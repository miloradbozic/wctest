package service_test

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
	"wctest/app/config"
	"wctest/app/db"
	"wctest/app/model"
	"wctest/app/mocks"
	"wctest/app/service"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestEmployeeService_GetOrganizationTree(t *testing.T) {
	tests := []struct {
		name          string
		mockSetup     func(*mocks.MockEmployeeRepository)
		expectedTree  []db.EmployeeNode
		expectedError error
	}{
		{
			name: "successful retrieval",
			mockSetup: func(m *mocks.MockEmployeeRepository) {
				expectedTree := []db.EmployeeNode{
					{
						Employee: model.Employee{
							ID:   1,
							Name: "John Doe",
						},
						Reports: []db.EmployeeNode{
							{
								Employee: model.Employee{
									ID:   2,
									Name: "Jane Smith",
								},
							},
						},
					},
				}
				m.EXPECT().GetEmployeeTree().Return(expectedTree, nil)
			},
			expectedTree: []db.EmployeeNode{
				{
					Employee: model.Employee{
						ID:   1,
						Name: "John Doe",
					},
					Reports: []db.EmployeeNode{
						{
							Employee: model.Employee{
								ID:   2,
								Name: "Jane Smith",
							},
						},
					},
				},
			},
			expectedError: nil,
		},
		{
			name: "repository error",
			mockSetup: func(m *mocks.MockEmployeeRepository) {
				m.EXPECT().GetEmployeeTree().Return(nil, errors.New("database error"))
			},
			expectedTree:  nil,
			expectedError: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mocks.NewMockEmployeeRepository(ctrl)
			tt.mockSetup(mockRepo)

			service := service.NewEmployeeService(mockRepo, config.NewConfig())
			tree, err := service.GetOrganizationTree()

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, len(tt.expectedTree), len(tree))
			
			if len(tree) > 0 {
				assert.Equal(t, len(tt.expectedTree[0].Reports), len(tree[0].Reports))
				
				// Check root employee name
				assert.Equal(t, tt.expectedTree[0].Employee.Name, tree[0].Employee.Name)
				
				// Check reports' names
				for i, report := range tree[0].Reports {
					assert.Equal(t, tt.expectedTree[0].Reports[i].Employee.Name, report.Employee.Name)
				}
			}
		})
	}
}

func TestEmployeeService_InitializeData(t *testing.T) {
	tests := []struct {
		name          string
		mockSetup     func(*mocks.MockEmployeeRepository)
		httpSetup     func()
		expectedError error
	}{
		{
			name: "database not empty",
			mockSetup: func(m *mocks.MockEmployeeRepository) {
				m.EXPECT().IsEmpty().Return(false, nil)
			},
			httpSetup:     func() {},
			expectedError: nil,
		},
		{
			name: "successful initialization",
			mockSetup: func(m *mocks.MockEmployeeRepository) {
				m.EXPECT().IsEmpty().Return(true, nil)
				m.EXPECT().Create(gomock.Any()).Return(nil)
				m.EXPECT().Create(gomock.Any()).Return(nil)
			},
			httpSetup: func() {
				originalHTTPGet := service.HTTPGet
				service.HTTPGet = func(url string) (*http.Response, error) {
					return &http.Response{
						Body: io.NopCloser(strings.NewReader(`[
							{"id": 1, "name": "John Doe", "title": "CEO", "manager_id": null},
							{"id": 2, "name": "Jane Smith", "title": "CTO", "manager_id": 1}
						]`)),
					}, nil
				}
				defer func() { service.HTTPGet = originalHTTPGet }()
			},
			expectedError: nil,
		},
		{
			name: "database check error",
			mockSetup: func(m *mocks.MockEmployeeRepository) {
				m.EXPECT().IsEmpty().Return(false, errors.New("database error"))
			},
			httpSetup:     func() {},
			expectedError: errors.New("database error"),
		},
		{
			name: "HTTP request error",
			mockSetup: func(m *mocks.MockEmployeeRepository) {
				m.EXPECT().IsEmpty().Return(true, nil)
			},
			httpSetup: func() {
				originalHTTPGet := service.HTTPGet
				service.HTTPGet = func(url string) (*http.Response, error) {
					return nil, errors.New("HTTP error")
				}
				defer func() { service.HTTPGet = originalHTTPGet }()
			},
			expectedError: errors.New("HTTP error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mocks.NewMockEmployeeRepository(ctrl)
			tt.mockSetup(mockRepo)
			tt.httpSetup()

			service := service.NewEmployeeService(mockRepo, config.NewConfig())
			err := service.InitializeData()

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				return
			}

			assert.NoError(t, err)
		})
	}
}
