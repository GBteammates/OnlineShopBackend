/*
 * Backend for Online Shop
 *
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

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
