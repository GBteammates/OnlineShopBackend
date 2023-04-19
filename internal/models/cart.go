package models

import (
	"time"

	"github.com/google/uuid"
)

type Cart struct {
	Id       uuid.UUID
	UserId   uuid.UUID
	Items    []ItemWithQuantity
	ExpireAt time.Time
}
