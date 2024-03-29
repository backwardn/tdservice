/*
 * Copyright (C) 2019 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */
package types

import "time"

type HostInfo struct {
	Hostname string `json:"hostname,omitempty"`
	Version  string `json:"version"`
	Build    string `json:"build"`
	OS       string `json:"os"`
}

type Host struct {
	ID        string    `json:"id" gorm:"primary_key;type:uuid"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
	// embed
	HostInfo
	Status string `json:"status"`
}

type HostCreateResponse struct {
	Host
	User  string `json:user`
	Token string `json:token`
}
