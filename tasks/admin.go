package tasks

import (
	"errors"
	"flag"
	"fmt"
	"intel/isecl/lib/common/setup"
	"intel/isecl/lib/common/validation"
        consts "intel/isecl/tdservice/constants"
	"intel/isecl/tdservice/repository"
	"intel/isecl/tdservice/types"
	"io"

	"golang.org/x/crypto/bcrypt"

	log "github.com/sirupsen/logrus"
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
        envRegHostUser, _ := c.GetenvString("TDS_REG_HOST_USERNAME", "Username for register-host-user authentication")
        envRegHostPass, _ := c.GetenvSecret("TDS_REG_HOST_PASSWORD", "Password for register-host-user authentication")
	fs := flag.NewFlagSet("admin", flag.ContinueOnError)
	force := fs.Bool("force", false, "force creation")
	username := fs.String("admin-user", envUser, "Username for admin authentication")
	password := fs.String("admin-pass", envPass, "Password for admin authentication")
        regHostUsername := fs.String("reg-host-user", envRegHostUser,  "Username for RegisterHostUser authentication")
        regHostUserPassword := fs.String("reg-host-password", envRegHostPass,  "Password for RegisterHostUser authentication")
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
        if *regHostUsername == "" {
                return errors.New("admin setup: Username for RegisterHostUser cannot be empty")
        }
        if *regHostUserPassword == "" {
                return errors.New("admin setup: Password for for RegisterHostUser cannot be empty")
        }
	valid_err := validation.ValidateAccount(*username, *password)
	if valid_err != nil {
		return valid_err
	}
	valid_err = validation.ValidateAccount(*regHostUsername, *regHostUserPassword)
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
                db.UserRepository().Delete(types.User{Name: *regHostUsername})
	}

    adminRole := types.Role{Name: consts.AdminGroupName}
	db.UserRepository().Create(types.User{Name:  *username, 
                                              PasswordHash: hash,
                                              Roles: []types.Role{adminRole},
                                             })

	registerHostRole := types.Role{Name: consts.RegisterHostGroupName}
	hash, _ = bcrypt.GenerateFromPassword([]byte(*regHostUserPassword), bcrypt.DefaultCost)
	db.UserRepository().Create(types.User{Name: *regHostUsername,
					      PasswordHash: hash,
					      Roles: []types.Role{registerHostRole},
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
