/*
 * Copyright (C) 2019 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */
package context

import (
	"intel/isecl/tdservice/types"

	"net/http"

	"github.com/gorilla/context"
)

type key int

const (
	rolesKey = iota
)

func SetUserRoles(r *http.Request, val types.Roles) {
	context.Set(r, rolesKey, val)
}

func GetUserRoles(r *http.Request) types.Roles {
	if rv := context.Get(r, rolesKey); rv != nil {
		return rv.(types.Roles)
	}
	return nil
}
