package blob

import (
	"context"
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
)

type blobHandler interface {
	handleBlob(path string)
}

type directoryWalker struct {
	name            string
	logPrefix       string
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
		logPrefix:       fmt.Sprintf("[%s]", name),
		containerClient: containerClient,
		context:         context.Background(),
		handler:         handler,
	}
}

func (w *directoryWalker) walk(dir string) error {
	log.Printf("%s Starting directory \"%s\"\n", w.logPrefix, dir)

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
			log.Printf("%s Error has occurred while walking in directory \"%s\": %v", w.logPrefix, dir, err)
			return err
		}
		for _, blob := range resp.Segment.BlobItems {
			w.totalFound++
			w.handler.handleBlob(*blob.Name)
		}
	}

	log.Printf("%s Finished directory \"%s\"\n", w.logPrefix, dir)

	return nil
}

func (w *directoryWalker) getTotalFound() int {
	return w.totalFound
}
