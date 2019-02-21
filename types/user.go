package types

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           string     `gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	CreatedAt    time.Time  `json:"-"`
	UpdatedAt    time.Time  `json:"-"`
	DeletedAt    *time.Time `json:"-"`
	Name         string
	PasswordHash []byte
}

func (u *User) CheckPassword(password []byte) error {
	return bcrypt.CompareHashAndPassword(u.PasswordHash, password)
}
