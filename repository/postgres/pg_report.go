package postgres

import (
	"intel/isecl/tdservice/repository"
	"intel/isecl/tdservice/types"

	"github.com/jinzhu/gorm"
)

type PostgresReportRepository struct {
	db *gorm.DB
}

func (r *PostgresReportRepository) Create(report types.Report) (*types.Report, error) {
	var err error
	err = r.db.Create(&report).Error
	return &report, err
}

func (r *PostgresReportRepository) Retrieve(report types.Report) (*types.Report, error) {
	err := r.db.Preload("Host").First(&report).Error
	return &report, err
}

func (r *PostgresReportRepository) RetrieveAll(report types.Report) ([]types.Report, error) {
	var reports []types.Report
	err := r.db.Where(&report).Find(&reports).Error
	return reports, err
}

func (r *PostgresReportRepository) RetrieveByFilterCriteria(filter repository.ReportFilter) ([]types.Report, error) {
	var reports []types.Report
	query := r.db.Where(&filter.Report)

	// Add From
	if !filter.From.IsZero() {
		query = query.Where("(detection_json->>'timestamp')::int >= ?", filter.From.Unix())
	}

	// Add To
	if !filter.To.IsZero() {
		query = query.Where("(detection_json->>'timestamp')::int <= ?", filter.To.Unix())
	}

	// Add Hostname
	if filter.Hostname != "" {
		query = query.Joins("JOIN hosts ON hosts.id = reports.host_id").Where("hosts.hostname = ?", filter.Hostname)
	}

	err := query.Find(&reports).Error
	return reports, err
}

func (r *PostgresReportRepository) Update(report types.Report) error {
	return r.db.Save(report).Error
}

func (r *PostgresReportRepository) Delete(report types.Report) error {
	return r.db.Delete(report).Error
}
