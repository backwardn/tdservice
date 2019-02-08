package postgres

import (
	"intel/isecl/threat-detection-service/repository"

	"github.com/jinzhu/gorm"
)

type PostgresDatabase struct {
	DB *gorm.DB
}

func (pd *PostgresDatabase) Migrate() error {
	pd.DB.AutoMigrate()
	return nil
}

func (pd *PostgresDatabase) HostRepository() repository.HostRepository {
	return nil
}

func (pd *PostgresDatabase) ReportRepository() repository.ReportRepository {
	return nil
}
