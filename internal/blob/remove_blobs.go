package blob

import (
	"fmt"
	"sync"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
)

type removingJob struct {
	containerClient *container.Client
	dirs            chan string
	blobs           chan string
	walkerGroup     sync.WaitGroup
	processorGroup  sync.WaitGroup
	walkers         []*directoryWalker
	processors      []blobProcessor
}

func (client *ContainerClient) RemoveBlobsInList(listFile string, walkers int, processors int) error {

	var job = removingJob{
		containerClient: client.containerClient,
		dirs:            make(chan string, walkers),
		blobs:           make(chan string, processors),
		walkers:         make([]*directoryWalker, 0, walkers),
		processors:      make([]blobProcessor, 0, processors),
	}

	defer func() {
		job.close()
		job.printSummary()
	}()

	return job.doJob(listFile, walkers, processors)
}

func (job *removingJob) doJob(listFile string, walkers int, processors int) error {

	parser, err := newListParser(listFile, job)
	if err != nil {
		return err
	}
	defer parser.close()

	job.startProcessors(processors)
	job.startWalkers(walkers)

	return parser.parseAll()
}

func (job *removingJob) handleDirectory(path string) {
	job.dirs <- path
}

func (job *removingJob) handleBlob(path string) {
	job.blobs <- path
}

func (job *removingJob) startProcessors(count int) {

	for i := range count {
		var name = fmt.Sprintf("processor-%03d", i+1)
		var processor blobProcessor = newRemovingBlobProcessor(name, job.containerClient)
		job.processors = append(job.processors, processor)
		job.processorGroup.Add(1)
		go func() {
			defer job.processorGroup.Done()
			for blob := range job.blobs {
				processor.processBlob(blob)
			}
		}()
	}
}

func (job *removingJob) startWalkers(count int) {

	for i := range count {
		var name = fmt.Sprintf("walker-%03d", i+1)
		var walker = newDirectoryWalker(name, job.containerClient, job)
		job.walkers = append(job.walkers, walker)
		job.walkerGroup.Add(1)
		go func() {
			defer job.walkerGroup.Done()
			for dir := range job.dirs {
				walker.walk(dir)
			}
		}()
	}
}

func (job *removingJob) close() {

	// Waits for directory walkers to complete
	close(job.dirs)
	job.walkerGroup.Wait()

	// Waits for blob walkers to complete
	close(job.blobs)
	job.processorGroup.Wait()
}

func (job *removingJob) printSummary() {
	found := job.getFoundBlobs()
	successful, failed := job.getProcessedBlobs()

	fmt.Println()
	fmt.Printf("Summary: blobs found: %d, successful: %d, failed: %d\n", found, successful, failed)
}

func (job *removingJob) getFoundBlobs() int {
	var total = 0
	for _, walker := range job.walkers {
		total += walker.getTotalFound()
	}
	return total
}

func (job *removingJob) getProcessedBlobs() (int, int) {
	var totalSuccessful = 0
	var totalFailed = 0
	for _, processor := range job.processors {
		successful, failed := processor.getTotalProcessed()
		totalSuccessful += successful
		totalFailed += failed
	}
	return totalSuccessful, totalFailed
}
