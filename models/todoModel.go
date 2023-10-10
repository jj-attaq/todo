package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Todo struct {
	// ID        string `gorm:"default:uuid_generate_v4()"`
    ID       uuid.UUID `gorm:"primary_key; unique; type:uuid; column:id; default:uuid_generate_v4()"`
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt gorm.DeletedAt `gorm:"index"`
    // UserID   string
    // UserPW   string
    Title       string
    Description string
    Status      bool
    // UniqueID string
}

/* func (user *Todo) BeforeCreate(db *gorm.DB) error {
	user.ID = uuid.New().String()
	return nil
} */
