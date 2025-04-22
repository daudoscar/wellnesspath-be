package main

import (
	"log"
	"wellnesspath/config"
	"wellnesspath/routes"
)

func main() {
	// Load configuration
	config.LoadConfig()

	// Connect to the database
	config.ConnectDatabase()

	// Initialize router
	router := routes.SetupRouter()

	// Start the server using the loaded ADDR and PORT
	log.Printf("Server is running on %s:%s", config.ENV.Addr, config.ENV.Port)
	router.Run(config.ENV.Addr + ":" + config.ENV.Port)
}
