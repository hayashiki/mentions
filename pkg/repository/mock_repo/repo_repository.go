// Code generated by MockGen. DO NOT EDIT.
// Source: repo_repository.go

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	model "github.com/hayashiki/mentions/pkg/model"
	reflect "reflect"
)

// MockRepoRepository is a mock of RepoRepository interface
type MockRepoRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepoRepositoryMockRecorder
}

// MockRepoRepositoryMockRecorder is the mock recorder for MockRepoRepository
type MockRepoRepositoryMockRecorder struct {
	mock *MockRepoRepository
}

// NewMockRepoRepository creates a new mock instance
func NewMockRepoRepository(ctrl *gomock.Controller) *MockRepoRepository {
	mock := &MockRepoRepository{ctrl: ctrl}
	mock.recorder = &MockRepoRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockRepoRepository) EXPECT() *MockRepoRepositoryMockRecorder {
	return m.recorder
}

// List mocks base method
func (m *MockRepoRepository) List(ctx context.Context, cursor string, limit int) ([]*model.Repo, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", ctx, cursor, limit)
	ret0, _ := ret[0].([]*model.Repo)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// List indicates an expected call of List
func (mr *MockRepoRepositoryMockRecorder) List(ctx, cursor, limit interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockRepoRepository)(nil).List), ctx, cursor, limit)
}

// Get mocks base method
func (m *MockRepoRepository) Get(ctx context.Context, id int64) (*model.Repo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, id)
	ret0, _ := ret[0].(*model.Repo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get
func (mr *MockRepoRepositoryMockRecorder) Get(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockRepoRepository)(nil).Get), ctx, id)
}

// Put mocks base method
func (m *MockRepoRepository) Put(ctx context.Context, repo *model.Repo) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Put", ctx, repo)
	ret0, _ := ret[0].(error)
	return ret0
}

// Put indicates an expected call of Put
func (mr *MockRepoRepositoryMockRecorder) Put(ctx, repo interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Put", reflect.TypeOf((*MockRepoRepository)(nil).Put), ctx, repo)
}

// Delete mocks base method
func (m *MockRepoRepository) Delete(ctx context.Context, id int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete
func (mr *MockRepoRepositoryMockRecorder) Delete(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockRepoRepository)(nil).Delete), ctx, id)
}
