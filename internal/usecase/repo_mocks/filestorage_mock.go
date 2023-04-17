// Code generated by MockGen. DO NOT EDIT.
// Source: internal/usecase/fs_interface.go

// Package mocks is a generated GoMock package.
package mocks

import (
	models "OnlineShopBackend/internal/models"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockFileStorager is a mock of FileStorager interface.
type MockFileStorager struct {
	ctrl     *gomock.Controller
	recorder *MockFileStoragerMockRecorder
}

// MockFileStoragerMockRecorder is the mock recorder for MockFileStorager.
type MockFileStoragerMockRecorder struct {
	mock *MockFileStorager
}

// NewMockFileStorager creates a new mock instance.
func NewMockFileStorager(ctrl *gomock.Controller) *MockFileStorager {
	mock := &MockFileStorager{ctrl: ctrl}
	mock.recorder = &MockFileStoragerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockFileStorager) EXPECT() *MockFileStoragerMockRecorder {
	return m.recorder
}

// DeleteCategoryImage mocks base method.
func (m *MockFileStorager) DeleteCategoryImage(id, filename string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteCategoryImage", id, filename)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteCategoryImage indicates an expected call of DeleteCategoryImage.
func (mr *MockFileStoragerMockRecorder) DeleteCategoryImage(id, filename interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteCategoryImage", reflect.TypeOf((*MockFileStorager)(nil).DeleteCategoryImage), id, filename)
}

// DeleteCategoryImageById mocks base method.
func (m *MockFileStorager) DeleteCategoryImageById(id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteCategoryImageById", id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteCategoryImageById indicates an expected call of DeleteCategoryImageById.
func (mr *MockFileStoragerMockRecorder) DeleteCategoryImageById(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteCategoryImageById", reflect.TypeOf((*MockFileStorager)(nil).DeleteCategoryImageById), id)
}

// DeleteItemImage mocks base method.
func (m *MockFileStorager) DeleteItemImage(id, filename string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteItemImage", id, filename)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteItemImage indicates an expected call of DeleteItemImage.
func (mr *MockFileStoragerMockRecorder) DeleteItemImage(id, filename interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteItemImage", reflect.TypeOf((*MockFileStorager)(nil).DeleteItemImage), id, filename)
}

// DeleteItemImagesFolderById mocks base method.
func (m *MockFileStorager) DeleteItemImagesFolderById(id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteItemImagesFolderById", id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteItemImagesFolderById indicates an expected call of DeleteItemImagesFolderById.
func (mr *MockFileStoragerMockRecorder) DeleteItemImagesFolderById(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteItemImagesFolderById", reflect.TypeOf((*MockFileStorager)(nil).DeleteItemImagesFolderById), id)
}

// GetCategoriesImagesList mocks base method.
func (m *MockFileStorager) GetCategoriesImagesList() ([]*models.FileInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCategoriesImagesList")
	ret0, _ := ret[0].([]*models.FileInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCategoriesImagesList indicates an expected call of GetCategoriesImagesList.
func (mr *MockFileStoragerMockRecorder) GetCategoriesImagesList() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCategoriesImagesList", reflect.TypeOf((*MockFileStorager)(nil).GetCategoriesImagesList))
}

// GetItemsImagesList mocks base method.
func (m *MockFileStorager) GetItemsImagesList() ([]*models.FileInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetItemsImagesList")
	ret0, _ := ret[0].([]*models.FileInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetItemsImagesList indicates an expected call of GetItemsImagesList.
func (mr *MockFileStoragerMockRecorder) GetItemsImagesList() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetItemsImagesList", reflect.TypeOf((*MockFileStorager)(nil).GetItemsImagesList))
}

// PutCategoryImage mocks base method.
func (m *MockFileStorager) PutCategoryImage(id, filename string, file []byte) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PutCategoryImage", id, filename, file)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PutCategoryImage indicates an expected call of PutCategoryImage.
func (mr *MockFileStoragerMockRecorder) PutCategoryImage(id, filename, file interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PutCategoryImage", reflect.TypeOf((*MockFileStorager)(nil).PutCategoryImage), id, filename, file)
}

// PutItemImage mocks base method.
func (m *MockFileStorager) PutItemImage(id, filename string, file []byte) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PutItemImage", id, filename, file)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PutItemImage indicates an expected call of PutItemImage.
func (mr *MockFileStoragerMockRecorder) PutItemImage(id, filename, file interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PutItemImage", reflect.TypeOf((*MockFileStorager)(nil).PutItemImage), id, filename, file)
}
