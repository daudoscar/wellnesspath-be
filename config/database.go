package config

import (
	"context"
	"fmt"
	"log"
	"syncspend/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	if ENV == nil {
		LoadConfig()
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		ENV.DBUser, ENV.DBPassword, ENV.DBHost, ENV.DBPort, ENV.DBName)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	log.Println("Database connection successful")

	err = DB.AutoMigrate(
		&models.User{},
		&models.Exercise{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate database tables: %v", err)
	}
	log.Println("Database tables migrated successfully")

	err = TestBlobStorageConnection()
	if err != nil {
		log.Fatalf("Failed to connect to Azure Blob Storage: %v", err)
	}
}

func TestBlobStorageConnection() error {
	client := BlobClient
	containerName := ENV.AzureContainerName

	ctx := context.Background()
	pager := client.NewListBlobsFlatPager(containerName, nil)

	if pager.More() {
		_, err := pager.NextPage(ctx)
		if err != nil {
			return err
		}
	}

	log.Println("Azure Blob Storage connection successful.")
	return nil
}
