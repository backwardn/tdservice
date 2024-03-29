/*
 * Copyright (C) 2019 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */
package repository

type TDSDatabase interface {
	Migrate() error
	HostRepository() HostRepository
	ReportRepository() ReportRepository
	UserRepository() UserRepository
	RoleRepository() RoleRepository
	Close()
}
