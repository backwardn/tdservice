package repository

import "intel/isecl/tdservice/types"

type HostRepository interface {
	// Create should return a pointer to Host
	Create(host types.Host) (*types.Host, error)
	Retrieve(host types.Host) (*types.Host, error)
	RetrieveAll(host types.Host) ([]types.Host, error)
	Update(host types.Host) error
	Delete(host types.Host) error
}
