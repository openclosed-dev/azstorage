package blob

import (
	"context"
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
)

type processorStats struct {
	total      int
	successful int
	failed     int
}

type blobProcessor interface {
	processBlob(blobName string) error
	getStats() processorStats
}

type removingBlobProcessor struct {
	name            string
	logPrefix       string
	containerClient *container.Client
	context         context.Context
	total           int
	successful      int
	failed          int
}

func newRemovingBlobProcessor(
	name string,
	containerClient *container.Client) *removingBlobProcessor {
	return &removingBlobProcessor{
		name:            name,
		logPrefix:       fmt.Sprintf("[%s]", name),
		containerClient: containerClient,
		context:         context.Background(),
	}
}

func (p *removingBlobProcessor) processBlob(blobName string) error {

	log.Printf("%s Deleting a blob \"%s\"\n", p.logPrefix, blobName)
	p.total++

	var blobClient = p.containerClient.NewBlobClient(blobName)
	_, err := blobClient.Delete(p.context, &blob.DeleteOptions{})
	if err != nil {
		log.Printf("%s Failed to delete a blob \"%s\": %v", p.logPrefix, blobName, err)
		p.failed++
		return err
	}

	log.Printf("%s Deleted a blob \"%s\" successfully\n", p.logPrefix, blobName)
	p.successful++

	return nil
}

func (p *removingBlobProcessor) getStats() processorStats {
	return processorStats{
		p.total, p.successful, p.failed,
	}
}
