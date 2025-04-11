package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"wctest/app/config"
	"wctest/app/db"
	"wctest/app/service"
)

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		next.ServeHTTP(w, r)
	})
}

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

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

	// API handler
	apiHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

	// Serve static files
	fs := http.FileServer(http.Dir(filepath.Join("frontend", "build")))
	
	mux := http.NewServeMux()
	mux.Handle("/api/employees", enableCORS(apiHandler))
	mux.Handle("/", fs)

	addr := fmt.Sprintf(":%d", cfg.Port)
	log.Printf("Server starting on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}