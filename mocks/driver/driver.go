// Code generated by mockery v2.2.1. DO NOT EDIT.

package mocks

import (
	context "context"

	driver "github.com/hashicorp/consul-nia/driver"
	mock "github.com/stretchr/testify/mock"
)

// Driver is an autogenerated mock type for the Driver type
type Driver struct {
	mock.Mock
}

// ApplyTask provides a mock function with given fields: ctx
func (_m *Driver) ApplyTask(ctx context.Context) error {
	ret := _m.Called(ctx)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Init provides a mock function with given fields: ctx
func (_m *Driver) Init(ctx context.Context) error {
	ret := _m.Called(ctx)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// InitTask provides a mock function with given fields: task, force
func (_m *Driver) InitTask(task driver.Task, force bool) error {
	ret := _m.Called(task, force)

	var r0 error
	if rf, ok := ret.Get(0).(func(driver.Task, bool) error); ok {
		r0 = rf(task, force)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Version provides a mock function with given fields:
func (_m *Driver) Version() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}
