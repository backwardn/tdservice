package mock

import (
	"intel/isecl/tdservice/types"
)

type MockUserRepository struct {
	CreateFunc   func(types.User) (*types.User, error)
	RetrieveFunc func(types.User) (*types.User, error)
	UpdateFunc   func(types.User) error
	DeleteFunc   func(types.User) error
}

func (m *MockUserRepository) Create(user types.User) (*types.User, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(user)
	}
	return nil, nil
}

func (m *MockUserRepository) Retrieve(user types.User) (*types.User, error) {
	if m.RetrieveFunc != nil {
		return m.RetrieveFunc(user)
	}
	return nil, nil
}

func (m *MockUserRepository) Update(user types.User) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(user)
	}
	return nil
}

func (m *MockUserRepository) Delete(user types.User) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(user)
	}
	return nil
}
