package tasks

import (
	"errors"
	"flag"
	"fmt"
	"intel/isecl/lib/common/setup"
	"intel/isecl/lib/common/validation"
	"intel/isecl/tdservice/repository"
	"intel/isecl/tdservice/types"
	"io"

	"golang.org/x/crypto/bcrypt"

	log "github.com/sirupsen/logrus"
)

const (

ADMIN_ROLE="ADMINISTRATOR"

)
type Admin struct {
	Flags           []string
	DatabaseFactory func() (repository.TDSDatabase, error)
	ConsoleWriter   io.Writer
}

func (a Admin) Run(c setup.Context) error {
	fmt.Fprintln(a.ConsoleWriter, "Running admin setup...")
	envUser, _ := c.GetenvString("TDS_ADMIN_USERNAME", "Username for admin authentication")
	envPass, _ := c.GetenvSecret("TDS_ADMIN_PASSWORD", "Password for admin authentication")
	fs := flag.NewFlagSet("admin", flag.ContinueOnError)
	force := fs.Bool("force", false, "force creation")
	username := fs.String("admin-user", envUser, "Username for admin authentication")
	password := fs.String("admin-pass", envPass, "Password for admin authentication")
         
	err := fs.Parse(a.Flags)
	if err != nil {
		return err
	}
	if a.Validate(c) == nil && *force == false {
		return nil
	}
	if *username == "" {
		return errors.New("admin setup: Username cannot be empty")
	}
	if *password == "" {
		return errors.New("admin setup: Password cannot be empty")
	}
	valid_err := validation.ValidateAccount(*username, *password)
	if valid_err != nil {
		return valid_err
	}
	db, err := a.DatabaseFactory()
	if err != nil {
		log.WithError(err).Error("failed to open database")
		return err
	}
	defer db.Close()
	hash, _ := bcrypt.GenerateFromPassword([]byte(*password), bcrypt.DefaultCost)
	if *force {
		// if force, delete any users with the name
		db.UserRepository().Delete(types.User{Name: *username})
	}
        admin_role := types.Role{Name: ADMIN_ROLE}
	db.UserRepository().Create(types.User{Name:  *username, 
                                              PasswordHash: hash,
                                              Roles: []types.Role{admin_role},
                                             })

	return nil
}

func (a Admin) Validate(c setup.Context) error {
	// check if Users is not empty, eventually will be check if there is an admin user via a permission column
	// but as it stands only type of user is an admin user
	db, err := a.DatabaseFactory()
	if err != nil {
		log.WithError(err).Error("admin setup: failed to open database")
		return err
	}
	defer db.Close()
	_, err = db.UserRepository().Retrieve(types.User{}) // passing in a Zero struct finds First record
	return err
}
