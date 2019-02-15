package mock

import (
	"intel/isecl/threat-detection-service/repository"
)

type MockDatabase struct {
	MockHostRepository   MockHostRepository
	MockReportRepository MockReportRepository
	MockUserRepository   MockUserRepository
}

func (m *MockDatabase) Migrate() error {
	return nil
}

func (m *MockDatabase) HostRepository() repository.HostRepository {
	return &m.MockHostRepository
}

func (m *MockDatabase) ReportRepository() repository.ReportRepository {
	return &m.MockReportRepository
}

func (m *MockDatabase) UserRepository() repository.UserRepository {
	return &m.MockUserRepository
}

func (m *MockDatabase) Close() {

}
