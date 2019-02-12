package types

import (
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           string `gorm:"primary_key;type:uuid;"`
	Name         string
	PasswordHash []byte
}

func (u *User) CheckPassword(password []byte) error {
	return bcrypt.CompareHashAndPassword(u.PasswordHash, password)
}
