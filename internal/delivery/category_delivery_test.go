package delivery

import (
	"OnlineShopBackend/internal/delivery/category"
	fs "OnlineShopBackend/internal/repository/mocks"
	"OnlineShopBackend/internal/models"
	"OnlineShopBackend/internal/usecase/mocks"
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
	"github.com/golang-module/carbon/v2"
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
	testModelsCategoryWithId2 = &models.Category{
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
	testCategoryWithImage = models.Category{
		Id:          testId,
		Name:        "testName",
		Description: "testDescription",
		Image:       "testImagePath",
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

func TestCreateCategory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	logger := zap.L()
	itemUsecase := mocks.NewMockIItemUsecase(ctrl)
	categoryUsecase := mocks.NewMockICategoryUsecase(ctrl)

	cartUsecase := mocks.NewMockICartUsecase(ctrl)
	filestorage := fs.NewMockFileStorager(ctrl)
	orderUsecase := mocks.NewMockIOrderUsecase(ctrl)
	delivery := NewDelivery(itemUsecase, nil, categoryUsecase, cartUsecase, logger, filestorage, orderUsecase)

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

	cartUsecase := mocks.NewMockICartUsecase(ctrl)
	filestorage := fs.NewMockFileStorager(ctrl)
	orderUsecase := mocks.NewMockIOrderUsecase(ctrl)
	delivery := NewDelivery(itemUsecase, nil, categoryUsecase, cartUsecase, logger, filestorage, orderUsecase)

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

	cartUsecase := mocks.NewMockICartUsecase(ctrl)
	filestorage := fs.NewMockFileStorager(ctrl)
	orderUsecase := mocks.NewMockIOrderUsecase(ctrl)
	delivery := NewDelivery(itemUsecase, nil, categoryUsecase, cartUsecase, logger, filestorage, orderUsecase)

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

	cartUsecase := mocks.NewMockICartUsecase(ctrl)
	filestorage := fs.NewMockFileStorager(ctrl)
	orderUsecase := mocks.NewMockIOrderUsecase(ctrl)
	delivery := NewDelivery(itemUsecase, nil, categoryUsecase, cartUsecase, logger, filestorage, orderUsecase)

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

	cartUsecase := mocks.NewMockICartUsecase(ctrl)
	filestorage := fs.NewMockFileStorager(ctrl)
	orderUsecase := mocks.NewMockIOrderUsecase(ctrl)
	delivery := NewDelivery(itemUsecase, nil, categoryUsecase, cartUsecase, logger, filestorage, orderUsecase)

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

	cartUsecase := mocks.NewMockICartUsecase(ctrl)
	filestorage := fs.NewMockFileStorager(ctrl)
	orderUsecase := mocks.NewMockIOrderUsecase(ctrl)
	delivery := NewDelivery(itemUsecase, nil, categoryUsecase, cartUsecase, logger, filestorage, orderUsecase)

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
	filestorage.EXPECT().PutCategoryImage(testId.String(), carbon.Now().ToShortDateTimeString()+".jpeg", testFile).Return("", fmt.Errorf("error"))
	delivery.UploadCategoryImage(c)
	require.Equal(t, 507, w.Code)

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
	filestorage.EXPECT().PutCategoryImage(testId.String(), carbon.Now().ToShortDateTimeString()+".png", testFile).Return("testImagePath", nil)
	categoryUsecase.EXPECT().GetCategory(ctx, testId).Return(&testEmptyModelsCategory, fmt.Errorf("error"))
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
	MockCatFile(c, "png", testFile)
	filestorage.EXPECT().PutCategoryImage(testId.String(), carbon.Now().ToShortDateTimeString()+".png", testFile).Return("testImagePath", nil)
	categoryUsecase.EXPECT().GetCategory(ctx, testId).Return(&testEmptyModelsCategory, models.ErrorNotFound{})
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
	MockCatFile(c, "png", testFile)
	filestorage.EXPECT().PutCategoryImage(testId.String(), carbon.Now().ToShortDateTimeString()+".png", testFile).Return("testImagePath", nil)
	categoryUsecase.EXPECT().GetCategory(ctx, testId).Return(testModelsCategoryWithId, nil)
	categoryUsecase.EXPECT().UpdateCategory(ctx, &testCategoryWithImage).Return(fmt.Errorf("error"))
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
	MockCatFile(c, "png", testFile)
	filestorage.EXPECT().PutCategoryImage(testId.String(), carbon.Now().ToShortDateTimeString()+".png", testFile).Return("testImagePath", nil)
	categoryUsecase.EXPECT().GetCategory(ctx, testId).Return(testModelsCategoryWithId, nil)
	categoryUsecase.EXPECT().UpdateCategory(ctx, &testCategoryWithImage).Return(nil)
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

	cartUsecase := mocks.NewMockICartUsecase(ctrl)
	filestorage := fs.NewMockFileStorager(ctrl)
	orderUsecase := mocks.NewMockIOrderUsecase(ctrl)
	delivery := NewDelivery(itemUsecase, nil, categoryUsecase, cartUsecase, logger, filestorage, orderUsecase)

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
	filestorage.EXPECT().DeleteCategoryImage(testId.String(), "testName").Return(fmt.Errorf("error"))
	delivery.DeleteCategoryImage(c)
	require.Equal(t, 500, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Request.URL, _ = url.Parse(fmt.Sprintf("?id=%s&name=testName", testId.String()+"l"))
	filestorage.EXPECT().DeleteCategoryImage(testId.String()+"l", "testName").Return(nil)
	delivery.DeleteCategoryImage(c)
	require.Equal(t, 400, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Request.URL, _ = url.Parse(fmt.Sprintf("?id=%s&name=testName", testId.String()))
	filestorage.EXPECT().DeleteCategoryImage(testId.String(), "testName").Return(fmt.Errorf("error"))
	delivery.DeleteCategoryImage(c)
	require.Equal(t, 500, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Request.URL, _ = url.Parse(fmt.Sprintf("?id=%s&name=testName", testId.String()))
	filestorage.EXPECT().DeleteCategoryImage(testId.String(), "testName").Return(nil)
	categoryUsecase.EXPECT().GetCategory(ctx, testId).Return(&testEmptyModelsCategory, fmt.Errorf("error"))
	delivery.DeleteCategoryImage(c)
	require.Equal(t, 500, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Request.URL, _ = url.Parse(fmt.Sprintf("?id=%s&name=testName", testId.String()))
	filestorage.EXPECT().DeleteCategoryImage(testId.String(), "testName").Return(nil)
	categoryUsecase.EXPECT().GetCategory(ctx, testId).Return(&testEmptyModelsCategory, models.ErrorNotFound{})
	delivery.DeleteCategoryImage(c)
	require.Equal(t, 404, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Request.URL, _ = url.Parse(fmt.Sprintf("?id=%s&name=testImagePath", testId.String()))
	filestorage.EXPECT().DeleteCategoryImage(testId.String(), "testImagePath").Return(nil)
	categoryUsecase.EXPECT().GetCategory(ctx, testId).Return(&testCategoryWithImage, nil)
	categoryUsecase.EXPECT().UpdateCategory(ctx, testModelsCategoryWithId2).Return(fmt.Errorf("error"))
	delivery.DeleteCategoryImage(c)
	require.Equal(t, 500, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Request.URL, _ = url.Parse(fmt.Sprintf("?id=%s&name=testImagePath", testId.String()))
	filestorage.EXPECT().DeleteCategoryImage(testId.String(), "testImagePath").Return(nil)
	categoryUsecase.EXPECT().GetCategory(ctx, testId).Return(&testCategoryWithImage, nil)
	categoryUsecase.EXPECT().UpdateCategory(ctx, testModelsCategoryWithId2).Return(nil)
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

	cartUsecase := mocks.NewMockICartUsecase(ctrl)
	filestorage := fs.NewMockFileStorager(ctrl)
	orderUsecase := mocks.NewMockIOrderUsecase(ctrl)
	delivery := NewDelivery(itemUsecase, nil, categoryUsecase, cartUsecase, logger, filestorage, orderUsecase)

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
	categoryUsecase.EXPECT().DeleteCategoryCash(ctx, testCategoryWithImage2.Name).Return(fmt.Errorf("error"))
	filestorage.EXPECT().DeleteCategoryImageById(testId.String()).Return(nil)
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
	categoryUsecase.EXPECT().DeleteCategoryCash(ctx, testCategoryWithImage2.Name).Return(nil)
	filestorage.EXPECT().DeleteCategoryImageById(testId.String()).Return(nil)
	categoryUsecase.EXPECT().GetCategoryByName(ctx, "NoCategory").Return(&testNoCategoryWithId, nil)
	itemUsecase.EXPECT().UpdateItem(ctx, &testModelsItemNoCat).Return(fmt.Errorf("error"))
	itemUsecase.EXPECT().UpdateItemsInCategoryCash(ctx, &testModelsItemNoCat, "create").Return(fmt.Errorf("error"))
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
	categoryUsecase.EXPECT().DeleteCategoryCash(ctx, testCategoryWithImage2.Name).Return(fmt.Errorf("error"))
	filestorage.EXPECT().DeleteCategoryImageById(testId.String()).Return(nil)
	categoryUsecase.EXPECT().GetCategoryByName(ctx, "NoCategory").Return(&models.Category{}, models.ErrorNotFound{})
	categoryUsecase.EXPECT().CreateCategory(ctx, &testNoCategory).Return(testId, nil)
	itemUsecase.EXPECT().UpdateItem(ctx, &testModelsItemNoCat).Return(fmt.Errorf("error"))
	itemUsecase.EXPECT().UpdateItemsInCategoryCash(ctx, &testModelsItemNoCat, "create").Return(fmt.Errorf("error"))
	delivery.DeleteCategory(c)
	require.Equal(t, 200, w.Code)
}

func TestDeleteCategory2(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	logger := zap.L()
	itemUsecase := mocks.NewMockIItemUsecase(ctrl)
	categoryUsecase := mocks.NewMockICategoryUsecase(ctrl)

	cartUsecase := mocks.NewMockICartUsecase(ctrl)
	filestorage := fs.NewMockFileStorager(ctrl)
	orderUsecase := mocks.NewMockIOrderUsecase(ctrl)
	delivery := NewDelivery(itemUsecase, nil, categoryUsecase, cartUsecase, logger, filestorage, orderUsecase)

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
