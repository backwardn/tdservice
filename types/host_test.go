package types

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncoding(t *testing.T) {
	h := Host{
		ID: "1234",
		HostInfo: HostInfo{
			Version: "v1.0",
			Build:   "12313131",
			OS:      "Linux",
		},
		Status: "offline",
	}
	j, _ := json.Marshal(h)
	t.Log(string(j))
	assert.Equal(t, `{"id":"1234","hostname":"","version":"v1.0","build":"12313131","os":"Linux","status":"offline"}`, string(j))
}

func TestDecoding(t *testing.T) {
	var h Host
	json.Unmarshal([]byte(`{"id":"1234","hostname":"","version":"v1.0","build":"12313131","os":"Linux","status":"offline"}`), &h)
	assert.Equal(t, "1234", h.ID)
	assert.Equal(t, "v1.0", h.Version)
}