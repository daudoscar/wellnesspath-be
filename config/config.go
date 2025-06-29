package config

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/spf13/viper"
)

type Config struct {
	DBUser               string
	DBPassword           string
	DBHost               string
	DBPort               string
	DBName               string
	JWTSecret            string
	Addr                 string
	Port                 string
	AccessTTL            time.Duration
	RefreshTTL           time.Duration
	AzureStorageAccount  string
	AzureStorageKey      string
	AzureContainerName   string
	AzureStorageEndpoint string
	Environment          string
	Queue                string
}

var (
	ENV        *Config
	BlobClient *azblob.Client
)

func LoadConfig() *Config {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Println("No .env file found, reading environment variables")
	}

	accessTTL, err := time.ParseDuration(viper.GetString("ACCESS_TTL"))
	if err != nil {
		accessTTL = 24 * time.Hour
	}

	refreshTTL, err := time.ParseDuration(viper.GetString("REFRESH_TTL"))
	if err != nil {
		refreshTTL = 7 * 24 * time.Hour
	}

	ENV = &Config{
		DBUser:               viper.GetString("DB_USER"),
		DBPassword:           viper.GetString("DB_PASSWORD"),
		DBHost:               viper.GetString("DB_HOST"),
		DBPort:               viper.GetString("DB_PORT"),
		DBName:               viper.GetString("DB_NAME"),
		JWTSecret:            viper.GetString("JWT_SECRET"),
		Addr:                 viper.GetString("ADDR"),
		Port:                 viper.GetString("PORT"),
		AccessTTL:            accessTTL,
		RefreshTTL:           refreshTTL,
		AzureStorageAccount:  viper.GetString("AZURE_STORAGE_ACCOUNT_NAME"),
		AzureStorageKey:      viper.GetString("AZURE_STORAGE_ACCOUNT_KEY"),
		AzureContainerName:   viper.GetString("AZURE_STORAGE_CONTAINER_NAME"),
		AzureStorageEndpoint: viper.GetString("AZURE_STORAGE_ENDPOINT"),
		Environment:          viper.GetString("ENVIRONMENT"),
		Queue:                viper.GetString("QUEUE"),
	}

	err = InitBlobClient()
	if err != nil {
		log.Fatalf("Failed to initialize Azure Blob Storage client: %v", err)
	}

	return ENV
}

func InitBlobClient() error {
	// 1. Inisialisasi kredensial
	cred, err := azblob.NewSharedKeyCredential(ENV.AzureStorageAccount, ENV.AzureStorageKey)
	if err != nil {
		return fmt.Errorf("failed to create Azure shared key credential: %v", err)
	}

	// 2. Buat blob client
	BlobClient, err = azblob.NewClientWithSharedKeyCredential(ENV.AzureStorageEndpoint, cred, nil)
	if err != nil {
		return fmt.Errorf("failed to create blob client: %v", err)
	}

	log.Println("Azure Blob Storage client initialized successfully.")

	// 3. Auto-create container jika belum ada
	ctx := context.Background()
	containerClient := BlobClient.ServiceClient().NewContainerClient(ENV.AzureContainerName)
	_, err = containerClient.Create(ctx, nil)
	if err != nil {
		// Jika container sudah ada, abaikan error
		if !strings.Contains(err.Error(), "ContainerAlreadyExists") {
			return fmt.Errorf("failed to create container '%s': %v", ENV.AzureContainerName, err)
		}
	}

	log.Printf("✅ Container '%s' is ready.\n", ENV.AzureContainerName)
	return nil
}
