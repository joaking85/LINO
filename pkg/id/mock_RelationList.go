// Code generated by mockery v1.0.0. DO NOT EDIT.

package id

import mock "github.com/stretchr/testify/mock"

// MockRelationList is an autogenerated mock type for the RelationList type
type MockRelationList struct {
	mock.Mock
}

// Contains provides a mock function with given fields: _a0
func (_m *MockRelationList) Contains(_a0 string) bool {
	ret := _m.Called(_a0)

	var r0 bool
	if rf, ok := ret.Get(0).(func(string) bool); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// Len provides a mock function with given fields:
func (_m *MockRelationList) Len() uint {
	ret := _m.Called()

	var r0 uint
	if rf, ok := ret.Get(0).(func() uint); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint)
	}

	return r0
}

// Relation provides a mock function with given fields: idx
func (_m *MockRelationList) Relation(idx uint) Relation {
	ret := _m.Called(idx)

	var r0 Relation
	if rf, ok := ret.Get(0).(func(uint) Relation); ok {
		r0 = rf(idx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(Relation)
		}
	}

	return r0
}

// String provides a mock function with given fields:
func (_m *MockRelationList) String() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}
