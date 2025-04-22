package helpers

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"path/filepath"
	"strings"

	"syncspend/config"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blockblob"
)

func UploadProfileImage(header *multipart.FileHeader, id int) (string, error) {
	client := config.BlobClient
	containerName := config.ENV.AzureContainerName
	folder := "profile"

	ext := strings.ToLower(filepath.Ext(header.Filename))
	blobName := fmt.Sprintf("%s/profile_%d%s", folder, id, ext)

	file, err := OpenFileFromMultipartHeader(header)
	if err != nil {
		return "", fmt.Errorf("failed to open profile image file: %v", err)
	}
	defer file.Close()

	fileContent := make([]byte, header.Size)
	_, err = file.Read(fileContent)
	if err != nil {
		return "", fmt.Errorf("failed to read file content: %v", err)
	}

	contentType := getContentTypeByExtension(ext)

	ctx := context.Background()

	settings := blockblob.UploadBufferOptions{
		Concurrency: 1,
		HTTPHeaders: &blob.HTTPHeaders{
			BlobContentType: to.Ptr(contentType),
		},
	}

	_, err = client.UploadBuffer(ctx, containerName, blobName, fileContent, &settings)
	if err != nil {
		return "", fmt.Errorf("failed to upload blob: %v", err)
	}

	url := fmt.Sprintf("%s/%s/%s", strings.TrimRight(config.ENV.AzureStorageEndpoint, "/"), containerName, blobName)
	log.Printf("Image uploaded successfully to: %s", url)

	return url, nil
}

func OpenFileFromMultipartHeader(header *multipart.FileHeader) (multipart.File, error) {
	file, err := header.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	return file, nil
}

func getContentTypeByExtension(ext string) string {
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".bmp":
		return "image/bmp"
	case ".webp":
		return "image/webp"
	default:
		return "application/octet-stream"
	}
}
