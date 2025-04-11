package server

import (
	"encoding/json"
	"log"
	"net/http"
	"wctest/app/service"
)

type Server struct {
	employeeService *service.EmployeeService
}

func NewServer(employeeService *service.EmployeeService) *Server {
	return &Server{
		employeeService: employeeService,
	}
}

func (s *Server) Start(addr string) error {
	http.HandleFunc("/employees", s.handleGetEmployees)
	log.Printf("Server starting on %s", addr)
	return http.ListenAndServe(addr, nil)
}

func (s *Server) handleGetEmployees(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	tree, err := s.employeeService.GetOrganizationTree()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(tree); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
} 