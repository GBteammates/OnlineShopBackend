// Code generated by MockGen. DO NOT EDIT.
// Source: internal/usecase/usecase_interface.go

// Package mocks is a generated GoMock package.
package mocks

import (
	user "OnlineShopBackend/internal/delivery/user"
	models "OnlineShopBackend/internal/models"
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
)

// MockIItemUsecase is a mock of IItemUsecase interface.
type MockIItemUsecase struct {
	ctrl     *gomock.Controller
	recorder *MockIItemUsecaseMockRecorder
}

// MockIItemUsecaseMockRecorder is the mock recorder for MockIItemUsecase.
type MockIItemUsecaseMockRecorder struct {
	mock *MockIItemUsecase
}

// NewMockIItemUsecase creates a new mock instance.
func NewMockIItemUsecase(ctrl *gomock.Controller) *MockIItemUsecase {
	mock := &MockIItemUsecase{ctrl: ctrl}
	mock.recorder = &MockIItemUsecaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIItemUsecase) EXPECT() *MockIItemUsecaseMockRecorder {
	return m.recorder
}

// AddFavouriteItem mocks base method.
func (m *MockIItemUsecase) AddFavouriteItem(ctx context.Context, userId, itemId uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddFavouriteItem", ctx, userId, itemId)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddFavouriteItem indicates an expected call of AddFavouriteItem.
func (mr *MockIItemUsecaseMockRecorder) AddFavouriteItem(ctx, userId, itemId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddFavouriteItem", reflect.TypeOf((*MockIItemUsecase)(nil).AddFavouriteItem), ctx, userId, itemId)
}

// CreateItem mocks base method.
func (m *MockIItemUsecase) CreateItem(ctx context.Context, item *models.Item) (uuid.UUID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateItem", ctx, item)
	ret0, _ := ret[0].(uuid.UUID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateItem indicates an expected call of CreateItem.
func (mr *MockIItemUsecaseMockRecorder) CreateItem(ctx, item interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateItem", reflect.TypeOf((*MockIItemUsecase)(nil).CreateItem), ctx, item)
}

// DeleteFavouriteItem mocks base method.
func (m *MockIItemUsecase) DeleteFavouriteItem(ctx context.Context, userId, itemId uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteFavouriteItem", ctx, userId, itemId)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteFavouriteItem indicates an expected call of DeleteFavouriteItem.
func (mr *MockIItemUsecaseMockRecorder) DeleteFavouriteItem(ctx, userId, itemId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteFavouriteItem", reflect.TypeOf((*MockIItemUsecase)(nil).DeleteFavouriteItem), ctx, userId, itemId)
}

// DeleteItem mocks base method.
func (m *MockIItemUsecase) DeleteItem(ctx context.Context, id uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteItem", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteItem indicates an expected call of DeleteItem.
func (mr *MockIItemUsecaseMockRecorder) DeleteItem(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteItem", reflect.TypeOf((*MockIItemUsecase)(nil).DeleteItem), ctx, id)
}

// GetFavouriteItems mocks base method.
func (m *MockIItemUsecase) GetFavouriteItems(ctx context.Context, userId uuid.UUID, limitOptions map[string]int, sortOptions map[string]string) ([]models.Item, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFavouriteItems", ctx, userId, limitOptions, sortOptions)
	ret0, _ := ret[0].([]models.Item)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFavouriteItems indicates an expected call of GetFavouriteItems.
func (mr *MockIItemUsecaseMockRecorder) GetFavouriteItems(ctx, userId, limitOptions, sortOptions interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFavouriteItems", reflect.TypeOf((*MockIItemUsecase)(nil).GetFavouriteItems), ctx, userId, limitOptions, sortOptions)
}

// GetFavouriteItemsId mocks base method.
func (m *MockIItemUsecase) GetFavouriteItemsId(ctx context.Context, userId uuid.UUID) (*map[uuid.UUID]uuid.UUID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFavouriteItemsId", ctx, userId)
	ret0, _ := ret[0].(*map[uuid.UUID]uuid.UUID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFavouriteItemsId indicates an expected call of GetFavouriteItemsId.
func (mr *MockIItemUsecaseMockRecorder) GetFavouriteItemsId(ctx, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFavouriteItemsId", reflect.TypeOf((*MockIItemUsecase)(nil).GetFavouriteItemsId), ctx, userId)
}

// GetItem mocks base method.
func (m *MockIItemUsecase) GetItem(ctx context.Context, id uuid.UUID) (*models.Item, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetItem", ctx, id)
	ret0, _ := ret[0].(*models.Item)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetItem indicates an expected call of GetItem.
func (mr *MockIItemUsecaseMockRecorder) GetItem(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetItem", reflect.TypeOf((*MockIItemUsecase)(nil).GetItem), ctx, id)
}

// GetItemsByCategory mocks base method.
func (m *MockIItemUsecase) GetItemsByCategory(ctx context.Context, categoryName string, limitOptions map[string]int, sortOptions map[string]string) ([]models.Item, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetItemsByCategory", ctx, categoryName, limitOptions, sortOptions)
	ret0, _ := ret[0].([]models.Item)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetItemsByCategory indicates an expected call of GetItemsByCategory.
func (mr *MockIItemUsecaseMockRecorder) GetItemsByCategory(ctx, categoryName, limitOptions, sortOptions interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetItemsByCategory", reflect.TypeOf((*MockIItemUsecase)(nil).GetItemsByCategory), ctx, categoryName, limitOptions, sortOptions)
}

// ItemsList mocks base method.
func (m *MockIItemUsecase) ItemsList(ctx context.Context, limitOptions map[string]int, sortOptions map[string]string) ([]models.Item, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ItemsList", ctx, limitOptions, sortOptions)
	ret0, _ := ret[0].([]models.Item)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ItemsList indicates an expected call of ItemsList.
func (mr *MockIItemUsecaseMockRecorder) ItemsList(ctx, limitOptions, sortOptions interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ItemsList", reflect.TypeOf((*MockIItemUsecase)(nil).ItemsList), ctx, limitOptions, sortOptions)
}

// ItemsQuantity mocks base method.
func (m *MockIItemUsecase) ItemsQuantity(ctx context.Context) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ItemsQuantity", ctx)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ItemsQuantity indicates an expected call of ItemsQuantity.
func (mr *MockIItemUsecaseMockRecorder) ItemsQuantity(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ItemsQuantity", reflect.TypeOf((*MockIItemUsecase)(nil).ItemsQuantity), ctx)
}

// ItemsQuantityInCategory mocks base method.
func (m *MockIItemUsecase) ItemsQuantityInCategory(ctx context.Context, categoryName string) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ItemsQuantityInCategory", ctx, categoryName)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ItemsQuantityInCategory indicates an expected call of ItemsQuantityInCategory.
func (mr *MockIItemUsecaseMockRecorder) ItemsQuantityInCategory(ctx, categoryName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ItemsQuantityInCategory", reflect.TypeOf((*MockIItemUsecase)(nil).ItemsQuantityInCategory), ctx, categoryName)
}

// ItemsQuantityInFavourite mocks base method.
func (m *MockIItemUsecase) ItemsQuantityInFavourite(ctx context.Context, userId uuid.UUID) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ItemsQuantityInFavourite", ctx, userId)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ItemsQuantityInFavourite indicates an expected call of ItemsQuantityInFavourite.
func (mr *MockIItemUsecaseMockRecorder) ItemsQuantityInFavourite(ctx, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ItemsQuantityInFavourite", reflect.TypeOf((*MockIItemUsecase)(nil).ItemsQuantityInFavourite), ctx, userId)
}

// ItemsQuantityInSearch mocks base method.
func (m *MockIItemUsecase) ItemsQuantityInSearch(ctx context.Context, search string) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ItemsQuantityInSearch", ctx, search)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ItemsQuantityInSearch indicates an expected call of ItemsQuantityInSearch.
func (mr *MockIItemUsecaseMockRecorder) ItemsQuantityInSearch(ctx, search interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ItemsQuantityInSearch", reflect.TypeOf((*MockIItemUsecase)(nil).ItemsQuantityInSearch), ctx, search)
}

// SearchLine mocks base method.
func (m *MockIItemUsecase) SearchLine(ctx context.Context, param string, limitOptions map[string]int, sortOptions map[string]string) ([]models.Item, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchLine", ctx, param, limitOptions, sortOptions)
	ret0, _ := ret[0].([]models.Item)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SearchLine indicates an expected call of SearchLine.
func (mr *MockIItemUsecaseMockRecorder) SearchLine(ctx, param, limitOptions, sortOptions interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchLine", reflect.TypeOf((*MockIItemUsecase)(nil).SearchLine), ctx, param, limitOptions, sortOptions)
}

// SortItems mocks base method.
func (m *MockIItemUsecase) SortItems(items []models.Item, sortType, sortOrder string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SortItems", items, sortType, sortOrder)
}

// SortItems indicates an expected call of SortItems.
func (mr *MockIItemUsecaseMockRecorder) SortItems(items, sortType, sortOrder interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SortItems", reflect.TypeOf((*MockIItemUsecase)(nil).SortItems), items, sortType, sortOrder)
}

// UpdateCash mocks base method.
func (m *MockIItemUsecase) UpdateCash(ctx context.Context, id uuid.UUID, op string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateCash", ctx, id, op)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateCash indicates an expected call of UpdateCash.
func (mr *MockIItemUsecaseMockRecorder) UpdateCash(ctx, id, op interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateCash", reflect.TypeOf((*MockIItemUsecase)(nil).UpdateCash), ctx, id, op)
}

// UpdateFavouriteItemsCash mocks base method.
func (m *MockIItemUsecase) UpdateFavouriteItemsCash(ctx context.Context, userId, itemId uuid.UUID, op string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "UpdateFavouriteItemsCash", ctx, userId, itemId, op)
}

// UpdateFavouriteItemsCash indicates an expected call of UpdateFavouriteItemsCash.
func (mr *MockIItemUsecaseMockRecorder) UpdateFavouriteItemsCash(ctx, userId, itemId, op interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateFavouriteItemsCash", reflect.TypeOf((*MockIItemUsecase)(nil).UpdateFavouriteItemsCash), ctx, userId, itemId, op)
}

// UpdateItem mocks base method.
func (m *MockIItemUsecase) UpdateItem(ctx context.Context, item *models.Item) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateItem", ctx, item)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateItem indicates an expected call of UpdateItem.
func (mr *MockIItemUsecaseMockRecorder) UpdateItem(ctx, item interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateItem", reflect.TypeOf((*MockIItemUsecase)(nil).UpdateItem), ctx, item)
}

// UpdateItemsInCategoryCash mocks base method.
func (m *MockIItemUsecase) UpdateItemsInCategoryCash(ctx context.Context, newItem *models.Item, op string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateItemsInCategoryCash", ctx, newItem, op)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateItemsInCategoryCash indicates an expected call of UpdateItemsInCategoryCash.
func (mr *MockIItemUsecaseMockRecorder) UpdateItemsInCategoryCash(ctx, newItem, op interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateItemsInCategoryCash", reflect.TypeOf((*MockIItemUsecase)(nil).UpdateItemsInCategoryCash), ctx, newItem, op)
}

// MockICategoryUsecase is a mock of ICategoryUsecase interface.
type MockICategoryUsecase struct {
	ctrl     *gomock.Controller
	recorder *MockICategoryUsecaseMockRecorder
}

// MockICategoryUsecaseMockRecorder is the mock recorder for MockICategoryUsecase.
type MockICategoryUsecaseMockRecorder struct {
	mock *MockICategoryUsecase
}

// NewMockICategoryUsecase creates a new mock instance.
func NewMockICategoryUsecase(ctrl *gomock.Controller) *MockICategoryUsecase {
	mock := &MockICategoryUsecase{ctrl: ctrl}
	mock.recorder = &MockICategoryUsecaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockICategoryUsecase) EXPECT() *MockICategoryUsecaseMockRecorder {
	return m.recorder
}

// CreateCategory mocks base method.
func (m *MockICategoryUsecase) CreateCategory(ctx context.Context, category *models.Category) (uuid.UUID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateCategory", ctx, category)
	ret0, _ := ret[0].(uuid.UUID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateCategory indicates an expected call of CreateCategory.
func (mr *MockICategoryUsecaseMockRecorder) CreateCategory(ctx, category interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateCategory", reflect.TypeOf((*MockICategoryUsecase)(nil).CreateCategory), ctx, category)
}

// DeleteCategory mocks base method.
func (m *MockICategoryUsecase) DeleteCategory(ctx context.Context, id uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteCategory", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteCategory indicates an expected call of DeleteCategory.
func (mr *MockICategoryUsecaseMockRecorder) DeleteCategory(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteCategory", reflect.TypeOf((*MockICategoryUsecase)(nil).DeleteCategory), ctx, id)
}

// DeleteCategoryCash mocks base method.
func (m *MockICategoryUsecase) DeleteCategoryCash(ctx context.Context, name string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteCategoryCash", ctx, name)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteCategoryCash indicates an expected call of DeleteCategoryCash.
func (mr *MockICategoryUsecaseMockRecorder) DeleteCategoryCash(ctx, name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteCategoryCash", reflect.TypeOf((*MockICategoryUsecase)(nil).DeleteCategoryCash), ctx, name)
}

// GetCategory mocks base method.
func (m *MockICategoryUsecase) GetCategory(ctx context.Context, id uuid.UUID) (*models.Category, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCategory", ctx, id)
	ret0, _ := ret[0].(*models.Category)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCategory indicates an expected call of GetCategory.
func (mr *MockICategoryUsecaseMockRecorder) GetCategory(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCategory", reflect.TypeOf((*MockICategoryUsecase)(nil).GetCategory), ctx, id)
}

// GetCategoryByName mocks base method.
func (m *MockICategoryUsecase) GetCategoryByName(ctx context.Context, name string) (*models.Category, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCategoryByName", ctx, name)
	ret0, _ := ret[0].(*models.Category)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCategoryByName indicates an expected call of GetCategoryByName.
func (mr *MockICategoryUsecaseMockRecorder) GetCategoryByName(ctx, name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCategoryByName", reflect.TypeOf((*MockICategoryUsecase)(nil).GetCategoryByName), ctx, name)
}

// GetCategoryList mocks base method.
func (m *MockICategoryUsecase) GetCategoryList(ctx context.Context) ([]models.Category, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCategoryList", ctx)
	ret0, _ := ret[0].([]models.Category)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCategoryList indicates an expected call of GetCategoryList.
func (mr *MockICategoryUsecaseMockRecorder) GetCategoryList(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCategoryList", reflect.TypeOf((*MockICategoryUsecase)(nil).GetCategoryList), ctx)
}

// UpdateCash mocks base method.
func (m *MockICategoryUsecase) UpdateCash(ctx context.Context, id uuid.UUID, op string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateCash", ctx, id, op)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateCash indicates an expected call of UpdateCash.
func (mr *MockICategoryUsecaseMockRecorder) UpdateCash(ctx, id, op interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateCash", reflect.TypeOf((*MockICategoryUsecase)(nil).UpdateCash), ctx, id, op)
}

// UpdateCategory mocks base method.
func (m *MockICategoryUsecase) UpdateCategory(ctx context.Context, category *models.Category) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateCategory", ctx, category)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateCategory indicates an expected call of UpdateCategory.
func (mr *MockICategoryUsecaseMockRecorder) UpdateCategory(ctx, category interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateCategory", reflect.TypeOf((*MockICategoryUsecase)(nil).UpdateCategory), ctx, category)
}

// MockIOrderUsecase is a mock of IOrderUsecase interface.
type MockIOrderUsecase struct {
	ctrl     *gomock.Controller
	recorder *MockIOrderUsecaseMockRecorder
}

// MockIOrderUsecaseMockRecorder is the mock recorder for MockIOrderUsecase.
type MockIOrderUsecaseMockRecorder struct {
	mock *MockIOrderUsecase
}

// NewMockIOrderUsecase creates a new mock instance.
func NewMockIOrderUsecase(ctrl *gomock.Controller) *MockIOrderUsecase {
	mock := &MockIOrderUsecase{ctrl: ctrl}
	mock.recorder = &MockIOrderUsecaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIOrderUsecase) EXPECT() *MockIOrderUsecaseMockRecorder {
	return m.recorder
}

// ChangeAddress mocks base method.
func (m *MockIOrderUsecase) ChangeAddress(ctx context.Context, order *models.Order, newAddress models.UserAddress) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ChangeAddress", ctx, order, newAddress)
	ret0, _ := ret[0].(error)
	return ret0
}

// ChangeAddress indicates an expected call of ChangeAddress.
func (mr *MockIOrderUsecaseMockRecorder) ChangeAddress(ctx, order, newAddress interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChangeAddress", reflect.TypeOf((*MockIOrderUsecase)(nil).ChangeAddress), ctx, order, newAddress)
}

// ChangeStatus mocks base method.
func (m *MockIOrderUsecase) ChangeStatus(ctx context.Context, order *models.Order, newStatus models.Status) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ChangeStatus", ctx, order, newStatus)
	ret0, _ := ret[0].(error)
	return ret0
}

// ChangeStatus indicates an expected call of ChangeStatus.
func (mr *MockIOrderUsecaseMockRecorder) ChangeStatus(ctx, order, newStatus interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChangeStatus", reflect.TypeOf((*MockIOrderUsecase)(nil).ChangeStatus), ctx, order, newStatus)
}

// DeleteOrder mocks base method.
func (m *MockIOrderUsecase) DeleteOrder(ctx context.Context, order *models.Order) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteOrder", ctx, order)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteOrder indicates an expected call of DeleteOrder.
func (mr *MockIOrderUsecaseMockRecorder) DeleteOrder(ctx, order interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteOrder", reflect.TypeOf((*MockIOrderUsecase)(nil).DeleteOrder), ctx, order)
}

// GetOrder mocks base method.
func (m *MockIOrderUsecase) GetOrder(ctx context.Context, id uuid.UUID) (*models.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOrder", ctx, id)
	ret0, _ := ret[0].(*models.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOrder indicates an expected call of GetOrder.
func (mr *MockIOrderUsecaseMockRecorder) GetOrder(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrder", reflect.TypeOf((*MockIOrderUsecase)(nil).GetOrder), ctx, id)
}

// GetOrdersForUser mocks base method.
func (m *MockIOrderUsecase) GetOrdersForUser(ctx context.Context, user *models.User) ([]models.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOrdersForUser", ctx, user)
	ret0, _ := ret[0].([]models.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOrdersForUser indicates an expected call of GetOrdersForUser.
func (mr *MockIOrderUsecaseMockRecorder) GetOrdersForUser(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrdersForUser", reflect.TypeOf((*MockIOrderUsecase)(nil).GetOrdersForUser), ctx, user)
}

// PlaceOrder mocks base method.
func (m *MockIOrderUsecase) PlaceOrder(ctx context.Context, cart *models.Cart, user models.User, address models.UserAddress) (*models.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PlaceOrder", ctx, cart, user, address)
	ret0, _ := ret[0].(*models.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PlaceOrder indicates an expected call of PlaceOrder.
func (mr *MockIOrderUsecaseMockRecorder) PlaceOrder(ctx, cart, user, address interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PlaceOrder", reflect.TypeOf((*MockIOrderUsecase)(nil).PlaceOrder), ctx, cart, user, address)
}

// MockICartUsecase is a mock of ICartUsecase interface.
type MockICartUsecase struct {
	ctrl     *gomock.Controller
	recorder *MockICartUsecaseMockRecorder
}

// MockICartUsecaseMockRecorder is the mock recorder for MockICartUsecase.
type MockICartUsecaseMockRecorder struct {
	mock *MockICartUsecase
}

// NewMockICartUsecase creates a new mock instance.
func NewMockICartUsecase(ctrl *gomock.Controller) *MockICartUsecase {
	mock := &MockICartUsecase{ctrl: ctrl}
	mock.recorder = &MockICartUsecaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockICartUsecase) EXPECT() *MockICartUsecaseMockRecorder {
	return m.recorder
}

// AddItemToCart mocks base method.
func (m *MockICartUsecase) AddItemToCart(ctx context.Context, cartId, itemId uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddItemToCart", ctx, cartId, itemId)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddItemToCart indicates an expected call of AddItemToCart.
func (mr *MockICartUsecaseMockRecorder) AddItemToCart(ctx, cartId, itemId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddItemToCart", reflect.TypeOf((*MockICartUsecase)(nil).AddItemToCart), ctx, cartId, itemId)
}

// Create mocks base method.
func (m *MockICartUsecase) Create(ctx context.Context, userId uuid.UUID) (uuid.UUID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, userId)
	ret0, _ := ret[0].(uuid.UUID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockICartUsecaseMockRecorder) Create(ctx, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockICartUsecase)(nil).Create), ctx, userId)
}

// DeleteCart mocks base method.
func (m *MockICartUsecase) DeleteCart(ctx context.Context, cartId uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteCart", ctx, cartId)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteCart indicates an expected call of DeleteCart.
func (mr *MockICartUsecaseMockRecorder) DeleteCart(ctx, cartId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteCart", reflect.TypeOf((*MockICartUsecase)(nil).DeleteCart), ctx, cartId)
}

// DeleteItemFromCart mocks base method.
func (m *MockICartUsecase) DeleteItemFromCart(ctx context.Context, cartId, itemId uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteItemFromCart", ctx, cartId, itemId)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteItemFromCart indicates an expected call of DeleteItemFromCart.
func (mr *MockICartUsecaseMockRecorder) DeleteItemFromCart(ctx, cartId, itemId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteItemFromCart", reflect.TypeOf((*MockICartUsecase)(nil).DeleteItemFromCart), ctx, cartId, itemId)
}

// GetCart mocks base method.
func (m *MockICartUsecase) GetCart(ctx context.Context, cartId uuid.UUID) (*models.Cart, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCart", ctx, cartId)
	ret0, _ := ret[0].(*models.Cart)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCart indicates an expected call of GetCart.
func (mr *MockICartUsecaseMockRecorder) GetCart(ctx, cartId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCart", reflect.TypeOf((*MockICartUsecase)(nil).GetCart), ctx, cartId)
}

// GetCartByUserId mocks base method.
func (m *MockICartUsecase) GetCartByUserId(ctx context.Context, userId uuid.UUID) (*models.Cart, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCartByUserId", ctx, userId)
	ret0, _ := ret[0].(*models.Cart)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCartByUserId indicates an expected call of GetCartByUserId.
func (mr *MockICartUsecaseMockRecorder) GetCartByUserId(ctx, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCartByUserId", reflect.TypeOf((*MockICartUsecase)(nil).GetCartByUserId), ctx, userId)
}

// MockIUserUsecase is a mock of IUserUsecase interface.
type MockIUserUsecase struct {
	ctrl     *gomock.Controller
	recorder *MockIUserUsecaseMockRecorder
}

// MockIUserUsecaseMockRecorder is the mock recorder for MockIUserUsecase.
type MockIUserUsecaseMockRecorder struct {
	mock *MockIUserUsecase
}

// NewMockIUserUsecase creates a new mock instance.
func NewMockIUserUsecase(ctrl *gomock.Controller) *MockIUserUsecase {
	mock := &MockIUserUsecase{ctrl: ctrl}
	mock.recorder = &MockIUserUsecaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIUserUsecase) EXPECT() *MockIUserUsecaseMockRecorder {
	return m.recorder
}

// CreateRights mocks base method.
func (m *MockIUserUsecase) CreateRights(ctx context.Context, rights *models.Rights) (uuid.UUID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateRights", ctx, rights)
	ret0, _ := ret[0].(uuid.UUID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateRights indicates an expected call of CreateRights.
func (mr *MockIUserUsecaseMockRecorder) CreateRights(ctx, rights interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateRights", reflect.TypeOf((*MockIUserUsecase)(nil).CreateRights), ctx, rights)
}

// CreateUser mocks base method.
func (m *MockIUserUsecase) CreateUser(ctx context.Context, user *user.CreateUserData) (*models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", ctx, user)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockIUserUsecaseMockRecorder) CreateUser(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockIUserUsecase)(nil).CreateUser), ctx, user)
}

// GetRightsId mocks base method.
func (m *MockIUserUsecase) GetRightsId(ctx context.Context, name string) (*models.Rights, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRightsId", ctx, name)
	ret0, _ := ret[0].(*models.Rights)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRightsId indicates an expected call of GetRightsId.
func (mr *MockIUserUsecaseMockRecorder) GetRightsId(ctx, name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRightsId", reflect.TypeOf((*MockIUserUsecase)(nil).GetRightsId), ctx, name)
}

// GetRightsList mocks base method.
func (m *MockIUserUsecase) GetRightsList(ctx context.Context) ([]models.Rights, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRightsList", ctx)
	ret0, _ := ret[0].([]models.Rights)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRightsList indicates an expected call of GetRightsList.
func (mr *MockIUserUsecaseMockRecorder) GetRightsList(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRightsList", reflect.TypeOf((*MockIUserUsecase)(nil).GetRightsList), ctx)
}

// GetUserByEmail mocks base method.
func (m *MockIUserUsecase) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByEmail", ctx, email)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByEmail indicates an expected call of GetUserByEmail.
func (mr *MockIUserUsecaseMockRecorder) GetUserByEmail(ctx, email interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByEmail", reflect.TypeOf((*MockIUserUsecase)(nil).GetUserByEmail), ctx, email)
}

// UpdateUserData mocks base method.
func (m *MockIUserUsecase) UpdateUserData(ctx context.Context, id uuid.UUID, user *user.CreateUserData) (*models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUserData", ctx, id, user)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateUserData indicates an expected call of UpdateUserData.
func (mr *MockIUserUsecaseMockRecorder) UpdateUserData(ctx, id, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUserData", reflect.TypeOf((*MockIUserUsecase)(nil).UpdateUserData), ctx, id, user)
}

// UpdateUserRole mocks base method.
func (m *MockIUserUsecase) UpdateUserRole(ctx context.Context, roleId uuid.UUID, email string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUserRole", ctx, roleId, email)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateUserRole indicates an expected call of UpdateUserRole.
func (mr *MockIUserUsecaseMockRecorder) UpdateUserRole(ctx, roleId, email interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUserRole", reflect.TypeOf((*MockIUserUsecase)(nil).UpdateUserRole), ctx, roleId, email)
}
