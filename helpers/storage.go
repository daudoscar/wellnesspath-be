package helpers

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
	"wellnesspath/config"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blockblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/sas"
)

// GenerateSASURL generates a signed URL (SAS Token) for a blob with read-only access and expiration time
func GenerateSASURL(filename string, expiry time.Duration) (string, error) {
	accountName := config.ENV.AzureStorageAccount
	accountKey := config.ENV.AzureStorageKey
	containerName := config.ENV.AzureContainerName
	environment := config.ENV.Environment // "local" atau "Hosted"

	// Tentukan ekstensi file dan path folder
	ext := strings.ToLower(filepath.Ext(filename))
	var blobPath string
	switch ext {
	case ".png", ".jpg", ".jpeg":
		var base string
		if ext == ".jpeg" {
			base = filename[:len(filename)-5]
		} else {
			base = filename[:len(filename)-4]
		}

		blobPaths := []string{
			fmt.Sprintf("%s.jpg", base),
			fmt.Sprintf("%s.png", base),
			fmt.Sprintf("%s.jpeg", base),
		}

		found := false
		var blobfound string

		for _, path := range blobPaths {
			fmt.Printf("Checking: %s\n", path)
			if BlobExists(path) {
				found = true
				blobfound = path
				break
			}
		}

		if found {
			blobPath = blobfound
		} else {
			blobPath = "images/placeholder.png"
			filename = "placeholder.png"
		}
	case ".mp4", ".mov", ".avi":
		if !BlobExists(filename) {
			blobPath = "videos/default.mp4"
			filename = "placeholder.png"
		} else {
			blobPath = filename
		}
	default:
		return "", fmt.Errorf("unsupported file extension: %s", ext)
	}

	// Fallback ke default image jika tidak ditemukan
	// if !BlobExists(filename) {
	// 	if filename[len(filename)-3:] == "mp4" {
	// 		blobPath = "videos/default.mp4"
	// 	} else {
	// 		blobPath = "images/placeholder.png"
	// 	}

	// 	filename = "placeholder.png"
	// }

	// Jika environment lokal → buat URL lokal tanpa SAS
	if strings.ToLower(environment) == "local" {
		url := fmt.Sprintf("http://127.0.0.1:10000/%s/%s?temp=true&exp=%d",
			containerName, blobPath, time.Now().Add(expiry).Unix())
		return url, nil
	}

	// Jika environment Hosted → buat SAS token valid
	cred, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		return "", fmt.Errorf("failed to create credentials: %w", err)
	}

	expiresOn := time.Now().UTC().Add(expiry)
	permissions := sas.BlobPermissions{Read: true}

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

	url := fmt.Sprintf("https://%s.blob.core.windows.net/%s/%s?%s",
		accountName, containerName, blobPath, sasQueryParams.Encode())

	return url, nil
}

func GenerateSASURLAds(filename string, expiry time.Duration) (string, error) {
	accountName := config.ENV.AzureStorageAccount
	accountKey := config.ENV.AzureStorageKey
	containerName := config.ENV.AzureContainerName

	// Jika environment Hosted → buat SAS token valid
	cred, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		return "", fmt.Errorf("failed to create credentials: %w", err)
	}

	expiresOn := time.Now().UTC().Add(expiry)
	permissions := sas.BlobPermissions{Read: true}

	sasQueryParams, err := sas.BlobSignatureValues{
		ContainerName: containerName,
		BlobName:      filename,
		Permissions:   permissions.String(),
		StartTime:     time.Now().UTC(),
		ExpiryTime:    expiresOn,
	}.SignWithSharedKey(cred)

	if err != nil {
		return "", fmt.Errorf("failed to sign SAS: %w", err)
	}

	url := fmt.Sprintf("https://%s.blob.core.windows.net/%s/%s?%s",
		accountName, containerName, filename, sasQueryParams.Encode())

	return url, nil
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

	blobPath := filename

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

// LOCAL

func GenerateLocalSASURL(container, blobName string, expiry time.Duration) (string, error) {
	cred, err := azblob.NewSharedKeyCredential("devstoreaccount1", "Eby8vdM02xNOcqFlqUwJPLlmEtlCDXJ1OUzFT50uSRZ6IFsuFq2UVErCz4I6tq/K1SZFPTOtr/KBHBeksoGMGw==")
	if err != nil {
		return "", err
	}

	sasParams, err := sas.BlobSignatureValues{
		Protocol:      sas.ProtocolHTTPSandHTTP,
		ContainerName: container,
		BlobName:      blobName,
		Permissions:   (&sas.BlobPermissions{Read: true}).String(),
		StartTime:     time.Now().UTC(),
		ExpiryTime:    time.Now().UTC().Add(expiry),
	}.SignWithSharedKey(cred)

	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("http://127.0.0.1:10000/devstoreaccount1/%s/%s?%s", container, blobName, sasParams.Encode())
	return url, nil
}

func UploadDefaultImageToAzurite() error {
	client := config.BlobClient
	containerName := config.ENV.AzureContainerName
	blobName := "images/default.jpg"
	localPath := filepath.Join("__blobstorage__", "images", "default.jpg")

	// Baca file dari local filesystem
	fileContent, err := os.ReadFile(localPath)
	if err != nil {
		return fmt.Errorf("failed to read local file: %w", err)
	}

	// Tentukan content type berdasarkan ekstensi
	contentType := getContentType(".jpg")

	ctx := context.Background()

	settings := blockblob.UploadBufferOptions{
		Concurrency: 1,
		HTTPHeaders: &blob.HTTPHeaders{
			BlobContentType: to.Ptr(contentType),
		},
	}

	// Upload ke Azurite
	_, err = client.UploadBuffer(ctx, containerName, blobName, fileContent, &settings)
	if err != nil {
		return fmt.Errorf("failed to upload blob: %w", err)
	}

	url := fmt.Sprintf("%s/%s/%s", strings.TrimRight(config.ENV.AzureStorageEndpoint, "/"), containerName, blobName)
	log.Printf("✅ Uploaded default.jpg to: %s", url)

	return nil
}

// Helper content-type
func getContentType(ext string) string {
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	default:
		return "application/octet-stream"
	}
}
