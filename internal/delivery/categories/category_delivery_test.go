package delivery

import (
	"OnlineShopBackend/internal/delivery/categories/category"
	"OnlineShopBackend/internal/models"
	mocks "OnlineShopBackend/internal/usecase/usecase_mocks"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"go.uber.org/zap"
)

type WrongStruct struct {
	Name        int   `json:"name"`
	Description int32 `json:"description"`
}

var (
	err              = fmt.Errorf("error")
	testId           = uuid.New()
	testCategoryNoId = &models.Category{
		Name:        "testName",
		Description: "testDescription",
	}
	testShortCategory = category.ShortCategory{
		Name:        "testName",
		Description: "testDescription",
	}
	testShortNoCategory = category.ShortCategory{
		Name:        "NoCategory",
		Description: "testDescription",
	}
	testCategoryId = category.CategoryId{
		Value: testId.String(),
	}
	testCategoryWithId = category.Category{
		Id:          testId.String(),
		Name:        "testName",
		Description: "testDescription",
	}
	testCatNoCategoryWithId = category.Category{
		Id:          testId.String(),
		Name:        "NoCategory",
		Description: "testDescription",
	}
	testModelsCategoryWithId = &models.Category{
		Id:          testId,
		Name:        "testName",
		Description: "testDescription",
	}
	testEmptyCategory = category.ShortCategory{
		Name:        "",
		Description: "",
	}
	testEmptyModelsCategory = models.Category{
		Id:          uuid.Nil,
		Name:        "",
		Description: "",
	}
	testWrong = WrongStruct{
		Name:        5,
		Description: 6,
	}
	testList    = []models.Category{*testModelsCategoryWithId}
	testOutList = []category.Category{
		testCategoryWithId,
	}
	testNoCategory = models.Category{
		Name:        "NoCategory",
		Description: "Category for items from deleting categories",
	}
	testNoCategoryWithId = models.Category{
		Id:          testId,
		Name:        "NoCategory",
		Description: "Category for items from deleting categories",
	}
	testCategoryWithImage2 = &models.Category{
		Id:          testId,
		Name:        "testName",
		Description: "testDescription",
		Image:       "testImagePath",
	}
	testModelsItemNoCat = models.Item{
		Id:          testId,
		Title:       "testTitle",
		Description: "testDescription",
		Category:    testNoCategoryWithId,
		Price:       10,
		Vendor:      "testVendor",
	}
	post                 = "POST"
	put                  = "PUT"
	testModelsItemWithId = &models.Item{
		Id:          testId,
		Title:       "testTitle",
		Description: "testDescription",
		Category: models.Category{
			Id:          testId,
			Name:        "testName",
			Description: "testDescription",
		},
		Price:  10,
		Vendor: "testVendor",
	}
	testFile = []byte{0xff, 0xd8, 0xff, 0xe0, 0x0, 0x10, 0x4a, 0x46, 0x49, 0x46, 0x0, 0x1, 0x1, 0x1, 0x0, 0x48, 0x0, 0x48, 0x0, 0x0, 0xff, 0xe1, 0x0, 0x22, 0x45, 0x78, 0x69, 0x66, 0x0, 0x0, 0x4d, 0x4d, 0x0, 0x2a, 0x0, 0x0, 0x0, 0x8, 0x0, 0x1, 0x1, 0x12, 0x0, 0x3, 0x0, 0x0, 0x0, 0x1, 0x0, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xff, 0xfe, 0x0, 0xd, 0x53, 0x65, 0x63, 0x6c, 0x75, 0x62, 0x2e, 0x6f, 0x72, 0x67, 0x0, 0xff, 0xdb, 0x0, 0x43, 0x0, 0x2, 0x1, 0x1, 0x2, 0x1, 0x1, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x3, 0x5, 0x3, 0x3, 0x3, 0x3, 0x3, 0x6, 0x4, 0x4, 0x3, 0x5, 0x7, 0x6, 0x7, 0x7, 0x7, 0x6, 0x7, 0x7, 0x8, 0x9, 0xb, 0x9, 0x8, 0x8, 0xa, 0x8, 0x7, 0x7, 0xa, 0xd, 0xa, 0xa, 0xb, 0xc, 0xc, 0xc, 0xc, 0x7, 0x9, 0xe, 0xf, 0xd, 0xc, 0xe, 0xb, 0xc, 0xc, 0xc, 0xff, 0xdb, 0x0, 0x43, 0x1, 0x2, 0x2, 0x2, 0x3, 0x3, 0x3, 0x6, 0x3, 0x3, 0x6, 0xc, 0x8, 0x7, 0x8, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xff, 0xc0, 0x0, 0x11, 0x8, 0x0, 0x1, 0x0, 0x1, 0x3, 0x1, 0x22, 0x0, 0x2, 0x11, 0x1, 0x3, 0x11, 0x1, 0xff, 0xc4, 0x0, 0x1f, 0x0, 0x0, 0x1, 0x5, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7, 0x8, 0x9, 0xa, 0xb, 0xff, 0xc4, 0x0, 0xb5, 0x10, 0x0, 0x2, 0x1, 0x3, 0x3, 0x2, 0x4, 0x3, 0x5, 0x5, 0x4, 0x4, 0x0, 0x0, 0x1, 0x7d, 0x1, 0x2, 0x3, 0x0, 0x4, 0x11, 0x5, 0x12, 0x21, 0x31, 0x41, 0x6, 0x13, 0x51, 0x61, 0x7, 0x22, 0x71, 0x14, 0x32, 0x81, 0x91, 0xa1, 0x8, 0x23, 0x42, 0xb1, 0xc1, 0x15, 0x52, 0xd1, 0xf0, 0x24, 0x33, 0x62, 0x72, 0x82, 0x9, 0xa, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2a, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x3a, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48, 0x49, 0x4a, 0x53, 0x54, 0x55, 0x56, 0x57, 0x58, 0x59, 0x5a, 0x63, 0x64, 0x65, 0x66, 0x67, 0x68, 0x69, 0x6a, 0x73, 0x74, 0x75, 0x76, 0x77, 0x78, 0x79, 0x7a, 0x83, 0x84, 0x85, 0x86, 0x87, 0x88, 0x89, 0x8a, 0x92,
		0x93, 0x94, 0x95, 0x96, 0x97, 0x98, 0x99, 0x9a, 0xa2, 0xa3, 0xa4, 0xa5, 0xa6, 0xa7, 0xa8, 0xa9, 0xaa, 0xb2, 0xb3, 0xb4, 0xb5, 0xb6, 0xb7, 0xb8, 0xb9, 0xba, 0xc2, 0xc3, 0xc4, 0xc5, 0xc6, 0xc7, 0xc8, 0xc9, 0xca, 0xd2, 0xd3, 0xd4, 0xd5, 0xd6, 0xd7, 0xd8, 0xd9, 0xda, 0xe1, 0xe2, 0xe3, 0xe4, 0xe5, 0xe6, 0xe7, 0xe8, 0xe9, 0xea, 0xf1, 0xf2, 0xf3, 0xf4, 0xf5, 0xf6, 0xf7, 0xf8, 0xf9, 0xfa, 0xff, 0xc4, 0x0, 0x1f, 0x1, 0x0, 0x3, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7, 0x8, 0x9, 0xa, 0xb, 0xff, 0xc4, 0x0, 0xb5, 0x11, 0x0, 0x2, 0x1, 0x2, 0x4, 0x4, 0x3, 0x4, 0x7, 0x5, 0x4, 0x4, 0x0, 0x1, 0x2, 0x77, 0x0, 0x1, 0x2, 0x3, 0x11, 0x4, 0x5, 0x21, 0x31, 0x6, 0x12, 0x41, 0x51, 0x7, 0x61, 0x71, 0x13, 0x22, 0x32, 0x81, 0x8, 0x14, 0x42, 0x91, 0xa1, 0xb1, 0xc1, 0x9, 0x23, 0x33, 0x52, 0xf0, 0x15, 0x62, 0x72, 0xd1, 0xa, 0x16, 0x24, 0x34, 0xe1, 0x25, 0xf1, 0x17, 0x18, 0x19, 0x1a, 0x26, 0x27, 0x28, 0x29, 0x2a, 0x35, 0x36, 0x37, 0x38, 0x39, 0x3a, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48, 0x49, 0x4a, 0x53, 0x54, 0x55, 0x56, 0x57, 0x58, 0x59, 0x5a, 0x63, 0x64, 0x65, 0x66, 0x67, 0x68, 0x69, 0x6a, 0x73, 0x74, 0x75, 0x76, 0x77, 0x78, 0x79, 0x7a, 0x82, 0x83, 0x84, 0x85, 0x86, 0x87, 0x88, 0x89, 0x8a, 0x92, 0x93, 0x94, 0x95, 0x96, 0x97, 0x98, 0x99, 0x9a, 0xa2, 0xa3, 0xa4, 0xa5, 0xa6, 0xa7, 0xa8, 0xa9, 0xaa, 0xb2, 0xb3, 0xb4, 0xb5, 0xb6, 0xb7, 0xb8, 0xb9, 0xba, 0xc2, 0xc3, 0xc4, 0xc5, 0xc6, 0xc7, 0xc8, 0xc9, 0xca, 0xd2, 0xd3, 0xd4, 0xd5, 0xd6, 0xd7, 0xd8, 0xd9, 0xda, 0xe2, 0xe3, 0xe4, 0xe5, 0xe6, 0xe7, 0xe8, 0xe9, 0xea, 0xf2, 0xf3, 0xf4, 0xf5, 0xf6, 0xf7, 0xf8, 0xf9, 0xfa, 0xff, 0xda, 0x0, 0xc, 0x3, 0x1, 0x0, 0x2, 0x11, 0x3, 0x11, 0x0, 0x3f, 0x0, 0xfc, 0x8b, 0xa2, 0x8a, 0x2b, 0xf3, 0xb3, 0xf6, 0x83, 0xff, 0xd9}
)

func MockCatJson(c *gin.Context, content interface{}, method string) {
	if method == "POST" {
		c.Request.Method = "POST"
	}
	if method == "PUT" {
		c.Request.Method = "PUT"
	}

	c.Request.Header.Set("Content-Type", "application/json")

	jsonbytes, err := json.Marshal(content)
	if err != nil {
		panic(err)
	}
	c.Request.Body = io.NopCloser(bytes.NewBuffer(jsonbytes))
}

func MockCatFile(c *gin.Context, fileType string, file []byte) {
	c.Request.Method = "POST"
	if fileType == "jpeg" {
		c.Request.Header.Set("Content-Type", "image/jpeg")
	} else {
		c.Request.Header.Set("Content-Type", "image/png")
	}
	c.Request.Body = io.NopCloser(bytes.NewBuffer(file))
}

func MockJson(c *gin.Context, content interface{}, method string) {
	if method == "POST" {
		c.Request.Method = "POST"
	}
	if method == "PUT" {
		c.Request.Method = "PUT"
	}

	c.Request.Header.Set("Content-Type", "application/json")

	jsonbytes, err := json.Marshal(content)
	if err != nil {
		panic(err)
	}
	c.Request.Body = io.NopCloser(bytes.NewBuffer(jsonbytes))
}

func TestCreateCategory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	logger := zap.L()
	itemUsecase := mocks.NewMockIItemUsecase(ctrl)
	categoryUsecase := mocks.NewMockICategoryUsecase(ctrl)
	delivery := NewCategoryDelivery(categoryUsecase, itemUsecase, logger.Sugar())

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}

	MockCatJson(c, testShortCategory, post)
	bytesRes, _ := json.Marshal(&testCategoryId)
	categoryUsecase.EXPECT().CreateCategory(ctx, testCategoryNoId).Return(testId, nil)
	delivery.CreateCategory(c)
	require.Equal(t, 201, w.Code)
	require.Equal(t, bytesRes, w.Body.Bytes())

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	c.Request = &http.Request{
		Header: make(http.Header),
	}
	MockCatJson(c, testWrong, post)
	c.Request.Header.Set("Content-Type", "application/text")
	delivery.CreateCategory(c)
	require.Equal(t, 400, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	c.Request = &http.Request{
		Header: make(http.Header),
	}
	MockCatJson(c, testEmptyCategory, post)
	delivery.CreateCategory(c)
	require.Equal(t, 400, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	c.Request = &http.Request{
		Header: make(http.Header),
	}
	MockJson(c, testEmptyCategory, post)
	delivery.CreateCategory(c)
	require.Equal(t, 400, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	c.Request = &http.Request{
		Header: make(http.Header),
	}
	MockJson(c, testShortCategory, post)
	categoryUsecase.EXPECT().CreateCategory(ctx, testCategoryNoId).Return(uuid.Nil, fmt.Errorf("error"))
	delivery.CreateCategory(c)
	require.Equal(t, 500, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	c.Request = &http.Request{
		Header: make(http.Header),
	}
	MockJson(c, testShortNoCategory, post)
	delivery.CreateCategory(c)
	require.Equal(t, 400, w.Code)
}

func TestUpdateCategory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	logger := zap.L()
	itemUsecase := mocks.NewMockIItemUsecase(ctrl)
	categoryUsecase := mocks.NewMockICategoryUsecase(ctrl)
	delivery := NewCategoryDelivery(categoryUsecase, itemUsecase, logger.Sugar())

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Params = []gin.Param{
		{
			Key:   "categoryID",
			Value: testId.String(),
		},
	}
	MockCatJson(c, testCategoryWithId, put)
	categoryUsecase.EXPECT().UpdateCategory(ctx, testModelsCategoryWithId).Return(nil)
	delivery.UpdateCategory(c)
	require.Equal(t, 200, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Params = []gin.Param{
		{
			Key:   "categoryID",
			Value: testId.String(),
		},
	}
	MockCatJson(c, testCatNoCategoryWithId, put)
	delivery.UpdateCategory(c)
	require.Equal(t, 400, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Params = []gin.Param{
		{
			Key:   "categoryID",
			Value: testId.String(),
		},
	}
	MockCatJson(c, testCategoryWithId, put)
	categoryUsecase.EXPECT().UpdateCategory(ctx, testModelsCategoryWithId).Return(fmt.Errorf("error"))
	delivery.UpdateCategory(c)
	require.Equal(t, 500, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Params = []gin.Param{
		{
			Key:   "categoryID",
			Value: testId.String(),
		},
	}
	MockCatJson(c, testCategoryWithId, put)
	categoryUsecase.EXPECT().UpdateCategory(ctx, testModelsCategoryWithId).Return(models.ErrorNotFound{})
	delivery.UpdateCategory(c)
	require.Equal(t, 404, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Params = []gin.Param{
		{
			Key:   "categoryID",
			Value: testId.String() + "l",
		},
	}
	MockCatJson(c, testCategoryWithId, put)
	delivery.UpdateCategory(c)
	require.Equal(t, 400, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Params = []gin.Param{
		{
			Key:   "categoryID",
			Value: testId.String(),
		},
	}
	MockCatJson(c, testWrong, put)
	delivery.UpdateCategory(c)
	require.Equal(t, 400, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	delivery.UpdateCategory(c)
	require.Equal(t, 400, w.Code)
}

func TestGetCategoryList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	logger := zap.L()
	itemUsecase := mocks.NewMockIItemUsecase(ctrl)
	categoryUsecase := mocks.NewMockICategoryUsecase(ctrl)
	delivery := NewCategoryDelivery(categoryUsecase, itemUsecase, logger.Sugar())

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	categoryUsecase.EXPECT().GetCategoryList(ctx).Return([]models.Category{}, fmt.Errorf("error"))
	delivery.GetCategoryList(c)
	require.Equal(t, 500, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}

	testBytes, _ := json.Marshal(&testOutList)
	categoryUsecase.EXPECT().GetCategoryList(ctx).Return(testList, nil)
	delivery.GetCategoryList(c)
	require.Equal(t, 200, w.Code)
	require.Equal(t, testBytes, w.Body.Bytes())
}

func TestGetCategoryList2(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	logger := zap.L()
	itemUsecase := mocks.NewMockIItemUsecase(ctrl)
	categoryUsecase := mocks.NewMockICategoryUsecase(ctrl)
	delivery := NewCategoryDelivery(categoryUsecase, itemUsecase, logger.Sugar())

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	categoryUsecase.EXPECT().GetCategoryList(ctx).Return([]models.Category{{Name: "NoCategory"}}, nil)
	itemUsecase.EXPECT().ItemsQuantityInCategory(ctx, "NoCategory").Return(0, err)
	delivery.GetCategoryList(c)
	require.Equal(t, 200, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	categoryUsecase.EXPECT().GetCategoryList(ctx).Return([]models.Category{{Name: "NoCategory"}}, nil)
	itemUsecase.EXPECT().ItemsQuantityInCategory(ctx, "NoCategory").Return(0, nil)
	delivery.GetCategoryList(c)
	require.Equal(t, 200, w.Code)
}

func TestGetCategory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	logger := zap.L()
	itemUsecase := mocks.NewMockIItemUsecase(ctrl)
	categoryUsecase := mocks.NewMockICategoryUsecase(ctrl)
	delivery := NewCategoryDelivery(categoryUsecase, itemUsecase, logger.Sugar())

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	delivery.GetCategory(c)
	require.Equal(t, 400, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Params = []gin.Param{
		{
			Key:   "categoryID",
			Value: testId.String(),
		},
	}
	categoryUsecase.EXPECT().GetCategory(ctx, testId).Return(&testEmptyModelsCategory, fmt.Errorf("error"))
	delivery.GetCategory(c)
	require.Equal(t, 500, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Params = []gin.Param{
		{
			Key:   "categoryID",
			Value: testId.String(),
		},
	}
	categoryUsecase.EXPECT().GetCategory(ctx, testId).Return(&testEmptyModelsCategory, models.ErrorNotFound{})
	delivery.GetCategory(c)
	require.Equal(t, 404, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Params = []gin.Param{
		{
			Key:   "categoryID",
			Value: testId.String() + "n",
		},
	}
	delivery.GetCategory(c)
	require.Equal(t, 400, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Params = []gin.Param{
		{
			Key:   "categoryID",
			Value: testId.String(),
		},
	}
	testBytes, _ := json.Marshal(&testCategoryWithId)
	categoryUsecase.EXPECT().GetCategory(ctx, testId).Return(testModelsCategoryWithId, nil)
	delivery.GetCategory(c)
	require.Equal(t, 200, w.Code)
	require.Equal(t, testBytes, w.Body.Bytes())
}

func TestUploadCategoryImage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	logger := zap.L()
	itemUsecase := mocks.NewMockIItemUsecase(ctrl)
	categoryUsecase := mocks.NewMockICategoryUsecase(ctrl)
	delivery := NewCategoryDelivery(categoryUsecase, itemUsecase, logger.Sugar())

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	delivery.UploadCategoryImage(c)
	require.Equal(t, 400, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Params = []gin.Param{
		{
			Key:   "categoryID",
			Value: testId.String(),
		},
	}

	delivery.UploadCategoryImage(c)
	require.Equal(t, 415, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Params = []gin.Param{
		{
			Key:   "categoryID",
			Value: testId.String() + "l",
		},
	}

	delivery.UploadCategoryImage(c)
	require.Equal(t, 400, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Params = []gin.Param{
		{
			Key:   "categoryID",
			Value: testId.String(),
		},
	}
	MockCatFile(c, "jpeg", testFile)
	categoryUsecase.EXPECT().UploadCategoryImage(ctx, testId, gomock.Any(), testFile).Return(fmt.Errorf("error"))
	delivery.UploadCategoryImage(c)
	require.Equal(t, 500, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Params = []gin.Param{
		{
			Key:   "categoryID",
			Value: testId.String(),
		},
	}
	MockCatFile(c, "jpeg", testFile)
	categoryUsecase.EXPECT().UploadCategoryImage(ctx, testId, gomock.Any(), testFile).Return(models.ErrorNotFound{})
	delivery.UploadCategoryImage(c)
	require.Equal(t, 404, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Params = []gin.Param{
		{
			Key:   "categoryID",
			Value: testId.String(),
		},
	}
	MockCatFile(c, "jpeg", testFile)
	categoryUsecase.EXPECT().UploadCategoryImage(ctx, testId, gomock.Any(), testFile).Return(nil)
	delivery.UploadCategoryImage(c)
	require.Equal(t, 201, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Params = []gin.Param{
		{
			Key:   "categoryID",
			Value: testId.String(),
		},
	}
	MockCatFile(c, "png", testFile)
	categoryUsecase.EXPECT().UploadCategoryImage(ctx, testId, gomock.Any(), testFile).Return(nil)
	delivery.UploadCategoryImage(c)
	require.Equal(t, 201, w.Code)
}

func TestDeleteCategoryImage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	logger := zap.L()
	itemUsecase := mocks.NewMockIItemUsecase(ctrl)
	categoryUsecase := mocks.NewMockICategoryUsecase(ctrl)
	delivery := NewCategoryDelivery(categoryUsecase, itemUsecase, logger.Sugar())

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	MockCatJson(c, testWrong, post)
	delivery.DeleteCategoryImage(c)
	require.Equal(t, 400, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	delivery.DeleteCategoryImage(c)
	require.Equal(t, 400, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Request.URL, _ = url.Parse(fmt.Sprintf("?id=%s&name=testName", testId.String()))
	categoryUsecase.EXPECT().DeleteCategoryImage(ctx, testId, "testName").Return(fmt.Errorf("error"))
	delivery.DeleteCategoryImage(c)
	require.Equal(t, 500, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Request.URL, _ = url.Parse(fmt.Sprintf("?id=%s&name=testName", testId.String()))
	categoryUsecase.EXPECT().DeleteCategoryImage(ctx, testId, "testName").Return(models.ErrorNotFound{})
	delivery.DeleteCategoryImage(c)
	require.Equal(t, 404, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Request.URL, _ = url.Parse(fmt.Sprintf("?id=%s&name=testName", testId.String()))
	categoryUsecase.EXPECT().DeleteCategoryImage(ctx, testId, "testName").Return(nil)
	delivery.DeleteCategoryImage(c)
	require.Equal(t, 200, w.Code)
}

func TestDeleteCategory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	logger := zap.L()
	itemUsecase := mocks.NewMockIItemUsecase(ctrl)
	categoryUsecase := mocks.NewMockICategoryUsecase(ctrl)
	delivery := NewCategoryDelivery(categoryUsecase, itemUsecase, logger.Sugar())

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	delivery.DeleteCategory(c)
	require.Equal(t, 400, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Params = []gin.Param{
		{
			Key: "categoryID",

			Value: testId.String() + "l",
		},
	}
	delivery.DeleteCategory(c)
	require.Equal(t, 400, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Params = []gin.Param{
		{
			Key:   "categoryID",
			Value: testId.String(),
		},
	}
	categoryUsecase.EXPECT().GetCategory(ctx, testId).Return(&testEmptyModelsCategory, fmt.Errorf("error"))
	delivery.DeleteCategory(c)
	require.Equal(t, 500, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Params = []gin.Param{
		{
			Key:   "categoryID",
			Value: testId.String(),
		},
	}
	categoryUsecase.EXPECT().GetCategory(ctx, testId).Return(&testEmptyModelsCategory, models.ErrorNotFound{})
	delivery.DeleteCategory(c)
	require.Equal(t, 404, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Params = []gin.Param{
		{
			Key:   "categoryID",
			Value: testId.String(),
		},
	}
	categoryUsecase.EXPECT().GetCategory(ctx, testId).Return(testModelsCategoryWithId, nil)
	itemUsecase.EXPECT().ItemsQuantityInCategory(ctx, testModelsCategoryWithId.Name).Return(0, nil)
	categoryUsecase.EXPECT().DeleteCategory(ctx, testId).Return(fmt.Errorf("error"))
	delivery.DeleteCategory(c)
	require.Equal(t, 500, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Params = []gin.Param{
		{
			Key:   "categoryID",
			Value: testId.String(),
		},
	}
	categoryUsecase.EXPECT().GetCategory(ctx, testId).Return(testModelsCategoryWithId, nil)
	itemUsecase.EXPECT().ItemsQuantityInCategory(ctx, testModelsCategoryWithId.Name).Return(-1, fmt.Errorf("error"))
	delivery.DeleteCategory(c)
	require.Equal(t, 500, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Params = []gin.Param{
		{
			Key:   "categoryID",
			Value: testId.String(),
		},
	}
	categoryUsecase.EXPECT().GetCategory(ctx, testId).Return(testCategoryWithImage2, nil)
	itemUsecase.EXPECT().ItemsQuantityInCategory(ctx, testCategoryWithImage2.Name).Return(0, nil)
	categoryUsecase.EXPECT().DeleteCategory(ctx, testId).Return(nil)
	categoryUsecase.EXPECT().DeleteCategoryCache(ctx, testCategoryWithImage2.Name).Return(fmt.Errorf("error"))
	categoryUsecase.EXPECT().DeleteCategoryImageById(ctx, testId).Return(fmt.Errorf("error"))
	delivery.DeleteCategory(c)
	require.Equal(t, 200, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Params = []gin.Param{
		{
			Key:   "categoryID",
			Value: testId.String(),
		},
	}
	limitOptions := map[string]int{"offset": 0, "limit": 1}
	sortOptions := map[string]string{"sortType": "name", "sortOrder": "asc"}

	categoryUsecase.EXPECT().GetCategory(ctx, testId).Return(testCategoryWithImage2, nil)
	itemUsecase.EXPECT().ItemsQuantityInCategory(ctx, testCategoryWithImage2.Name).Return(1, nil)
	itemUsecase.EXPECT().GetItemsByCategory(ctx, testCategoryWithImage2.Name, limitOptions, sortOptions).Return([]models.Item{*testModelsItemWithId}, nil)
	categoryUsecase.EXPECT().DeleteCategory(ctx, testId).Return(nil)
	categoryUsecase.EXPECT().DeleteCategoryCache(ctx, testCategoryWithImage2.Name).Return(nil)
	categoryUsecase.EXPECT().DeleteCategoryImageById(ctx, testId).Return(nil)
	categoryUsecase.EXPECT().GetCategoryByName(ctx, "NoCategory").Return(&testNoCategoryWithId, nil)
	itemUsecase.EXPECT().UpdateItem(ctx, &testModelsItemNoCat).Return(fmt.Errorf("error"))
	itemUsecase.EXPECT().UpdateItemsInCategoryCache(ctx, &testModelsItemNoCat, "create").Return(fmt.Errorf("error"))
	delivery.DeleteCategory(c)
	require.Equal(t, 200, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Params = []gin.Param{
		{
			Key:   "categoryID",
			Value: testId.String(),
		},
	}
	categoryUsecase.EXPECT().GetCategory(ctx, testId).Return(testCategoryWithImage2, nil)
	itemUsecase.EXPECT().ItemsQuantityInCategory(ctx, testCategoryWithImage2.Name).Return(1, nil)
	itemUsecase.EXPECT().GetItemsByCategory(ctx, testCategoryWithImage2.Name, limitOptions, sortOptions).Return(nil, fmt.Errorf("error"))
	delivery.DeleteCategory(c)
	require.Equal(t, 500, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Params = []gin.Param{
		{
			Key:   "categoryID",
			Value: testId.String(),
		},
	}
	categoryUsecase.EXPECT().GetCategory(ctx, testId).Return(testCategoryWithImage2, nil)
	itemUsecase.EXPECT().ItemsQuantityInCategory(ctx, testCategoryWithImage2.Name).Return(1, nil)
	itemUsecase.EXPECT().GetItemsByCategory(ctx, testCategoryWithImage2.Name, limitOptions, sortOptions).Return([]models.Item{*testModelsItemWithId}, nil)
	categoryUsecase.EXPECT().DeleteCategory(ctx, testId).Return(nil)
	categoryUsecase.EXPECT().DeleteCategoryCache(ctx, testCategoryWithImage2.Name).Return(fmt.Errorf("error"))
	categoryUsecase.EXPECT().DeleteCategoryImageById(ctx, testId).Return(nil)
	categoryUsecase.EXPECT().GetCategoryByName(ctx, "NoCategory").Return(&models.Category{}, models.ErrorNotFound{})
	categoryUsecase.EXPECT().CreateCategory(ctx, &testNoCategory).Return(testId, nil)
	itemUsecase.EXPECT().UpdateItem(ctx, &testModelsItemNoCat).Return(fmt.Errorf("error"))
	itemUsecase.EXPECT().UpdateItemsInCategoryCache(ctx, &testModelsItemNoCat, "create").Return(fmt.Errorf("error"))
	delivery.DeleteCategory(c)
	require.Equal(t, 200, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Params = []gin.Param{
		{
			Key:   "categoryID",
			Value: testId.String(),
		},
	}
	categoryUsecase.EXPECT().GetCategory(ctx, testId).Return(testCategoryWithImage2, nil)
	itemUsecase.EXPECT().ItemsQuantityInCategory(ctx, testCategoryWithImage2.Name).Return(1, nil)
	itemUsecase.EXPECT().GetItemsByCategory(ctx, testCategoryWithImage2.Name, limitOptions, sortOptions).Return([]models.Item{*testModelsItemWithId}, nil)
	categoryUsecase.EXPECT().DeleteCategory(ctx, testId).Return(nil)
	categoryUsecase.EXPECT().DeleteCategoryCache(ctx, testCategoryWithImage2.Name).Return(fmt.Errorf("error"))
	categoryUsecase.EXPECT().DeleteCategoryImageById(ctx, testId).Return(nil)
	categoryUsecase.EXPECT().GetCategoryByName(ctx, "NoCategory").Return(&models.Category{}, models.ErrorNotFound{})
	categoryUsecase.EXPECT().CreateCategory(ctx, &testNoCategory).Return(uuid.Nil, fmt.Errorf("error"))
	delivery.DeleteCategory(c)
	require.Equal(t, 500, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Params = []gin.Param{
		{
			Key:   "categoryID",
			Value: testId.String(),
		},
	}
	categoryUsecase.EXPECT().GetCategory(ctx, testId).Return(testCategoryWithImage2, nil)
	itemUsecase.EXPECT().ItemsQuantityInCategory(ctx, testCategoryWithImage2.Name).Return(1, nil)
	itemUsecase.EXPECT().GetItemsByCategory(ctx, testCategoryWithImage2.Name, limitOptions, sortOptions).Return([]models.Item{*testModelsItemWithId}, nil)
	categoryUsecase.EXPECT().DeleteCategory(ctx, testId).Return(nil)
	categoryUsecase.EXPECT().DeleteCategoryCache(ctx, testCategoryWithImage2.Name).Return(fmt.Errorf("error"))
	categoryUsecase.EXPECT().DeleteCategoryImageById(ctx, testId).Return(nil)
	categoryUsecase.EXPECT().GetCategoryByName(ctx, "NoCategory").Return(&models.Category{}, fmt.Errorf("error"))
	delivery.DeleteCategory(c)
	require.Equal(t, 500, w.Code)
}

func TestDeleteCategory2(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	logger := zap.L()
	itemUsecase := mocks.NewMockIItemUsecase(ctrl)
	categoryUsecase := mocks.NewMockICategoryUsecase(ctrl)
	delivery := NewCategoryDelivery(categoryUsecase, itemUsecase, logger.Sugar())

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Params = []gin.Param{
		{
			Key:   "categoryID",
			Value: testId.String(),
		},
	}
	categoryUsecase.EXPECT().GetCategory(ctx, testId).Return(&testNoCategoryWithId, nil)
	delivery.DeleteCategory(c)
	require.Equal(t, 400, w.Code)
}

func TestGetCategoriesImagesList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	logger := zap.L()
	itemUsecase := mocks.NewMockIItemUsecase(ctrl)
	categoryUsecase := mocks.NewMockICategoryUsecase(ctrl)
	delivery := NewCategoryDelivery(categoryUsecase, itemUsecase, logger.Sugar())

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}

	categoryUsecase.EXPECT().GetCategoriesImagesList(ctx).Return(nil, fmt.Errorf("error"))
	delivery.GetCategoriesImagesList(c)
	require.Equal(t, 500, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}

	categoryUsecase.EXPECT().GetCategoriesImagesList(ctx).Return([]*models.FileInfo{}, nil)
	delivery.GetCategoriesImagesList(c)
	require.Equal(t, 200, w.Code)
}