// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// Bcrypt is an autogenerated mock type for the Bcrypt type
type Bcrypt struct {
	mock.Mock
}

// CompareHashAndPassword provides a mock function with given fields: hashedPassword, password
func (_m *Bcrypt) CompareHashAndPassword(hashedPassword []byte, password []byte) error {
	ret := _m.Called(hashedPassword, password)

	var r0 error
	if rf, ok := ret.Get(0).(func([]byte, []byte) error); ok {
		r0 = rf(hashedPassword, password)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewBcrypt interface {
	mock.TestingT
	Cleanup(func())
}

// NewBcrypt creates a new instance of Bcrypt. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewBcrypt(t mockConstructorTestingTNewBcrypt) *Bcrypt {
	mock := &Bcrypt{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
