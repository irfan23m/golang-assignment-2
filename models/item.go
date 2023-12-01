package models

import "gorm.io/gorm"

type Item struct {
	gorm.Model
	ItemCode    string `json:"itemCode" gorm:"not null"`
	Description string `json:"description" gorm:"not null"`
	Quantity    int    `json:"quantity" gorm:"not null"`
	OrderID     int    `json:"orderId"`

	Order *Order `json:"-"`
}

type UpdateItem struct {
	ItemCode    string `json:"itemCode"`
	Description string `json:"description"`
	Quantity    int    `json:"quantity"`
}
