package blob

import (
	"context"
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
)

type blobHandler interface {
	handleBlob(path string)
}

type directoryWalker struct {
	name            string
	logger          *log.Logger
	containerClient *container.Client
	context         context.Context
	handler         blobHandler
	totalFound      int
}

func newDirectoryWalker(
	name string,
	containerClient *container.Client,
	handler blobHandler) *directoryWalker {
	return &directoryWalker{
		name:            name,
		logger:          log.New(os.Stdout, "["+name+"] ", 0),
		containerClient: containerClient,
		context:         context.Background(),
		handler:         handler,
	}
}

func (w *directoryWalker) walk(dir string) error {
	w.logger.Printf("Starting directory \"%s\"\n", dir)

	var options = azblob.ListBlobsFlatOptions{
		Include: container.ListBlobsInclude{Deleted: false, Versions: false},
	}

	if dir != "/" {
		options.Prefix = &dir
	}

	pager := w.containerClient.NewListBlobsFlatPager(&options)

	for pager.More() {
		resp, err := pager.NextPage(w.context)
		if err != nil {
			w.logger.Printf("Error has occurred while walking \"%s\": %v", dir, err)
			return err
		}
		for _, blob := range resp.Segment.BlobItems {
			w.totalFound++
			w.handler.handleBlob(*blob.Name)
		}
	}

	w.logger.Printf("Finished directory \"%s\"\n", dir)

	return nil
}

func (w *directoryWalker) getTotalFound() int {
	return w.totalFound
}
