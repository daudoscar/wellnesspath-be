package config

import (
	"context"
	"fmt"
	"log"
	"wellnesspath/models"

	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	if ENV == nil {
		LoadConfig()
	}

	// SQL Server DSN format
	dsn := fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s",
		ENV.DBUser, ENV.DBPassword, ENV.DBHost, ENV.DBPort, ENV.DBName)

	var err error
	DB, err = gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	log.Println("✅ Database connection successful")

	err = DB.AutoMigrate(
		&models.User{},
		&models.Profile{},
		&models.Exercise{},
		&models.WorkoutPlan{},
		&models.WorkoutPlanDay{},
		&models.WorkoutPlanExercise{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate database tables: %v", err)
	}
	log.Println("✅ Database tables migrated successfully")

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

	log.Println("✅ Azure Blob Storage connection successful.")
	return nil
}
