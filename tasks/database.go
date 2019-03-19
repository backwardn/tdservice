package tasks

import (
	"errors"
	"flag"
	"fmt"
	"intel/isecl/lib/common/setup"
	"intel/isecl/lib/common/validation"
	"intel/isecl/tdservice/config"
	"io"
)

type Database struct {
	Flags         []string
	Config        *config.Configuration
	ConsoleWriter io.Writer
}

func (db Database) Run(c setup.Context) error {
	fmt.Fprintln(db.ConsoleWriter, "Running database setup...")
	envHost, _ := c.GetenvString("TDS_DB_HOSTNAME", "Database Hostname")
	envPort, _ := c.GetenvInt("TDS_DB_PORT", "Database Port")
	envUser, _ := c.GetenvString("TDS_DB_USERNAME", "Database Username")
	envPass, _ := c.GetenvSecret("TDS_DB_PASSWORD", "Database Password")
	envDB, _ := c.GetenvString("TDS_DB_NAME", "Database Name")

	fs := flag.NewFlagSet("database", flag.ContinueOnError)
	fs.StringVar(&db.Config.Postgres.Hostname, "db-host", envHost, "Database Hostname")
	fs.IntVar(&db.Config.Postgres.Port, "db-port", envPort, "Database Port")
	fs.StringVar(&db.Config.Postgres.Username, "db-user", envUser, "Database Username")
	fs.StringVar(&db.Config.Postgres.Password, "db-pass", envPass, "Database Password")
	fs.StringVar(&db.Config.Postgres.DBName, "db-name", envDB, "Database Name")
	err := fs.Parse(db.Flags)
	if err != nil {
		return err
	}

	var valid_err error

	valid_err = validation.ValidateHostname(db.Config.Postgres.Hostname)
	if valid_err != nil {
		return fmt.Errorf("Validation fail: %s", valid_err.Error())
	}
	valid_err = validation.ValidateAccount(db.Config.Postgres.Username, db.Config.Postgres.Password)
	if valid_err != nil {
		return fmt.Errorf("Validation fail: %s", valid_err.Error())
	}
	valid_err = validation.ValidateIdentifier(db.Config.Postgres.DBName)
	if valid_err != nil {
		return fmt.Errorf("Validation fail: %s", valid_err.Error())
	}

	return db.Config.Save()
}

func (db Database) Validate(c setup.Context) error {
	if db.Config.Postgres.Hostname == "" {
		return errors.New("database setup: Hostname is not set")
	}
	if db.Config.Postgres.Port == 0 {
		return errors.New("database setup: Port is not set")
	}
	if db.Config.Postgres.Username == "" {
		return errors.New("database setup: Username is not set")
	}
	if db.Config.Postgres.Password == "" {
		return errors.New("database setup: Password is not set")
	}
	if db.Config.Postgres.DBName == "" {
		return errors.New("database: Schema is not set")
	}
	return nil
}
