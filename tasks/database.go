package tasks

import (
	"errors"
	"flag"
	"intel/isecl/lib/common/setup"
	"intel/isecl/threat-detection-service/config"
)

type Database struct {
	Flags  []string
	Config *config.Configuration
}

func (db Database) Run(c setup.Context) error {
	envHost, _ := c.GetenvString("TDS_DB_HOSTNAME", "Database Hostname")
	envPort, _ := c.GetenvInt("TDS_DB_PORT", "Database Port")
	envUser, _ := c.GetenvString("TDS_DB_USERNAME", "Database Username")
	envPass, _ := c.GetenvSecret("TDS_DB_PASSWORD", "Database Password")
	envDB, _ := c.GetenvString("TDS_DB_NAME", "Database Name")

	fs := flag.NewFlagSet("database", flag.ContinueOnError)
	fs.StringVar(&db.Config.Postgres.Hostname, "db-host", envHost, "Database Hostname")
	fs.IntVar(&db.Config.Postgres.Port, "db-port", envPort, "Database Port")
	fs.StringVar(&db.Config.Postgres.Username, "db-username", envUser, "Database Username")
	fs.StringVar(&db.Config.Postgres.Password, "db-password", envPass, "Database Password")
	fs.StringVar(&db.Config.Postgres.DBName, "db-name", envDB, "Database Name")
	err := fs.Parse(db.Flags)
	if err != nil {
		return err
	}
	return db.Config.Save()
}

func (db Database) Validate(c setup.Context) error {
	if db.Config.Postgres.Hostname == "" {
		return errors.New("Database: Hostname is not set")
	}
	if db.Config.Postgres.Port == 0 {
		return errors.New("Database: Port is not set")
	}
	if db.Config.Postgres.Username == "" {
		return errors.New("Database: Username is not set")
	}
	if db.Config.Postgres.Password == "" {
		return errors.New("Database: Password is not set")
	}
	if db.Config.Postgres.DBName == "" {
		return errors.New("Database: Schema is not set")
	}
	return nil
}
