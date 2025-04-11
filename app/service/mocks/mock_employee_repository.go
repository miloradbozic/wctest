package mocks

import (
	"reflect"
	"wctest/app/db"
	"wctest/app/model"

	"github.com/golang/mock/gomock"
)

// MockEmployeeRepository is a mock of EmployeeRepository interface
type MockEmployeeRepository struct {
	ctrl     *gomock.Controller
	recorder *MockEmployeeRepositoryMockRecorder
}

// MockEmployeeRepositoryMockRecorder is the mock recorder for MockEmployeeRepository
type MockEmployeeRepositoryMockRecorder struct {
	mock *MockEmployeeRepository
}

// NewMockEmployeeRepository creates a new mock instance
func NewMockEmployeeRepository(ctrl *gomock.Controller) *MockEmployeeRepository {
	mock := &MockEmployeeRepository{ctrl: ctrl}
	mock.recorder = &MockEmployeeRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockEmployeeRepository) EXPECT() *MockEmployeeRepositoryMockRecorder {
	return m.recorder
}

// GetEmployeeTree mocks base method
func (m *MockEmployeeRepository) GetEmployeeTree() ([]db.EmployeeNode, error) {
	ret := m.ctrl.Call(m, "GetEmployeeTree")
	ret0, _ := ret[0].([]db.EmployeeNode)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetEmployeeTree indicates an expected call of GetEmployeeTree
func (mr *MockEmployeeRepositoryMockRecorder) GetEmployeeTree() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEmployeeTree", reflect.TypeOf((*MockEmployeeRepository)(nil).GetEmployeeTree))
}

// IsEmpty mocks base method
func (m *MockEmployeeRepository) IsEmpty() (bool, error) {
	ret := m.ctrl.Call(m, "IsEmpty")
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsEmpty indicates an expected call of IsEmpty
func (mr *MockEmployeeRepositoryMockRecorder) IsEmpty() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsEmpty", reflect.TypeOf((*MockEmployeeRepository)(nil).IsEmpty))
}

// Create mocks base method
func (m *MockEmployeeRepository) Create(employee *model.Employee) error {
	ret := m.ctrl.Call(m, "Create", employee)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create
func (mr *MockEmployeeRepositoryMockRecorder) Create(employee interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockEmployeeRepository)(nil).Create), employee)
} 