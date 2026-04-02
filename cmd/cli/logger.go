package cli

import (
	"os"

	"github.com/charmbracelet/log"
)

var Logger = log.New(os.Stderr)

func LoggerFromFlag(level string, _ error) {
	switch level {
	case "ERR":
		Logger.SetLevel(log.ErrorLevel)
	case "INFO":
		Logger.SetLevel(log.InfoLevel)
	case "DEBUG":
		Logger.SetLevel(log.DebugLevel)
	default:
		Logger.SetLevel(log.WarnLevel)
	}
}
