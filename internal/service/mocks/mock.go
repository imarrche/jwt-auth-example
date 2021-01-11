// Code generated by MockGen. DO NOT EDIT.
// Source: interface.go

// Package mock_service is a generated GoMock package.
package mock_service

import (
	gomock "github.com/golang/mock/gomock"
	model "github.com/imarrche/jwt-auth-example/internal/model"
	service "github.com/imarrche/jwt-auth-example/internal/service"
	reflect "reflect"
)

// MockService is a mock of Service interface
type MockService struct {
	ctrl     *gomock.Controller
	recorder *MockServiceMockRecorder
}

// MockServiceMockRecorder is the mock recorder for MockService
type MockServiceMockRecorder struct {
	mock *MockService
}

// NewMockService creates a new mock instance
func NewMockService(ctrl *gomock.Controller) *MockService {
	mock := &MockService{ctrl: ctrl}
	mock.recorder = &MockServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockService) EXPECT() *MockServiceMockRecorder {
	return m.recorder
}

// Auth mocks base method
func (m *MockService) Auth() service.Auth {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Auth")
	ret0, _ := ret[0].(service.Auth)
	return ret0
}

// Auth indicates an expected call of Auth
func (mr *MockServiceMockRecorder) Auth() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Auth", reflect.TypeOf((*MockService)(nil).Auth))
}

// MockAuth is a mock of Auth interface
type MockAuth struct {
	ctrl     *gomock.Controller
	recorder *MockAuthMockRecorder
}

// MockAuthMockRecorder is the mock recorder for MockAuth
type MockAuthMockRecorder struct {
	mock *MockAuth
}

// NewMockAuth creates a new mock instance
func NewMockAuth(ctrl *gomock.Controller) *MockAuth {
	mock := &MockAuth{ctrl: ctrl}
	mock.recorder = &MockAuthMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockAuth) EXPECT() *MockAuthMockRecorder {
	return m.recorder
}

// SignUp mocks base method
func (m *MockAuth) SignUp(arg0 model.User) (model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SignUp", arg0)
	ret0, _ := ret[0].(model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SignUp indicates an expected call of SignUp
func (mr *MockAuthMockRecorder) SignUp(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignUp", reflect.TypeOf((*MockAuth)(nil).SignUp), arg0)
}

// SignIn mocks base method
func (m *MockAuth) SignIn(arg0, arg1 string) (string, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SignIn", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// SignIn indicates an expected call of SignIn
func (mr *MockAuthMockRecorder) SignIn(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignIn", reflect.TypeOf((*MockAuth)(nil).SignIn), arg0, arg1)
}

// ValidateJWT mocks base method
func (m *MockAuth) ValidateJWT(arg0, arg1 string) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidateJWT", arg0, arg1)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ValidateJWT indicates an expected call of ValidateJWT
func (mr *MockAuthMockRecorder) ValidateJWT(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateJWT", reflect.TypeOf((*MockAuth)(nil).ValidateJWT), arg0, arg1)
}

// RefreshAccessJWT mocks base method
func (m *MockAuth) RefreshAccessJWT(arg0 string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RefreshAccessJWT", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RefreshAccessJWT indicates an expected call of RefreshAccessJWT
func (mr *MockAuthMockRecorder) RefreshAccessJWT(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RefreshAccessJWT", reflect.TypeOf((*MockAuth)(nil).RefreshAccessJWT), arg0)
}
