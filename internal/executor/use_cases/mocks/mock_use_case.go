// Code generated by MockGen. DO NOT EDIT.
// Source: C:\Dev\Calculator\internal\executor\use_cases\use_case.go
//
// Generated by this command:
//
//	mockgen -source=C:\Dev\Calculator\internal\executor\use_cases\use_case.go -destination=C:\Dev\Calculator\internal\executor\use_cases\mocks\mock_use_case.go -package=mocks
//

// Package mocks is a generated GoMock package.
package mocks

import (
	executor "Calculator/internal/executor"
	dto "Calculator/internal/executor/dto"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockCommunicationService is a mock of CommunicationService interface.
type MockCommunicationService struct {
	ctrl     *gomock.Controller
	recorder *MockCommunicationServiceMockRecorder
	isgomock struct{}
}

// MockCommunicationServiceMockRecorder is the mock recorder for MockCommunicationService.
type MockCommunicationServiceMockRecorder struct {
	mock *MockCommunicationService
}

// NewMockCommunicationService creates a new mock instance.
func NewMockCommunicationService(ctrl *gomock.Controller) *MockCommunicationService {
	mock := &MockCommunicationService{ctrl: ctrl}
	mock.recorder = &MockCommunicationServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCommunicationService) EXPECT() *MockCommunicationServiceMockRecorder {
	return m.recorder
}

// ConsumeResults mocks base method.
func (m *MockCommunicationService) ConsumeResults(queue string, rp executor.ResultProcessor) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ConsumeResults", queue, rp)
	ret0, _ := ret[0].(error)
	return ret0
}

// ConsumeResults indicates an expected call of ConsumeResults.
func (mr *MockCommunicationServiceMockRecorder) ConsumeResults(queue, rp any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ConsumeResults", reflect.TypeOf((*MockCommunicationService)(nil).ConsumeResults), queue, rp)
}

// DeclareResultsQueue mocks base method.
func (m *MockCommunicationService) DeclareResultsQueue(queueName string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeclareResultsQueue", queueName)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeclareResultsQueue indicates an expected call of DeclareResultsQueue.
func (mr *MockCommunicationServiceMockRecorder) DeclareResultsQueue(queueName any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeclareResultsQueue", reflect.TypeOf((*MockCommunicationService)(nil).DeclareResultsQueue), queueName)
}

// RequestCalculation mocks base method.
func (m *MockCommunicationService) RequestCalculation(cd *dto.CalculationData) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RequestCalculation", cd)
}

// RequestCalculation indicates an expected call of RequestCalculation.
func (mr *MockCommunicationServiceMockRecorder) RequestCalculation(cd any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RequestCalculation", reflect.TypeOf((*MockCommunicationService)(nil).RequestCalculation), cd)
}
