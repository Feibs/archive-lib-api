// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	context "context"
	entity "archive_lib/entity"

	mock "github.com/stretchr/testify/mock"
)

// BorrowRepo is an autogenerated mock type for the BorrowRepo type
type BorrowRepo struct {
	mock.Mock
}

// GetBookByBorrowId provides a mock function with given fields: ctx, id
func (_m *BorrowRepo) GetBookByBorrowId(ctx context.Context, id int) (int, error) {
	ret := _m.Called(ctx, id)

	var r0 int
	if rf, ok := ret.Get(0).(func(context.Context, int) int); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(int)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IsBorrowExisted provides a mock function with given fields: ctx, id
func (_m *BorrowRepo) IsBorrowExisted(ctx context.Context, id int) (bool, error) {
	ret := _m.Called(ctx, id)

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context, int) bool); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IsReturned provides a mock function with given fields: ctx, id
func (_m *BorrowRepo) IsReturned(ctx context.Context, id int) (bool, error) {
	ret := _m.Called(ctx, id)

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context, int) bool); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IsUserAuthorized provides a mock function with given fields: ctx, record_id, user_id
func (_m *BorrowRepo) IsUserAuthorized(ctx context.Context, record_id int, user_id int) (bool, error) {
	ret := _m.Called(ctx, record_id, user_id)

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context, int, int) bool); ok {
		r0 = rf(ctx, record_id, user_id)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int, int) error); ok {
		r1 = rf(ctx, record_id, user_id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Record provides a mock function with given fields: ctx, borrow
func (_m *BorrowRepo) Record(ctx context.Context, borrow *entity.Borrow) (*entity.Borrow, error) {
	ret := _m.Called(ctx, borrow)

	var r0 *entity.Borrow
	if rf, ok := ret.Get(0).(func(context.Context, *entity.Borrow) *entity.Borrow); ok {
		r0 = rf(ctx, borrow)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.Borrow)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *entity.Borrow) error); ok {
		r1 = rf(ctx, borrow)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Return provides a mock function with given fields: ctx, borrow
func (_m *BorrowRepo) Return(ctx context.Context, borrow *entity.Borrow) (*entity.Borrow, error) {
	ret := _m.Called(ctx, borrow)

	var r0 *entity.Borrow
	if rf, ok := ret.Get(0).(func(context.Context, *entity.Borrow) *entity.Borrow); ok {
		r0 = rf(ctx, borrow)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.Borrow)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *entity.Borrow) error); ok {
		r1 = rf(ctx, borrow)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewBorrowRepo interface {
	mock.TestingT
	Cleanup(func())
}

// NewBorrowRepo creates a new instance of BorrowRepo. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewBorrowRepo(t mockConstructorTestingTNewBorrowRepo) *BorrowRepo {
	mock := &BorrowRepo{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
