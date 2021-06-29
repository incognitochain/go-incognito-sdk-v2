package incclient

import (
	"fmt"
	"log"
	"os"
)

// IncLogger implements a logger for the incclient package.
type IncLogger struct {
	Log      *log.Logger
	IsEnable bool
}

// NewLogger creates a new IncLogger. If isEnable = true, it will do logging.
// If logFile is set, it will store logging information into the given logFile.
func NewLogger(isEnable bool, logFile ...string) *IncLogger {
	writer := os.Stdout
	if len(logFile) != 0 {
		var err error
		writer, err = os.OpenFile(logFile[0], os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			fmt.Println("Error opening file:", err)
			os.Exit(1)
		}
	}
	Log := log.New(writer, "", log.Ldate|log.Ltime|log.Lshortfile)

	return &IncLogger{
		Log:      Log,
		IsEnable: isEnable,
	}
}

var incLogger = NewLogger(false)
