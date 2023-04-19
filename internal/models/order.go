package models

import (
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	StatusCreated    Status = "order created"
	StatusProcessing Status = "order processing"
	StatusProcessed  Status = "order processed"
	StatusReady      Status = "ready for shipment"
	StatusCourier    Status = "picked by courier"
	StatusShipped    Status = "delivered"

	StandardShipmentPeriod  time.Duration = 24 * 3 * time.Hour
	ProlongedShipmentPeriod time.Duration = 24 * 7 * time.Hour
)

type Order struct {
	Id           uuid.UUID
	CreatedAt    time.Time
	ShipmentTime time.Time
	User         User
	Address      UserAddress
	Status       Status
	Items        []ItemWithQuantity
}
