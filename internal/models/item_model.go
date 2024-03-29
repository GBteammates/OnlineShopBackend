/*
 * Backend for Online Shop
 *
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package models

import "github.com/google/uuid"

type Item struct {
	Id          uuid.UUID
	Title       string
	Description string
	Price       int32
	Category    Category
	Vendor      string
	Images      []string
}

type ItemWithQuantity struct {
	Item
	Quantity int
}