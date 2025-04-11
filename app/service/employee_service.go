package service

import (
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

func (s *EmployeeService) InitializeSampleData() error {
	if err := s.repo.Cleanup(); err != nil {
		return err
	}

	// Create CEO
	ceo := &model.Employee{
		FirstName: "Michael",
		LastName:  "Chen",
		Title:     "CEO",
	}
	if err := s.repo.Create(ceo); err != nil {
		return err
	}

	// Create CTO
	cto := &model.Employee{
		FirstName: "Barrett",
		LastName:  "Glasauer",
		Title:     "CTO",
		ReportsTo: ceo,
	}
	if err := s.repo.Create(cto); err != nil {
		return err
	}

	// Create COO
	coo := &model.Employee{
		FirstName: "Andres",
		LastName:  "Green",
		Title:     "COO",
		ReportsTo: ceo,
	}
	if err := s.repo.Create(coo); err != nil {
		return err
	}

	return nil
} 