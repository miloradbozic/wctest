package main

import (
	"fmt"
	"log"
	"wctest/app/config"
	"wctest/app/db"
	"wctest/app/server"
	"wctest/app/service"
)

func main() {
	// Load configuration
	cfg := config.NewConfig()

	// Initialize database
	database, err := db.InitDB(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	// Create repository
	repo := db.NewEmployeeRepository(database)

	// Create service
	employeeService := service.NewEmployeeService(repo, cfg)

	// Initialize data if needed
	if err := employeeService.InitializeData(); err != nil {
		log.Fatalf("Failed to initialize data: %v", err)
	}

	// Create and start server
	srv := server.NewServer(employeeService)
	addr := fmt.Sprintf(":%d", cfg.Port)
	log.Printf("Server starting on %s", addr)
	if err := srv.Start(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}