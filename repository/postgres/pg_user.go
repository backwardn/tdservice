package postgres

import (
	"intel/isecl/tdservice/types"

	"github.com/jinzhu/gorm"
)

type PostgresUserRepository struct {
	db *gorm.DB
}

func (r *PostgresUserRepository) Create(u types.User) (*types.User, error) {
	err := r.db.Create(&u).Error
	return &u, err
}

func (r *PostgresUserRepository) Retrieve(u types.User) (*types.User, error) {
	err := r.db.First(&u).Error
	return &u, err
}

func (r *PostgresUserRepository) Update(u types.User) error {
	return r.db.Save(&u).Error
}

func (r *PostgresUserRepository) Delete(u types.User) error {
	return r.db.Delete(&u).Error
}
