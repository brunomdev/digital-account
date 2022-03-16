// Code generated by MockGen. DO NOT EDIT.
// Source: contract.go

// Package mock_transaction is a generated GoMock package.
package mock_transaction

import (
	context "context"
	reflect "reflect"

	entity "github.com/brunomdev/digital-account/entity"
	gomock "github.com/golang/mock/gomock"
)

// MockService is a mock of Service interface.
type MockService struct {
	ctrl     *gomock.Controller
	recorder *MockServiceMockRecorder
}

// MockServiceMockRecorder is the mock recorder for MockService.
type MockServiceMockRecorder struct {
	mock *MockService
}

// NewMockService creates a new mock instance.
func NewMockService(ctrl *gomock.Controller) *MockService {
	mock := &MockService{ctrl: ctrl}
	mock.recorder = &MockServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockService) EXPECT() *MockServiceMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockService) Create(ctx context.Context, accountID, operationTypeID int, amount float64) (*entity.Transaction, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, accountID, operationTypeID, amount)
	ret0, _ := ret[0].(*entity.Transaction)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockServiceMockRecorder) Create(ctx, accountID, operationTypeID, amount interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockService)(nil).Create), ctx, accountID, operationTypeID, amount)
}

// MockRepository is a mock of Repository interface.
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder
}

// MockRepositoryMockRecorder is the mock recorder for MockRepository.
type MockRepositoryMockRecorder struct {
	mock *MockRepository
}

// NewMockRepository creates a new mock instance.
func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	mock := &MockRepository{ctrl: ctrl}
	mock.recorder = &MockRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepository) EXPECT() *MockRepositoryMockRecorder {
	return m.recorder
}

// Save mocks base method.
func (m *MockRepository) Save(ctx context.Context, accountID, operationTypeID int, amount float64) (*entity.Transaction, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", ctx, accountID, operationTypeID, amount)
	ret0, _ := ret[0].(*entity.Transaction)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Save indicates an expected call of Save.
func (mr *MockRepositoryMockRecorder) Save(ctx, accountID, operationTypeID, amount interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockRepository)(nil).Save), ctx, accountID, operationTypeID, amount)
}