package main

import (
	"intel/isecl/tdservice/middleware"
	"strings"
	"intel/isecl/tdservice/tasks"
	"intel/isecl/lib/common/setup"
	"flag"
	"path"
	"context"
	"fmt"
	"intel/isecl/tdservice/config"
	"intel/isecl/tdservice/constants"
	"intel/isecl/tdservice/repository"
	"intel/isecl/tdservice/repository/postgres"
	"intel/isecl/tdservice/resource"
	"io"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
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
	DataDir 		string
	ExecutablePath string
	Config         *config.Configuration
	ConsoleWriter  io.Writer
	LogWriter      io.Writer
	HTTPLogWriter  io.Writer
}

func (a *App) printUsage() {
	fmt.Fprintln(a.consoleWriter(), "Help message here")
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
	return path.Join(constants.ExecutableDir, "tdservice")
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

func (a *App) Run(args []string) error {
	if len(args) < 2 {
		a.printUsage()
		os.Exit(1)
	}
	//bin := args[0]
	cmd := args[1]
	switch cmd {
	default:
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
	case "setup": 
		if len(args) <= 2 {
			fmt.Fprintln(os.Stdout, "Available setup tasks:\n- database\n- admin\n- server\n- tls\n-----------------\n- [all]")
			os.Exit(1)
		}
		task := strings.ToLower(args[2])
		flags := args[2:]
		setupRunner := &setup.Runner {
			Tasks: []setup.Task{
				tasks.Database{
					Flags: flags,
					Config: a.configuration(),
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
					Flags: flags,
					Config: a.configuration(),
					ConsoleWriter: os.Stdout,
				},
				tasks.TLS{
					Flags: flags,
					TLSCertFile: path.Join(a.configDir(), constants.TLSCertFile),
					TLSKeyFile: path.Join(a.configDir(), constants.TLSKeyFile), 
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
	if s, _ := a.state(); s == constants.Stopped {
		cmd := exec.Command(a.executablePath(), "run")
		err := cmd.Start()
		if err != nil {
			log.WithError(err).Error("Failed to start tdservice as a daemon")
			return err
		}
		pidFile := path.Join(a.runDir(), constants.PIDFile)
		err = config.WritePidFile(pidFile, cmd.Process.Pid)
		cmd.Process.Release()
		if err != nil {
			log.WithError(err).Error("failed to write pid file")
			return err
		}
		fmt.Fprintln(a.consoleWriter(), "Started Threat Detection Service")
	} else {
		fmt.Fprintln(a.consoleWriter(), "Threat Detection Service is already running")
	}
	return nil
}

func (a *App) stop() error {
	if s, pid := a.state(); s == constants.Running {
		if err := syscall.Kill(pid, syscall.SIGTERM); err != nil {
			log.WithError(err).Error("Failed to terminate Threat Detection Service with signal SIGTERM")
			fmt.Fprintln(a.consoleWriter(), "Failed to stop Threat Detection Service")
			return err
		}
		fmt.Fprintln(a.consoleWriter(), "Threat Detection Service stopped")
	} else {
		fmt.Fprintln(a.consoleWriter(), "Threat Detection Service is already stopped")
	}
	return nil
}

func (a *App) state() (state constants.State, pid int) {
	pidFile := path.Join(a.runDir(), constants.PIDFile)
	pid, err := config.CheckPidFile(pidFile)
	if err != nil {
		log.WithError(err).Debug("failed to check pid file")
		return
	}
	p, err := os.FindProcess(pid)
	if err != nil {
		log.WithError(err).Error("failed to find process")
		return
	}
	if err := p.Signal(syscall.Signal(0)); err != nil {
		log.WithError(err).Error("failed to signal process")
		return
	}
	state = constants.Running
	return
}

func (a *App) status() error {
	s, pid := a.state()
	if s == constants.Running {
		fmt.Fprintf(a.consoleWriter(), "Threat Detection Service is running (PID: %d)\n", pid)
	} else {
		fmt.Fprintln(a.consoleWriter(), "Threat Detection Service is not running")
	}
	return nil
}

func (a *App) uninstall(keepConfig bool) {
	a.stop()
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
}
