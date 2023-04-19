package models

import "github.com/google/uuid"

type User struct {
	Id       uuid.UUID   `json:"id"`
	Firstname string      `json:"firstname,omitempty"`
	Lastname  string      `json:"lastname,omitempty"`
	Password  string      `json:"password,omitempty"`
	Email     string      `json:"email,omitempty"`
	Address   UserAddress `json:"address,omitempty"`
	Rights    Rights      `json:"rights"`
}
