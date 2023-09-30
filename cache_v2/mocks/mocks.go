// Code generated by MockGen. DO NOT EDIT.
// Source: cache_v2/cache.go

// Package mock_cache is a generated GoMock package.
package mock_cache

import (
	context "context"
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
)

// MockCacher is a mock of Cacher interface.
type MockCacher struct {
	ctrl     *gomock.Controller
	recorder *MockCacherMockRecorder
}

// MockCacherMockRecorder is the mock recorder for MockCacher.
type MockCacherMockRecorder struct {
	mock *MockCacher
}

// NewMockCacher creates a new mock instance.
func NewMockCacher(ctrl *gomock.Controller) *MockCacher {
	mock := &MockCacher{ctrl: ctrl}
	mock.recorder = &MockCacherMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCacher) EXPECT() *MockCacherMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockCacher) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockCacherMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockCacher)(nil).Close))
}

// Delete mocks base method.
func (m *MockCacher) Delete(ctx context.Context, key string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, key)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockCacherMockRecorder) Delete(ctx, key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockCacher)(nil).Delete), ctx, key)
}

// Get mocks base method.
func (m *MockCacher) Get(ctx context.Context, key string, data interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, key, data)
	ret0, _ := ret[0].(error)
	return ret0
}

// Get indicates an expected call of Get.
func (mr *MockCacherMockRecorder) Get(ctx, key, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockCacher)(nil).Get), ctx, key, data)
}

// Set mocks base method.
func (m *MockCacher) Set(ctx context.Context, key string, data interface{}, ttl ...time.Duration) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, key, data}
	for _, a := range ttl {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Set", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Set indicates an expected call of Set.
func (mr *MockCacherMockRecorder) Set(ctx, key, data interface{}, ttl ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, key, data}, ttl...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Set", reflect.TypeOf((*MockCacher)(nil).Set), varargs...)
}
