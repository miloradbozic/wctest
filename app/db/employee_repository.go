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
		INSERT INTO employees (id, name, title, manager_id)
		VALUES (?, ?, ?, ?)
	`

	_, err := r.db.Exec(query, 
		employee.ID,
		employee.Name, 
		employee.Title, 
		employee.ManagerID,
	)
	return err
}

// EmployeeNode represents an employee in the tree structure
type EmployeeNode struct {
	Employee model.Employee `json:"employee"`
	Reports  []EmployeeNode `json:"reports"`
}

func (r *EmployeeRepository) GetEmployeeTree() ([]EmployeeNode, error) {
	// Get all employees
	employees, err := r.GetAll()
	if err != nil {
		return nil, err
	}

	// Create a map of employees by ID
	employeeMap := make(map[int]*model.Employee)
	for i := range employees {
		employeeMap[employees[i].ID] = &employees[i]
	}

	// Create a map of reports for each employee
	reportsMap := make(map[int][]*model.Employee)
	for i := range employees {
		emp := &employees[i]
		if emp.ManagerID != nil {
			managerID := *emp.ManagerID
			reportsMap[managerID] = append(reportsMap[managerID], emp)
		}
	}

	// Sort reports by name
	for _, reports := range reportsMap {
		sort.Slice(reports, func(i, j int) bool {
			return reports[i].Name < reports[j].Name
		})
	}

	// Build the tree structure
	var roots []EmployeeNode
	for i := range employees {
		emp := &employees[i]
		// If the employee has no manager, they are a root
		if emp.ManagerID == nil {
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
		SELECT id, name, title, manager_id
		FROM employees
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var employees []model.Employee
	for rows.Next() {
		var emp model.Employee
		var managerID sql.NullInt64

		err := rows.Scan(
			&emp.ID,
			&emp.Name,
			&emp.Title,
			&managerID,
		)
		if err != nil {
			return nil, err
		}

		if managerID.Valid {
			id := int(managerID.Int64)
			emp.ManagerID = &id
		}

		employees = append(employees, emp)
	}

	return employees, nil
}

func (r *EmployeeRepository) GetByID(id int) (*model.Employee, error) {
	query := `
		SELECT id, name, title, manager_id
		FROM employees
		WHERE id = ?
	`

	var emp model.Employee
	var managerID sql.NullInt64

	err := r.db.QueryRow(query, id).Scan(
		&emp.ID,
		&emp.Name,
		&emp.Title,
		&managerID,
	)
	if err != nil {
		return nil, err
	}

	if managerID.Valid {
		id := int(managerID.Int64)
		emp.ManagerID = &id
	}

	return &emp, nil
}

func (r *EmployeeRepository) IsEmpty() (bool, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM employees").Scan(&count)
	if err != nil {
		return false, err
	}
	return count == 0, nil
} 