// +build integration

package postgres

import (
	"intel/isecl/tdservice/repository"
	"intel/isecl/tdservice/types"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestReportCreateRetrieve(t *testing.T) {
	db := dialDatabase(t)
	assert := assert.New(t)
	// If you somehow run this on production, god bless your poor soul
	db.DB.DropTableIfExists(types.Host{}, types.Report{})
	db.Migrate()

	// create a host to associate it with
	created, err := createHost("10.0.0.1", db.HostRepository())

	report := types.Report{}
	report.Detection.PID = 1
	report.HostID = created.ID

	createdReport, err := db.ReportRepository().Create(report)
	assert.NoError(err)
	assert.NotEmpty(createdReport.ID)
	assert.NotEmpty(createdReport.HostID)

	retrieved, err := db.ReportRepository().Retrieve(types.Report{ID: createdReport.ID})
	assert.Equal(createdReport.ID, retrieved.ID)
	assert.NotEmpty(retrieved.Host.ID)
	assert.Equal(created.ID, retrieved.Host.ID)

	// create another report

	report2 := report
	report2.Detection.PID = 2

	_, err = db.ReportRepository().Create(report2)

	all, err := db.ReportRepository().RetrieveAll(types.Report{})
	assert.NoError(err)
	assert.Len(all, 2)
}

func TestReportRetrieveByFilter(t *testing.T) {
	db := dialDatabase(t)
	assert := assert.New(t)
	// If you somehow run this on production, god bless your poor soul
	db.DB.DropTableIfExists(types.Host{}, types.Report{})
	db.Migrate()

	// create a host to associate it with
	created, _ := createHost("10.0.0.1", db.HostRepository())

	report := types.Report{}
	report.Detection.PID = 1
	now := time.Now()
	report.Detection.Timestamp = now.Unix()
	report.HostID = created.ID

	report2 := report
	report2.Detection.Timestamp = now.Add(-time.Hour).Unix()

	report3 := report
	report3.Detection.Timestamp = now.Add(time.Hour).Unix()
	db.ReportRepository().Create(report)
	db.ReportRepository().Create(report2)
	db.ReportRepository().Create(report3)

	filter := repository.ReportFilter{
		From: now.Add(-2 * time.Hour),
	}
	all, err := db.ReportRepository().RetrieveByFilterCriteria(filter)
	assert.NoError(err)
	assert.Len(all, 3)

	filter = repository.ReportFilter{
		Hostname: "10.0.0.1",
		From:     now.Add(-2 * time.Hour),
		To:       now.Add(time.Minute),
	}
	all, err = db.ReportRepository().RetrieveByFilterCriteria(filter)
	assert.NoError(err)
	assert.Len(all, 2)
}