package resource

import (
	"bytes"
	"encoding/json"
	"intel/isecl/threat-detection-service/repository"
	"intel/isecl/threat-detection-service/repository/mock"
	"intel/isecl/threat-detection-service/types"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreateReport(t *testing.T) {
	assert := assert.New(t)
	db := new(mock.MockDatabase)
	var reportCreated bool

	r := setupRouter(db)
	recorder := httptest.NewRecorder()

	report := types.Report{
		HostID: "2",
		Detection: types.Detection{
			Description:      "description",
			PID:              1,
			TID:              2,
			ProcessName:      "process.name",
			ProcessImagePath: "/usr/bin/process.name",
			ProcessCmdLine:   "/usr/bin/process.name -f",
			Timestamp:        1234,
			Severity:         10,
			CVEIDs:           []string{"CVE-1", "CVE-2"},
			ThreatClass:      "spectre variant 1",
		},
	}
	db.MockReportRepository.CreateFunc = func(r types.Report) (*types.Report, error) {
		reportCreated = true
		r.ID = "1"
		r.Host = types.Host{ID: "2"}
		assert.Equal(report.Detection, r.Detection)
		return &r, nil
	}
	reportJSON, _ := json.Marshal(report)

	req := httptest.NewRequest("POST", "/tds/reports", bytes.NewBuffer(reportJSON))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(recorder, req)
	assert.Equal(http.StatusCreated, recorder.Code)
	assert.Equal("application/json", recorder.Header().Get("Content-Type"))
	assert.True(reportCreated)

	var out types.Report
	json.Unmarshal(recorder.Body.Bytes(), &out)
	report.ID = "1"
	assert.Equal(report, out)
}

func TestRetrieveReport(t *testing.T) {
	assert := assert.New(t)
	db := new(mock.MockDatabase)
	report := types.Report{
		ID:     "1",
		HostID: "2",
		Detection: types.Detection{
			Description:      "description",
			PID:              1,
			TID:              2,
			ProcessName:      "process.name",
			ProcessImagePath: "/usr/bin/process.name",
			ProcessCmdLine:   "/usr/bin/process.name -f",
			Timestamp:        1234,
			Severity:         10,
			CVEIDs:           []string{"CVE-1", "CVE-2"},
			ThreatClass:      "spectre variant 1",
		},
	}
	var reportRetrieved bool
	db.MockReportRepository.RetrieveFunc = func(h types.Report) (*types.Report, error) {
		reportRetrieved = true
		assert.Equal("1", h.ID)
		return &report, nil
	}
	r := setupRouter(db)
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/tds/reports/1", nil)
	r.ServeHTTP(recorder, req)
	assert.Equal(http.StatusOK, recorder.Code)
	assert.Equal("application/json", recorder.Header().Get("Content-Type"))
	assert.True(reportRetrieved)
	assert.NotEmpty(recorder.Body.String())
}

func TestQueryReports(t *testing.T) {
	assert := assert.New(t)
	db := new(mock.MockDatabase)
	report := types.Report{
		HostID: "2",
		Detection: types.Detection{
			Description:      "description",
			PID:              1,
			TID:              2,
			ProcessName:      "process.name",
			ProcessImagePath: "/usr/bin/process.name",
			ProcessCmdLine:   "/usr/bin/process.name -f",
			Timestamp:        1234,
			Severity:         10,
			CVEIDs:           []string{"CVE-1", "CVE-2"},
			ThreatClass:      "spectre variant 1",
		},
	}

	db.MockReportRepository.RetrieveByFilterCriteriaFunc = func(f repository.ReportFilter) ([]types.Report, error) {
		exFrom, _ := time.Parse(time.RFC3339, "2002-10-02T15:00:00Z")
		exTo, _ := time.Parse(time.RFC3339, "2020-10-02T15:00:00Z")
		assert.Equal("10.1.2.3", f.Hostname)
		assert.Equal(exFrom, f.From)
		assert.Equal(exTo, f.To)
		return []types.Report{report}, nil
	}
	r := setupRouter(db)
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/tds/reports?hostname=10.1.2.3&from=2002-10-02T15:00:00Z&to=2020-10-02T15:00:00Z", nil)
	r.ServeHTTP(recorder, req)
	assert.Equal(http.StatusOK, recorder.Code)
}
