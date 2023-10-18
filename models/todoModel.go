package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Todo struct {
    ID       uuid.UUID `gorm:"primary_key; unique; type:uuid; column:id; default:uuid_generate_v4()"`
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt gorm.DeletedAt `gorm:"index"`
    UserID  uuid.UUID  `gorm:"required"` // ID from user table this is this table's 'Foreign key'
    Title       string
    Description string
    Status      bool
}
// https://launchschool.com/books/sql/read/table_relationships
// Raw sql for foreign key, find gorm version.
// FOREIGN KEY (fk_col_name)
// REFERENCES target_table_name (pk_col_name);
