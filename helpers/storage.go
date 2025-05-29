package helpers

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"
	"wellnesspath/config"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/sas"
)

// GenerateSASURL generates a signed URL (SAS Token) for a blob with read-only access and expiration time
func GenerateSASURL(filename string, expiry time.Duration) (string, error) {
	accountName := config.ENV.AzureStorageAccount
	accountKey := config.ENV.AzureStorageKey
	containerName := config.ENV.AzureContainerName

	cred, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		return "", fmt.Errorf("failed to create credentials: %w", err)
	}

	ext := strings.ToLower(filepath.Ext(filename))
	var blobPath string
	switch ext {
	case ".png", ".jpg", ".jpeg", ".gif":
		blobPath = fmt.Sprintf("picture/%s", filename)
	case ".mp4", ".mov", ".avi":
		blobPath = fmt.Sprintf("video/%s", filename)
	default:
		return "", fmt.Errorf("unsupported file extension: %s", ext)
	}

	// Fallback to default image if the requested blob does not exist
	if !BlobExists(filename) {
		blobPath = "picture/default.jpg"
	}

	expiresOn := time.Now().UTC().Add(expiry)
	permissions := sas.BlobPermissions{
		Read: true,
	}

	sasQueryParams, err := sas.BlobSignatureValues{
		ContainerName: containerName,
		BlobName:      blobPath,
		Permissions:   permissions.String(),
		StartTime:     time.Now().UTC(),
		ExpiryTime:    expiresOn,
	}.SignWithSharedKey(cred)
	if err != nil {
		return "", fmt.Errorf("failed to sign SAS: %w", err)
	}

	blobURL := fmt.Sprintf("https://%s.blob.core.windows.net/%s/%s?%s",
		accountName,
		containerName,
		blobPath,
		sasQueryParams.Encode())

	return blobURL, nil
}

// BlobExists checks whether a blob file exists in Azure Storage
func BlobExists(filename string) bool {
	accountName := config.ENV.AzureStorageAccount
	accountKey := config.ENV.AzureStorageKey
	containerName := config.ENV.AzureContainerName

	cred, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		return false
	}

	ext := strings.ToLower(filepath.Ext(filename))
	var blobPath string
	switch ext {
	case ".png", ".jpg", ".jpeg", ".gif":
		blobPath = fmt.Sprintf("picture/%s", filename)
	case ".mp4", ".mov", ".avi":
		blobPath = fmt.Sprintf("video/%s", filename)
	default:
		return false
	}

	blobClient, err := blob.NewClientWithSharedKeyCredential(
		fmt.Sprintf("https://%s.blob.core.windows.net/%s/%s", accountName, containerName, blobPath), cred, nil)
	if err != nil {
		return false
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = blobClient.GetProperties(ctx, nil)
	return err == nil
}
