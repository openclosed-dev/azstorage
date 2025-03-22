package blob

import (
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/service"
)

func newServiceClient(accountName string) (*service.Client, error) {

	serviceURL := fmt.Sprintf("https://%s.blob.core.windows.net/", accountName)

	accountKey, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_KEY")
	if ok {

		credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
		if err != nil {
			return nil, fmt.Errorf("authentication failure with account key: %w", err)
		}

		return service.NewClientWithSharedKeyCredential(serviceURL, credential, nil)

	} else {

		credential, err := azidentity.NewDefaultAzureCredential(nil)
		if err != nil {
			return nil, err
		}

		return service.NewClient(serviceURL, credential, nil)
	}
}

func newContainerClient(accountName string, containerName string) (*container.Client, error) {

	serviceClient, err := newServiceClient(accountName)
	if err != nil {
		return nil, err
	}

	return serviceClient.NewContainerClient(containerName), nil
}
