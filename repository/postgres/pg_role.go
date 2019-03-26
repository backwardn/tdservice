package postgres

import (
	"intel/isecl/tdservice/types"

	"github.com/jinzhu/gorm"
)

type PostgresRoleRepository struct {
	db *gorm.DB
}

func (r *PostgresRoleRepository) Create(role types.Role) (*types.Role, error) {
	err := r.db.Create(&role).Error
	return &role, err
}

func (r *PostgresRoleRepository) Retrieve(role types.Role) (*types.Role, error) {
	err := r.db.First(&role).Error
	return &role, err
}

func (r *PostgresRoleRepository) Update(role types.Role) error {
	return r.db.Save(&role).Error
}

func (r *PostgresRoleRepository) Delete(role types.Role) error {
	return r.db.Delete(&role).Error
}