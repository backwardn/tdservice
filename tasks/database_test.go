package tasks

import (
	"intel/isecl/lib/common/setup"
	"intel/isecl/threat-detection-service/config"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDatabaseSetup(t *testing.T) {
	assert := assert.New(t)
	c := config.Configuration{}
	s := Database{
		Flags:  []string{"-db-host=hostname", "-db-port=5432", "-db-user=user", "-db-pass=password", "-db-name=tds_db"},
		Config: &c,
	}
	ctx := setup.Context{}
	err := s.Run(ctx)
	assert.NoError(err)
	assert.Equal("hostname", c.Postgres.Hostname)
	assert.Equal(5432, c.Postgres.Port)
	assert.Equal("user", c.Postgres.Username)
	assert.Equal("password", c.Postgres.Password)
	assert.Equal("tds_db", c.Postgres.DBName)
}

func TestDatabaseSetupEnv(t *testing.T) {
	assert := assert.New(t)
	os.Setenv("TDS_DB_HOSTNAME", "hostname")
	os.Setenv("TDS_DB_PORT", "5432")
	os.Setenv("TDS_DB_USERNAME", "user")
	os.Setenv("TDS_DB_PASSWORD", "password")
	os.Setenv("TDS_DB_NAME", "tds_db")
	c := config.Configuration{}
	s := Database{
		Flags:  nil,
		Config: &c,
	}
	ctx := setup.Context{}
	err := s.Run(ctx)
	assert.NoError(err)
	assert.Equal("hostname", c.Postgres.Hostname)
	assert.Equal(5432, c.Postgres.Port)
	assert.Equal("user", c.Postgres.Username)
	assert.Equal("password", c.Postgres.Password)
	assert.Equal("tds_db", c.Postgres.DBName)
}
