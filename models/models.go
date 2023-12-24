package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Todo struct {
	ID        uuid.UUID `gorm:"primary_key; unique; type:uuid; column:id; default:uuid_generate_v4()"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	UserID      uuid.UUID `gorm:"foreignKey:ID,constraint:OnUpdate:CASCADE,OnDelete:SET NULL"` // ID from user table this is this table's 'Foreign key' // JOIN
	Title       string
	Description string
	Status      bool
}
type TeamTodo struct {
    Todo
    TeamID uuid.UUID `gorm:"foreignKey:ID,constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
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
	// Role      string         `json:"role"`
    Team uuid.UUID // foreign key
    Role string
	Username  string         `gorm:"unique" json:"userName"`
	Email     string         `gorm:"unique" json:"email" binding:"required,email"`
	Password  string         `json:"-" binding:"required,gte=6,lte=30"` // not working. pw can be 4 letters long atm.
	Todos     []Todo         // `gorm:"foreignKey:UserID,constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	//LoggedIn bool
}

type Team struct {
	ID        uuid.UUID `gorm:"primary_key; unique; type:uuid; column:id; default:uuid_generate_v4()"` // this needs to be associated to the todoIDs so we know which todo corresponds to whom.
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
    TeamName string
    TeamTodos []TeamTodo
    Members []*User `gorm:"-"` // MAYBE?
    // Members []User `gorm:"foreignKey:ID,constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	// Todos     []Todo
}

type Session struct {
	ID        uuid.UUID `gorm:"primary_key; unique; type:uuid; column:id; "` //default:uuid_generate_v4()"` // TOO LONG FOR bcrypt package to encrypt, 72 bits max!!!
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	//
	UserID uuid.UUID // partial index at ../migrate/migrate.go var sessionPartial // allows for soft delete
	Email  string
	Expiry time.Time
}

type Identifiable interface {
	// change models above from uuid, and THEN use this interface (after modding output of func) as input for func responsible for encrypting IDs stored on the DB?
	GetID() uuid.UUID
}

func (t *Todo) GetID() uuid.UUID {
	return t.ID
}

func (u *User) GetID() uuid.UUID {
	return u.ID
}

func (t *Team) GetID() uuid.UUID {
    return t.ID
}

func (s *Session) GetID() uuid.UUID {
	return s.ID
}
func (s Session) IsExpired() bool {
	return s.Expiry.Before(time.Now())
}
