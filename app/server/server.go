package server

import (
	"encoding/json"
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

func (s *Server) getEmployees(w http.ResponseWriter, r *http.Request) {
	employees, err := s.employeeService.GetOrganizationTree()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(employees)
}

func (s *Server) Start(port int) error {
	http.HandleFunc("/employees", s.getEmployees)
	return http.ListenAndServe(":8080", nil)
} 