// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	context "context"
	dto "archive_lib/dto"

	mock "github.com/stretchr/testify/mock"
)

// BorrowUsecase is an autogenerated mock type for the BorrowUsecase type
type BorrowUsecase struct {
	mock.Mock
}

// Record provides a mock function with given fields: ctx, borrowRequest
func (_m *BorrowUsecase) Record(ctx context.Context, borrowRequest *dto.BorrowRequest) (*dto.BorrowResponse, error) {
	ret := _m.Called(ctx, borrowRequest)

	var r0 *dto.BorrowResponse
	if rf, ok := ret.Get(0).(func(context.Context, *dto.BorrowRequest) *dto.BorrowResponse); ok {
		r0 = rf(ctx, borrowRequest)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*dto.BorrowResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *dto.BorrowRequest) error); ok {
		r1 = rf(ctx, borrowRequest)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Return provides a mock function with given fields: ctx, returnRequest
func (_m *BorrowUsecase) Return(ctx context.Context, returnRequest *dto.ReturnRequest) (*dto.BorrowResponse, error) {
	ret := _m.Called(ctx, returnRequest)

	var r0 *dto.BorrowResponse
	if rf, ok := ret.Get(0).(func(context.Context, *dto.ReturnRequest) *dto.BorrowResponse); ok {
		r0 = rf(ctx, returnRequest)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*dto.BorrowResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *dto.ReturnRequest) error); ok {
		r1 = rf(ctx, returnRequest)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewBorrowUsecase interface {
	mock.TestingT
	Cleanup(func())
}

// NewBorrowUsecase creates a new instance of BorrowUsecase. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewBorrowUsecase(t mockConstructorTestingTNewBorrowUsecase) *BorrowUsecase {
	mock := &BorrowUsecase{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}