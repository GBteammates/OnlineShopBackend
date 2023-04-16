package delivery

import (
	"OnlineShopBackend/internal/delivery/category"
	"OnlineShopBackend/internal/delivery/item"
	"OnlineShopBackend/internal/models"
	fs "OnlineShopBackend/internal/repository/mocks"
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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.uber.org/zap"
)

type WrongShortItem struct {
	Title       int
	Description int32
	Category    []int
	Price       string
	Vendor      int
}

type WrongAddFav struct {
	UserId int
	ItemId int32
}

type WrongInItem struct {
	Id          int
	Title       int
	Description int32
	Category    []int
	Price       string
	Vendor      int
}

var (
	testId        = uuid.New()
	testId2       = uuid.New()
	testItemId    = item.ItemId{Value: testId.String()}
	testShortItem = item.ShortItem{
		Title:       "testTitle",
		Description: "testDescription",
		Category:    testId.String(),
		Price:       10,
		Vendor:      "testVendor",
	}
	testShortItemWithoutCat = item.ShortItem{
		Title:       "testTitle",
		Description: "testDescription",
		Price:       10,
		Vendor:      "testVendor",
	}
	wrongShortItem = WrongShortItem{
		Title:       10,
		Description: 11,
		Category:    []int{1},
		Price:       "5",
		Vendor:      5,
	}
	wrongInItem = WrongInItem{
		Id:          5,
		Title:       10,
		Description: 11,
		Category:    []int{1},
		Price:       "5",
		Vendor:      5,
	}
	testModelsItemWithoutId = &models.Item{
		Title:       "testTitle",
		Description: "testDescription",
		Category: models.Category{
			Id: testId,
		},
		Price:  10,
		Vendor: "testVendor",
	}
	testInItem = item.InItem{
		Id:          testId.String(),
		Title:       "testTitle",
		Description: "testDescription",
		Category:    testId.String(),
		Price:       10,
		Vendor:      "testVendor",
	}
	testInItemWithWrongId = item.InItem{
		Id:          testId.String() + "1",
		Title:       "testTitle",
		Description: "testDescription",
		Category:    testId.String(),
		Price:       10,
		Vendor:      "testVendor",
	}
	testInItemWithWrongCatId = item.InItem{
		Id:          testId.String(),
		Title:       "testTitle",
		Description: "testDescription",
		Category:    testId.String() + "1",
		Price:       10,
		Vendor:      "testVendor",
	}
	testOutItem = item.OutItem{
		Id:          testId.String(),
		Title:       "testTitle",
		Description: "testDescription",
		Category: category.Category{
			Id:          testInItem.Category,
			Name:        "testName",
			Description: "testDescription",
		},
		Price:  10,
		Vendor: "testVendor",
	}
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

	testModelsItemWithIdAndOtherCatId = &models.Item{
		Id:          testId,
		Title:       "testTitle",
		Description: "testDescription",
		Category: models.Category{
			Id:          testId2,
			Name:        "testName",
			Description: "testDescription",
		},
		Price:  10,
		Vendor: "testVendor",
	}
	testModelsItemWithId2 = &models.Item{
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
	testShortModelsItemWithIdWithEmptyImage2 = &models.Item{
		Id:          testId,
		Title:       "testTitle",
		Description: "testDescription",
		Category:    models.Category{},
		Price:       10,
		Vendor:      "testVendor",
		Images:      []string{""},
	}
	testModelsItemWithImage = models.Item{
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
		Images: []string{"testName"},
	}
	testModelsItemWithImage2 = models.Item{
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
		Images: []string{"testName.jpeg"},
	}
	testEmptyItem = item.ShortItem{
		Title:       "",
		Description: "",
		Category:    testId.String(),
		Price:       0,
		Vendor:      "",
	}
	testAddFav = item.AddFavItem{
		UserId: testId.String(),
		ItemId: testId2.String(),
	}
	testWrongAddFav = WrongAddFav{
		UserId: 1,
		ItemId: 2,
	}
	post         = "POST"
	put          = "PUT"
	testItems    = []models.Item{*testModelsItemWithId}
	testOutItems = item.ItemsList{
		List: []item.OutItem{testOutItem},
	}
	testFile = []byte{0xff, 0xd8, 0xff, 0xe0, 0x0, 0x10, 0x4a, 0x46, 0x49, 0x46, 0x0, 0x1, 0x1, 0x1, 0x0, 0x48, 0x0, 0x48, 0x0, 0x0, 0xff, 0xe1, 0x0, 0x22, 0x45, 0x78, 0x69, 0x66, 0x0, 0x0, 0x4d, 0x4d, 0x0, 0x2a, 0x0, 0x0, 0x0, 0x8, 0x0, 0x1, 0x1, 0x12, 0x0, 0x3, 0x0, 0x0, 0x0, 0x1, 0x0, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xff, 0xfe, 0x0, 0xd, 0x53, 0x65, 0x63, 0x6c, 0x75, 0x62, 0x2e, 0x6f, 0x72, 0x67, 0x0, 0xff, 0xdb, 0x0, 0x43, 0x0, 0x2, 0x1, 0x1, 0x2, 0x1, 0x1, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x3, 0x5, 0x3, 0x3, 0x3, 0x3, 0x3, 0x6, 0x4, 0x4, 0x3, 0x5, 0x7, 0x6, 0x7, 0x7, 0x7, 0x6, 0x7, 0x7, 0x8, 0x9, 0xb, 0x9, 0x8, 0x8, 0xa, 0x8, 0x7, 0x7, 0xa, 0xd, 0xa, 0xa, 0xb, 0xc, 0xc, 0xc, 0xc, 0x7, 0x9, 0xe, 0xf, 0xd, 0xc, 0xe, 0xb, 0xc, 0xc, 0xc, 0xff, 0xdb, 0x0, 0x43, 0x1, 0x2, 0x2, 0x2, 0x3, 0x3, 0x3, 0x6, 0x3, 0x3, 0x6, 0xc, 0x8, 0x7, 0x8, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xff, 0xc0, 0x0, 0x11, 0x8, 0x0, 0x1, 0x0, 0x1, 0x3, 0x1, 0x22, 0x0, 0x2, 0x11, 0x1, 0x3, 0x11, 0x1, 0xff, 0xc4, 0x0, 0x1f, 0x0, 0x0, 0x1, 0x5, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7, 0x8, 0x9, 0xa, 0xb, 0xff, 0xc4, 0x0, 0xb5, 0x10, 0x0, 0x2, 0x1, 0x3, 0x3, 0x2, 0x4, 0x3, 0x5, 0x5, 0x4, 0x4, 0x0, 0x0, 0x1, 0x7d, 0x1, 0x2, 0x3, 0x0, 0x4, 0x11, 0x5, 0x12, 0x21, 0x31, 0x41, 0x6, 0x13, 0x51, 0x61, 0x7, 0x22, 0x71, 0x14, 0x32, 0x81, 0x91, 0xa1, 0x8, 0x23, 0x42, 0xb1, 0xc1, 0x15, 0x52, 0xd1, 0xf0, 0x24, 0x33, 0x62, 0x72, 0x82, 0x9, 0xa, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2a, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x3a, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48, 0x49, 0x4a, 0x53, 0x54, 0x55, 0x56, 0x57, 0x58, 0x59, 0x5a, 0x63, 0x64, 0x65, 0x66, 0x67, 0x68, 0x69, 0x6a, 0x73, 0x74, 0x75, 0x76, 0x77, 0x78, 0x79, 0x7a, 0x83, 0x84, 0x85, 0x86, 0x87, 0x88, 0x89, 0x8a, 0x92,
		0x93, 0x94, 0x95, 0x96, 0x97, 0x98, 0x99, 0x9a, 0xa2, 0xa3, 0xa4, 0xa5, 0xa6, 0xa7, 0xa8, 0xa9, 0xaa, 0xb2, 0xb3, 0xb4, 0xb5, 0xb6, 0xb7, 0xb8, 0xb9, 0xba, 0xc2, 0xc3, 0xc4, 0xc5, 0xc6, 0xc7, 0xc8, 0xc9, 0xca, 0xd2, 0xd3, 0xd4, 0xd5, 0xd6, 0xd7, 0xd8, 0xd9, 0xda, 0xe1, 0xe2, 0xe3, 0xe4, 0xe5, 0xe6, 0xe7, 0xe8, 0xe9, 0xea, 0xf1, 0xf2, 0xf3, 0xf4, 0xf5, 0xf6, 0xf7, 0xf8, 0xf9, 0xfa, 0xff, 0xc4, 0x0, 0x1f, 0x1, 0x0, 0x3, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7, 0x8, 0x9, 0xa, 0xb, 0xff, 0xc4, 0x0, 0xb5, 0x11, 0x0, 0x2, 0x1, 0x2, 0x4, 0x4, 0x3, 0x4, 0x7, 0x5, 0x4, 0x4, 0x0, 0x1, 0x2, 0x77, 0x0, 0x1, 0x2, 0x3, 0x11, 0x4, 0x5, 0x21, 0x31, 0x6, 0x12, 0x41, 0x51, 0x7, 0x61, 0x71, 0x13, 0x22, 0x32, 0x81, 0x8, 0x14, 0x42, 0x91, 0xa1, 0xb1, 0xc1, 0x9, 0x23, 0x33, 0x52, 0xf0, 0x15, 0x62, 0x72, 0xd1, 0xa, 0x16, 0x24, 0x34, 0xe1, 0x25, 0xf1, 0x17, 0x18, 0x19, 0x1a, 0x26, 0x27, 0x28, 0x29, 0x2a, 0x35, 0x36, 0x37, 0x38, 0x39, 0x3a, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48, 0x49, 0x4a, 0x53, 0x54, 0x55, 0x56, 0x57, 0x58, 0x59, 0x5a, 0x63, 0x64, 0x65, 0x66, 0x67, 0x68, 0x69, 0x6a, 0x73, 0x74, 0x75, 0x76, 0x77, 0x78, 0x79, 0x7a, 0x82, 0x83, 0x84, 0x85, 0x86, 0x87, 0x88, 0x89, 0x8a, 0x92, 0x93, 0x94, 0x95, 0x96, 0x97, 0x98, 0x99, 0x9a, 0xa2, 0xa3, 0xa4, 0xa5, 0xa6, 0xa7, 0xa8, 0xa9, 0xaa, 0xb2, 0xb3, 0xb4, 0xb5, 0xb6, 0xb7, 0xb8, 0xb9, 0xba, 0xc2, 0xc3, 0xc4, 0xc5, 0xc6, 0xc7, 0xc8, 0xc9, 0xca, 0xd2, 0xd3, 0xd4, 0xd5, 0xd6, 0xd7, 0xd8, 0xd9, 0xda, 0xe2, 0xe3, 0xe4, 0xe5, 0xe6, 0xe7, 0xe8, 0xe9, 0xea, 0xf2, 0xf3, 0xf4, 0xf5, 0xf6, 0xf7, 0xf8, 0xf9, 0xfa, 0xff, 0xda, 0x0, 0xc, 0x3, 0x1, 0x0, 0x2, 0x11, 0x3, 0x11, 0x0, 0x3f, 0x0, 0xfc, 0x8b, 0xa2, 0x8a, 0x2b, 0xf3, 0xb3, 0xf6, 0x83, 0xff, 0xd9}
)

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

func MockFile(c *gin.Context, fileType string, file []byte) {
	c.Request.Method = "POST"
	if fileType == "jpeg" {
		c.Request.Header.Set("Content-Type", "image/jpeg")
	} else {
		c.Request.Header.Set("Content-Type", "image/png")
	}
	c.Request.Body = io.NopCloser(bytes.NewBuffer(file))
}
func TestCreateItem(t *testing.T) {
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
	MockJson(c, testShortItem, post)
	bytesRes, _ := json.Marshal(&testItemId)
	itemUsecase.EXPECT().CreateItem(ctx, testModelsItemWithoutId).Return(testId, nil)
	delivery.CreateItem(c)
	require.Equal(t, 201, w.Code)
	require.Equal(t, bytesRes, w.Body.Bytes())

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	c.Request = &http.Request{
		Header: make(http.Header),
	}
	MockJson(c, wrongShortItem, post)
	delivery.CreateItem(c)
	assert.Equal(t, 400, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	c.Request = &http.Request{
		Header: make(http.Header),
	}
	MockJson(c, testEmptyItem, post)
	delivery.CreateItem(c)
	require.Equal(t, 400, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	c.Request = &http.Request{
		Header: make(http.Header),
	}
	MockJson(c, testShortItem, post)
	itemUsecase.EXPECT().CreateItem(ctx, testModelsItemWithoutId).Return(uuid.Nil, fmt.Errorf("error"))
	delivery.CreateItem(c)
	require.Equal(t, 500, w.Code)

	testNoCategory := models.Category{
		Name:        "NoCategory",
		Description: "Category for items without categories",
	}
	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	c.Request = &http.Request{
		Header: make(http.Header),
	}
	MockJson(c, testShortItemWithoutCat, post)
	categoryUsecase.EXPECT().GetCategoryByName(ctx, "NoCategory").Return(nil, models.ErrorNotFound{})
	categoryUsecase.EXPECT().CreateCategory(ctx, &testNoCategory).Return(uuid.Nil, fmt.Errorf("error"))
	delivery.CreateItem(c)
	require.Equal(t, 500, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	c.Request = &http.Request{
		Header: make(http.Header),
	}
	MockJson(c, testShortItemWithoutCat, post)
	categoryUsecase.EXPECT().GetCategoryByName(ctx, "NoCategory").Return(nil, err)
	delivery.CreateItem(c)
	require.Equal(t, 500, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	c.Request = &http.Request{
		Header: make(http.Header),
	}
	MockJson(c, testShortItemWithoutCat, post)
	categoryUsecase.EXPECT().GetCategoryByName(ctx, "NoCategory").Return(nil, models.ErrorNotFound{})
	categoryUsecase.EXPECT().CreateCategory(ctx, &testNoCategory).Return(testId, nil)
	itemUsecase.EXPECT().CreateItem(ctx, testModelsItemWithoutId).Return(testId, nil)
	delivery.CreateItem(c)
	require.Equal(t, 201, w.Code)
}

func TestGetItem(t *testing.T) {
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
			Key:   "itemID",
			Value: testId.String(),
		},
	}
	bytesRes, _ := json.Marshal(&testOutItem)
	itemUsecase.EXPECT().GetItem(ctx, testId).Return(testModelsItemWithId, nil)
	delivery.GetItem(c)
	require.Equal(t, 200, w.Code)
	require.Equal(t, bytesRes, w.Body.Bytes())

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	delivery.GetItem(c)
	require.Equal(t, 400, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Params = []gin.Param{
		{
			Key:   "itemID",
			Value: testId.String(),
		},
	}
	itemUsecase.EXPECT().GetItem(ctx, testId).Return(&models.Item{}, fmt.Errorf("error"))
	delivery.GetItem(c)
	require.Equal(t, 500, w.Code)
}

func TestUpdateItem(t *testing.T) {
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
	MockJson(c, wrongInItem, put)
	delivery.UpdateItem(c)
	require.Equal(t, 400, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	MockJson(c, testInItemWithWrongId, put)
	delivery.UpdateItem(c)
	require.Equal(t, 400, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	MockJson(c, testInItemWithWrongCatId, put)
	delivery.UpdateItem(c)
	require.Equal(t, 400, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	MockJson(c, testInItem, put)
	itemUsecase.EXPECT().GetItem(ctx, testId).Return(nil, fmt.Errorf("error"))
	delivery.UpdateItem(c)
	require.Equal(t, 500, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	MockJson(c, testInItem, put)
	itemUsecase.EXPECT().GetItem(ctx, testId).Return(nil, models.ErrorNotFound{})
	delivery.UpdateItem(c)
	require.Equal(t, 404, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	MockJson(c, testInItem, put)
	itemUsecase.EXPECT().GetItem(ctx, testId).Return(testModelsItemWithIdAndOtherCatId, nil)
	categoryUsecase.EXPECT().GetCategory(ctx, testModelsItemWithId.Category.Id).Return(&models.Category{}, nil)
	itemUsecase.EXPECT().UpdateItem(ctx, testShortModelsItemWithIdWithEmptyImage2).Return(fmt.Errorf("error"))
	delivery.UpdateItem(c)
	require.Equal(t, 500, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	MockJson(c, testInItem, put)
	itemUsecase.EXPECT().GetItem(ctx, testId).Return(nil, models.ErrorNotFound{})
	delivery.UpdateItem(c)
	require.Equal(t, 404, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	MockJson(c, testInItem, put)
	itemUsecase.EXPECT().GetItem(ctx, testId).Return(testModelsItemWithIdAndOtherCatId, nil)
	categoryUsecase.EXPECT().GetCategory(ctx, testModelsItemWithId.Category.Id).Return(&models.Category{}, err)
	delivery.UpdateItem(c)
	require.Equal(t, 500, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	MockJson(c, testInItem, put)
	itemUsecase.EXPECT().GetItem(ctx, testId).Return(testModelsItemWithIdAndOtherCatId, nil)
	categoryUsecase.EXPECT().GetCategory(ctx, testModelsItemWithId.Category.Id).Return(&models.Category{}, models.ErrorNotFound{})
	delivery.UpdateItem(c)
	require.Equal(t, 404, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	MockJson(c, testInItem, put)
	itemUsecase.EXPECT().GetItem(ctx, testId).Return(testModelsItemWithIdAndOtherCatId, nil)
	categoryUsecase.EXPECT().GetCategory(ctx, testModelsItemWithId.Category.Id).Return(&models.Category{}, nil)
	itemUsecase.EXPECT().UpdateItem(ctx, testShortModelsItemWithIdWithEmptyImage2).Return(models.ErrorNotFound{})
	delivery.UpdateItem(c)
	require.Equal(t, 404, w.Code)
}

func TestItemsList(t *testing.T) {
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
	c.Request.URL, _ = url.Parse("?offset=0&limit=1&sortType=name&sortOrder=asc")

	testOutItems.Quantity = 1
	bytesRes, _ := json.Marshal(&testOutItems)
	testLimitOptions := map[string]int{"offset": 0, "limit": 1}
	testSortOptions := map[string]string{"sortType": "name", "sortOrder": "asc"}
	itemUsecase.EXPECT().ItemsList(ctx, testLimitOptions, testSortOptions).Return(testItems, nil)
	itemUsecase.EXPECT().ItemsQuantity(ctx).Return(1, nil)
	delivery.ItemsList(c)
	require.Equal(t, 200, w.Code)
	require.Equal(t, bytesRes, w.Body.Bytes())

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Request.URL, _ = url.Parse("?offset=0&limit=1")

	itemUsecase.EXPECT().ItemsList(ctx, testLimitOptions, testSortOptions).Return([]models.Item{}, fmt.Errorf("error"))
	delivery.ItemsList(c)
	require.Equal(t, 500, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Request.URL, _ = url.Parse("?offset=0&limit=k")

	delivery.ItemsList(c)
	require.Equal(t, 400, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}

	testOutItems.Quantity = 1
	bytesRes, _ = json.Marshal(&testOutItems)
	itemUsecase.EXPECT().ItemsQuantity(ctx).Return(1, nil)
	itemUsecase.EXPECT().ItemsList(ctx, testLimitOptions, testSortOptions).Return(testItems, nil)
	itemUsecase.EXPECT().ItemsQuantity(ctx).Return(1, nil)
	delivery.ItemsList(c)
	require.Equal(t, 200, w.Code)
	require.Equal(t, bytesRes, w.Body.Bytes())

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}

	testOutItems.Quantity = 100
	bytesRes, _ = json.Marshal(&testOutItems)
	itemUsecase.EXPECT().ItemsQuantity(ctx).Return(100, nil)
	itemUsecase.EXPECT().ItemsList(ctx, map[string]int{"offset": 0, "limit": 10}, testSortOptions).Return(testItems, nil)
	itemUsecase.EXPECT().ItemsQuantity(ctx).Return(100, nil)
	delivery.ItemsList(c)
	require.Equal(t, 200, w.Code)
	require.Equal(t, bytesRes, w.Body.Bytes())

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}

	itemUsecase.EXPECT().ItemsQuantity(ctx).Return(100, nil)
	itemUsecase.EXPECT().ItemsList(ctx, map[string]int{"offset": 0, "limit": 10}, testSortOptions).Return(testItems, nil)
	itemUsecase.EXPECT().ItemsQuantity(ctx).Return(-1, err)
	delivery.ItemsList(c)
	require.Equal(t, 500, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}

	itemUsecase.EXPECT().ItemsQuantity(ctx).Return(0, nil)
	delivery.ItemsList(c)
	require.Equal(t, 200, w.Code)
}

func TestItemsQuantity(t *testing.T) {
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

	testQuantity := item.ItemsQuantity{
		Quantity: 1,
	}
	bytesRes, _ := json.Marshal(&testQuantity)
	itemUsecase.EXPECT().ItemsQuantity(ctx).Return(1, nil)
	delivery.ItemsQuantity(c)
	require.Equal(t, 200, w.Code)
	require.Equal(t, bytesRes, w.Body.Bytes())

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	itemUsecase.EXPECT().ItemsQuantity(ctx).Return(-1, fmt.Errorf("error"))
	delivery.ItemsQuantity(c)
	require.Equal(t, 500, w.Code)
}

func TestSearchLine(t *testing.T) {
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
	c.Request.URL, _ = url.Parse("?param=test&offset=0&limit=1")
	testLimitOptions := map[string]int{"offset": 0, "limit": 1}
	testSortOptions := map[string]string{"sortType": "name", "sortOrder": "asc"}

	testOutItems.Quantity = 1
	bytesRes, _ := json.Marshal(&testOutItems)
	itemUsecase.EXPECT().SearchLine(ctx, "test", testLimitOptions, testSortOptions).Return(testItems, nil)
	itemUsecase.EXPECT().ItemsQuantityInSearch(ctx, "test").Return(1, nil)
	delivery.SearchLine(c)
	require.Equal(t, 200, w.Code)
	require.Equal(t, bytesRes, w.Body.Bytes())

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Request.URL, _ = url.Parse("?param=test&offset=0&limit=1")

	itemUsecase.EXPECT().SearchLine(ctx, "test", testLimitOptions, testSortOptions).Return([]models.Item{}, fmt.Errorf("error"))
	delivery.SearchLine(c)
	require.Equal(t, 500, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Request.URL, _ = url.Parse("?param=test&offset=0&limit=k")

	delivery.SearchLine(c)
	require.Equal(t, 400, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Request.URL, _ = url.Parse("?offset=0&limit=1")

	delivery.SearchLine(c)
	require.Equal(t, 400, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Request.URL, _ = url.Parse("?param=test&offset=0&limit=0")

	itemUsecase.EXPECT().SearchLine(ctx, "test", map[string]int{"offset": 0, "limit": 10}, testSortOptions).Return(testItems, nil)
	itemUsecase.EXPECT().ItemsQuantityInSearch(ctx, "test").Return(1, nil)
	delivery.SearchLine(c)
	require.Equal(t, 200, w.Code)
	require.Equal(t, bytesRes, w.Body.Bytes())

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Request.URL, _ = url.Parse("?param=test&offset=0&limit=0")

	itemUsecase.EXPECT().SearchLine(ctx, "test", map[string]int{"offset": 0, "limit": 10}, testSortOptions).Return(testItems, nil)
	itemUsecase.EXPECT().ItemsQuantityInSearch(ctx, "test").Return(-1, err)
	delivery.SearchLine(c)
	require.Equal(t, 500, w.Code)
}

func TestGetItemsByCategory(t *testing.T) {
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
	c.Request.URL, _ = url.Parse("?param=test&offset=0&limit=1")
	testLimitOptions := map[string]int{"offset": 0, "limit": 1}
	testSortOptions := map[string]string{"sortType": "name", "sortOrder": "asc"}

	testOutItems.Quantity = 1
	bytesRes, _ := json.Marshal(&testOutItems)
	itemUsecase.EXPECT().GetItemsByCategory(ctx, "test", testLimitOptions, testSortOptions).Return(testItems, nil)
	itemUsecase.EXPECT().ItemsQuantityInCategory(ctx, "test").Return(1, nil)
	delivery.GetItemsByCategory(c)
	require.Equal(t, 200, w.Code)
	require.Equal(t, bytesRes, w.Body.Bytes())

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Request.URL, _ = url.Parse("?param=test&offset=0&limit=1")

	itemUsecase.EXPECT().GetItemsByCategory(ctx, "test", testLimitOptions, testSortOptions).Return([]models.Item{}, fmt.Errorf("error"))
	delivery.GetItemsByCategory(c)
	require.Equal(t, 500, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Request.URL, _ = url.Parse("?param=test&offset=0&limit=k")

	delivery.GetItemsByCategory(c)
	require.Equal(t, 400, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Request.URL, _ = url.Parse("?offset=0&limit=1")

	delivery.GetItemsByCategory(c)
	require.Equal(t, 400, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Request.URL, _ = url.Parse("?param=test&offset=0&limit=0")

	itemUsecase.EXPECT().GetItemsByCategory(ctx, "test", map[string]int{"offset": 0, "limit": 10}, testSortOptions).Return(testItems, nil)
	itemUsecase.EXPECT().ItemsQuantityInCategory(ctx, "test").Return(1, nil)
	delivery.GetItemsByCategory(c)
	require.Equal(t, 200, w.Code)
	require.Equal(t, bytesRes, w.Body.Bytes())

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Request.URL, _ = url.Parse("?param=test&offset=0&limit=0")

	itemUsecase.EXPECT().GetItemsByCategory(ctx, "test", map[string]int{"offset": 0, "limit": 10}, testSortOptions).Return(testItems, nil)
	itemUsecase.EXPECT().ItemsQuantityInCategory(ctx, "test").Return(-1, err)
	delivery.GetItemsByCategory(c)
	require.Equal(t, 500, w.Code)
}

func TestUploadItemImage(t *testing.T) {
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
	delivery.UploadItemImage(c)
	require.Equal(t, 400, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Params = []gin.Param{
		{
			Key:   "itemID",
			Value: testId.String(),
		},
	}
	delivery.UploadItemImage(c)
	require.Equal(t, 415, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Params = []gin.Param{
		{
			Key:   "itemID",
			Value: testId.String(),
		},
	}
	MockFile(c, "jpeg", testFile)
	itemUsecase.EXPECT().GetItem(ctx, testId).Return(&models.Item{}, fmt.Errorf("error"))
	delivery.UploadItemImage(c)
	require.Equal(t, 500, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Params = []gin.Param{
		{
			Key:   "itemID",
			Value: testId.String(),
		},
	}
	MockFile(c, "jpeg", testFile)
	itemUsecase.EXPECT().GetItem(ctx, testId).Return(&models.Item{}, models.ErrorNotFound{})
	delivery.UploadItemImage(c)
	require.Equal(t, 404, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Params = []gin.Param{
		{
			Key:   "itemID",
			Value: testId.String(),
		},
	}
	MockFile(c, "jpeg", testFile)
	itemUsecase.EXPECT().GetItem(ctx, testId).Return(&models.Item{}, nil)
	filestorage.EXPECT().PutItemImage(testId.String(), carbon.Now().ToShortDateTimeString()+".jpeg", testFile).Return("", fmt.Errorf("error"))
	delivery.UploadItemImage(c)
	require.Equal(t, 507, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Params = []gin.Param{
		{
			Key:   "itemID",
			Value: testId.String(),
		},
	}
	MockFile(c, "png", testFile)
	itemUsecase.EXPECT().GetItem(ctx, testId).Return(&models.Item{}, nil)
	filestorage.EXPECT().PutItemImage(testId.String(), carbon.Now().ToShortDateTimeString()+".png", testFile).Return("", fmt.Errorf("error"))
	delivery.UploadItemImage(c)
	require.Equal(t, 507, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Params = []gin.Param{
		{
			Key:   "itemID",
			Value: testId.String(),
		},
	}
	MockFile(c, "jpeg", testFile)
	itemUsecase.EXPECT().GetItem(ctx, testId).Return(testModelsItemWithId, nil)
	filestorage.EXPECT().PutItemImage(testId.String(), carbon.Now().ToShortDateTimeString()+".jpeg", testFile).Return("testName", nil)
	itemUsecase.EXPECT().UpdateItem(ctx, &testModelsItemWithImage).Return(fmt.Errorf("error"))
	delivery.UploadItemImage(c)
	require.Equal(t, 500, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Params = []gin.Param{
		{
			Key:   "itemID",
			Value: testId.String(),
		},
	}
	MockFile(c, "jpeg", testFile)
	filestorage.EXPECT().PutItemImage(testId.String(), carbon.Now().ToShortDateTimeString()+".jpeg", testFile).Return("testName", nil)
	itemUsecase.EXPECT().GetItem(ctx, testId).Return(testModelsItemWithId2, nil)
	itemUsecase.EXPECT().UpdateItem(ctx, &testModelsItemWithImage).Return(nil)
	delivery.UploadItemImage(c)
	require.Equal(t, 201, w.Code)
}

func TestDeleteItemImage(t *testing.T) {
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
	delivery.DeleteItemImage(c)
	require.Equal(t, 400, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Request.URL, _ = url.Parse(fmt.Sprintf("?id=%s&name=testName.jpg", testId.String()))
	itemUsecase.EXPECT().GetItem(ctx, testId).Return(&models.Item{}, models.ErrorNotFound{})
	delivery.DeleteItemImage(c)
	require.Equal(t, 404, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Request.URL, _ = url.Parse(fmt.Sprintf("?id=%s&name=testName.jpg", testId.String()))
	itemUsecase.EXPECT().GetItem(ctx, testId).Return(&models.Item{}, err)
	delivery.DeleteItemImage(c)
	require.Equal(t, 500, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Request.URL, _ = url.Parse(fmt.Sprintf("?id=%s&name=testName.jpg", testId.String()))
	itemUsecase.EXPECT().GetItem(ctx, testId).Return(&models.Item{}, nil)
	filestorage.EXPECT().DeleteItemImage(testId.String(), "testName.jpg").Return(fmt.Errorf("error"))
	delivery.DeleteItemImage(c)
	require.Equal(t, 500, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Request.URL, _ = url.Parse(fmt.Sprintf("?id=%s&name=testName.jpeg", testId.String()))
	itemUsecase.EXPECT().GetItem(ctx, testId).Return(&testModelsItemWithImage2, nil)
	filestorage.EXPECT().DeleteItemImage(testId.String(), "testName.jpeg").Return(nil)
	testModelsItemWithImage2.Images = []string{}
	itemUsecase.EXPECT().UpdateItem(ctx, &testModelsItemWithImage2).Return(fmt.Errorf("error"))
	delivery.DeleteItemImage(c)
	require.Equal(t, 500, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Request.URL, _ = url.Parse(fmt.Sprintf("?id=%s&name=testName.jpeg", testId.String()))
	itemUsecase.EXPECT().GetItem(ctx, testId).Return(&testModelsItemWithImage2, nil)
	filestorage.EXPECT().DeleteItemImage(testId.String(), "testName.jpeg").Return(nil)
	testModelsItemWithImage2.Images = []string{}
	itemUsecase.EXPECT().UpdateItem(ctx, &testModelsItemWithImage2).Return(nil)
	delivery.DeleteItemImage(c)
	require.Equal(t, 200, w.Code)
}

func TestDeleteItem(t *testing.T) {
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
	delivery.DeleteItem(c)
	require.Equal(t, 400, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Params = []gin.Param{
		{
			Key:   "itemID",
			Value: testId.String() + "l",
		},
	}
	delivery.DeleteItem(c)
	require.Equal(t, 400, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Params = []gin.Param{
		{
			Key:   "itemID",
			Value: testId.String(),
		},
	}
	itemUsecase.EXPECT().GetItem(ctx, testId).Return(nil, fmt.Errorf("error"))
	delivery.DeleteItem(c)
	require.Equal(t, 500, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Params = []gin.Param{
		{
			Key:   "itemID",
			Value: testId.String(),
		},
	}
	itemUsecase.EXPECT().GetItem(ctx, testId).Return(testModelsItemWithId, nil)
	itemUsecase.EXPECT().DeleteItem(ctx, testId).Return(fmt.Errorf("error"))
	delivery.DeleteItem(c)
	require.Equal(t, 500, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Params = []gin.Param{
		{
			Key:   "itemID",
			Value: testId.String(),
		},
	}
	itemUsecase.EXPECT().GetItem(ctx, testId).Return(&testModelsItemWithImage, nil)
	itemUsecase.EXPECT().DeleteItem(ctx, testId).Return(nil)
	itemUsecase.EXPECT().UpdateItemsInCategoryCache(ctx, &testModelsItemWithImage, "delete").Return(fmt.Errorf("error"))
	filestorage.EXPECT().DeleteItemImagesFolderById(testId.String()).Return(fmt.Errorf("error"))
	delivery.DeleteItem(c)
	require.Equal(t, 200, w.Code)
}

func TestItemsQuantityInCategory(t *testing.T) {
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
			Key:   "ctegoryName",
			Value: "test",
		},
	}
	delivery.ItemsQuantityInCategory(c)
	require.Equal(t, 400, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Params = []gin.Param{
		{
			Key:   "categoryName",
			Value: "test",
		},
	}

	itemUsecase.EXPECT().ItemsQuantityInCategory(ctx, "test").Return(-1, err)
	delivery.ItemsQuantityInCategory(c)
	require.Equal(t, 500, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Params = []gin.Param{
		{
			Key:   "categoryName",
			Value: "test",
		},
	}

	itemUsecase.EXPECT().ItemsQuantityInCategory(ctx, "test").Return(1, nil)
	delivery.ItemsQuantityInCategory(c)
	require.Equal(t, 200, w.Code)
}

func TestItemsQuantityInFavourite(t *testing.T) {
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
			Key:   "ctegoryName",
			Value: "test",
		},
	}
	delivery.ItemsQuantityInFavourite(c)
	require.Equal(t, 400, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Params = []gin.Param{
		{
			Key:   "userID",
			Value: testId.String(),
		},
	}

	itemUsecase.EXPECT().ItemsQuantityInFavourite(ctx, testId).Return(-1, err)
	delivery.ItemsQuantityInFavourite(c)
	require.Equal(t, 500, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Params = []gin.Param{
		{
			Key:   "userID",
			Value: testId.String(),
		},
	}

	itemUsecase.EXPECT().ItemsQuantityInFavourite(ctx, testId).Return(1, nil)
	delivery.ItemsQuantityInFavourite(c)
	require.Equal(t, 200, w.Code)
}

func TestAddFavouriteItem(t *testing.T) {
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
	MockJson(c, testWrongAddFav, post)
	delivery.AddFavouriteItem(c)
	require.Equal(t, 400, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	MockJson(c, testAddFav, post)
	itemUsecase.EXPECT().AddFavouriteItem(ctx, testId, testId2).Return(models.ErrorNotFound{})
	delivery.AddFavouriteItem(c)
	require.Equal(t, 404, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	MockJson(c, testAddFav, post)
	itemUsecase.EXPECT().AddFavouriteItem(ctx, testId, testId2).Return(err)
	delivery.AddFavouriteItem(c)
	require.Equal(t, 500, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	MockJson(c, testAddFav, post)
	itemUsecase.EXPECT().AddFavouriteItem(ctx, testId, testId2).Return(nil)
	delivery.AddFavouriteItem(c)
	require.Equal(t, 200, w.Code)
}

func TestDeleteFavouriteItem(t *testing.T) {
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
			Key:   "userID",
			Value: "test",
		},
	}
	delivery.DeleteFavouriteItem(c)
	require.Equal(t, 400, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Params = []gin.Param{
		{
			Key:   "userID",
			Value: testId.String(),
		},
		{
			Key:   "itemID",
			Value: "test",
		},
	}
	delivery.DeleteFavouriteItem(c)
	require.Equal(t, 400, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Params = []gin.Param{
		{
			Key:   "userID",
			Value: testId.String(),
		},
		{
			Key:   "itemID",
			Value: testId2.String(),
		},
	}
	itemUsecase.EXPECT().DeleteFavouriteItem(ctx, testId, testId2).Return(err)
	delivery.DeleteFavouriteItem(c)
	require.Equal(t, 500, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Params = []gin.Param{
		{
			Key:   "userID",
			Value: testId.String(),
		},
		{
			Key:   "itemID",
			Value: testId2.String(),
		},
	}
	itemUsecase.EXPECT().DeleteFavouriteItem(ctx, testId, testId2).Return(models.ErrorNotFound{})
	delivery.DeleteFavouriteItem(c)
	require.Equal(t, 404, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Params = []gin.Param{
		{
			Key:   "userID",
			Value: testId.String(),
		},
		{
			Key:   "itemID",
			Value: testId2.String(),
		},
	}
	itemUsecase.EXPECT().DeleteFavouriteItem(ctx, testId, testId2).Return(nil)
	delivery.DeleteFavouriteItem(c)
	require.Equal(t, 200, w.Code)
}

func TestGetFavouriteItems(t *testing.T) {
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
	c.Request.URL, _ = url.Parse(fmt.Sprintf("?param=%s&offset=0&limit=1", testId.String()))
	testLimitOptions := map[string]int{"offset": 0, "limit": 1}
	testSortOptions := map[string]string{"sortType": "name", "sortOrder": "asc"}

	itemUsecase.EXPECT().GetFavouriteItems(ctx, testId, testLimitOptions, testSortOptions).Return(testItems, nil)
	itemUsecase.EXPECT().ItemsQuantityInFavourite(ctx, testId).Return(1, nil)
	delivery.GetFavouriteItems(c)
	require.Equal(t, 200, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Request.URL, _ = url.Parse(fmt.Sprintf("?param=%s&offset=0&limit=1", testId.String()))

	itemUsecase.EXPECT().GetFavouriteItems(ctx, testId, testLimitOptions, testSortOptions).Return([]models.Item{}, fmt.Errorf("error"))
	delivery.GetFavouriteItems(c)
	require.Equal(t, 500, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Request.URL, _ = url.Parse("?param=test&offset=0&limit=k")

	delivery.GetFavouriteItems(c)
	require.Equal(t, 400, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Request.URL, _ = url.Parse("?offset=0&limit=1")

	delivery.GetFavouriteItems(c)
	require.Equal(t, 400, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Request.URL, _ = url.Parse(fmt.Sprintf("?param=%s&offset=0&limit=0", testId.String()))

	itemUsecase.EXPECT().GetFavouriteItems(ctx, testId, map[string]int{"offset": 0, "limit": 10}, testSortOptions).Return(testItems, nil)
	itemUsecase.EXPECT().ItemsQuantityInFavourite(ctx, testId).Return(-1, err)
	delivery.GetFavouriteItems(c)
	require.Equal(t, 500, w.Code)
}

func TestItemsQuantityInSearch(t *testing.T) {
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
			Key:   "searchRequet",
			Value: "test",
		},
	}
	delivery.ItemsQuantityInSearch(c)
	require.Equal(t, 400, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Params = []gin.Param{
		{
			Key:   "searchRequest",
			Value: "test",
		},
	}

	itemUsecase.EXPECT().ItemsQuantityInSearch(ctx, "test").Return(-1, err)
	delivery.ItemsQuantityInSearch(c)
	require.Equal(t, 500, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Params = []gin.Param{
		{
			Key:   "searchRequest",
			Value: "test",
		},
	}

	itemUsecase.EXPECT().ItemsQuantityInSearch(ctx, "test").Return(1, nil)
	delivery.ItemsQuantityInSearch(c)
	require.Equal(t, 200, w.Code)
}
