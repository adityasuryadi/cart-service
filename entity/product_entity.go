package entity

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	Id          uuid.UUID `gorm:"primaryKey;type:uuid;" column:"id"`
	Name        string    `gorm:"column:name"`
	Price       float64   `gorm:"column:price"`
	Qty         int       `gorm:"column:qty"`
	Description string    `gorm:"column:description"`
	CreatedAt   time.Time `gorm:"column:created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at"`
}

func (Product) TableName() string {
	return "products"
}
