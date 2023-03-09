// Code generated by MockGen. DO NOT EDIT.
// Source: database/treatments.go

// Package mock_database is a generated GoMock package.
package mock_database

import (
	reflect "reflect"
	database "treatment-scheme-organizer/database"

	gomock "github.com/golang/mock/gomock"
)

// MockTreatmentsRepository is a mock of TreatmentsRepository interface.
type MockTreatmentsRepository struct {
	ctrl     *gomock.Controller
	recorder *MockTreatmentsRepositoryMockRecorder
}

// MockTreatmentsRepositoryMockRecorder is the mock recorder for MockTreatmentsRepository.
type MockTreatmentsRepositoryMockRecorder struct {
	mock *MockTreatmentsRepository
}

// NewMockTreatmentsRepository creates a new mock instance.
func NewMockTreatmentsRepository(ctrl *gomock.Controller) *MockTreatmentsRepository {
	mock := &MockTreatmentsRepository{ctrl: ctrl}
	mock.recorder = &MockTreatmentsRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTreatmentsRepository) EXPECT() *MockTreatmentsRepositoryMockRecorder {
	return m.recorder
}

// Add mocks base method.
func (m *MockTreatmentsRepository) Add(illnessID uint, begunAt, endedAt string, schemes []database.TreatmentScheme) (uint, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Add", illnessID, begunAt, endedAt, schemes)
	ret0, _ := ret[0].(uint)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Add indicates an expected call of Add.
func (mr *MockTreatmentsRepositoryMockRecorder) Add(illnessID, begunAt, endedAt, schemes interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Add", reflect.TypeOf((*MockTreatmentsRepository)(nil).Add), illnessID, begunAt, endedAt, schemes)
}

// ByID mocks base method.
func (m *MockTreatmentsRepository) ByID(id uint) (database.Treatment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ByID", id)
	ret0, _ := ret[0].(database.Treatment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ByID indicates an expected call of ByID.
func (mr *MockTreatmentsRepositoryMockRecorder) ByID(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ByID", reflect.TypeOf((*MockTreatmentsRepository)(nil).ByID), id)
}

// ByIllness mocks base method.
func (m *MockTreatmentsRepository) ByIllness(illnessID uint, limit, offset int) ([]database.Treatment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ByIllness", illnessID, limit, offset)
	ret0, _ := ret[0].([]database.Treatment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ByIllness indicates an expected call of ByIllness.
func (mr *MockTreatmentsRepositoryMockRecorder) ByIllness(illnessID, limit, offset interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ByIllness", reflect.TypeOf((*MockTreatmentsRepository)(nil).ByIllness), illnessID, limit, offset)
}