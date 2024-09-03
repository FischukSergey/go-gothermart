// Code generated by MockGen. DO NOT EDIT.
// Source: internal/app/handlers/withdraw/withdraw.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	models "github.com/FischukSergey/go-gothermart.git/internal/models"
	gomock "github.com/golang/mock/gomock"
)

// MockOrderBalanceWithdraw is a mock of OrderBalanceWithdraw interface.
type MockOrderBalanceWithdraw struct {
	ctrl     *gomock.Controller
	recorder *MockOrderBalanceWithdrawMockRecorder
}

// MockOrderBalanceWithdrawMockRecorder is the mock recorder for MockOrderBalanceWithdraw.
type MockOrderBalanceWithdrawMockRecorder struct {
	mock *MockOrderBalanceWithdraw
}

// NewMockOrderBalanceWithdraw creates a new mock instance.
func NewMockOrderBalanceWithdraw(ctrl *gomock.Controller) *MockOrderBalanceWithdraw {
	mock := &MockOrderBalanceWithdraw{ctrl: ctrl}
	mock.recorder = &MockOrderBalanceWithdrawMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockOrderBalanceWithdraw) EXPECT() *MockOrderBalanceWithdrawMockRecorder {
	return m.recorder
}

// CreateOrderWithdraw mocks base method.
func (m *MockOrderBalanceWithdraw) CreateOrderWithdraw(ctx context.Context, order models.Order) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateOrderWithdraw", ctx, order)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateOrderWithdraw indicates an expected call of CreateOrderWithdraw.
func (mr *MockOrderBalanceWithdrawMockRecorder) CreateOrderWithdraw(ctx, order interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateOrderWithdraw", reflect.TypeOf((*MockOrderBalanceWithdraw)(nil).CreateOrderWithdraw), ctx, order)
}
