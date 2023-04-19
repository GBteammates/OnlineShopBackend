package router

import (
	carts "OnlineShopBackend/internal/delivery/carts"
	categories "OnlineShopBackend/internal/delivery/categories"
	items "OnlineShopBackend/internal/delivery/items"
	orders "OnlineShopBackend/internal/delivery/orders"
	"OnlineShopBackend/internal/delivery/swagger/docs"
	users "OnlineShopBackend/internal/delivery/users"
	"net/http"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
)

// Route is the information for every URI.
type Route struct {
	// Name is the name of this Route.
	Name string
	// Method is the string for the HTTP method. ex) GET, POST etc..
	Method string
	// Pattern is the pattern of the URI.
	Pattern string
	//
	Middleware gin.HandlerFunc
	// HandlerFunc is the handler function of this route.
	HandlerFunc gin.HandlerFunc
}

// Routes is the list of the generated Route.
type Routes []Route

type Router struct {
	*gin.Engine
	itemDelivery     *items.ItemDelivery
	categoryDelivery *categories.CategoryDelivery
	cartDelivery     *carts.CartDelivery
	orderDelivery    *orders.OrderDelivery
	userDelivery     *users.UserDelivery
	logger           *zap.Logger
}

// NewRouter returns a new router.
func NewRouter(itemDelivery *items.ItemDelivery,
	categoryDelivery *categories.CategoryDelivery,
	cartDelivery *carts.CartDelivery,
	orderDelivery *orders.OrderDelivery,
	userDelivery *users.UserDelivery,
	logger *zap.Logger) *Router {

	logger.Debug("Enter in NewRouter()")
	gin := gin.Default()
	gin.Use(CORSMiddleware())
	gin.Use(ginzap.RecoveryWithZap(logger, true))
	gin.Static("/files", "./static/files")
	docs.SwaggerInfo.BasePath = "/"
	gin.Group("/docs").Any("/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router := &Router{
		itemDelivery:  itemDelivery,
		categoryDelivery: categoryDelivery,
		cartDelivery:  cartDelivery,
		orderDelivery: orderDelivery,
		userDelivery:  userDelivery,
		logger:        logger,
	}

	routes := Routes{
		// -------------------------CATEGORY----------------------------------------------------------------------------
		{
			"CreateCategory",
			http.MethodPost,
			"/categories/create",
			AdminAuth(),
			categoryDelivery.CreateCategory,
		},
		{
			"GetCategory",
			http.MethodGet,
			"/categories/:categoryID",
			noOpMiddleware,
			categoryDelivery.GetCategory,
		},
		{
			"GetCategoryList",
			http.MethodGet,
			"/categories/list",
			noOpMiddleware,
			categoryDelivery.GetCategoryList,
		},
		{
			"UpdateCategory",
			http.MethodPut,
			"/categories/:categoryID",
			AdminAuth(),
			categoryDelivery.UpdateCategory,
		},
		{
			"UploadCategoryImage",
			http.MethodPost,
			"/categories/image/upload/:categoryID",
			AdminAuth(),
			categoryDelivery.UploadCategoryImage,
		},
		{
			"DeleteCategoryImage",
			http.MethodDelete,
			"/categories/image/delete", //?id=25f32441-587a-452d-af8c-b3876ae29d45&name=20221209194557.jpeg
			AdminAuth(),
			categoryDelivery.DeleteCategoryImage,
		},
		{
			"DeleteCategory",
			http.MethodDelete,
			"/categories/delete/:categoryID",
			AdminAuth(),
			categoryDelivery.DeleteCategory,
		},
		{
			"GetCategoriesImagesList",
			http.MethodGet,
			"/categories/images/list",
			AdminAuth(),
			categoryDelivery.GetCategoriesImagesList,
		},
		// -------------------------ITEM--------------------------------------------------------------------------------
		{
			"CreateItem",
			http.MethodPost,
			"/items/create",
			AdminAuth(),
			itemDelivery.CreateItem,
		},
		{
			"GetItem",
			http.MethodGet,
			"/items/:itemID",
			noOpMiddleware,
			itemDelivery.GetItem,
		},
		{
			"GetItemsByCategory",
			http.MethodGet,
			"/items/", //?param=categoryName&offset=20&limit=10&sort_type=name&sort_order=asc (sort_type == name or price, sort_order == asc or desc)
			noOpMiddleware,
			itemDelivery.GetItemsByCategory,
		},
		{
			"UpdateItem",
			http.MethodPut,
			"/items/update",
			AdminAuth(),
			itemDelivery.UpdateItem,
		},
		{
			"UploadItemImage",
			http.MethodPost,
			"/items/image/upload/:itemID",
			AdminAuth(),
			itemDelivery.UploadItemImage,
		},
		{
			"DeleteItemImage",
			http.MethodDelete,
			"/items/image/delete", //?id=25f32441-587a-452d-af8c-b3876ae29d45&name=20221209194557.jpeg
			AdminAuth(),
			itemDelivery.DeleteItemImage,
		},
		{
			"ItemsQuantity",
			http.MethodGet,
			"/items/quantity",
			noOpMiddleware,
			itemDelivery.ItemsQuantity,
		},
		{
			"ItemsQuantityInCategory",
			http.MethodGet,
			"/items/quantityCat/:categoryName",
			noOpMiddleware,
			itemDelivery.ItemsQuantityInCategory,
		},
		{
			"ItemsQuantityInFavourite",
			http.MethodGet,
			"/items/quantityFav/:userID",
			noOpMiddleware,
			itemDelivery.ItemsQuantityInFavourite,
		},
		{
			"ItemsQuantityInSearch",
			http.MethodGet,
			"/items/quantitySearch/:searchRequest",
			noOpMiddleware,
			itemDelivery.ItemsQuantityInSearch,
		},
		{
			"ItemsList",
			http.MethodGet,
			"/items/list", //?offset=20&limit=10&sort_type=name&sort_order=asc (sort_type == name or price, sort_order == asc or desc)
			noOpMiddleware,
			itemDelivery.ItemsList,
		},
		{
			"SearchLine",
			http.MethodGet,
			"/items/search/", //?param=searchRequest&offset=20&limit=10&sort_type=name&sort_order=asc (sort_type == name or price, sort_order == asc or desc)
			noOpMiddleware,
			itemDelivery.SearchLine,
		},
		{
			"DeleteItem",
			http.MethodDelete,
			"/items/delete/:itemID",
			AdminAuth(),
			itemDelivery.DeleteItem,
		},
		{
			"AddFavouriteItem",
			http.MethodPost,
			"/items/addFavItem",
			UserAuth(),
			itemDelivery.AddFavouriteItem,
		},
		{
			"DeleteFavouriteItem",
			http.MethodDelete,
			"/items/deleteFav/:userID/:itemID",
			UserAuth(),
			itemDelivery.DeleteFavouriteItem,
		},
		{
			"GetFavouriteItems",
			http.MethodGet,
			"/items/favList/", //?param=userIDt&offset=20&limit=10&sort_type=name&sort_order=asc (sort_type == name or price, sort_order == asc or desc)
			UserAuth(),
			itemDelivery.GetFavouriteItems,
		},
		{
			"GetItemsImagesList",
			http.MethodGet,
			"/items/images/list",
			AdminAuth(),
			itemDelivery.GetItemsImagesList,
		},
		// -------------------------CART--------------------------------------------------------------------------------
		{
			"GetCart",
			http.MethodGet,
			"/cart/:cartID",
			UserAuth(),
			cartDelivery.GetCart,
		},
		{
			"GetCartByUserId",
			http.MethodGet,
			"/cart/byUser/:userID",
			UserAuth(),
			cartDelivery.GetCartByUserId,
		},
		{
			"CreateCart",
			http.MethodPost,
			"/cart/create/:userID",
			UserAuth(),
			cartDelivery.CreateCart,
		},
		{
			"AddItemToCart",
			http.MethodPut,
			"/cart/addItem",
			UserAuth(),
			cartDelivery.AddItemToCart,
		},
		{
			"DeleteItemFromCart",
			http.MethodDelete,
			"/cart/delete/:cartID/:itemID",
			UserAuth(),
			cartDelivery.DeleteItemFromCart,
		},
		{
			"DeleteCart",
			http.MethodDelete,
			"/cart/delete/:cartID",
			UserAuth(),
			cartDelivery.DeleteCart,
		},
		// -------------------------USER--------------------------------------------------------------------------------
		{
			"CreateUser",
			http.MethodPost,
			"/user/create",
			noOpMiddleware,
			userDelivery.CreateUser,
		},
		{
			"LoginUser",
			http.MethodPost,
			"/user/login",
			noOpMiddleware,
			userDelivery.LoginUser,
		},
		{
			"LogoutUser",
			http.MethodGet,
			"/user/logout",
			noOpMiddleware,
			userDelivery.LogoutUser,
		},
		{
			"LoginUserGoogle",
			http.MethodGet,
			"/user/login/google",
			noOpMiddleware,
			userDelivery.LoginUserGoogle,
		},
		{
			"callbackGoogle",
			http.MethodGet,
			"/user/callbackGoogle",
			noOpMiddleware,
			userDelivery.CallbackGoogle,
		},

		{
			"userProfile",
			http.MethodGet,
			"/user/profile",
			UserAuth(),
			userDelivery.GetUserProfile,
		},
		{
			"userProfileUpdate",
			http.MethodPut,
			"/user/profile/edit",
			UserAuth(),
			userDelivery.UpdateUserData,
		},
		{
			"ChangeRole",
			http.MethodPut,
			"/user/role/update",
			AdminAuth(),
			userDelivery.ChangeRole,
		},
		{
			"UserRolesList",
			http.MethodGet,
			"/user/rights/list",
			AdminAuth(),
			userDelivery.RolesList,
		},
		{
			"CreateRights",
			http.MethodPost,
			"/user/createRights",
			AdminAuth(),
			userDelivery.CreateRights,
		},
		// -------------------------ORDER--------------------------------------------------------------------------------
		{
			"CreateOrder",
			http.MethodPost,
			"/order/create",
			UserAuth(),
			orderDelivery.CreateOrder,
		},
		{
			"GetOrder",
			http.MethodGet,
			"/order/:orderID",
			UserAuth(),
			orderDelivery.GetOrder,
		},
		{
			"GetOrdersForUsers",
			http.MethodGet,
			"/order/list/:userID",
			UserAuth(),
			orderDelivery.GetOrdersForUser,
		},
		{
			"DeleteOrder",
			http.MethodDelete,
			"/order/delete/:orderID",
			AdminAuth(),
			orderDelivery.DeleteOrder,
		},
		{
			"ChangeAddress",
			http.MethodPatch,
			"/order/changeaddress",
			UserAuth(),
			orderDelivery.ChangeAddress,
		},
		{
			"ChangeStatus",
			http.MethodPatch,
			"/order/changestatus",
			AdminAuth(),
			orderDelivery.ChangeStatus,
		},
	}

	for _, route := range routes {
		switch route.Method {
		case http.MethodGet:
			gin.GET(route.Pattern, route.Middleware, route.HandlerFunc)
		case http.MethodPost:
			gin.POST(route.Pattern, route.Middleware, route.HandlerFunc)
		case http.MethodPut:
			gin.PUT(route.Pattern, route.Middleware, route.HandlerFunc)
		case http.MethodPatch:
			gin.PATCH(route.Pattern, route.Middleware, route.HandlerFunc)
		case http.MethodDelete:
			gin.DELETE(route.Pattern, route.Middleware, route.HandlerFunc)
		}
	}
	router.Engine = gin
	return router
}
