// Code generated by MockGen. DO NOT EDIT.
// Source: database/procedures.go

// Package mock_database is a generated GoMock package.
package mock_database

import (
	reflect "reflect"
	database "treatment-scheme-organizer/database"

	gomock "github.com/golang/mock/gomock"
)

// MockProceduresRepository is a mock of ProceduresRepository interface.
type MockProceduresRepository struct {
	ctrl     *gomock.Controller
	recorder *MockProceduresRepositoryMockRecorder
}

// MockProceduresRepositoryMockRecorder is the mock recorder for MockProceduresRepository.
type MockProceduresRepositoryMockRecorder struct {
	mock *MockProceduresRepository
}

// NewMockProceduresRepository creates a new mock instance.
func NewMockProceduresRepository(ctrl *gomock.Controller) *MockProceduresRepository {
	mock := &MockProceduresRepository{ctrl: ctrl}
	mock.recorder = &MockProceduresRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockProceduresRepository) EXPECT() *MockProceduresRepositoryMockRecorder {
	return m.recorder
}

// Add mocks base method.
func (m *MockProceduresRepository) Add(title string) (uint, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Add", title)
	ret0, _ := ret[0].(uint)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Add indicates an expected call of Add.
func (mr *MockProceduresRepositoryMockRecorder) Add(title interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Add", reflect.TypeOf((*MockProceduresRepository)(nil).Add), title)
}

// All mocks base method.
func (m *MockProceduresRepository) All(limit, offset int) ([]database.Procedure, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "All", limit, offset)
	ret0, _ := ret[0].([]database.Procedure)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// All indicates an expected call of All.
func (mr *MockProceduresRepositoryMockRecorder) All(limit, offset interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "All", reflect.TypeOf((*MockProceduresRepository)(nil).All), limit, offset)
}

// ByID mocks base method.
func (m *MockProceduresRepository) ByID(id uint) (database.Procedure, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ByID", id)
	ret0, _ := ret[0].(database.Procedure)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ByID indicates an expected call of ByID.
func (mr *MockProceduresRepositoryMockRecorder) ByID(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ByID", reflect.TypeOf((*MockProceduresRepository)(nil).ByID), id)
}
