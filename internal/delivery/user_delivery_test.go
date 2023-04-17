package delivery

/*import (
	fs "OnlineShopBackend/internal/filestorage/mocks"
	"OnlineShopBackend/internal/models"
	"OnlineShopBackend/internal/usecase/mocks"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	//"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

var testModelNewUser = &models.User{
	Firstname: "Test",
	Password:  "password",
	Email:     "test@gmail.com",
}

//var testNewUserId = uuid.New()

//var testTokenSet = &jwtauth.Token{
//	AccessToken:  "",
//	RefreshToken: "",
//}

func TestCreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	logger := zap.L()
	itemUsecase := mocks.NewMockIItemUsecase(ctrl)
	categoryUsecase := mocks.NewMockICategoryUsecase(ctrl)
	userUsecase := mocks.NewMockIUserUsecase(ctrl)
	cartUsecase := mocks.NewMockICartUsecase(ctrl)
	filestorage := fs.NewMockFileStorager(ctrl)
	delivery := NewDelivery(itemUsecase, nil, categoryUsecase, cartUsecase, logger, filestorage)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}

	userUsecase.EXPECT().CreateUser(ctx, testModelNewUser).Return(testModelNewUser, nil)
	delivery.CreateUser(c)
	require.Equal(t, 201, w.Code)

}*/

import (
	"OnlineShopBackend/internal/delivery/user"
	"OnlineShopBackend/internal/models"
	mocks "OnlineShopBackend/internal/usecase/usecase_mocks"
	"bytes"
	"context"
	"encoding/json"
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

type WrongShortRights struct {
	Name int
}
type WrongRights struct {
	Id   int
	Name int
}

var (
	testShortRights = user.ShortRights{
		Name: "test",
	}
	testModelRightsNoId = models.Rights{
		Name: "test",
	}
	testSREmptyName = user.ShortRights{
		Name: "",
	}
	testWrongShortRights = WrongShortRights{
		Name: 5,
	}
)

func MockRightsJson(c *gin.Context, content interface{}, method string) {
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

func TestCreateRights(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := zap.L()
	itemUsecase := mocks.NewMockIItemUsecase(ctrl)
	categoryUsecase := mocks.NewMockICategoryUsecase(ctrl)
	userUsecase := mocks.NewMockIUserUsecase(ctrl)
	cartUsecase := mocks.NewMockICartUsecase(ctrl)
	orderUsecase := mocks.NewMockIOrderUsecase(ctrl)
	delivery := NewDelivery(itemUsecase, userUsecase, categoryUsecase, cartUsecase, logger, orderUsecase)
	ctx := context.Background()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	MockRightsJson(c, testWrongShortRights, post)
	delivery.CreateRights(c)
	require.Equal(t, 400, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	MockRightsJson(c, testSREmptyName, post)
	delivery.CreateRights(c)
	require.Equal(t, 400, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	MockRightsJson(c, testShortRights, post)
	userUsecase.EXPECT().CreateRights(ctx, &testModelRightsNoId).Return(uuid.Nil, err)
	delivery.CreateRights(c)
	require.Equal(t, 500, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	MockRightsJson(c, testShortRights, post)
	userUsecase.EXPECT().CreateRights(ctx, &testModelRightsNoId).Return(testId, nil)
	delivery.CreateRights(c)
	require.Equal(t, 201, w.Code)
}
