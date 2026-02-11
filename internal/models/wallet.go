package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Wallet struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID    uuid.UUID      `gorm:"type:uuid;uniqueIndex;not null" json:"user_id"`
	Balance   float64        `gorm:"type:decimal(15,2);default:0;not null" json:"balance"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate hook to generate UUID
func (w *Wallet) BeforeCreate(tx *gorm.DB) error {
	if w.ID == uuid.Nil {
		w.ID = uuid.New()
	}
	return nil
}
