package repository

import "intel/isecl/tdservice/types"

type RoleRepository interface {
	Create(types.Role) (*types.Role, error)
	Retrieve(types.Role) (*types.Role, error)
	Update(types.Role) error
	Delete(types.Role) error
}