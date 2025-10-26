package logger_factory

import (
	"log"
	"os"
)

const filePermission = 0666

/////////////////////////////////////////////////////////////////////////////////////////////
//
//  @Brief			GetLogger creates a new logger that writes to the specified file.
//
//	@Description	The log file is truncated if it exists, created otherwise
//	@Description	Log entries include "TSM" prefix, date, time, and file:line information
//
//	@Param			filename	string			Path to the log file
//
//	@Return			*log.Logger	Configured logger instance
//
/////////////////////////////////////////////////////////////////////////////////////////////
func GetLogger(filename string) *log.Logger {
	logfile, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, filePermission)

	if err != nil {
		panic("No can do for logging")
	}

	return log.New(logfile, "TSM", log.Ldate|log.Ltime|log.Lshortfile)
}
