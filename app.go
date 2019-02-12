package main

import (
	"intel/isecl/threat-detection-service/middleware"
	"strings"
	"intel/isecl/threat-detection-service/tasks"
	"intel/isecl/lib/common/setup"
	"flag"
	"path"
	"context"
	"fmt"
	"intel/isecl/threat-detection-service/config"
	"intel/isecl/threat-detection-service/constants"
	"intel/isecl/threat-detection-service/repository"
	"intel/isecl/threat-detection-service/repository/postgres"
	"intel/isecl/threat-detection-service/resource"
	"io"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	stdlog "log"

	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type App struct {
	ConfigDir      string
	RunDir         string
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
	return config.Global
}

func (a *App) Run(args []string) error {
	if len(args) < 2 {
		a.printUsage()
		os.Exit(1)
	}
	//bin := args[0]
	cmd := args[1]
	switch cmd {
	case "run":
		return a.startServer()
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
		if len(args) < 2 {
			fmt.Fprintln(os.Stdout, "Available setup tasks:\n- server\n- tls\n- all")
		}
		task := strings.ToLower(args[1])
		flags := args[2:]
		setupRunner := &setup.Runner {
			Tasks: []setup.Task{
				tasks.Server{
					Flags: flags,
					Config: a.configuration(),
				},
				tasks.TLS{
					Flags: flags,
					TLSCertFile: path.Join(a.ConfigDir, constants.TLSCertFile),
					TLSKeyFile: path.Join(a.ConfigDir, constants.TLSKeyFile),
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
	// start webserver
	var sslMode string
	c := a.configuration()
	if a.configuration().Postgres.SSLMode {
		sslMode = "enable"
	} else {
		sslMode = "disable"
	}
	var db *gorm.DB
	var dbErr error
	const numAttempts = 4
	for i := 0; i < numAttempts; i = i + 1 {
		const retryTime = 5
		db, dbErr = gorm.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=%s",
			c.Postgres.Hostname, c.Postgres.Port, c.Postgres.Username, c.Postgres.DBName, c.Postgres.Password, sslMode))
		if dbErr != nil {
			log.WithError(dbErr).Info("Failed to connect to DB, retrying")
			fmt.Fprintf(a.consoleWriter(), "Failed to connect to DB, retrying in %d seconds...\n", retryTime)
		} else {
			break
		}
		time.Sleep(retryTime * time.Second)
	}
	defer db.Close()
	if dbErr != nil {
		log.WithError(dbErr).Infof("Failed to connect to db after %d attempts\n", numAttempts)
		return dbErr
	}
	tdsDB := &postgres.PostgresDatabase{DB: db}
	tdsDB.Migrate()
	r := mux.NewRouter().PathPrefix("/tds").Subrouter()
	r.Use(middleware.NewBasicAuth(tdsDB.UserRepository()))
	func(setters ...func(*mux.Router, repository.TDSDatabase)) {
		for _, s := range setters {
			s(r, tdsDB)
		}
	}(resource.SetHosts, resource.SetReports)

	stop := make(chan os.Signal)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	httpLog := stdlog.New(a.httpLogWriter(), "", 0)
	h := &http.Server{
		Addr:     fmt.Sprintf(":%d", c.Port),
		Handler:  handlers.RecoveryHandler(handlers.RecoveryLogger(httpLog), handlers.PrintRecoveryStack(true))(handlers.CombinedLoggingHandler(a.httpLogWriter(), r)),
		ErrorLog: httpLog,
	}
	go func() {
		tlsCert := path.Join(a.ConfigDir, constants.TLSCertFile)
		tlsKey := path.Join(a.ConfigDir, constants.TLSKeyFile)
		if err := h.ListenAndServeTLS(tlsCert, tlsKey); err != nil {
			log.WithError(err).Info("Failed to start HTTPS server")
			stop <- syscall.SIGTERM
		}
	}()
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
		cmd := exec.Command(a.ExecutablePath, "run")
		err := cmd.Start()
		if err != nil {
			log.WithError(err).Error("Failed to start tdagentd")
			return err
		}
		cmd.Process.Release()
		fmt.Fprintln(a.consoleWriter(), "Started Threat Detection Agent")
	} else {
		fmt.Fprintln(a.consoleWriter(), "Threat Detection Agent is already running")
	}
	return nil
}

func (a *App) stop() error {
	if s, pid := a.state(); s == constants.Running {
		if err := syscall.Kill(pid, syscall.SIGTERM); err != nil {
			log.WithError(err).Error("Failed to terminate Threat Detection Agent with signal SIGTERM")
			fmt.Fprintln(a.consoleWriter(), "Failed to stop Threat Detection Agent")
			return err
		}
		fmt.Fprintln(a.consoleWriter(), "Threat Detection agent stopped")
	} else {
		fmt.Fprintln(a.consoleWriter(), "Threat Detection Agent is already stopped")
	}
	return nil
}

func (a *App) state() (state constants.State, pid int) {
	pidFile := path.Join(a.RunDir, constants.PIDFile)
	pid, err := config.CheckPidFile(pidFile)
	if err != nil {
		return
	}
	p, err := os.FindProcess(pid)
	if err != nil {
		return
	}
	if err := p.Signal(syscall.Signal(0)); err != nil {
		return
	}
	state = constants.Running
	return
}

func (a *App) status() error {
	s, pid := a.state()
	if s == constants.Running {
		fmt.Fprintf(a.consoleWriter(), "Threat Detection Agent is running (PID: %d)\n", pid)
	} else {
		fmt.Fprintln(a.consoleWriter(), "Threat Detection Agent is not running")
	}
	return nil
}

func (a *App) uninstall(keepConfig bool) {
	os.Remove(a.ExecutablePath)
	if !keepConfig {
		os.RemoveAll(a.ConfigDir)
	}
	os.RemoveAll(a.RunDir)
}
