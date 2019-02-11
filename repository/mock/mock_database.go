package mock

import (
	"intel/isecl/threat-detection-service/repository"
)

type MockDatabase struct {
	MockHostRepository MockHostRepository
}

func (m *MockDatabase) Migrate() error {
	return nil
}

func (m *MockDatabase) HostRepository() repository.HostRepository {
	return nil
}

func (pd *MockDatabase) ReportRepository() repository.ReportRepository {
	return nil
}
