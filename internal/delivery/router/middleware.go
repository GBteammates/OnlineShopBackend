package router

import (
	"OnlineShopBackend/internal/delivery/users/user/jwtauth"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

const (
	authorizationHeader = "Authorization"
	admin               = "Admin"
	customer            = "Customer"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-Auth-Token, Set-Cookie")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
func JWTMiddleware(c *gin.Context) {

	tokenString := c.GetHeader(authorizationHeader)

	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization header is empty"})
		c.Abort()
		return
	}

	headerSplit := strings.Split(tokenString, " ")
	if len(headerSplit) != 2 || headerSplit[0] != "Bearer" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "header issue"})
		return
	}
	if len(headerSplit[1]) == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "empty token"})
		return
	}

	jwtKey, err := jwtauth.NewJWTKeyConfig()
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "empty config"})
		return
	}

	claims := &jwtauth.Payload{}
	token, err := jwt.ParseWithClaims(headerSplit[1], claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtKey.Key), nil //TODO
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid signature"})
			c.Abort()
			return
		}
	}

	if !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"message": tokenString})
		c.Abort()
		return
	}

	c.Set("claims", claims)
}

// AdminAuth method grants permission only to users with the role 'Admin'
func AdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		JWTMiddleware(c)
		userCr, ok := c.MustGet("claims").(*jwtauth.Payload)
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "incorrect claims"})
			return
		}
		if userCr.Role != admin {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "not permitted"})
			return
		}
		c.Next()
	}
}

// CustomerAuth method grants permission only to users with the role 'Customer'
func CustomerAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		JWTMiddleware(c)
		userCr, ok := c.MustGet("claims").(*jwtauth.Payload)
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "incorrect claims"})
			return
		}
		if userCr.Role != customer {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "not permitted"})
			return
		}
		c.Next()
	}
}

// UserAuth method confirms that user is authorized
func UserAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		JWTMiddleware(c)
		userCr, ok := c.MustGet("claims").(*jwtauth.Payload)
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "incorrect claims"})
			return
		}
		if userCr.UserId == uuid.Nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user unauthorized 333"})
			return
		}
		c.Next()
	}
}

// noOpMiddleware is a dummy method of middleware
func noOpMiddleware(c *gin.Context) {
	c.Next()
}
