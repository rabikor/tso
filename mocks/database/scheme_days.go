// Code generated by MockGen. DO NOT EDIT.
// Source: database/scheme_days.go

// Package mock_database is a generated GoMock package.
package mock_database

import (
	reflect "reflect"
	database "treatment-scheme-organizer/database"

	gomock "github.com/golang/mock/gomock"
)

// MockSchemeDaysRepository is a mock of SchemeDaysRepository interface.
type MockSchemeDaysRepository struct {
	ctrl     *gomock.Controller
	recorder *MockSchemeDaysRepositoryMockRecorder
}

// MockSchemeDaysRepositoryMockRecorder is the mock recorder for MockSchemeDaysRepository.
type MockSchemeDaysRepositoryMockRecorder struct {
	mock *MockSchemeDaysRepository
}

// NewMockSchemeDaysRepository creates a new mock instance.
func NewMockSchemeDaysRepository(ctrl *gomock.Controller) *MockSchemeDaysRepository {
	mock := &MockSchemeDaysRepository{ctrl: ctrl}
	mock.recorder = &MockSchemeDaysRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSchemeDaysRepository) EXPECT() *MockSchemeDaysRepositoryMockRecorder {
	return m.recorder
}

// Add mocks base method.
func (m *MockSchemeDaysRepository) Add(schemeID, procedureID, drugID, order, times, frequency uint) (uint, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Add", schemeID, procedureID, drugID, order, times, frequency)
	ret0, _ := ret[0].(uint)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Add indicates an expected call of Add.
func (mr *MockSchemeDaysRepositoryMockRecorder) Add(schemeID, procedureID, drugID, order, times, frequency interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Add", reflect.TypeOf((*MockSchemeDaysRepository)(nil).Add), schemeID, procedureID, drugID, order, times, frequency)
}

// ByScheme mocks base method.
func (m *MockSchemeDaysRepository) ByScheme(schemeID, limit, offset int) ([]database.SchemeDay, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ByScheme", schemeID, limit, offset)
	ret0, _ := ret[0].([]database.SchemeDay)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ByScheme indicates an expected call of ByScheme.
func (mr *MockSchemeDaysRepositoryMockRecorder) ByScheme(schemeID, limit, offset interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ByScheme", reflect.TypeOf((*MockSchemeDaysRepository)(nil).ByScheme), schemeID, limit, offset)
}
