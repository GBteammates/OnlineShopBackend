package user

import (
	"OnlineShopBackend/internal/delivery/users/user/jwtauth"
	"OnlineShopBackend/internal/models"

	"github.com/google/uuid"
)

type LoginResponseData struct {
	CartId uuid.UUID `json:"cartId"`
	Token jwtauth.Token `json:"token"`

}

type CreateUserData struct {
	ID        uuid.UUID          `json:"id,omitempty"`
	Firstname string             `json:"firstname,omitempty" binding:"required" example:"Jane"`
	Lastname  string             `json:"lastname,omitempty" binding:"required" example:"Doe"`
	Password  string             `json:"password,omitempty"`
	Email     string             `json:"email,omitempty"`
	Address   models.UserAddress `json:"address,omitempty" binding:"required"`
	Rights    models.Rights      `json:"rights,omitempty"`
}

type ShortRights struct {
	Name  string   `json:"name" binding:"required" example:"admin"`
	Rules []string `json:"rules,omitempty"`
}

type RightsId struct {
	Value string `json:"id" uri:"itemID" binding:"required,uuid" example:"00000000-0000-0000-0000-000000000000" format:"uuid"`
}

type Credentials struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}