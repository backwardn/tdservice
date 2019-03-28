package types

import "time"

type Role struct {
	ID        string    `json:"id" gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`

	// Name: UpdateHost
	Name string `json:"rolename" gorm:"not null"`
	// 1234-88769876-28768
	Domain string `json:"roledomain,omitempty"`
}

type Roles []Role
