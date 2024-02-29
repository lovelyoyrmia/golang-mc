// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/Foedie/foedie-server-v2/user/domain/clients (interfaces: AuthServiceClient)

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	pb "github.com/Foedie/foedie-server-v2/user/domain/pb"
	gomock "github.com/golang/mock/gomock"
)

// MockAuthServiceClient is a mock of AuthServiceClient interface.
type MockAuthServiceClient struct {
	ctrl     *gomock.Controller
	recorder *MockAuthServiceClientMockRecorder
}

// MockAuthServiceClientMockRecorder is the mock recorder for MockAuthServiceClient.
type MockAuthServiceClientMockRecorder struct {
	mock *MockAuthServiceClient
}

// NewMockAuthServiceClient creates a new mock instance.
func NewMockAuthServiceClient(ctrl *gomock.Controller) *MockAuthServiceClient {
	mock := &MockAuthServiceClient{ctrl: ctrl}
	mock.recorder = &MockAuthServiceClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuthServiceClient) EXPECT() *MockAuthServiceClientMockRecorder {
	return m.recorder
}

// VerifyEmail mocks base method.
func (m *MockAuthServiceClient) VerifyEmail(arg0 string) (*pb.VerifyEmailResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "VerifyEmail", arg0)
	ret0, _ := ret[0].(*pb.VerifyEmailResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// VerifyEmail indicates an expected call of VerifyEmail.
func (mr *MockAuthServiceClientMockRecorder) VerifyEmail(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VerifyEmail", reflect.TypeOf((*MockAuthServiceClient)(nil).VerifyEmail), arg0)
}
