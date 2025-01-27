// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// TransactionRepo is an autogenerated mock type for the TransactionRepo type
type TransactionRepo struct {
	mock.Mock
}

// WithinTransaction provides a mock function with given fields: ctx, fn
func (_m *TransactionRepo) WithinTransaction(ctx context.Context, fn func(context.Context) error) error {
	ret := _m.Called(ctx, fn)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, func(context.Context) error) error); ok {
		r0 = rf(ctx, fn)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewTransactionRepo interface {
	mock.TestingT
	Cleanup(func())
}

// NewTransactionRepo creates a new instance of TransactionRepo. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewTransactionRepo(t mockConstructorTestingTNewTransactionRepo) *TransactionRepo {
	mock := &TransactionRepo{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
