package blob

import (
	"context"
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/service"
)

func newServiceClient(accountName string) (*service.Client, error) {

	serviceURL := fmt.Sprintf("https://%s.blob.core.windows.net/", accountName)

	var client *service.Client
	accountKey, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_KEY")
	if ok {

		credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
		if err != nil {
			return nil, fmt.Errorf("authentication failure with account key: %w", err)
		}

		client, err = service.NewClientWithSharedKeyCredential(serviceURL, credential, nil)
		if err != nil {
			return nil, err
		}

	} else {

		credential, err := azidentity.NewDefaultAzureCredential(nil)
		if err != nil {
			return nil, err
		}

		client, err = service.NewClient(serviceURL, credential, nil)
		if err != nil {
			return nil, err
		}
	}

	if err := verifyServiceClient(client); err != nil {
		return nil, fmt.Errorf("failed to access the blob service: %w", err)
	}

	return client, nil
}

func newContainerClient(accountName string, containerName string) (*container.Client, error) {

	serviceClient, err := newServiceClient(accountName)
	if err != nil {
		return nil, err
	}

	client := serviceClient.NewContainerClient(containerName)

	if err := verifyContainerClient(client); err != nil {
		return nil, fmt.Errorf("failed to access the blob container: %w", err)
	}

	return client, nil
}

func verifyServiceClient(client *service.Client) error {
	_, err := client.GetProperties(context.Background(), nil)
	return err
}

func verifyContainerClient(client *container.Client) error {
	_, err := client.GetProperties(context.Background(), nil)
	return err
}
