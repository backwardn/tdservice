package tasks

import (
	"flag"
	"intel/isecl/lib/common/setup"
	"intel/isecl/threat-detection-service/config"
)

type Server struct {
	Flags  []string
	Config *config.Configuration
}

func (s Server) Run(c setup.Context) error {
	defaultPort, err := c.GetenvInt("TDS_PORT", "threat detection service http port")
	if err != nil {
		defaultPort = 8443
	}
	fs := flag.NewFlagSet("server", flag.ContinueOnError)
	fs.IntVar(&s.Config.Port, "port", defaultPort, "threat detection service http port")
	err = fs.Parse(s.Flags)
	if err != nil {
		return err
	}
	return s.Config.Save()
}

func (s Server) Validate(c setup.Context) error {
	return nil
}
