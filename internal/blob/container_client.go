package blob

import (
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/service"
)

type ContainerClient struct {
	accountName     string
	containerClient *container.Client
}

func NewContainerClient(accountName string, containerName string) (*ContainerClient, error) {

	serviceURL := fmt.Sprintf("https://%s.blob.core.windows.net/", accountName)

	var serviceClient *service.Client

	accountKey, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_KEY")
	if ok {

		credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
		if err != nil {
			return nil, fmt.Errorf("authentication failure with account key: %w", err)
		}

		serviceClient, err = service.NewClientWithSharedKeyCredential(serviceURL, credential, nil)
		if err != nil {
			return nil, err
		}

	} else {

		credential, err := azidentity.NewDefaultAzureCredential(nil)
		if err != nil {
			return nil, err
		}

		serviceClient, err = service.NewClient(serviceURL, credential, nil)
		if err != nil {
			return nil, err
		}
	}

	containerClient := serviceClient.NewContainerClient(containerName)

	return &ContainerClient{
		accountName:     accountName,
		containerClient: containerClient,
	}, nil
}

func (client *ContainerClient) RemoveBlobsInList(listFile string, walkers int, processors int) error {
	var job = newBlobRemovingJob(walkers, processors, client.containerClient)
	return job.doJob(listFile)
}
