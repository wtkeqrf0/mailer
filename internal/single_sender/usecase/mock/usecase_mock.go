// Code generated by MockGen. DO NOT EDIT.
// Source: ../usecase.go

// Package mock is a generated GoMock package.
package mock

import (
	consumer "mailer/internal/consumer"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockUseCase is a mock of UseCase interface.
type MockUseCase struct {
	ctrl     *gomock.Controller
	recorder *MockUseCaseMockRecorder
}

// MockUseCaseMockRecorder is the mock recorder for MockUseCase.
type MockUseCaseMockRecorder struct {
	mock *MockUseCase
}

// NewMockUseCase creates a new mock instance.
func NewMockUseCase(ctrl *gomock.Controller) *MockUseCase {
	mock := &MockUseCase{ctrl: ctrl}
	mock.recorder = &MockUseCaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUseCase) EXPECT() *MockUseCaseMockRecorder {
	return m.recorder
}

// SendEmail mocks base method.
func (m *MockUseCase) SendEmail(queueEmail consumer.Email) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendEmail", queueEmail)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SendEmail indicates an expected call of SendEmail.
func (mr *MockUseCaseMockRecorder) SendEmail(queueEmail interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendEmail", reflect.TypeOf((*MockUseCase)(nil).SendEmail), queueEmail)
}
