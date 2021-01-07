// Code generated by MockGen. DO NOT EDIT.
// Source: user_repository.go

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	model "github.com/hayashiki/mentions/pkg/model"
	reflect "reflect"
)

// MockUserRepository is a mock of UserRepository interface
type MockUserRepository struct {
	ctrl     *gomock.Controller
	recorder *MockUserRepositoryMockRecorder
}

// MockUserRepositoryMockRecorder is the mock recorder for MockUserRepository
type MockUserRepositoryMockRecorder struct {
	mock *MockUserRepository
}

// NewMockUserRepository creates a new mock instance
func NewMockUserRepository(ctrl *gomock.Controller) *MockUserRepository {
	mock := &MockUserRepository{ctrl: ctrl}
	mock.recorder = &MockUserRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockUserRepository) EXPECT() *MockUserRepositoryMockRecorder {
	return m.recorder
}

// List mocks base method
func (m *MockUserRepository) List(ctx context.Context, team *model.Team, cursor string, limit int) ([]*model.User, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", ctx, team, cursor, limit)
	ret0, _ := ret[0].([]*model.User)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// List indicates an expected call of List
func (mr *MockUserRepositoryMockRecorder) List(ctx, team, cursor, limit interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockUserRepository)(nil).List), ctx, team, cursor, limit)
}

// Put mocks base method
func (m *MockUserRepository) Put(ctx context.Context, team *model.Team, user *model.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Put", ctx, team, user)
	ret0, _ := ret[0].(error)
	return ret0
}

// Put indicates an expected call of Put
func (mr *MockUserRepositoryMockRecorder) Put(ctx, team, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Put", reflect.TypeOf((*MockUserRepository)(nil).Put), ctx, team, user)
}

// FindByGithubID mocks base method
func (m *MockUserRepository) FindByGithubID(ctx context.Context, githubID string) (*model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByGithubID", ctx, githubID)
	ret0, _ := ret[0].(*model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByGithubID indicates an expected call of FindByGithubID
func (mr *MockUserRepositoryMockRecorder) FindByGithubID(ctx, githubID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByGithubID", reflect.TypeOf((*MockUserRepository)(nil).FindByGithubID), ctx, githubID)
}

// FindBySlackID mocks base method
func (m *MockUserRepository) FindBySlackID(ctx context.Context, githubID string) (*model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindBySlackID", ctx, githubID)
	ret0, _ := ret[0].(*model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindBySlackID indicates an expected call of FindBySlackID
func (mr *MockUserRepositoryMockRecorder) FindBySlackID(ctx, githubID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindBySlackID", reflect.TypeOf((*MockUserRepository)(nil).FindBySlackID), ctx, githubID)
}

// Delete mocks base method
func (m *MockUserRepository) Delete(ctx context.Context, team *model.Team, id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, team, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete
func (mr *MockUserRepositoryMockRecorder) Delete(ctx, team, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockUserRepository)(nil).Delete), ctx, team, id)
}