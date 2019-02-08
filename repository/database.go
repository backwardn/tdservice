package repository

type TDSDatabase interface {
	Migrate() error
	HostRepository() HostRepository
	ReportRepository() ReportRepository
}
