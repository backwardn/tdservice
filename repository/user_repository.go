package repository

import "intel/isecl/tdservice/types"

type UserRepository interface {
	Create(types.User) (*types.User, error)
	Retrieve(types.User) (*types.User, error)
	Update(types.User) error
	Delete(types.User) error
}
