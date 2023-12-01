package models

import (
	"time"

	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	CustomerName string    `json:"customerName" gorm:"not null"`
	OrderedAt    time.Time `json:"orderedAt"`
	Items        []Item    `json:"items" gorm:"foreignKey:OrderID"`
}

type UpdateOrder struct {
	ID           uint         `json:"id"`
	CustomerName string       `json:"customerName"`
	UpdateItem   []UpdateItem `json:"items"`
}
