package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Todo struct {
	ID        uuid.UUID `gorm:"primary_key; unique; type:uuid; column:id; default:uuid_generate_v4()"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	// User User
	UserID      uuid.UUID `gorm:"foreignKey:ID,constraint:OnUpdate:CASCADE,OnDelete:SET NULL"` // ID from user table this is this table's 'Foreign key' // JOIN
	Title       string
	Description string
	Status      bool
}

// https://launchschool.com/books/sql/read/table_relationships
// Raw sql for foreign key, find gorm version.
// FOREIGN KEY (fk_col_name)
// REFERENCES target_table_name (pk_col_name);

type User struct {
	ID        uuid.UUID `gorm:"primary_key; unique; type:uuid; column:id; default:uuid_generate_v4()"` // this needs to be associated to the todoIDs so we know which todo corresponds to whom.
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Role      string         `json:"role"`
	Username  string         `gorm:"unique" json:"userName"`
	Email     string         `gorm:"unique" json:"email" binding:"required,email"`
	Password  string         `json:"-" binding:"required,gte=6,lte=30"`
	Todos     []Todo         // `gorm:"foreignKey:UserID,constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
    //LoggedIn bool
}

type Session struct {
	ID        uuid.UUID `gorm:"primary_key; unique; type:uuid; column:id; "` //default:uuid_generate_v4()"`
	CreatedAt time.Time
	UpdatedAt time.Time
    DeletedAt gorm.DeletedAt `gorm:"index"`
    //
    UserID uuid.UUID // partial index at ../migrate/migrate.go var sessionPartial // allows for soft delete
    Email string
    Expiry time.Time
}

func (s Session) IsExpired() bool {
    return s.Expiry.Before(time.Now())
}
