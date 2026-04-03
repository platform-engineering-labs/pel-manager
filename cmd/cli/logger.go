package cli

import (
	"os"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

var Logger = log.New(os.Stderr)

func LoggerFromCmd(cmd *cobra.Command) {
	level, _ := cmd.Flags().GetString("level")

	switch level {
	case "ERR":
		Logger.SetLevel(log.ErrorLevel)
	case "INFO":
		Logger.SetLevel(log.InfoLevel)
	case "DEBUG":
		Logger.SetLevel(log.DebugLevel)
	default:
		if cmd.Name() == "pelmgr" {
			Logger.SetLevel(log.WarnLevel)
		} else {
			Logger.SetLevel(log.InfoLevel)
		}
	}
}
