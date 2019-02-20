package config

import (
	"intel/isecl/tdservice/constants"
	"os"
	"path"
	"sync"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

// should move this into lib common, as its duplicated across TDS and TDA

// Configuration is the global configuration struct that is marshalled/unmarshaled to a persisted yaml file
// Probably should embed a config generic struct
type Configuration struct {
	configFile string
	Port       int
	Postgres   struct {
		DBName   string
		Username string
		Password string
		Hostname string
		Port     int
		SSLMode  bool
	}
	LogLevel log.Level
}

var mu sync.Mutex

var Global = &Configuration{}

func (c *Configuration) Save() error {
	if c.configFile == "" {
		return nil
	}
	file, err := os.OpenFile(c.configFile, os.O_RDWR, 0)
	if err != nil {
		// we have an error
		if os.IsNotExist(err) {
			// error is that the config doesnt yet exist, create it
			file, err = os.Create(c.configFile)
			if err != nil {
				return err
			}
		} else {
			// someother I/O related error
			return err
		}
	}
	defer file.Close()
	return yaml.NewEncoder(file).Encode(c)
}

func Load(path string) (*Configuration, error) {
	var c Configuration
	file, err := os.Open(path)
	if err == nil {
		defer file.Close()
		yaml.NewDecoder(file).Decode(&c)
		c.configFile = path
		return &c, nil
	}
	return nil, err
}

func init() {
	// load from config
	g, err := Load(path.Join(constants.ConfigDir, constants.ConfigFile))
	if err == nil {
		Global = g
	}
}
