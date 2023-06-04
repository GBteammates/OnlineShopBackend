package models

import "github.com/google/uuid"

type Category struct {
	Id          uuid.UUID
	Name        string
	Description string
	Image       string
}

const (
	CategoriesList = "CategoriesList"
)
