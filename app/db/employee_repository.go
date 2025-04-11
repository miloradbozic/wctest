package db

import (
	"database/sql"
	"wctest/app/model"
)

type EmployeeRepository struct {
	db *sql.DB
}

func NewEmployeeRepository() *EmployeeRepository {
	return &EmployeeRepository{
		db: DB,
	}
}

func (r *EmployeeRepository) Cleanup() error {
	_, err := r.db.Exec("DELETE FROM employees")
	return err
}

func (r *EmployeeRepository) Create(employee *model.Employee) error {
	query := `
		INSERT INTO employees (first_name, last_name, title, reports_to_id)
		VALUES (?, ?, ?, ?)
	`

	var reportsToID *int
	if employee.ReportsTo != nil {
		reportsToID = &employee.ReportsTo.ID
	}

	result, err := r.db.Exec(query, employee.FirstName, employee.LastName, employee.Title, reportsToID)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	employee.ID = int(id)
	return nil
}

func (r *EmployeeRepository) GetAll() ([]model.Employee, error) {
	query := `
		SELECT e.id, e.first_name, e.last_name, e.title, 
		       m.id, m.first_name, m.last_name, m.title
		FROM employees e
		LEFT JOIN employees m ON e.reports_to_id = m.id
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var employees []model.Employee
	employeeMap := make(map[int]*model.Employee)

	for rows.Next() {
		var emp model.Employee
		var managerID sql.NullInt64
		var managerFirstName, managerLastName, managerTitle sql.NullString

		err := rows.Scan(
			&emp.ID,
			&emp.FirstName,
			&emp.LastName,
			&emp.Title,
			&managerID,
			&managerFirstName,
			&managerLastName,
			&managerTitle,
		)
		if err != nil {
			return nil, err
		}

		if managerID.Valid {
			manager := &model.Employee{
				ID:        int(managerID.Int64),
				FirstName: managerFirstName.String,
				LastName:  managerLastName.String,
				Title:     managerTitle.String,
			}
			emp.ReportsTo = manager
		}

		employeeMap[emp.ID] = &emp
		employees = append(employees, emp)
	}

	return employees, nil
}

func (r *EmployeeRepository) GetByID(id int) (*model.Employee, error) {
	query := `
		SELECT e.id, e.first_name, e.last_name, e.title,
		       m.id, m.first_name, m.last_name, m.title
		FROM employees e
		LEFT JOIN employees m ON e.reports_to_id = m.id
		WHERE e.id = ?
	`

	var emp model.Employee
	var managerID sql.NullInt64
	var managerFirstName, managerLastName, managerTitle sql.NullString

	err := r.db.QueryRow(query, id).Scan(
		&emp.ID,
		&emp.FirstName,
		&emp.LastName,
		&emp.Title,
		&managerID,
		&managerFirstName,
		&managerLastName,
		&managerTitle,
	)
	if err != nil {
		return nil, err
	}

	if managerID.Valid {
		manager := &model.Employee{
			ID:        int(managerID.Int64),
			FirstName: managerFirstName.String,
			LastName:  managerLastName.String,
			Title:     managerTitle.String,
		}
		emp.ReportsTo = manager
	}

	return &emp, nil
} 