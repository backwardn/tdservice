package tasks

import (
	"intel/isecl/lib/common/setup"
	"intel/isecl/tdservice/config"
	"os"

	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServerSetup(t *testing.T) {
	c := config.Configuration{}
	s := Server{
		Flags:  []string{"-port=1337"},
		Config: &c,
	}
	ctx := setup.Context{}
	err := s.Run(ctx)
	assert.NoError(t, err)
	assert.Equal(t, 1337, c.Port)
}

func TestServerSetupEnv(t *testing.T) {
	os.Setenv("TDS_PORT", "1337")
	c := config.Configuration{}
	s := Server{
		Flags:  nil,
		Config: &c,
	}
	ctx := setup.Context{}
	err := s.Run(ctx)
	assert.NoError(t, err)
	assert.Equal(t, 1337, c.Port)
}
