package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"wctest/app/model"
)

func getEmployees(w http.ResponseWriter, r *http.Request) {
	
	ceo := &model.Employee{
		FirstName: "Michael",
		LastName:  "Chen",
		Title:     "CEO",
	}
	cto := &model.Employee{
		FirstName: "Barrett",
		LastName:  "Glasauer",
		Title:     "CTO",
		ReportsTo: ceo,
	}

	coo := &model.Employee{
		FirstName: "Andres",
		LastName:  "Green",
		Title:     "COO",
		ReportsTo: ceo,
	}

	employees := []model.Employee{
		*ceo,
		*cto,
		*coo,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(employees)
}

func main() {
	http.HandleFunc("/employees", getEmployees)
	
	fmt.Println("Server starting on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
