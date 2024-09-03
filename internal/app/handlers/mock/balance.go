// Code generated by MockGen. DO NOT EDIT.
// Source: internal/app/handlers/balance/balance.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockBalancer is a mock of Balancer interface.
type MockBalancer struct {
	ctrl     *gomock.Controller
	recorder *MockBalancerMockRecorder
}

// MockBalancerMockRecorder is the mock recorder for MockBalancer.
type MockBalancerMockRecorder struct {
	mock *MockBalancer
}

// NewMockBalancer creates a new mock instance.
func NewMockBalancer(ctrl *gomock.Controller) *MockBalancer {
	mock := &MockBalancer{ctrl: ctrl}
	mock.recorder = &MockBalancerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBalancer) EXPECT() *MockBalancerMockRecorder {
	return m.recorder
}

// GetUserBalance mocks base method.
func (m *MockBalancer) GetUserBalance(ctx context.Context, userID int) (float32, float32, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserBalance", ctx, userID)
	ret0, _ := ret[0].(float32)
	ret1, _ := ret[1].(float32)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetUserBalance indicates an expected call of GetUserBalance.
func (mr *MockBalancerMockRecorder) GetUserBalance(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserBalance", reflect.TypeOf((*MockBalancer)(nil).GetUserBalance), ctx, userID)
}
