package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Cart struct {
	Id           uuid.UUID `gorm:"primaryKey;type:uuid;" column:"id"`
	ProductId    uuid.UUID `gorm:"column:product_id;type:uuid"`
	UserId       uuid.UUID `gorm:"column:user_id;type:uuid"`
	ProductName  string    `gorm:"column:product_name"`
	ProductPrice float64   `gorm:"column:product_price"`
	Qty          int       `gorm:"column:qty"`
	Email        string    `gorm:"column:email"`
	TotalPrice   float64   `gorm:"column:total_price"`
	CreatedAt    time.Time `gorm:"column:created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at"`
}

func (Cart) TableName() string {
	return "carts"
}

func (entity *Cart) BeforeCreate(db *gorm.DB) error {
	entity.Id = uuid.New()
	entity.CreatedAt = time.Now().Local()
	return nil
}

func (entity *Cart) BeforeUpdate(db *gorm.DB) error {
	entity.UpdatedAt = time.Now().Local()
	return nil
}
