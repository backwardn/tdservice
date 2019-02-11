package repository

import "intel/isecl/threat-detection-service/types"

type ReportFilter struct {
}

type ReportRepository interface {
	Create(report types.Report) error
	Retrieve(report types.Report) (*types.Report, error)
	RetrieveByFilterCriteria(filter ReportFilter) ([]types.Report, error)
	Update(report types.Report) error
	Delete(report types.Report) error
}
