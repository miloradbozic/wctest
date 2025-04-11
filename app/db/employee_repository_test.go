package db_test

import (
	"testing"
	"wctest/app/db"
	"wctest/app/model"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestEmployeeRepository_Create(t *testing.T) {
	tests := []struct {
		name        string
		employee    *model.Employee
		expectedErr error
	}{
		{
			name: "successful creation",
			employee: &model.Employee{
				ID:        1,
				Name:      "John Johnson",
				Title:     "CEO",
				ManagerID: nil,
			},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer mockDB.Close()

			// Create repository with mock DB
			repo := db.NewEmployeeRepository(&db.DB{DB: mockDB})

			// Set up expectations
			mock.ExpectExec("INSERT INTO employees").
				WithArgs(tt.employee.ID, tt.employee.Name, tt.employee.Title, tt.employee.ManagerID).
				WillReturnResult(sqlmock.NewResult(1, 1))

			err = repo.Create(tt.employee)

			assert.Equal(t, tt.expectedErr, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestEmployeeRepository_GetEmployeeTree(t *testing.T) {
	tests := []struct {
		name          string
		rows          *sqlmock.Rows
		expectedTree  []db.EmployeeNode
		expectedError error
	}{
		{
			name: "single employee",
			rows: sqlmock.NewRows([]string{"id", "name", "title", "manager_id"}).
				AddRow(1, "George Johnson", "CEO", nil),
			expectedTree: []db.EmployeeNode{
				{
					Employee: model.Employee{
						ID:        1,
						Name:      "George Johnson",
						Title:     "CEO",
						ManagerID: nil,
					},
					Reports: []db.EmployeeNode{},
				},
			},
			expectedError: nil,
		},
		{
			name: "employee with reports",
			rows: sqlmock.NewRows([]string{"id", "name", "title", "manager_id"}).
				AddRow(1, "George Johnson", "CEO", nil).
				AddRow(2, "Will Bronson", "CTO", 1),
			expectedTree: []db.EmployeeNode{
				{
					Employee: model.Employee{
						ID:        1,
						Name:      "George Johnson",
						Title:     "CEO",
						ManagerID: nil,
					},
					Reports: []db.EmployeeNode{
						{
							Employee: model.Employee{
								ID:        2,
								Name:      "Will Bronson",
								Title:     "CTO",
								ManagerID: intPtr(1),
							},
							Reports: []db.EmployeeNode{},
						},
					},
				},
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock database
			mockDB, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer mockDB.Close()

			// Create repository with mock DB
			repo := db.NewEmployeeRepository(&db.DB{DB: mockDB})

			// Set up expectations
			mock.ExpectQuery("SELECT id, name, title, manager_id FROM employees").
				WillReturnRows(tt.rows)

			// Execute the test
			tree, err := repo.GetEmployeeTree()

			// Assert results
			assert.Equal(t, tt.expectedError, err)
			assert.Equal(t, len(tt.expectedTree), len(tree))
			
			if len(tree) > 0 {
				assert.Equal(t, tt.expectedTree[0].Employee.Name, tree[0].Employee.Name)
				assert.Equal(t, len(tt.expectedTree[0].Reports), len(tree[0].Reports))
				
				if len(tree[0].Reports) > 0 {
					assert.Equal(t, tt.expectedTree[0].Reports[0].Employee.Name, tree[0].Reports[0].Employee.Name)
				}
			}
			
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

// Helper function to create int pointer
func intPtr(i int) *int {
	return &i
} 