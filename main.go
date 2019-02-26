package main

import (
	"intel/isecl/tdservice/constants"
	"os"
	"path"
)

func openLogFiles() (logFile *os.File, httpLogFile *os.File) {
	logFilePath := path.Join(constants.LogDir, constants.LogFile)
	logFile, _ = os.OpenFile(logFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	httpLogFilePath := path.Join(constants.LogDir, constants.HTTPLogFile)
	httpLogFile, _ = os.OpenFile(httpLogFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	return
}

func main() {
	l, h := openLogFiles()
	defer l.Close()
	defer h.Close()
	app := &App{
		LogWriter:     l,
		HTTPLogWriter: h,
	}
	err := app.Run(os.Args)
	if err != nil {
		os.Exit(1)
	}
}
