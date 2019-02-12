package repository

import "intel/isecl/threat-detection-service/types"

type UserRepository interface {
	Create(types.User) error
	Retrieve(types.User) (*types.User, error)
	Update(types.User) error
	Delete(types.User) error
}
