package main

import (
	"log"
	"wellnesspath/config"

	"wellnesspath/helpers"
	"wellnesspath/routes"
)

func main() {
	// Load configuration
	config.LoadConfig()

	// Connect to the database
	if config.ENV.Environment == "Hosted" {
		config.ConnectDatabase()
	} else {
		config.ConnectDatabase()
	}

	// if config.ENV.Environment != "Production" {
	// 	config.ResetEntireDatabase()
	// }

	err := helpers.UploadDefaultImageToAzurite()
	if err != nil {
		log.Fatalf("❌ Upload failed: %v", err)
	}

	// Initialize router
	router := routes.SetupRouter()

	// if err := seeder.SeedExercisesFromFile(); err != nil {
	// 	log.Fatal("Seeding failed:", err)
	// }

	// Start the server using the loaded ADDR and PORT
	log.Printf("Server is running on %s:%s", config.ENV.Addr, config.ENV.Port)
	router.Run(config.ENV.Addr + ":" + config.ENV.Port)
}
