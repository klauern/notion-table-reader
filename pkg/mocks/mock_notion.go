// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/klauern/notion-table-reader/pkg (interfaces: NotionClient)
//
// Generated by this command:
//
//	mockgen -destination=mocks/mock_notion.go -package=mocks . NotionClient
//

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	pkg "github.com/klauern/notion-table-reader/pkg"
	gomock "go.uber.org/mock/gomock"
)

// MockNotionClient is a mock of NotionClient interface.
type MockNotionClient struct {
	ctrl     *gomock.Controller
	recorder *MockNotionClientMockRecorder
}

// MockNotionClientMockRecorder is the mock recorder for MockNotionClient.
type MockNotionClientMockRecorder struct {
	mock *MockNotionClient
}

// NewMockNotionClient creates a new mock instance.
func NewMockNotionClient(ctrl *gomock.Controller) *MockNotionClient {
	mock := &MockNotionClient{ctrl: ctrl}
	mock.recorder = &MockNotionClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockNotionClient) EXPECT() *MockNotionClientMockRecorder {
	return m.recorder
}

// FetchPages mocks base method.
func (m *MockNotionClient) FetchPages(arg0 string, arg1 bool) ([]pkg.PageDetail, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FetchPages", arg0, arg1)
	ret0, _ := ret[0].([]pkg.PageDetail)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FetchPages indicates an expected call of FetchPages.
func (mr *MockNotionClientMockRecorder) FetchPages(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FetchPages", reflect.TypeOf((*MockNotionClient)(nil).FetchPages), arg0, arg1)
}

// TagPage mocks base method.
func (m *MockNotionClient) TagPage(arg0 string, arg1 []string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TagPage", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// TagPage indicates an expected call of TagPage.
func (mr *MockNotionClientMockRecorder) TagPage(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TagPage", reflect.TypeOf((*MockNotionClient)(nil).TagPage), arg0, arg1)
}
