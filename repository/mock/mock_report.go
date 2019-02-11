package mock

import (
	"intel/isecl/threat-detection-service/repository"
	"intel/isecl/threat-detection-service/types"
)

type MockReportRepository struct {
	CreateFunc                   func(types.Report) error
	RetrieveFunc                 func(types.Report) (*types.Report, error)
	RetrieveByFilterCriteriaFunc func(repository.ReportFilter) ([]types.Report, error)
	UpdateFunc                   func(types.Report) error
	DeleteFunc                   func(types.Report) error
}

func (m *MockReportRepository) Create(r types.Report) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(r)
	}
	return nil
}

func (m *MockReportRepository) Retrieve(r types.Report) (*types.Report, error) {
	if m.RetrieveFunc != nil {
		return m.RetrieveFunc(r)
	}
	return nil, nil
}

func (m *MockReportRepository) RetrieveByFilterCriteria(f repository.ReportFilter) ([]types.Report, error) {
	if m.RetrieveByFilterCriteriaFunc != nil {
		return m.RetrieveByFilterCriteriaFunc(f)
	}
	return nil, nil
}

func (m *MockReportRepository) Update(r types.Report) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(r)
	}
	return nil
}

func (m *MockReportRepository) Delete(r types.Report) error {
	if m.DeleteFunc != nil {
		return m.Delete(r)
	}
	return nil
}
