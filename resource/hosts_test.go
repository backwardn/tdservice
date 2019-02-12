package resource

import (
	"bytes"
	"intel/isecl/threat-detection-service/repository/mock"
	"intel/isecl/threat-detection-service/types"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jinzhu/gorm"

	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	assert := assert.New(t)
	db := new(mock.MockDatabase)
	var hostCreated bool
	db.MockHostRepository.CreateFunc = func(h types.Host) error {
		hostCreated = true
		assert.Equal("host.intel.com", h.Hostname)
		assert.Equal("v1.0", h.Version)
		assert.Equal("1234", h.Build)
		assert.Equal("linux", h.OS)
		return nil
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

func TestRetrieve(t *testing.T) {
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
			Status: "online",
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

func TestRetrieve404(t *testing.T) {
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
