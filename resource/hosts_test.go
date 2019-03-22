package resource

import (
	"bytes"
	"encoding/json"
	"intel/isecl/tdservice/repository/mock"
	"intel/isecl/tdservice/types"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jinzhu/gorm"

	"github.com/stretchr/testify/assert"
)

func TestCreateHost(t *testing.T) {
	assert := assert.New(t)
	db := new(mock.MockDatabase)
	var hostCreated bool
	db.MockHostRepository.CreateFunc = func(h types.Host) (*types.Host, error) {
		hostCreated = true
		h.ID = "12345"
		assert.Equal("host.intel.com", h.Hostname)
		assert.Equal("v1.0", h.Version)
		assert.Equal("1234", h.Build)
		assert.Equal("linux", h.OS)
		return &h, nil
	}
	r := setupRouter(db)
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/tds/hosts", bytes.NewBufferString(`{"hostname": "host.intel.com", "version": "v1.0", "build": "1234", "os":"linux"}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(recorder, req)
	assert.Equal(http.StatusCreated, recorder.Code)
	assert.Equal("application/json", recorder.Header().Get("Content-Type"))
	assert.True(hostCreated)
}

func TestRetrieveHost(t *testing.T) {
	assert := assert.New(t)
	db := new(mock.MockDatabase)
	var hostRetrieved bool
	db.MockHostRepository.RetrieveFunc = func(h types.Host) (*types.Host, error) {
		hostRetrieved = true
		assert.Equal("12345", h.ID)
		return &types.Host{
			ID: "12345",
			HostInfo: types.HostInfo{
				Version: "v1.0",
				Build:   "1234",
				OS:      "linux",
			},
			Status: "Reserve for future implementation",
		}, nil
	}
	r := setupRouter(db)
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/tds/hosts/12345", nil)
	r.ServeHTTP(recorder, req)
	assert.Equal(http.StatusOK, recorder.Code)
	assert.Equal("application/json", recorder.Header().Get("Content-Type"))
	assert.True(hostRetrieved)
	assert.NotEmpty(recorder.Body.String())
}

func TestRetrieveHost404(t *testing.T) {
	assert := assert.New(t)
	db := new(mock.MockDatabase)
	db.MockHostRepository.RetrieveFunc = func(h types.Host) (*types.Host, error) {
		return nil, gorm.ErrRecordNotFound
	}
	r := setupRouter(db)
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/tds/hosts/12345", nil)
	r.ServeHTTP(recorder, req)
	assert.Equal(http.StatusNotFound, recorder.Code)
}

func TestDeleteHost(t *testing.T) {
	assert := assert.New(t)
	db := new(mock.MockDatabase)
	id := "12345"
	var deleted bool
	db.MockHostRepository.DeleteFunc = func(h types.Host) error {
		assert.Equal(id, h.ID)
		deleted = true
		return nil
	}
	r := setupRouter(db)
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest("DELETE", "/tds/hosts/12345", nil)
	r.ServeHTTP(recorder, req)
	assert.Equal(http.StatusNoContent, recorder.Code)
	assert.True(deleted)
}

func TestDeleteHost404(t *testing.T) {
	assert := assert.New(t)
	db := new(mock.MockDatabase)
	db.MockHostRepository.DeleteFunc = func(h types.Host) error {
		return gorm.ErrRecordNotFound
	}
	r := setupRouter(db)
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest("DELETE", "/tds/hosts/12345", nil)
	r.ServeHTTP(recorder, req)
	assert.Equal(http.StatusNotFound, recorder.Code)
}

func TestQueryHosts(t *testing.T) {
	assert := assert.New(t)
	db := new(mock.MockDatabase)
	filter := types.Host{
		HostInfo: types.HostInfo{
			Hostname: "10.1.2.3",
			Version:  "1.0",
			Build:    "1234",
			OS:       "linux",
		},
		Status: "Reserve for future implementation",
	}
	db.MockHostRepository.RetrieveAllFunc = func(h types.Host) ([]types.Host, error) {
		assert.Equal(filter, h)
		h.ID = "12345"
		return []types.Host{h}, nil
	}
	r := setupRouter(db)
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/tds/hosts?hostname=10.1.2.3&version=1.0&build=1234&os=linux", nil)
	r.ServeHTTP(recorder, req)
	assert.Equal(http.StatusOK, recorder.Code)
	// setup expected
	filter.ID = "12345"
	expected := []types.Host{filter}
	var actual []types.Host
	json.Unmarshal(recorder.Body.Bytes(), &actual)
	assert.Equal(expected, actual)
}
