package db

import (
	"database/sql"
	"sort"
	"wctest/app/model"
)

type EmployeeRepository struct {
	db *DB
}

func NewEmployeeRepository(db *DB) *EmployeeRepository {
	return &EmployeeRepository{
		db: db,
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

type EmployeeNode struct {
	Employee model.Employee `json:"employee"`
	Reports  []EmployeeNode `json:"reports"`
}

func (r *EmployeeRepository) GetEmployeeTree() ([]EmployeeNode, error) {
	employees, err := r.GetAll()
	if err != nil {
		return nil, err
	}

	employeeMap := make(map[int]*model.Employee)
	for i := range employees {
		employeeMap[employees[i].ID] = &employees[i]
	}

	reportsMap := make(map[int][]*model.Employee)
	for i := range employees {
		emp := &employees[i]
		if emp.ReportsTo != nil {
			managerID := emp.ReportsTo.ID
			reportsMap[managerID] = append(reportsMap[managerID], emp)
		}
	}

	for _, reports := range reportsMap {
		sort.Slice(reports, func(i, j int) bool {
			return reports[i].LastName < reports[j].LastName
		})
	}

	var roots []EmployeeNode
	for i := range employees {
		emp := &employees[i]
		// If the employee has no manager, they are at the root node
		if emp.ReportsTo == nil || employeeMap[emp.ReportsTo.ID] == nil {
			root := EmployeeNode{
				Employee: *emp,
				Reports:  buildReports(emp.ID, reportsMap),
			}
			roots = append(roots, root)
		}
	}

	return roots, nil
}

func buildReports(managerID int, reportsMap map[int][]*model.Employee) []EmployeeNode {
	reports := reportsMap[managerID]
	if reports == nil {
		return nil
	}

	var nodes []EmployeeNode
	for _, emp := range reports {
		node := EmployeeNode{
			Employee: *emp,
			Reports:  buildReports(emp.ID, reportsMap),
		}
		nodes = append(nodes, node)
	}
	return nodes
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