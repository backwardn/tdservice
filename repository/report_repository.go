package repository

import "intel/isecl/threat-detection-service/types"

type ReportFilter struct {
}

type ReportRepository interface {
	// everything should take a non pointer struct and Create should return a pointer
	Create(report types.Report) (*types.Report, error)
	Retrieve(report types.Report) (*types.Report, error)
	RetrieveAll(report types.Report) ([]types.Report, error)
	RetrieveByFilterCriteria(filter ReportFilter) ([]types.Report, error)
	Update(report types.Report) error
	Delete(report types.Report) error
}
