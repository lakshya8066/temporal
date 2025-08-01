// Code generated by MockGen. DO NOT EDIT.
// Source: transaction_manager.go
//
// Generated by this command:
//
//	mockgen -package ndc -source transaction_manager.go -destination transaction_manager_mock.go
//

// Package ndc is a generated GoMock package.
package ndc

import (
	context "context"
	reflect "reflect"

	chasm "go.temporal.io/server/chasm"
	namespace "go.temporal.io/server/common/namespace"
	persistence "go.temporal.io/server/common/persistence"
	gomock "go.uber.org/mock/gomock"
)

// MockTransactionManager is a mock of TransactionManager interface.
type MockTransactionManager struct {
	ctrl     *gomock.Controller
	recorder *MockTransactionManagerMockRecorder
	isgomock struct{}
}

// MockTransactionManagerMockRecorder is the mock recorder for MockTransactionManager.
type MockTransactionManagerMockRecorder struct {
	mock *MockTransactionManager
}

// NewMockTransactionManager creates a new mock instance.
func NewMockTransactionManager(ctrl *gomock.Controller) *MockTransactionManager {
	mock := &MockTransactionManager{ctrl: ctrl}
	mock.recorder = &MockTransactionManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTransactionManager) EXPECT() *MockTransactionManagerMockRecorder {
	return m.recorder
}

// BackfillWorkflow mocks base method.
func (m *MockTransactionManager) BackfillWorkflow(ctx context.Context, targetWorkflow Workflow, targetWorkflowEventsSlice ...*persistence.WorkflowEvents) error {
	m.ctrl.T.Helper()
	varargs := []any{ctx, targetWorkflow}
	for _, a := range targetWorkflowEventsSlice {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "BackfillWorkflow", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// BackfillWorkflow indicates an expected call of BackfillWorkflow.
func (mr *MockTransactionManagerMockRecorder) BackfillWorkflow(ctx, targetWorkflow any, targetWorkflowEventsSlice ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, targetWorkflow}, targetWorkflowEventsSlice...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BackfillWorkflow", reflect.TypeOf((*MockTransactionManager)(nil).BackfillWorkflow), varargs...)
}

// CheckWorkflowExists mocks base method.
func (m *MockTransactionManager) CheckWorkflowExists(ctx context.Context, namespaceID namespace.ID, workflowID, runID string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckWorkflowExists", ctx, namespaceID, workflowID, runID)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckWorkflowExists indicates an expected call of CheckWorkflowExists.
func (mr *MockTransactionManagerMockRecorder) CheckWorkflowExists(ctx, namespaceID, workflowID, runID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckWorkflowExists", reflect.TypeOf((*MockTransactionManager)(nil).CheckWorkflowExists), ctx, namespaceID, workflowID, runID)
}

// CreateWorkflow mocks base method.
func (m *MockTransactionManager) CreateWorkflow(ctx context.Context, targetWorkflow Workflow) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateWorkflow", ctx, targetWorkflow)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateWorkflow indicates an expected call of CreateWorkflow.
func (mr *MockTransactionManagerMockRecorder) CreateWorkflow(ctx, targetWorkflow any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateWorkflow", reflect.TypeOf((*MockTransactionManager)(nil).CreateWorkflow), ctx, targetWorkflow)
}

// GetCurrentWorkflowRunID mocks base method.
func (m *MockTransactionManager) GetCurrentWorkflowRunID(ctx context.Context, namespaceID namespace.ID, workflowID string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCurrentWorkflowRunID", ctx, namespaceID, workflowID)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCurrentWorkflowRunID indicates an expected call of GetCurrentWorkflowRunID.
func (mr *MockTransactionManagerMockRecorder) GetCurrentWorkflowRunID(ctx, namespaceID, workflowID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCurrentWorkflowRunID", reflect.TypeOf((*MockTransactionManager)(nil).GetCurrentWorkflowRunID), ctx, namespaceID, workflowID)
}

// LoadWorkflow mocks base method.
func (m *MockTransactionManager) LoadWorkflow(ctx context.Context, namespaceID namespace.ID, workflowID, runID string, archetype chasm.Archetype) (Workflow, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LoadWorkflow", ctx, namespaceID, workflowID, runID, archetype)
	ret0, _ := ret[0].(Workflow)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LoadWorkflow indicates an expected call of LoadWorkflow.
func (mr *MockTransactionManagerMockRecorder) LoadWorkflow(ctx, namespaceID, workflowID, runID, archetype any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoadWorkflow", reflect.TypeOf((*MockTransactionManager)(nil).LoadWorkflow), ctx, namespaceID, workflowID, runID, archetype)
}

// UpdateWorkflow mocks base method.
func (m *MockTransactionManager) UpdateWorkflow(ctx context.Context, isWorkflowRebuilt bool, targetWorkflow, newWorkflow Workflow) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateWorkflow", ctx, isWorkflowRebuilt, targetWorkflow, newWorkflow)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateWorkflow indicates an expected call of UpdateWorkflow.
func (mr *MockTransactionManagerMockRecorder) UpdateWorkflow(ctx, isWorkflowRebuilt, targetWorkflow, newWorkflow any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateWorkflow", reflect.TypeOf((*MockTransactionManager)(nil).UpdateWorkflow), ctx, isWorkflowRebuilt, targetWorkflow, newWorkflow)
}
