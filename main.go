package main

import (
	"log"
	"wellnesspath/config"

	"wellnesspath/routes"
	seeder "wellnesspath/seeders"
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

	if config.ENV.Environment != "Production" {

		config.ResetEntireDatabase()

		if err := seeder.SeedExercisesFromFile(); err != nil {
			log.Fatal("Seeding failed:", err)
		}
	}

	// Initialize router
	router := routes.SetupRouter()

	// Start the server using the loaded ADDR and PORT
	log.Printf("Server is running on %s:%s", config.ENV.Addr, config.ENV.Port)
	router.Run(config.ENV.Addr + ":" + config.ENV.Port)
}
