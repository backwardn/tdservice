/*
 * Copyright (C) 2019 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */
package postgres

import (
	"fmt"
	"intel/isecl/tdservice/repository"
	"intel/isecl/tdservice/types"
	"time"

	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

type PostgresDatabase struct {
	DB *gorm.DB
}

func (pd *PostgresDatabase) Migrate() error {
	pd.DB.AutoMigrate(types.Host{}, types.Report{}, types.User{}, types.Role{})
	pd.DB.Table("reports").
		AddForeignKey("host_id", "hosts(id)", "SET NULL", "SET NULL")
	return nil
}

func (pd *PostgresDatabase) HostRepository() repository.HostRepository {
	return &PostgresHostRepository{db: pd.DB}
}

func (pd *PostgresDatabase) ReportRepository() repository.ReportRepository {
	return &PostgresReportRepository{db: pd.DB}
}

func (pd *PostgresDatabase) UserRepository() repository.UserRepository {
	return &PostgresUserRepository{db: pd.DB}
}

func (pd *PostgresDatabase) RoleRepository() repository.RoleRepository {
	return &PostgresRoleRepository{db: pd.DB}
}

func (pd *PostgresDatabase) Close() {
	if pd.DB != nil {
		pd.DB.Close()
	}
}

func Open(host string, port int, dbname string, user string, password string, ssl bool) (*PostgresDatabase, error) {
	var sslMode string
	if ssl {
		sslMode = "require"
	} else {
		sslMode = "disable"
	}
	var db *gorm.DB
	var dbErr error
	const numAttempts = 4
	for i := 0; i < numAttempts; i = i + 1 {
		const retryTime = 5
		db, dbErr = gorm.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=%s",
			host, port, user, dbname, password, sslMode))
		if dbErr != nil {
			log.WithError(dbErr).Infof("Failed to connect to DB, retrying attempt %d/%d", i, numAttempts)
		} else {
			break
		}
		time.Sleep(retryTime * time.Second)
	}
	if dbErr != nil {
		log.WithError(dbErr).Infof("Failed to connect to db after %d attempts\n", numAttempts)
		return nil, dbErr
	}
	return &PostgresDatabase{DB: db}, nil
}
