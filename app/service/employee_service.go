package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"wctest/app/config"
	"wctest/app/db"
	"wctest/app/model"
)

type EmployeeService struct {
	repo *db.EmployeeRepository
	cfg  *config.Config
}

func NewEmployeeService(repo *db.EmployeeRepository, cfg *config.Config) *EmployeeService {
	return &EmployeeService{
		repo: repo,
		cfg:  cfg,
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

	resp, err := http.Get(s.cfg.EmployeesURL)
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