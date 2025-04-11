package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"wctest/app/db"
	"wctest/app/model"
)

func getEmployees(w http.ResponseWriter, r *http.Request) {
	repo := db.NewEmployeeRepository()
	employees, err := repo.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(employees)
}

func main() {
	err := db.InitDB(true)
	if err != nil {
		log.Fatal(err)
	}

	repo := db.NewEmployeeRepository()
	
	err = repo.Cleanup()
	if err != nil {
		log.Fatal(err)
	}
	
	ceo := &model.Employee{
		FirstName: "Michael",
		LastName:  "Chen",
		Title:     "CEO",
	}
	err = repo.Create(ceo)
	if err != nil {
		log.Fatal(err)
	}

	cto := &model.Employee{
		FirstName: "Barrett",
		LastName:  "Glasauer",
		Title:     "CTO",
		ReportsTo: ceo,
	}
	err = repo.Create(cto)
	if err != nil {
		log.Fatal(err)
	}

	coo := &model.Employee{
		FirstName: "Andres",
		LastName:  "Green",
		Title:     "COO",
		ReportsTo: ceo,
	}
	err = repo.Create(coo)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/employees", getEmployees)
	
	fmt.Println("Server starting on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}