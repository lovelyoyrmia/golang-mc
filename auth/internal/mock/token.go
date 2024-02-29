// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/Foedie/foedie-server-v2/auth/pkg/token (interfaces: Maker)

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"
	time "time"

	token "github.com/Foedie/foedie-server-v2/auth/pkg/token"
	gomock "github.com/golang/mock/gomock"
)

// MockMaker is a mock of Maker interface.
type MockMaker struct {
	ctrl     *gomock.Controller
	recorder *MockMakerMockRecorder
}

// MockMakerMockRecorder is the mock recorder for MockMaker.
type MockMakerMockRecorder struct {
	mock *MockMaker
}

// NewMockMaker creates a new mock instance.
func NewMockMaker(ctrl *gomock.Controller) *MockMaker {
	mock := &MockMaker{ctrl: ctrl}
	mock.recorder = &MockMakerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMaker) EXPECT() *MockMakerMockRecorder {
	return m.recorder
}

// GenerateToken mocks base method.
func (m *MockMaker) GenerateToken(arg0 string, arg1 time.Duration) (string, *token.Payload, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateToken", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(*token.Payload)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GenerateToken indicates an expected call of GenerateToken.
func (mr *MockMakerMockRecorder) GenerateToken(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateToken", reflect.TypeOf((*MockMaker)(nil).GenerateToken), arg0, arg1)
}

// VerifyToken mocks base method.
func (m *MockMaker) VerifyToken(arg0 string) (*token.Payload, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "VerifyToken", arg0)
	ret0, _ := ret[0].(*token.Payload)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// VerifyToken indicates an expected call of VerifyToken.
func (mr *MockMakerMockRecorder) VerifyToken(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VerifyToken", reflect.TypeOf((*MockMaker)(nil).VerifyToken), arg0)
}
