package blob

import (
	"fmt"
	"log"
	"sync"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
)

type blobRemovingJob struct {
	containerClient *container.Client
	dirs            chan string
	blobs           chan string
	walkerGroup     sync.WaitGroup
	processorGroup  sync.WaitGroup
	walkers         []*directoryWalker
	processors      []blobProcessor
}

func RemoveBlobsInList(
	accountName string,
	containerName string,
	listFile string,
	walkers int,
	processors int) error {

	client, err := newContainerClient(accountName, containerName)
	if err != nil {
		return err
	}

	var job = blobRemovingJob{
		containerClient: client,
		dirs:            make(chan string, walkers),
		blobs:           make(chan string, processors),
		walkers:         make([]*directoryWalker, 0, walkers),
		processors:      make([]blobProcessor, 0, processors),
	}

	return job.doJob(listFile)
}

func (job *blobRemovingJob) doJob(listFile string) error {

	parser, err := newListParser(listFile, job)
	if err != nil {
		return err
	}
	defer parser.close()

	defer job.waitForCompletion()

	job.startProcessors(cap(job.processors))
	job.startWalkers(cap(job.walkers))

	return parser.parseAll()
}

func (job *blobRemovingJob) waitForCompletion() {
	job.close()
	job.printSummary()
}

func (job *blobRemovingJob) handleDirectory(path string) {
	job.dirs <- path
}

func (job *blobRemovingJob) handleBlob(path string) {
	job.blobs <- path
}

func (job *blobRemovingJob) startProcessors(count int) {

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

func (job *blobRemovingJob) startWalkers(count int) {

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

func (job *blobRemovingJob) close() {

	// Waits for directory walkers to complete
	close(job.dirs)
	job.walkerGroup.Wait()

	// Waits for blob walkers to complete
	close(job.blobs)
	job.processorGroup.Wait()
}

func (job *blobRemovingJob) printSummary() {
	stats := job.getStats()
	log.Printf("Summary: blobs processed: %d, successful: %d, failed: %d\n", stats.total, stats.successful, stats.failed)
}

func (job *blobRemovingJob) getStats() processorStats {

	var sum = processorStats{}

	for _, processor := range job.processors {
		stats := processor.getStats()
		sum.total += stats.total
		sum.successful += stats.successful
		sum.failed += stats.failed
	}

	return sum
}
