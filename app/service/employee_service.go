package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"wctest/app/db"
	"wctest/app/model"
)

type EmployeeService struct {
	repo *db.EmployeeRepository
}

func NewEmployeeService(repo *db.EmployeeRepository) *EmployeeService {
	return &EmployeeService{
		repo: repo,
	}
}

func (s *EmployeeService) GetOrganizationTree() ([]db.EmployeeNode, error) {
	return s.repo.GetEmployeeTree()
}

func (s *EmployeeService) InitializeData() error {
	isEmpty, err := s.repo.IsEmpty()
	if err != nil {
		return err
	}

	if !isEmpty {
		return nil // Database already has data
	}

	resp, err := http.Get("https://gist.githubusercontent.com/chancock09/6d2a5a4436dcd488b8287f3e3e4fc73d/raw/fa47d64c6d5fc860fabd3033a1a4e3c59336324e/employees.json")
	if err != nil {
		return fmt.Errorf("failed to fetch data: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %v", err)
	}

	var employees []model.Employee
	if err := json.Unmarshal(body, &employees); err != nil {
		return fmt.Errorf("failed to parse JSON: %v", err)
	}

	for _, emp := range employees {
		if err := s.repo.Create(&emp); err != nil {
			return fmt.Errorf("failed to create employee %s: %v", emp.Name, err)
		}
	}

	return nil
} 