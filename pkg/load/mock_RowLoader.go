// Code generated by mockery v1.0.0. DO NOT EDIT.

package load

import mock "github.com/stretchr/testify/mock"

// MockRowLoader is an autogenerated mock type for the RowLoader type
type MockRowLoader struct {
	mock.Mock
}

// Export provides a mock function with given fields: _a0
func (_m *MockRowLoader) Export(_a0 Row) *Error {
	ret := _m.Called(_a0)

	var r0 *Error
	if rf, ok := ret.Get(0).(func(Row) *Error); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*Error)
		}
	}

	return r0
}