// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/Foedie/foedie-server-v2/auth/internal/db (interfaces: Store)

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	db "github.com/Foedie/foedie-server-v2/auth/internal/db"
	gomock "github.com/golang/mock/gomock"
)

// MockStore is a mock of Store interface.
type MockStore struct {
	ctrl     *gomock.Controller
	recorder *MockStoreMockRecorder
}

// MockStoreMockRecorder is the mock recorder for MockStore.
type MockStoreMockRecorder struct {
	mock *MockStore
}

// NewMockStore creates a new mock instance.
func NewMockStore(ctrl *gomock.Controller) *MockStore {
	mock := &MockStore{ctrl: ctrl}
	mock.recorder = &MockStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStore) EXPECT() *MockStoreMockRecorder {
	return m.recorder
}

// CreateRecoveryAccount mocks base method.
func (m *MockStore) CreateRecoveryAccount(arg0 context.Context, arg1 db.CreateRecoveryAccountParams) (db.RecoverAccount, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateRecoveryAccount", arg0, arg1)
	ret0, _ := ret[0].(db.RecoverAccount)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateRecoveryAccount indicates an expected call of CreateRecoveryAccount.
func (mr *MockStoreMockRecorder) CreateRecoveryAccount(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateRecoveryAccount", reflect.TypeOf((*MockStore)(nil).CreateRecoveryAccount), arg0, arg1)
}

// CreateRefreshToken mocks base method.
func (m *MockStore) CreateRefreshToken(arg0 context.Context, arg1 db.CreateRefreshTokenParams) (db.Session, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateRefreshToken", arg0, arg1)
	ret0, _ := ret[0].(db.Session)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateRefreshToken indicates an expected call of CreateRefreshToken.
func (mr *MockStoreMockRecorder) CreateRefreshToken(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateRefreshToken", reflect.TypeOf((*MockStore)(nil).CreateRefreshToken), arg0, arg1)
}

// CreateUser mocks base method.
func (m *MockStore) CreateUser(arg0 context.Context, arg1 db.CreateUserParams) (db.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", arg0, arg1)
	ret0, _ := ret[0].(db.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockStoreMockRecorder) CreateUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockStore)(nil).CreateUser), arg0, arg1)
}

// CreateUserOTP mocks base method.
func (m *MockStore) CreateUserOTP(arg0 context.Context, arg1 db.CreateUserOTPParams) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUserOTP", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateUserOTP indicates an expected call of CreateUserOTP.
func (mr *MockStoreMockRecorder) CreateUserOTP(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUserOTP", reflect.TypeOf((*MockStore)(nil).CreateUserOTP), arg0, arg1)
}

// CreateUserTx mocks base method.
func (m *MockStore) CreateUserTx(arg0 context.Context, arg1 db.UserTxParams) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUserTx", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateUserTx indicates an expected call of CreateUserTx.
func (mr *MockStoreMockRecorder) CreateUserTx(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUserTx", reflect.TypeOf((*MockStore)(nil).CreateUserTx), arg0, arg1)
}

// CreateVerifyEmail mocks base method.
func (m *MockStore) CreateVerifyEmail(arg0 context.Context, arg1 db.CreateVerifyEmailParams) (db.VerifyEmail, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateVerifyEmail", arg0, arg1)
	ret0, _ := ret[0].(db.VerifyEmail)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateVerifyEmail indicates an expected call of CreateVerifyEmail.
func (mr *MockStoreMockRecorder) CreateVerifyEmail(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateVerifyEmail", reflect.TypeOf((*MockStore)(nil).CreateVerifyEmail), arg0, arg1)
}

// DeleteSessionUser mocks base method.
func (m *MockStore) DeleteSessionUser(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteSessionUser", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteSessionUser indicates an expected call of DeleteSessionUser.
func (mr *MockStoreMockRecorder) DeleteSessionUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteSessionUser", reflect.TypeOf((*MockStore)(nil).DeleteSessionUser), arg0, arg1)
}

// ExecTx mocks base method.
func (m *MockStore) ExecTx(arg0 context.Context, arg1 func(*db.Queries) error) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ExecTx", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// ExecTx indicates an expected call of ExecTx.
func (mr *MockStoreMockRecorder) ExecTx(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExecTx", reflect.TypeOf((*MockStore)(nil).ExecTx), arg0, arg1)
}

// GetRecoverAccount mocks base method.
func (m *MockStore) GetRecoverAccount(arg0 context.Context, arg1 db.GetRecoverAccountParams) (db.RecoverAccount, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRecoverAccount", arg0, arg1)
	ret0, _ := ret[0].(db.RecoverAccount)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRecoverAccount indicates an expected call of GetRecoverAccount.
func (mr *MockStoreMockRecorder) GetRecoverAccount(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRecoverAccount", reflect.TypeOf((*MockStore)(nil).GetRecoverAccount), arg0, arg1)
}

// GetRefreshToken mocks base method.
func (m *MockStore) GetRefreshToken(arg0 context.Context, arg1 string) (db.Session, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRefreshToken", arg0, arg1)
	ret0, _ := ret[0].(db.Session)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRefreshToken indicates an expected call of GetRefreshToken.
func (mr *MockStoreMockRecorder) GetRefreshToken(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRefreshToken", reflect.TypeOf((*MockStore)(nil).GetRefreshToken), arg0, arg1)
}

// GetUserByEmailAndUsername mocks base method.
func (m *MockStore) GetUserByEmailAndUsername(arg0 context.Context, arg1 db.GetUserByEmailAndUsernameParams) (db.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByEmailAndUsername", arg0, arg1)
	ret0, _ := ret[0].(db.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByEmailAndUsername indicates an expected call of GetUserByEmailAndUsername.
func (mr *MockStoreMockRecorder) GetUserByEmailAndUsername(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByEmailAndUsername", reflect.TypeOf((*MockStore)(nil).GetUserByEmailAndUsername), arg0, arg1)
}

// GetUserByEmailOrUsername mocks base method.
func (m *MockStore) GetUserByEmailOrUsername(arg0 context.Context, arg1 string) (db.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByEmailOrUsername", arg0, arg1)
	ret0, _ := ret[0].(db.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByEmailOrUsername indicates an expected call of GetUserByEmailOrUsername.
func (mr *MockStoreMockRecorder) GetUserByEmailOrUsername(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByEmailOrUsername", reflect.TypeOf((*MockStore)(nil).GetUserByEmailOrUsername), arg0, arg1)
}

// GetUserByUid mocks base method.
func (m *MockStore) GetUserByUid(arg0 context.Context, arg1 string) (db.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByUid", arg0, arg1)
	ret0, _ := ret[0].(db.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByUid indicates an expected call of GetUserByUid.
func (mr *MockStoreMockRecorder) GetUserByUid(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByUid", reflect.TypeOf((*MockStore)(nil).GetUserByUid), arg0, arg1)
}

// GetVerifyEmail mocks base method.
func (m *MockStore) GetVerifyEmail(arg0 context.Context, arg1 db.GetVerifyEmailParams) (db.VerifyEmail, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetVerifyEmail", arg0, arg1)
	ret0, _ := ret[0].(db.VerifyEmail)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetVerifyEmail indicates an expected call of GetVerifyEmail.
func (mr *MockStoreMockRecorder) GetVerifyEmail(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetVerifyEmail", reflect.TypeOf((*MockStore)(nil).GetVerifyEmail), arg0, arg1)
}

// UpdateRecoverAccount mocks base method.
func (m *MockStore) UpdateRecoverAccount(arg0 context.Context, arg1 db.UpdateRecoverAccountParams) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateRecoverAccount", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateRecoverAccount indicates an expected call of UpdateRecoverAccount.
func (mr *MockStoreMockRecorder) UpdateRecoverAccount(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateRecoverAccount", reflect.TypeOf((*MockStore)(nil).UpdateRecoverAccount), arg0, arg1)
}

// UpdateUserActive mocks base method.
func (m *MockStore) UpdateUserActive(arg0 context.Context, arg1 db.UpdateUserActiveParams) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUserActive", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateUserActive indicates an expected call of UpdateUserActive.
func (mr *MockStoreMockRecorder) UpdateUserActive(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUserActive", reflect.TypeOf((*MockStore)(nil).UpdateUserActive), arg0, arg1)
}

// UpdateUserLastLogin mocks base method.
func (m *MockStore) UpdateUserLastLogin(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUserLastLogin", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateUserLastLogin indicates an expected call of UpdateUserLastLogin.
func (mr *MockStoreMockRecorder) UpdateUserLastLogin(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUserLastLogin", reflect.TypeOf((*MockStore)(nil).UpdateUserLastLogin), arg0, arg1)
}

// UpdateUserPassword mocks base method.
func (m *MockStore) UpdateUserPassword(arg0 context.Context, arg1 db.UpdateUserPasswordParams) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUserPassword", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateUserPassword indicates an expected call of UpdateUserPassword.
func (mr *MockStoreMockRecorder) UpdateUserPassword(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUserPassword", reflect.TypeOf((*MockStore)(nil).UpdateUserPassword), arg0, arg1)
}

// UpdateUserVerified mocks base method.
func (m *MockStore) UpdateUserVerified(arg0 context.Context, arg1 db.UpdateUserVerifiedParams) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUserVerified", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateUserVerified indicates an expected call of UpdateUserVerified.
func (mr *MockStoreMockRecorder) UpdateUserVerified(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUserVerified", reflect.TypeOf((*MockStore)(nil).UpdateUserVerified), arg0, arg1)
}

// UpdateVerifyEmail mocks base method.
func (m *MockStore) UpdateVerifyEmail(arg0 context.Context, arg1 db.UpdateVerifyEmailParams) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateVerifyEmail", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateVerifyEmail indicates an expected call of UpdateVerifyEmail.
func (mr *MockStoreMockRecorder) UpdateVerifyEmail(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateVerifyEmail", reflect.TypeOf((*MockStore)(nil).UpdateVerifyEmail), arg0, arg1)
}

// ValidateEmailTx mocks base method.
func (m *MockStore) ValidateEmailTx(arg0 context.Context, arg1, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidateEmailTx", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// ValidateEmailTx indicates an expected call of ValidateEmailTx.
func (mr *MockStoreMockRecorder) ValidateEmailTx(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateEmailTx", reflect.TypeOf((*MockStore)(nil).ValidateEmailTx), arg0, arg1, arg2)
}

// ValidateRecoverAccountTx mocks base method.
func (m *MockStore) ValidateRecoverAccountTx(arg0 context.Context, arg1, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidateRecoverAccountTx", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// ValidateRecoverAccountTx indicates an expected call of ValidateRecoverAccountTx.
func (mr *MockStoreMockRecorder) ValidateRecoverAccountTx(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateRecoverAccountTx", reflect.TypeOf((*MockStore)(nil).ValidateRecoverAccountTx), arg0, arg1, arg2)
}
