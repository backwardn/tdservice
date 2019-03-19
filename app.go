package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"intel/isecl/lib/common/setup"
	"intel/isecl/lib/common/validation"
	"intel/isecl/tdservice/config"
	"intel/isecl/tdservice/constants"
	"intel/isecl/tdservice/middleware"
	"intel/isecl/tdservice/repository"
	"intel/isecl/tdservice/repository/postgres"
	"intel/isecl/tdservice/resource"
	"intel/isecl/tdservice/tasks"
	"intel/isecl/tdservice/version"
	"io"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"strings"
	"syscall"
	"time"

	stdlog "log"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	// Import driver for GORM
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type App struct {
	ConfigDir      string
	RunDir         string
	DataDir        string
	ExecutablePath string
	Config         *config.Configuration
	ConsoleWriter  io.Writer
	LogWriter      io.Writer
	HTTPLogWriter  io.Writer
}

func (a *App) printUsage() {
	fmt.Fprintln(a.consoleWriter(), "Usage:")
	fmt.Fprintln(a.consoleWriter(), "")
	fmt.Fprintln(a.consoleWriter(), "    tdservice <command> [arguments]")
	fmt.Fprintln(a.consoleWriter(), "")
	fmt.Fprintln(a.consoleWriter(), "Avaliable Commands:")
	fmt.Fprintln(a.consoleWriter(), "    help|-h|-help    Show this help message")
	fmt.Fprintln(a.consoleWriter(), "    setup [task]     Run setup task")
	fmt.Fprintln(a.consoleWriter(), "    start            Start tdservice")
	fmt.Fprintln(a.consoleWriter(), "    status           Show the status of tdservice")
	fmt.Fprintln(a.consoleWriter(), "    stop             Stop tdservice")
	fmt.Fprintln(a.consoleWriter(), "    uninstall        Uninstall tdservice")
	fmt.Fprintln(a.consoleWriter(), "    version          Show the version of tdservice")
	fmt.Fprintln(a.consoleWriter(), "")
	fmt.Fprintln(a.consoleWriter(), "Avaliable Tasks for setup:")
	fmt.Fprintln(a.consoleWriter(), "    tdservice setup database [-force] [--arguments=<argument_value>]")
	fmt.Fprintln(a.consoleWriter(), "        - Avaliable arguments are:")
	fmt.Fprintln(a.consoleWriter(), "            - db-host        alternatively, set environment variable TDS_DB_HOSTNAME")
	fmt.Fprintln(a.consoleWriter(), "            - db-port        alternatively, set environment variable TDS_DB_PORT")
	fmt.Fprintln(a.consoleWriter(), "            - db-username    alternatively, set environment variable TDS_DB_USERNAME")
	fmt.Fprintln(a.consoleWriter(), "            - db-password    alternatively, set environment variable TDS_DB_PASSWORD")
	fmt.Fprintln(a.consoleWriter(), "            - db-name        alternatively, set environment variable TDS_DB_NAME")
	fmt.Fprintln(a.consoleWriter(), "    tdservice setup server [--port=<port>]")
	fmt.Fprintln(a.consoleWriter(), "        - Setup http server on <port>")
	fmt.Fprintln(a.consoleWriter(), "        - Environment variable TDS_PORT=<port> can be set alternatively")
	fmt.Fprintln(a.consoleWriter(), "    tdservice setup tls [--force] [--host_names=<host_names>]")
	fmt.Fprintln(a.consoleWriter(), "        - Use the key and certificate provided in /etc/threat-detection if files exist")
	fmt.Fprintln(a.consoleWriter(), "        - Otherwise create its own self-signed TLS keypair in /etc/tdservice for quality of life")
	fmt.Fprintln(a.consoleWriter(), "        - Option [--force] overwrites any existing files, and always generate self-signed keypair")
	fmt.Fprintln(a.consoleWriter(), "        - Argument <host_names> is a list of host names used by local machine, seperated by comma")
	fmt.Fprintln(a.consoleWriter(), "        - Environment variable TDA_TLS_HOST_NAMES=<host_names> can be set alternatively")
	fmt.Fprintln(a.consoleWriter(), "    tdservice setup admin [--admin-user=<username>] [--admin-pass=<password>]")
	fmt.Fprintln(a.consoleWriter(), "        - Environment variable TDS_ADMIN_USERNAME=<username> can be set alternatively")
	fmt.Fprintln(a.consoleWriter(), "        - Environment variable TDS_ADMIN_PASSWORD=<password> can be set alternatively")
	fmt.Fprintln(a.consoleWriter(), "")
}

func (a *App) consoleWriter() io.Writer {
	if a.ConsoleWriter != nil {
		return a.ConsoleWriter
	}
	return os.Stdout
}

func (a *App) logWriter() io.Writer {
	if a.LogWriter != nil {
		return a.LogWriter
	}
	return os.Stderr
}

func (a *App) httpLogWriter() io.Writer {
	if a.HTTPLogWriter != nil {
		return a.HTTPLogWriter
	}
	return os.Stderr
}

func (a *App) configuration() *config.Configuration {
	if a.Config != nil {
		return a.Config
	}
	return config.Global()
}

func (a *App) executablePath() string {
	if a.ExecutablePath != "" {
		return a.ExecutablePath
	}
	exec, err := os.Executable()
	if err != nil {
		// if we can't find self-executable path, we're probably in a state that is panic() worthy
		panic(err)
	}
	return exec
}

func (a *App) configDir() string {
	if a.ConfigDir != "" {
		return a.ConfigDir
	}
	return constants.ConfigDir
}

func (a *App) dataDir() string {
	if a.DataDir != "" {
		return a.DataDir
	}
	return constants.DataDir
}

func (a *App) runDir() string {
	if a.RunDir != "" {
		return a.RunDir
	}
	return constants.RunDir
}

func (a *App) configureLogs() {
	log.SetOutput(io.MultiWriter(os.Stderr, a.logWriter()))
	log.SetLevel(a.configuration().LogLevel)

	// override golang logger
	w := log.StandardLogger().WriterLevel(a.configuration().LogLevel)
	stdlog.SetOutput(w)
}

func (a *App) Run(args []string) error {
	a.configureLogs()

	if len(args) < 2 {
		a.printUsage()
		os.Exit(1)
	}

	//bin := args[0]
	cmd := args[1]
	switch cmd {
	default:
		fmt.Println("Error: Unrecognized command: ", args[1])
		a.printUsage()
	case "run":
		return a.startServer()
	case "-help":
		fallthrough
	case "--h":
		fallthrough
	case "--help":
		fallthrough
	case "help":
		a.printUsage()
	case "start":
		return a.start()
	case "stop":
		return a.stop()
	case "status":
		return a.status()
	case "uninstall":
		var keepConfig bool
		flag.CommandLine.BoolVar(&keepConfig, "keep-config", false, "keep config when uninstalling")
		flag.CommandLine.Parse(args[2:])
		a.uninstall(keepConfig)
		os.Exit(0)
	case "version":
		fmt.Fprintf(a.consoleWriter(), "Threat Detection Service %s-%s\n", version.Version, version.GitHash)
	case "setup":

		if len(args) <= 2 {
			fmt.Fprintln(os.Stdout, "Available setup tasks:\n- database\n- admin\n- server\n- tls\n-----------------\n- [all]")
			os.Exit(1)
		}

		valid_err := validateSetupArgs(args[2], args[3:])
		if valid_err != nil {
			return valid_err
		}

		task := strings.ToLower(args[2])
		flags := args[3:]
		setupRunner := &setup.Runner{
			Tasks: []setup.Task{
				tasks.Database{
					Flags:         flags,
					Config:        a.configuration(),
					ConsoleWriter: os.Stdout,
				},
				tasks.Admin{
					Flags: flags,
					DatabaseFactory: func() (repository.TDSDatabase, error) {
						pg := &a.configuration().Postgres
						p, err := postgres.Open(pg.Hostname, pg.Port, pg.DBName, pg.Username, pg.Password, pg.SSLMode)
						if err != nil {
							log.WithError(err).Error("failed to open postgres connection for setup task")
							return nil, err
						}
						p.Migrate()
						return p, nil
					},
					ConsoleWriter: os.Stdout,
				},
				tasks.Server{
					Flags:         flags,
					Config:        a.configuration(),
					ConsoleWriter: os.Stdout,
				},
				tasks.TLS{
					Flags:         flags,
					TLSCertFile:   path.Join(a.configDir(), constants.TLSCertFile),
					TLSKeyFile:    path.Join(a.configDir(), constants.TLSKeyFile),
					ConsoleWriter: os.Stdout,
				},
			},
			AskInput: false,
		}
		var err error
		if task == "all" {
			err = setupRunner.RunTasks()
		} else {
			err = setupRunner.RunTasks(task)
		}
		if err != nil {
			log.WithError(err).Error("Error running setup")
			fmt.Println("Error running setup: ", err)
			return err
		}
	}
	return nil
}

func (a *App) startServer() error {
	c := a.configuration()

	// Open database
	tdsDB, err := postgres.Open(c.Postgres.Hostname, c.Postgres.Port, c.Postgres.DBName, c.Postgres.Username, c.Postgres.Password, c.Postgres.SSLMode)
	if err != nil {
		log.WithError(err).Error("failed to open Postgres database")
		return err
	}
	defer tdsDB.Close()
	log.Trace("Migrating Database")
	tdsDB.Migrate()

	// Create Router, set routes
	r := mux.NewRouter().PathPrefix("/tds").Subrouter()
	r.Use(middleware.NewBasicAuth(tdsDB.UserRepository()))
	func(setters ...func(*mux.Router, repository.TDSDatabase)) {
		for _, s := range setters {
			s(r, tdsDB)
		}
	}(resource.SetHosts, resource.SetReports, resource.SetVersion)

	// Setup signal handlers to gracefully handle termination
	stop := make(chan os.Signal)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	httpLog := stdlog.New(a.httpLogWriter(), "", 0)
	h := &http.Server{
		Addr:     fmt.Sprintf(":%d", c.Port),
		Handler:  handlers.RecoveryHandler(handlers.RecoveryLogger(httpLog), handlers.PrintRecoveryStack(true))(handlers.CombinedLoggingHandler(a.httpLogWriter(), r)),
		ErrorLog: httpLog,
	}

	// dispatch web server go routine
	go func() {
		tlsCert := path.Join(a.configDir(), constants.TLSCertFile)
		tlsKey := path.Join(a.configDir(), constants.TLSKeyFile)
		if err := h.ListenAndServeTLS(tlsCert, tlsKey); err != nil {
			log.WithError(err).Info("Failed to start HTTPS server")
			stop <- syscall.SIGTERM
		}
	}()

	fmt.Fprintln(a.consoleWriter(), "Threat Detection Service is running")
	// TODO dispatch Service status checker goroutine
	<-stop
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := h.Shutdown(ctx); err != nil {
		log.WithError(err).Info("Failed to gracefully shutdown webserver")
		return err
	}
	return nil
}

func (a *App) start() error {
	fmt.Fprintln(a.consoleWriter(), `Forwarding to "systemctl start tdservice"`)
	systemctl, err := exec.LookPath("systemctl")
	if err != nil {
		return err
	}
	return syscall.Exec(systemctl, []string{"systemctl", "start", "tdservice"}, os.Environ())
}

func (a *App) stop() error {
	fmt.Fprintln(a.consoleWriter(), `Forwarding to "systemctl stop tdservice"`)
	systemctl, err := exec.LookPath("systemctl")
	if err != nil {
		return err
	}
	return syscall.Exec(systemctl, []string{"systemctl", "stop", "tdservice"}, os.Environ())
}

func (a *App) status() error {
	fmt.Fprintln(a.consoleWriter(), `Forwarding to "systemctl status tdservice"`)
	systemctl, err := exec.LookPath("systemctl")
	if err != nil {
		return err
	}
	return syscall.Exec(systemctl, []string{"systemctl", "status", "tdservice"}, os.Environ())
}

func (a *App) uninstall(keepConfig bool) {
	err := os.Remove(a.executablePath())
	if err != nil {
		log.WithError(err).Error("error removing executable")
	}
	if !keepConfig {
		err = os.RemoveAll(a.configDir())
		if err != nil {
			log.WithError(err).Error("error removing config dir")
		}
	}
	err = os.RemoveAll(a.runDir())
	if err != nil {
		log.WithError(err).Error("error removing config dir")
	}
	err = os.RemoveAll(a.dataDir())
	if err != nil {
		log.WithError(err).Error("error removing data dir")
	}
	fmt.Fprintln(a.consoleWriter(), "Threat Detection Service uninstalled")
	a.stop()
}

func validateCmdAndEnv(env_names_cmd_opts map[string]string, flags *flag.FlagSet) error {

	env_names := make([]string, 0)
	for k, _ := range env_names_cmd_opts {
		env_names = append(env_names, k)
	}

	missing, valid_err := validation.ValidateEnvList(env_names)
	if valid_err != nil && missing != nil {
		for _, m := range missing {
			if cmd_f := flags.Lookup(env_names_cmd_opts[m]); cmd_f == nil {
				return errors.New("Insufficient arguments")
			}
		}
	}
	return nil
}

func validateSetupArgs(cmd string, args []string) error {

	var fs *flag.FlagSet

	switch cmd {
	default:
		return errors.New("Unknown command")

	case "database":

		env_names_cmd_opts := map[string]string{
			"TDS_DB_HOSTNAME": "db-host",
			"TDS_DB_PORT":     "db-port",
			"TDS_DB_USERNAME": "db-user",
			"TDS_DB_PASSWORD": "db-pass",
			"TDS_DB_NAME":     "db-name",
		}

		fs = flag.NewFlagSet("database", flag.ContinueOnError)
		fs.String("db-host", "", "Database Hostname")
		fs.Int("db-port", 0, "Database Port")
		fs.String("db-user", "", "Database Username")
		fs.String("db-pass", "", "Database Password")
		fs.String("db-name", "", "Database Name")

		err := fs.Parse(args)
		if err != nil {
			return fmt.Errorf("Fail to parse arguments: %s", err.Error())
		}
		return validateCmdAndEnv(env_names_cmd_opts, fs)

	case "admin":

		env_names_cmd_opts := map[string]string{
			"TDS_ADMIN_USERNAME": "admin-user",
			"TDS_ADMIN_PASSWORD": "admin-pass",
		}

		fs = flag.NewFlagSet("admin", flag.ContinueOnError)
		fs.String("admin-user", "", "Username for admin authentication")
		fs.String("admin-pass", "", "Password for admin authentication")

		err := fs.Parse(args)
		if err != nil {
			return fmt.Errorf("Fail to parse arguments: %s", err.Error())
		}
		return validateCmdAndEnv(env_names_cmd_opts, fs)

	case "server":
		// this has a default port value on 8443
		return nil

	case "tls":

		env_names_cmd_opts := map[string]string{
			"TDA_TLS_HOST_NAMES": "host_names",
		}

		fs = flag.NewFlagSet("tls", flag.ContinueOnError)
		fs.String("host_names", "", "comma separated list of hostnames to add to TLS self signed cert")

		err := fs.Parse(args)
		if err != nil {
			return fmt.Errorf("Fail to parse arguments: %s", err.Error())
		}
		return validateCmdAndEnv(env_names_cmd_opts, fs)

	case "all":
		if len(args) != 0 {
			return errors.New("Please setup the arguments with env")
		}
	}

	return nil
}
