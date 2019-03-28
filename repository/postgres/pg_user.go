package postgres

import (
	"fmt"
	"intel/isecl/tdservice/repository"

	"intel/isecl/tdservice/types"

	"github.com/jinzhu/gorm"
)

type PostgresUserRepository struct {
	db *gorm.DB
}

func (r *PostgresUserRepository) Create(u types.User) (*types.User, error) {

	uuid, err := repository.UUID()
	if err == nil {
		u.ID = uuid
	} else {
		return &u, err
	}
	err = r.db.Create(&u).Error
	return &u, err
}

func (r *PostgresUserRepository) Retrieve(u types.User) (*types.User, error) {
	err := r.db.Where(&u).First(&u).Error
	return &u, err
}

func (r *PostgresUserRepository) Update(u types.User) error {
	return r.db.Save(&u).Error
}

func (r *PostgresUserRepository) Delete(u types.User) error {
	return r.db.Delete(&u).Error
}

func (r *PostgresUserRepository) GetRoles(u types.User) (userRoles []types.Role, err error) {
	err = r.db.Select("roles.name, roles.domain").Joins("INNER JOIN user_roles on user_roles.role_id = roles.id INNER JOIN users on user_roles.user_id = users.id").Where(&u).Find(&userRoles).Error
	return userRoles, err
}

func (r *PostgresUserRepository) AddRoles(u types.User, roles []types.Role) error {

	// To be implemented later. We do not need this now as we are adding roles when the user is created.
	return fmt.Errorf("add roles function for user not implemented")
}
