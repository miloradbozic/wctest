package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"wctest/app/config"
	"wctest/app/db"
	"wctest/app/service"
)

func main() {
	cfg := config.NewConfig()
	database, err := db.InitDB(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	repo := db.NewEmployeeRepository(database)
	employeeService := service.NewEmployeeService(repo, cfg)
	if err := employeeService.InitializeData(); err != nil {
		log.Fatalf("Failed to initialize data: %v", err)
	}

	http.HandleFunc("/employees", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		tree, err := employeeService.GetOrganizationTree()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(tree); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	addr := fmt.Sprintf(":%d", cfg.Port)
	log.Printf("Server starting on %s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}