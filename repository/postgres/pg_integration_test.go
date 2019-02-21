package postgres

import (
	"fmt"
	"intel/isecl/tdservice/types"
	"testing"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/stretchr/testify/assert"
)

func dialDatabase(t *testing.T) *PostgresDatabase {
	g, err := gorm.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=%s",
		"localhost", 5432, "runner", "pgdb", "test", "disable"))
	if err != nil {
		t.SkipNow()
	}
	return &PostgresDatabase{DB: g}
}

func TestHostCreate(t *testing.T) {
	db := dialDatabase(t)
	assert := assert.New(t)
	// If you somehow run this on production, god bless your poor soul
	db.DB.DropTableIfExists(types.Host{})
	db.Migrate()
	host := types.Host{}
	host.Hostname = "10.0.0.1"
	host.OS = "linux"
	host.Status = "online"
	host.Version = "1.0"
	host.Build = "1234"
	created, err := db.HostRepository().Create(host)
	assert.NotEmpty(created.ID)
	assert.NoError(err)
}

func TestHostCreateDuplicate(t *testing.T) {
	db := dialDatabase(t)
	assert := assert.New(t)
	// If you somehow run this on production, god bless your poor soul
	db.DB.DropTableIfExists(types.Host{})
	db.Migrate()
	host := types.Host{}
	host.Hostname = "10.0.0.1"
	host.OS = "linux"
	host.Status = "online"
	host.Version = "1.0"
	host.Build = "1234"
	created, err := db.HostRepository().Create(host)
	assert.NotEmpty(created.ID)
	assert.NoError(err)

	_, err = db.HostRepository().Create(host)
	assert.Error(err)
}

func TestHostRetrieve(t *testing.T) {
	db := dialDatabase(t)
	assert := assert.New(t)
	// If you somehow run this on production, god bless your poor soul
	db.DB.DropTableIfExists(types.Host{})
	db.Migrate()
	host := types.Host{}
	host.Hostname = "10.0.0.1"
	host.OS = "linux"
	host.Status = "online"
	host.Version = "1.0"
	host.Build = "1234"
	created, err := db.HostRepository().Create(host)
	assert.NotEmpty(created.ID)
	assert.NoError(err)

	// fetch it
	fetched, err := db.HostRepository().Retrieve(types.Host{ID: created.ID})
	assert.NoError(err)
	assert.Equal(created.ID, fetched.ID)
}

func TestHostRetrieveAll(t *testing.T) {
	db := dialDatabase(t)
	assert := assert.New(t)
	// If you somehow run this on production, god bless your poor soul
	db.DB.DropTableIfExists(types.Host{})
	db.Migrate()
	host := types.Host{}
	host.Hostname = "10.0.0.1"
	host.OS = "linux"
	host.Status = "online"
	host.Version = "1.0"
	host.Build = "1234"
	created, err := db.HostRepository().Create(host)
	assert.NotEmpty(created.ID)
	assert.NoError(err)

	host.Hostname = "10.0.0.2"
	created2, err := db.HostRepository().Create(host)
	assert.NotEmpty(created2.ID)
	assert.NoError(err)

	// query all

	all, err := db.HostRepository().RetrieveAll(types.Host{})
	assert.NoError(err)
	assert.Len(all, 2)

	filter := types.Host{}
	filter.Hostname = "10.0.0.1"
	all, err = db.HostRepository().RetrieveAll(filter)
	assert.NoError(err)
	assert.Len(all, 1)
}

func TestHostUpdate(t *testing.T) {
	db := dialDatabase(t)
	assert := assert.New(t)
	// If you somehow run this on production, god bless your poor soul
	db.DB.DropTableIfExists(types.Host{})
	db.Migrate()
	host := types.Host{}
	host.Hostname = "10.0.0.1"
	host.OS = "linux"
	host.Status = "online"
	host.Version = "1.0"
	host.Build = "1234"
	created, err := db.HostRepository().Create(host)
	assert.NotEmpty(created.ID)
	assert.NoError(err)

	created.Hostname = "10.0.0.2"
	err = db.HostRepository().Update(*created)
	assert.NoError(err)
}

func TestHostDelete(t *testing.T) {
	db := dialDatabase(t)
	assert := assert.New(t)
	// If you somehow run this on production, god bless your poor soul
	db.DB.DropTableIfExists(types.Host{})
	db.Migrate()
	host := types.Host{}
	host.Hostname = "10.0.0.1"
	host.OS = "linux"
	host.Status = "online"
	host.Version = "1.0"
	host.Build = "1234"
	created, err := db.HostRepository().Create(host)
	assert.NotEmpty(created.ID)
	assert.NoError(err)

	err = db.HostRepository().Delete(*created)
	assert.NoError(err)

	all, err := db.HostRepository().RetrieveAll(types.Host{})
	assert.NoError(err)
	assert.Len(all, 0)
}
