package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type User struct {
    ID       uuid.UUID `gorm:"primary_key; unique; type:uuid; column:id; default:uuid_generate_v4()"` // this needs to be associated to the todoIDs so we know which todo corresponds to whom.
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt gorm.DeletedAt `gorm:"index"`
    Role string `json:"role"`
    Username string `gorm:"unique" json:"userName"`
    Email string `gorm:"unique" json:"email" binding:"required,email"`
    // FirstName string `json:"firstName"`
    // LastName string `json:"lastName"`
    Password string `json:"-" binding:"required,gte=6,lte=30"`
}

/* func ValidateUser(u *User) bool {
    return true
} */
