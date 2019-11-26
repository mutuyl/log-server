package main

import (
	log "github.com/alecthomas/log4go"
)

func initLogger() {
	var logFilename = "log.txt"
	logger = make(log.Logger)

	fileLogWriter := log.NewFileLogWriter(logFilename, false)
	consoleWriter := log.NewConsoleLogWriter()
	consoleWriter.SetFormat("[%T] (%S) %M")

	logger.AddFilter("stdout", log.FINEST, consoleWriter)
	logger.AddFilter("logfile", log.INFO, fileLogWriter)

	return
}
