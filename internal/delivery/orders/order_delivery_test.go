package delivery

/*import (
	"OnlineShopBackend/internal/delivery/cart"
	"OnlineShopBackend/internal/delivery/order"
	fs "OnlineShopBackend/internal/filestorage/mocks"
	"OnlineShopBackend/internal/models"
	"OnlineShopBackend/internal/usecase/mocks"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

var (
	testAddress = order.OrderAddress{
		Country: "Israel",
		City:    "Haifa",
		Zipcode: "40006",
		Street:  "Daniel 4",
	}

	testUser = order.UserForCart{
		Id:    uuid.New().String(),
		Email: "test@mail.ru",
		Role:  "user",
	}

	testOrderCart = cart.Cart{
		Id:     uuid.New().String(),
		UserId: testUser.Id,
		Items: []cart.CartItem{
			{
				Id:       uuid.NewString(),
				Title:    "test",
				Price:    300,
				Image:    "testurl",
				Quantity: 1,
			},
		},
	}

	testOrderCartWrongID = cart.Cart{
		Id:     "wrong",
		UserId: testUser.Id,
		Items: []cart.CartItem{
			{
				Id:       uuid.NewString(),
				Title:    "test",
				Price:    300,
				Image:    "testurl",
				Quantity: 1,
			},
		},
	}

	testCartUserAddress = order.CartAdressUser{
		Cart:    testOrderCart,
		Address: testAddress,
		User:    testUser,
	}

	testAddressWithUserAndId = order.AddressWithUserAndId{
		User:    testUser,
		Address: testAddress,
		OrderId: uuid.NewString(),
	}

	testStatusWithUSerAndId = order.StatusWithUserAndId{
		User:    testUser,
		Status:  string(models.StatusCourier),
		OrderId: uuid.NewString(),
	}
)

func MockCartUserAddressJson(c *gin.Context, content interface{}, method string) {
	if method == "DELETE" {
		c.Request.Method = "DELETE"
	}
	if method == "PUT" {
		c.Request.Method = "PUT"
	}
	if method == "POST" {
		c.Request.Method = "POST"
	}

	c.Request.Header.Set("Content-Type", "application/json")

	jsonbytes, err := json.Marshal(content)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(jsonbytes))
	c.Request.Body = io.NopCloser(bytes.NewBuffer(jsonbytes))
}

func TestCreateOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	// ctx := context.Background()
	// logger := zap.L()
	logger := zap.NewExample()
	itemUsecase := mocks.NewMockIItemUsecase(ctrl)
	categoryUsecase := mocks.NewMockICategoryUsecase(ctrl)
	cartUsecase := mocks.NewMockICartUsecase(ctrl)
	filestorage := fs.NewMockFileStorager(ctrl)
	orderUsecase := &mocks.OrderUsecaseMock{}
	delivery := NewDelivery(itemUsecase, categoryUsecase, cartUsecase, logger, filestorage, orderUsecase)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	MockCartUserAddressJson(c, testCartUserAddress, "POST")
	delivery.CreateOrder(c)
	require.Equal(t, 201, w.Code)
}

func TestCreateOrderInternalError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := zap.NewExample()
	itemUsecase := mocks.NewMockIItemUsecase(ctrl)
	categoryUsecase := mocks.NewMockICategoryUsecase(ctrl)
	cartUsecase := mocks.NewMockICartUsecase(ctrl)
	filestorage := fs.NewMockFileStorager(ctrl)
	orderUsecase := &mocks.OrderUsecaseMock{Err: fmt.Errorf("test")}
	delivery := NewDelivery(itemUsecase, categoryUsecase, cartUsecase, logger, filestorage, orderUsecase)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	MockCartUserAddressJson(c, testCartUserAddress, "POST")
	delivery.CreateOrder(c)
	require.Equal(t, 500, w.Code)
}

func TestCreateOrderBadRequestError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := zap.NewExample()
	itemUsecase := mocks.NewMockIItemUsecase(ctrl)
	categoryUsecase := mocks.NewMockICategoryUsecase(ctrl)
	cartUsecase := mocks.NewMockICartUsecase(ctrl)
	filestorage := fs.NewMockFileStorager(ctrl)
	orderUsecase := &mocks.OrderUsecaseMock{}
	delivery := NewDelivery(itemUsecase, categoryUsecase, cartUsecase, logger, filestorage, orderUsecase)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	testCartUserAddress.Cart = testOrderCartWrongID
	MockCartUserAddressJson(c, testCartUserAddress, "POST")
	delivery.CreateOrder(c)
	require.Equal(t, 400, w.Code)
}

func TestDeleteOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := zap.NewExample()
	itemUsecase := mocks.NewMockIItemUsecase(ctrl)
	categoryUsecase := mocks.NewMockICategoryUsecase(ctrl)
	cartUsecase := mocks.NewMockICartUsecase(ctrl)
	filestorage := fs.NewMockFileStorager(ctrl)
	orderUsecase := &mocks.OrderUsecaseMock{}
	delivery := NewDelivery(itemUsecase, categoryUsecase, cartUsecase, logger, filestorage, orderUsecase)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Params = []gin.Param{
		{
			Key:   "orderID",
			Value: uuid.NewString(),
		},
	}
	delivery.DeleteOrder(c)
	require.Equal(t, 200, w.Code)
}

func TestDeleteOrderBadRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := zap.NewExample()
	itemUsecase := mocks.NewMockIItemUsecase(ctrl)
	categoryUsecase := mocks.NewMockICategoryUsecase(ctrl)
	cartUsecase := mocks.NewMockICartUsecase(ctrl)
	filestorage := fs.NewMockFileStorager(ctrl)
	orderUsecase := &mocks.OrderUsecaseMock{}
	delivery := NewDelivery(itemUsecase, categoryUsecase, cartUsecase, logger, filestorage, orderUsecase)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Params = []gin.Param{
		{
			Key:   "orderID",
			Value: "wrong",
		},
	}
	delivery.DeleteOrder(c)
	require.Equal(t, 400, w.Code)
}

func TestDeleteOrderInternalError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := zap.NewExample()
	itemUsecase := mocks.NewMockIItemUsecase(ctrl)
	categoryUsecase := mocks.NewMockICategoryUsecase(ctrl)
	cartUsecase := mocks.NewMockICartUsecase(ctrl)
	filestorage := fs.NewMockFileStorager(ctrl)
	orderUsecase := &mocks.OrderUsecaseMock{Err: fmt.Errorf("internal error")}
	delivery := NewDelivery(itemUsecase, categoryUsecase, cartUsecase, logger, filestorage, orderUsecase)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Params = []gin.Param{
		{
			Key:   "orderID",
			Value: uuid.NewString(),
		},
	}
	delivery.DeleteOrder(c)
	require.Equal(t, 500, w.Code)
}

func TestGetOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := zap.NewExample()
	itemUsecase := mocks.NewMockIItemUsecase(ctrl)
	categoryUsecase := mocks.NewMockICategoryUsecase(ctrl)
	cartUsecase := mocks.NewMockICartUsecase(ctrl)
	filestorage := fs.NewMockFileStorager(ctrl)
	orderUsecase := &mocks.OrderUsecaseMock{}
	delivery := NewDelivery(itemUsecase, categoryUsecase, cartUsecase, logger, filestorage, orderUsecase)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Params = []gin.Param{
		{
			Key:   "orderID",
			Value: uuid.NewString(),
		},
	}
	delivery.GetOrder(c)
	require.Equal(t, 200, w.Code)
}

func TestGetOrderBadRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := zap.NewExample()
	itemUsecase := mocks.NewMockIItemUsecase(ctrl)
	categoryUsecase := mocks.NewMockICategoryUsecase(ctrl)
	cartUsecase := mocks.NewMockICartUsecase(ctrl)
	filestorage := fs.NewMockFileStorager(ctrl)
	orderUsecase := &mocks.OrderUsecaseMock{}
	delivery := NewDelivery(itemUsecase, categoryUsecase, cartUsecase, logger, filestorage, orderUsecase)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Params = []gin.Param{
		{
			Key:   "orderID",
			Value: "wrong",
		},
	}
	delivery.GetOrder(c)
	require.Equal(t, 400, w.Code)
}

func TestGetOrderInternal(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := zap.NewExample()
	itemUsecase := mocks.NewMockIItemUsecase(ctrl)
	categoryUsecase := mocks.NewMockICategoryUsecase(ctrl)
	cartUsecase := mocks.NewMockICartUsecase(ctrl)
	filestorage := fs.NewMockFileStorager(ctrl)
	orderUsecase := &mocks.OrderUsecaseMock{Err: fmt.Errorf("internal error")}
	delivery := NewDelivery(itemUsecase, categoryUsecase, cartUsecase, logger, filestorage, orderUsecase)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Params = []gin.Param{
		{
			Key:   "orderID",
			Value: uuid.NewString(),
		},
	}
	delivery.GetOrder(c)
	require.Equal(t, 500, w.Code)
}

func TestGetOrderForUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := zap.NewExample()
	itemUsecase := mocks.NewMockIItemUsecase(ctrl)
	categoryUsecase := mocks.NewMockICategoryUsecase(ctrl)
	cartUsecase := mocks.NewMockICartUsecase(ctrl)
	filestorage := fs.NewMockFileStorager(ctrl)
	orderUsecase := &mocks.OrderUsecaseMock{}
	delivery := NewDelivery(itemUsecase, categoryUsecase, cartUsecase, logger, filestorage, orderUsecase)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Params = []gin.Param{
		{
			Key:   "userID",
			Value: uuid.NewString(),
		},
	}
	delivery.GetOrdersForUser(c)
	require.Equal(t, 200, w.Code)
	var res []order.Order
	err = json.NewDecoder(w.Body).Decode(&res)
	assert.NoError(t, err)
	require.Equal(t, 2, len(res))
}

func TestGetOrderForUserBadRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := zap.NewExample()
	itemUsecase := mocks.NewMockIItemUsecase(ctrl)
	categoryUsecase := mocks.NewMockICategoryUsecase(ctrl)
	cartUsecase := mocks.NewMockICartUsecase(ctrl)
	filestorage := fs.NewMockFileStorager(ctrl)
	orderUsecase := &mocks.OrderUsecaseMock{}
	delivery := NewDelivery(itemUsecase, categoryUsecase, cartUsecase, logger, filestorage, orderUsecase)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Params = []gin.Param{
		{
			Key:   "userWrongID",
			Value: uuid.NewString(),
		},
	}
	delivery.GetOrdersForUser(c)
	require.Equal(t, 400, w.Code)

}

func TestGetOrderForUserInternal(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := zap.NewExample()
	itemUsecase := mocks.NewMockIItemUsecase(ctrl)
	categoryUsecase := mocks.NewMockICategoryUsecase(ctrl)
	cartUsecase := mocks.NewMockICartUsecase(ctrl)
	filestorage := fs.NewMockFileStorager(ctrl)
	orderUsecase := &mocks.OrderUsecaseMock{Err: fmt.Errorf("inetrnal error")}
	delivery := NewDelivery(itemUsecase, categoryUsecase, cartUsecase, logger, filestorage, orderUsecase)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Params = []gin.Param{
		{
			Key:   "userID",
			Value: uuid.NewString(),
		},
	}
	delivery.GetOrdersForUser(c)
	require.Equal(t, 500, w.Code)

}

func TestChangeAddress(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := zap.NewExample()
	itemUsecase := mocks.NewMockIItemUsecase(ctrl)
	categoryUsecase := mocks.NewMockICategoryUsecase(ctrl)
	cartUsecase := mocks.NewMockICartUsecase(ctrl)
	filestorage := fs.NewMockFileStorager(ctrl)
	orderUsecase := &mocks.OrderUsecaseMock{}
	delivery := NewDelivery(itemUsecase, categoryUsecase, cartUsecase, logger, filestorage, orderUsecase)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	testAddressWithUserAndId.User.Role = "admin"
	MockCartUserAddressJson(c, testAddressWithUserAndId, "PATCH")
	delivery.ChangeAddress(c)
	require.Equal(t, 200, w.Code)
}

func TestChangeAddressForbidden(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := zap.NewExample()
	itemUsecase := mocks.NewMockIItemUsecase(ctrl)
	categoryUsecase := mocks.NewMockICategoryUsecase(ctrl)
	cartUsecase := mocks.NewMockICartUsecase(ctrl)
	filestorage := fs.NewMockFileStorager(ctrl)
	orderUsecase := &mocks.OrderUsecaseMock{}
	delivery := NewDelivery(itemUsecase, categoryUsecase, cartUsecase, logger, filestorage, orderUsecase)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	testAddressWithUserAndId.User.Role = "user"
	MockCartUserAddressJson(c, testAddressWithUserAndId, "PATCH")
	delivery.ChangeAddress(c)
	require.Equal(t, 403, w.Code)
}

func TestChangeAddressBadRequst(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := zap.NewExample()
	itemUsecase := mocks.NewMockIItemUsecase(ctrl)
	categoryUsecase := mocks.NewMockICategoryUsecase(ctrl)
	cartUsecase := mocks.NewMockICartUsecase(ctrl)
	filestorage := fs.NewMockFileStorager(ctrl)
	orderUsecase := &mocks.OrderUsecaseMock{}
	delivery := NewDelivery(itemUsecase, categoryUsecase, cartUsecase, logger, filestorage, orderUsecase)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	testAddressWithUserAndId.User.Role = "admin"
	testAddressWithUserAndId.User.Id = "wrong"
	MockCartUserAddressJson(c, testAddressWithUserAndId, "PATCH")
	delivery.ChangeAddress(c)
	require.Equal(t, 400, w.Code)
}

func TestChangeAddressInternal(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := zap.NewExample()
	itemUsecase := mocks.NewMockIItemUsecase(ctrl)
	categoryUsecase := mocks.NewMockICategoryUsecase(ctrl)
	cartUsecase := mocks.NewMockICartUsecase(ctrl)
	filestorage := fs.NewMockFileStorager(ctrl)
	orderUsecase := &mocks.OrderUsecaseMock{Err: fmt.Errorf("Internal Error")}
	delivery := NewDelivery(itemUsecase, categoryUsecase, cartUsecase, logger, filestorage, orderUsecase)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	testAddressWithUserAndId.User.Role = "admin"
	MockCartUserAddressJson(c, testAddressWithUserAndId, "PATCH")
	delivery.ChangeAddress(c)
	require.Equal(t, 500, w.Code)
}

/*func TestChangeStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := zap.NewExample()
	itemUsecase := mocks.NewMockIItemUsecase(ctrl)
	categoryUsecase := mocks.NewMockICategoryUsecase(ctrl)
	cartUsecase := mocks.NewMockICartUsecase(ctrl)
	filestorage := fs.NewMockFileStorager(ctrl)
	orderUsecase := &mocks.OrderUsecaseMock{}
	delivery := NewDelivery(itemUsecase, categoryUsecase, cartUsecase, logger, filestorage, orderUsecase)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	testStatusWithUSerAndId.User.Role = "admin"
	MockCartUserAddressJson(c, testStatusWithUSerAndId, "PATCH")
	delivery.ChangeStatus(c)
	require.Equal(t, 200, w.Code)
}*/

/*func TestChangeStatusForbidden(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := zap.NewExample()
	itemUsecase := mocks.NewMockIItemUsecase(ctrl)
	categoryUsecase := mocks.NewMockICategoryUsecase(ctrl)
	cartUsecase := mocks.NewMockICartUsecase(ctrl)
	filestorage := fs.NewMockFileStorager(ctrl)
	orderUsecase := &mocks.OrderUsecaseMock{}
	delivery := NewDelivery(itemUsecase, categoryUsecase, cartUsecase, logger, filestorage, orderUsecase)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	testStatusWithUSerAndId.User.Role = "user"
	MockCartUserAddressJson(c, testStatusWithUSerAndId, "PATCH")
	delivery.ChangeStatus(c)
	require.Equal(t, 403, w.Code)
}

func TestChangeStatusBadRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := zap.NewExample()
	itemUsecase := mocks.NewMockIItemUsecase(ctrl)
	categoryUsecase := mocks.NewMockICategoryUsecase(ctrl)
	cartUsecase := mocks.NewMockICartUsecase(ctrl)
	filestorage := fs.NewMockFileStorager(ctrl)
	orderUsecase := &mocks.OrderUsecaseMock{}
	delivery := NewDelivery(itemUsecase, categoryUsecase, cartUsecase, logger, filestorage, orderUsecase)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}

	testStatusWithUSerAndId.User.Role = "admin"
	testStatusWithUSerAndId.OrderId = "wrong"
	defer func() {
		testStatusWithUSerAndId.OrderId = uuid.NewString()
	}()
	MockCartUserAddressJson(c, testStatusWithUSerAndId, "PATCH")
	delivery.ChangeStatus(c)
	require.Equal(t, 400, w.Code)
}

func TestChangeStatusInternal(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := zap.NewExample()
	itemUsecase := mocks.NewMockIItemUsecase(ctrl)
	categoryUsecase := mocks.NewMockICategoryUsecase(ctrl)
	cartUsecase := mocks.NewMockICartUsecase(ctrl)
	filestorage := fs.NewMockFileStorager(ctrl)
	orderUsecase := &mocks.OrderUsecaseMock{Err: fmt.Errorf("test error")}
	delivery := NewDelivery(itemUsecase, categoryUsecase, cartUsecase, logger, filestorage, orderUsecase)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	testStatusWithUSerAndId.User.Role = "admin"
	MockCartUserAddressJson(c, testStatusWithUSerAndId, "PATCH")
	delivery.ChangeStatus(c)
	require.Equal(t, 500, w.Code)
}*/
