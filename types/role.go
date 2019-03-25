package types

import "time"

type Role struct {
	ID        string    `json:"id" gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
	// embed
	Name   string `json:"rolename" gorm:"not null"`
	Domain string `json:"roledomain,omitempty"`
}

type Roles []Role
