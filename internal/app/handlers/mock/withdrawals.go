// Code generated by MockGen. DO NOT EDIT.
// Source: internal/app/handlers/withdrawals/withdrawals.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	models "github.com/FischukSergey/go-gothermart.git/internal/models"
	gomock "github.com/golang/mock/gomock"
)

// MockGetUserWithdrawAll is a mock of GetUserWithdrawAll interface.
type MockGetUserWithdrawAll struct {
	ctrl     *gomock.Controller
	recorder *MockGetUserWithdrawAllMockRecorder
}

// MockGetUserWithdrawAllMockRecorder is the mock recorder for MockGetUserWithdrawAll.
type MockGetUserWithdrawAllMockRecorder struct {
	mock *MockGetUserWithdrawAll
}

// NewMockGetUserWithdrawAll creates a new mock instance.
func NewMockGetUserWithdrawAll(ctrl *gomock.Controller) *MockGetUserWithdrawAll {
	mock := &MockGetUserWithdrawAll{ctrl: ctrl}
	mock.recorder = &MockGetUserWithdrawAllMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockGetUserWithdrawAll) EXPECT() *MockGetUserWithdrawAllMockRecorder {
	return m.recorder
}

// GetAllWithdraw mocks base method.
func (m *MockGetUserWithdrawAll) GetAllWithdraw(ctx context.Context, userID int) ([]models.GetAllWithdraw, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllWithdraw", ctx, userID)
	ret0, _ := ret[0].([]models.GetAllWithdraw)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllWithdraw indicates an expected call of GetAllWithdraw.
func (mr *MockGetUserWithdrawAllMockRecorder) GetAllWithdraw(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllWithdraw", reflect.TypeOf((*MockGetUserWithdrawAll)(nil).GetAllWithdraw), ctx, userID)
}
