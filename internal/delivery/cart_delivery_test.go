package delivery

import (
	"OnlineShopBackend/internal/delivery/cart"
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
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

type WrongShortCart struct {
	CartId int
	ItemId int
}

var (
	err           = fmt.Errorf("error")
	testUserId    = uuid.New()
	testCartId    = uuid.New()
	testModelCart = models.Cart{
		Id:     testCartId,
		UserId: testUserId,
		Items: []models.ItemWithQuantity{
			testModelItemWithQuantity,
		},
	}
	testModelItemWithQuantity = models.ItemWithQuantity{
		Quantity: 1,
	}
	testWrongShortCart = WrongShortCart{
		CartId: 5,
		ItemId: 5,
	}
	testShortCart = cart.ShortCart{
		CartId: testCartId.String(),
		ItemId: testId.String(),
	}
)

func MockCartJson(c *gin.Context, content interface{}, method string) {
	if method == "DELETE" {
		c.Request.Method = "DELETE"
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

func TestGetCart(t *testing.T) {
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
			Value: testId.String() + "n",
		},
	}
	delivery.GetCart(c)
	require.Equal(t, 400, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Params = []gin.Param{
		{
			Key:   "cartID",
			Value: testCartId.String(),
		},
	}
	cartUsecase.EXPECT().GetCart(ctx, testCartId).Return(nil, err)
	delivery.GetCart(c)
	require.Equal(t, 500, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Params = []gin.Param{
		{
			Key:   "cartID",
			Value: testCartId.String(),
		},
	}
	cartUsecase.EXPECT().GetCart(ctx, testCartId).Return(nil, models.ErrorNotFound{})
	delivery.GetCart(c)
	require.Equal(t, 404, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Params = []gin.Param{
		{
			Key:   "cartID",
			Value: testCartId.String(),
		},
	}
	cartUsecase.EXPECT().GetCart(ctx, testCartId).Return(&testModelCart, nil)
	delivery.GetCart(c)
	require.Equal(t, 200, w.Code)
}

func TestGetCartByUserId(t *testing.T) {
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
			Value: testId.String() + "n",
		},
	}
	delivery.GetCartByUserId(c)
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
	cartUsecase.EXPECT().GetCartByUserId(ctx, testId).Return(nil, err)
	delivery.GetCartByUserId(c)
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
	cartUsecase.EXPECT().GetCartByUserId(ctx, testId).Return(nil, models.ErrorNotFound{})
	delivery.GetCartByUserId(c)
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
	}
	cartUsecase.EXPECT().GetCartByUserId(ctx, testId).Return(&testModelCart, nil)
	delivery.GetCartByUserId(c)
	require.Equal(t, 200, w.Code)
}
func TestCreateCart(t *testing.T) {
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

	delivery.CreateCart(c)
	require.Equal(t, 400, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Params = []gin.Param{
		{
			Key:   "userID",
			Value: testUserId.String(),
		},
	}

	cartUsecase.EXPECT().Create(ctx, testUserId).Return(uuid.Nil, err)
	delivery.CreateCart(c)
	require.Equal(t, 500, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Params = []gin.Param{
		{
			Key:   "userID",
			Value: testUserId.String(),
		},
	}

	cartUsecase.EXPECT().Create(ctx, testUserId).Return(testCartId, nil)
	delivery.CreateCart(c)
	require.Equal(t, 201, w.Code)
}

func TestAddItemToCart(t *testing.T) {
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
	MockCartJson(c, testWrongShortCart, "PUT")
	delivery.AddItemToCart(c)
	require.Equal(t, 400, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	MockCartJson(c, testShortCart, "PUT")
	cartUsecase.EXPECT().AddItemToCart(ctx, testCartId, testId).Return(err)
	delivery.AddItemToCart(c)
	require.Equal(t, 500, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	MockCartJson(c, testShortCart, "PUT")
	cartUsecase.EXPECT().AddItemToCart(ctx, testCartId, testId).Return(nil)
	delivery.AddItemToCart(c)
	require.Equal(t, 200, w.Code)
}

func TestDeleteCart(t *testing.T) {
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
			Key:   "cartID",
			Value: testId.String() + "n",
		},
	}

	delivery.DeleteCart(c)
	require.Equal(t, 400, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Params = []gin.Param{
		{
			Key:   "cartID",
			Value: testCartId.String(),
		},
	}

	cartUsecase.EXPECT().DeleteCart(ctx, testCartId).Return(err)
	delivery.DeleteCart(c)
	require.Equal(t, 500, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Params = []gin.Param{
		{
			Key:   "cartID",
			Value: testCartId.String(),
		},
	}

	cartUsecase.EXPECT().DeleteCart(ctx, testCartId).Return(nil)
	delivery.DeleteCart(c)
	require.Equal(t, 200, w.Code)
}

func TestDeleteItemFromCart(t *testing.T) {
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
	delivery.DeleteItemFromCart(c)
	require.Equal(t, 400, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Params = []gin.Param{
		{
			Key:   "cartID",
			Value: testUserId.String(),
		},
		{
			Key:   "itemID",
			Value: testId.String() + "l",
		},
	}
	delivery.DeleteItemFromCart(c)
	require.Equal(t, 400, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Params = []gin.Param{
		{
			Key:   "cartID",
			Value: testUserId.String(),
		},
		{
			Key:   "itemID",
			Value: testId.String(),
		},
	}
	cartUsecase.EXPECT().DeleteItemFromCart(ctx, testUserId, testId).Return(err)
	delivery.DeleteItemFromCart(c)
	require.Equal(t, 500, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Params = []gin.Param{
		{
			Key:   "cartID",
			Value: testUserId.String(),
		},
		{
			Key:   "itemID",
			Value: testId.String(),
		},
	}

	cartUsecase.EXPECT().DeleteItemFromCart(ctx, testUserId, testId).Return(nil)
	delivery.DeleteItemFromCart(c)
	require.Equal(t, 200, w.Code)
}
