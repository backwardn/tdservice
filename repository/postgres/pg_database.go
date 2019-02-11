package postgres

import (
	"intel/isecl/threat-detection-service/repository"
	"intel/isecl/threat-detection-service/types"

	"github.com/jinzhu/gorm"
)

type PostgresDatabase struct {
	DB *gorm.DB
}

func (pd *PostgresDatabase) Migrate() error {
	pd.DB.AutoMigrate(types.Host{}, types.Report{})
	return nil
}

func (pd *PostgresDatabase) HostRepository() repository.HostRepository {
	return &PostgresHostRepository{db: pd.DB}
}

func (pd *PostgresDatabase) ReportRepository() repository.ReportRepository {
	return nil
}
