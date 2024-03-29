/*
 * Copyright (C) 2019 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */
package tasks

import (
	"errors"
	"intel/isecl/lib/common/setup"
	"intel/isecl/tdservice/repository"
	"intel/isecl/tdservice/repository/mock"
	"intel/isecl/tdservice/types"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateAdmin(t *testing.T) {
	m := &mock.MockDatabase{}
	m.MockUserRepository.RetrieveFunc = func(u types.User) (*types.User, error) {
		return &u, nil
	}
	task := Admin{
		Flags: []string{},
		DatabaseFactory: func() (repository.TDSDatabase, error) {
			return m, nil
		},
		ConsoleWriter: os.Stdout,
	}
	ctx := setup.Context{}
	err := task.Validate(ctx)
	assert.NoError(t, err)
}

func TestCreateAdmin(t *testing.T) {
	m := &mock.MockDatabase{}
	var user *types.User
	m.MockUserRepository.CreateFunc = func(u types.User) (*types.User, error) {
		user = &u
		return user, nil
	}
	m.MockUserRepository.RetrieveFunc = func(u types.User) (*types.User, error) {
		if user == nil {
			return nil, errors.New("Record not found")
		}
		return user, nil
	}
	task := Admin{
		Flags: []string{"-admin-user=admin", "-admin-pass=foobarfoobar", "-reg-host-user=regHostUser", "-reg-host-password=foobar"},
		DatabaseFactory: func() (repository.TDSDatabase, error) {
			return m, nil
		},
		ConsoleWriter: os.Stdout,
	}
	ctx := setup.Context{}
	err := task.Run(ctx)
	assert.NoError(t, err)
}

func TestCreateAdminForce(t *testing.T) {

}
