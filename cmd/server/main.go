package main

import (
	"log"
	"wctest/app/config"
	"wctest/app/db"
	"wctest/app/server"
	"wctest/app/service"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize database
	database, err := db.InitDB(cfg)
	if err != nil {
		log.Fatal(err)
	}


	repo := db.NewEmployeeRepository(database)
	employeeService := service.NewEmployeeService(repo)
	if err := employeeService.InitializeSampleData(); err != nil {
		log.Fatal(err)
	}

	srv := server.NewServer(employeeService)
	log.Printf("Server starting on port %d...", cfg.Port)
	log.Fatal(srv.Start(cfg.Port))
}