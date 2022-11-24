/*
 * Backend for Online Shop
 *
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package router

import (
	"OnlineShopBackend/internal/delivery"
	"log"

	"net/http"

	"github.com/gin-gonic/gin"
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
	// HandlerFunc is the handler function of this route.
	HandlerFunc gin.HandlerFunc
}

// Routes is the list of the generated Route.
type Routes []Route

type Router struct {
	router *gin.Engine
	del    *delivery.Delivery
	l      *zap.Logger
}

// NewRouter returns a new router.
func NewRouter(del *delivery.Delivery, logger *zap.Logger) *Router {
	log.Println("Enter in NewRouter()")
	router := gin.Default()
	routes := Routes{
		{
			"Index",
			http.MethodGet,
			"/",
			Index,
		},

		{
			"CreateCategory",
			http.MethodPost,
			"/categories/:category",
			del.CreateCategory,
		},

		{
			"CreateItem",
			http.MethodPost,
			"/items",
			del.CreateItem,
		},

		{
			"GetItem",
			http.MethodGet,
			"/items/:itemID",
			del.GetItem,
		},

		{
			"UpdateItem",
			http.MethodPut,
			"/items/:itemID",
			del.UpdateItem,
		},

		{
			"UploadFile",
			http.MethodPost,
			"/items/:itemID/upload",
			del.UploadFile,
		},

		{
			"GetCart",
			http.MethodGet,
			"/cart/:userID",
			delivery.GetCart,
		},

		{
			"GetCategoryList",
			http.MethodGet,
			"/items/categories/:category",
			del.GetCategoryList,
		},

		{
			"ItemsList",
			http.MethodGet,
			"/items",
			del.ItemsList,
		},

		{
			"SearchLine",
			http.MethodGet,
			"/search/:searchRequest",
			del.SearchLine,
		},

		{
			"CreateUser",
			http.MethodPost,
			"/user/create",
			delivery.CreateUser,
		},

		{
			"LoginUser",
			http.MethodPost,
			"/user/login",
			delivery.LoginUser,
		},

		{
			"LogoutUser",
			http.MethodPost,
			"/user/logout",
			delivery.LogoutUser,
		},
	}

	for _, route := range routes {
		switch route.Method {
		case http.MethodGet:
			router.GET(route.Pattern, route.HandlerFunc)
		case http.MethodPost:
			router.POST(route.Pattern, route.HandlerFunc)
		case http.MethodPut:
			router.PUT(route.Pattern, route.HandlerFunc)
		case http.MethodPatch:
			router.PATCH(route.Pattern, route.HandlerFunc)
		case http.MethodDelete:
			router.DELETE(route.Pattern, route.HandlerFunc)
		}
	}
	return &Router{router: router, del: del, l: logger}
}

// Index is the index handler.
func Index(c *gin.Context) {
	log.Println("Enter in Index")
	c.String(http.StatusOK, "Hello World!")
}

func (r *Router) Run(port string) error {
	return r.router.Run(port)
}
