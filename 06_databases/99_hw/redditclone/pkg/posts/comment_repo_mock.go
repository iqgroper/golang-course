// Code generated by MockGen. DO NOT EDIT.
// Source: posts/comment.go

// Package posts is a generated GoMock package.
package posts

import (
	user "redditclone/pkg/user"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockCommentRepo is a mock of CommentRepo interface.
type MockCommentRepo struct {
	ctrl     *gomock.Controller
	recorder *MockCommentRepoMockRecorder
}

// MockCommentRepoMockRecorder is the mock recorder for MockCommentRepo.
type MockCommentRepoMockRecorder struct {
	mock *MockCommentRepo
}

// NewMockCommentRepo creates a new mock instance.
func NewMockCommentRepo(ctrl *gomock.Controller) *MockCommentRepo {
	mock := &MockCommentRepo{ctrl: ctrl}
	mock.recorder = &MockCommentRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCommentRepo) EXPECT() *MockCommentRepoMockRecorder {
	return m.recorder
}

// Add mocks base method.
func (m *MockCommentRepo) Add(post_id, body string, user *user.User) (*Comment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Add", post_id, body, user)
	ret0, _ := ret[0].(*Comment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Add indicates an expected call of Add.
func (mr *MockCommentRepoMockRecorder) Add(post_id, body, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Add", reflect.TypeOf((*MockCommentRepo)(nil).Add), post_id, body, user)
}

// Delete mocks base method.
func (m *MockCommentRepo) Delete(post_id, comment_id string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", post_id, comment_id)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Delete indicates an expected call of Delete.
func (mr *MockCommentRepoMockRecorder) Delete(post_id, comment_id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockCommentRepo)(nil).Delete), post_id, comment_id)
}

// GetAll mocks base method.
func (m *MockCommentRepo) GetAll(post_id string) ([]*Comment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll", post_id)
	ret0, _ := ret[0].([]*Comment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll.
func (mr *MockCommentRepoMockRecorder) GetAll(post_id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockCommentRepo)(nil).GetAll), post_id)
}
