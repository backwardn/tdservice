/*
 * Copyright (C) 2019 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */
package mock

import (
	"intel/isecl/tdservice/types"
)

type MockRoleRepository struct {
	CreateFunc   func(types.Role) (*types.Role, error)
	RetrieveFunc func(types.Role) (*types.Role, error)
	UpdateFunc   func(types.Role) error
	DeleteFunc   func(types.Role) error
}

func (m *MockRoleRepository) Create(role types.Role) (*types.Role, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(role)
	}
	return nil, nil
}

func (m *MockRoleRepository) Retrieve(role types.Role) (*types.Role, error) {
	if m.RetrieveFunc != nil {
		return m.RetrieveFunc(role)
	}
	return nil, nil
}

func (m *MockRoleRepository) Update(role types.Role) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(role)
	}
	return nil
}
 
func (m *MockRoleRepository) Delete(role types.Role) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(role)
	}
	return nil
}