package repository

type TDSDatabase interface {
	Migrate() error
	HostRepository() interface{}
	ReportRepository() interface{}
}
