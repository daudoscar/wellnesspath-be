package main

import (
	"log"
	"wellnesspath/config"
	"wellnesspath/routes"
	// seeder "wellnesspath/seeders"
)

func main() {
	// Load configuration
	config.LoadConfig()

	// Connect to the database
	config.ConnectDatabase()

	// Initialize router
	router := routes.SetupRouter()

	// if err := seeder.SeedExercisesFromAzure(); err != nil {
	// 	log.Fatal("Seeding failed:", err)
	// }

	// Start the server using the loaded ADDR and PORT
	log.Printf("Server is running on %s:%s", config.ENV.Addr, config.ENV.Port)
	router.Run(config.ENV.Addr + ":" + config.ENV.Port)
}
