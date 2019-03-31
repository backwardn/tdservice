/*
 * Copyright (C) 2019 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */
package constants

import "crypto"

const (
	HomeDir              = "/opt/tdservice/"
	ConfigDir            = "/etc/tdservice/"
	ExecutableDir        = "/opt/tdservice/bin/"
	LogDir               = "/var/log/tdservice/"
	LogFile              = "tdservice.log"
	HTTPLogFile          = "http.log"
	ConfigFile           = "config.yml"
	TLSCertFile          = "cert.pem"
	TLSKeyFile           = "key.pem"
	PIDFile              = "tdservice.pid"
	ServiceRemoveCmd     = "systemctl disable tdservice"
	HashingAlgorithm     = crypto.SHA384
	PasswordRandomLength = 20
)

const (
	// privileges granted: GET_ANY_HOST, DELETE_ANY_HOST, QUERY_REPORT, VERSION, CREATE_HOST
	AdminGroupName = "Administrators"

	// privileges granted: CREATE_HOST
	RegisterHostGroupName = "RegisterHosts"

	// privileges granted: GET_HOST, POST_REPORT
	HostSelfUpdateGroupName = "HostSelfUpdate"
)

// State represents whether or not a daemon is running or not
type State bool

const (
	// Stopped is the default nil value, indicating not running
	Stopped State = false
	// Running means the daemon is active
	Running State = true
)
