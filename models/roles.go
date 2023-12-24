package models

import (
    "time"

    "github.com/google/uuid"
    "gorm.io/gorm"
)

type Permissions struct {
	ID        uuid.UUID `gorm:"primary_key; unique; type:uuid; column:id; "` //default:uuid_generate_v4()"`
	// ID        uuid.UUID `gorm:"primary_key; unique; type:uuid; column:id; "` //default:uuid_generate_v4()"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
    // UserID uuid.UUID `gorm:"foreignKey:ID,constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	// Team
	TeamCreate bool
	TeamRead   bool
	TeamUpdate bool
	TeamDelete bool
	//TeamTodos
	TTodoCreate bool
	TTodoRead   bool
	TTodoUpdate bool
	TTodoDelete bool
	// Todo
	TodoCreate bool
	TodoRead   bool
	TodoUpdate bool
	TodoDelete bool
}

var AdminPermissions = Permissions{
	// Team
	TeamCreate: true,
	TeamRead:   true,
	TeamUpdate: true,
	TeamDelete: true,
	// TeamTodos
	TTodoCreate: true,
	TTodoRead:   true,
	TTodoUpdate: true,
	TTodoDelete: true,
	// Todos
	TodoCreate: true,
	TodoRead:   true,
	TodoUpdate: true,
	TodoDelete: true,
}

var TeamLeadPermissions = Permissions{
	// Team
	TeamCreate: false,
	TeamRead:   true,
	TeamUpdate: true,
	TeamDelete: false,
	// TeamTodos
	TTodoCreate: true,
	TTodoRead:   true,
	TTodoUpdate: true,
	TTodoDelete: true,
	// Todos
	TodoCreate: true,
	TodoRead:   true,
	TodoUpdate: true,
	TodoDelete: true,
}

var MemberPermissions = Permissions{
	// Team
	TeamCreate: false,
	TeamRead:   true,
	TeamUpdate: false,
	TeamDelete: false,
	// TeamTodos
	TTodoCreate: false,
	TTodoRead:   true,
	TTodoUpdate: false,
	TTodoDelete: false,
	// Todos
	TodoCreate: true,
	TodoRead:   true,
	TodoUpdate: true,
	TodoDelete: true,
}

type Roles struct {
    Admin
    TeamLead
    Member
}

type TeamMember interface {
    // MkAdmin() Admin
    // MkTeamLead() TeamLead
    // MkMember() Member
	initialize()
	// GetPermissions() Permissions
	// GetAdminStatus() bool
}

type Admin struct {
	Permissions
	IsAdmin    bool
	IsTeamLead bool
}

func (a *Admin) initialize(u *User) {
    a.ID = u.ID
	a.Permissions = AdminPermissions
	a.IsAdmin = true
	a.IsTeamLead = true
}

type TeamLead struct {
	Permissions
	IsAdmin    bool
	IsTeamLead bool
}

func (tl *TeamLead) initialize(u *User) {
    tl.ID = u.ID
	tl.Permissions = TeamLeadPermissions
	tl.IsAdmin = false
	tl.IsTeamLead = true
}

type Member struct {
	Permissions
	IsAdmin    bool
	IsTeamLead bool
}

func (m *Member) initialize(u *User) {
    m.ID = u.ID
	m.Permissions = MemberPermissions
	m.IsAdmin = false
	m.IsTeamLead = false
}

// GetPermissions()
func (a *Admin) GetPermissions() Permissions {
	return a.Permissions
}
func (tl *TeamLead) GetPermissions() Permissions {
	return tl.Permissions
}
func (m *Member) GetPermissions() Permissions {
	return m.Permissions
}

// GetAdminStatus()
func (a *Admin) GetAdminStatus() bool {
	return a.IsAdmin
}
func (tl *TeamLead) GetAdminStatus() bool {
	return tl.IsAdmin
}
func (m *Member) GetAdminStatus() bool {
	return m.IsAdmin
}

// func (a *Admin) ChangeRole() {
// }

func (u *User) MkAdmin() Admin {
	var a Admin
	a.initialize(u)
	return a
}

func (u *User) MkTeamLead() TeamLead {
	var tl TeamLead
	tl.initialize(u)
	return tl
}

func (u *User) MkMember() Member {
	var m Member
	m.initialize(u)
	return m
}
