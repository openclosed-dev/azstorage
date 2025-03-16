package blob

import (
	"context"
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
)

type blobProcessor interface {
	processBlob(blobName string) error
	getTotalProcessed() (int, int)
}

type removingBlobProcessor struct {
	name            string
	logger          *log.Logger
	containerClient *container.Client
	context         context.Context
	successful      int
	failed          int
}

func newRemovingBlobProcessor(
	name string,
	containerClient *container.Client) *removingBlobProcessor {
	return &removingBlobProcessor{
		name:            name,
		logger:          log.New(os.Stdout, "["+name+"] ", 0),
		containerClient: containerClient,
		context:         context.Background(),
	}
}

func (p *removingBlobProcessor) processBlob(blobName string) error {
	p.logger.Printf("Deleting a blob \"%s\"\n", blobName)

	var blobClient = p.containerClient.NewBlobClient(blobName)
	_, err := blobClient.Delete(p.context, &blob.DeleteOptions{})
	if err != nil {
		p.logger.Printf("Failed to delete a blob \"%s\": %v", blobName, err)
		p.failed++
		return err
	}

	p.successful++
	p.logger.Printf("Deleted a blob \"%s\" successfully\n", blobName)

	return nil
}

func (p *removingBlobProcessor) getTotalProcessed() (int, int) {
	return p.successful, p.failed
}
