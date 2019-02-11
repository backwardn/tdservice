package postgres

import (
	"intel/isecl/threat-detection-service/repository"
	"intel/isecl/threat-detection-service/types"

	"github.com/jinzhu/gorm"
)

type PostgresReportRepository struct {
	db *gorm.DB
}

func (r *PostgresReportRepository) Create(report types.Report) error {
	return r.db.Create(&report).Error
}

func (r *PostgresReportRepository) Retrieve(report types.Report) (*types.Report, error) {
	err := r.db.First(&report).Error
	return &report, err
}

func (r *PostgresReportRepository) RetrieveByFilterCriteria(filter repository.ReportFilter) ([]types.Report, error) {
	return nil, nil
}

func (r *PostgresReportRepository) Update(report types.Report) error {
	return r.db.Save(&report).Error
}

func (r *PostgresReportRepository) Delete(report types.Report) error {
	return r.db.Delete(&report).Error
}
