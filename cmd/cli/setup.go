package cli

import (
	"fmt"
	"log/slog"

	"github.com/platform-engineering-labs/orbital/mgr"
	"github.com/platform-engineering-labs/pel-mananager/cmd/ui"
	"github.com/platform-engineering-labs/pel-mananager/sys"
	"github.com/platform-engineering-labs/pel-mananager/vals"
	"github.com/spf13/cobra"
)

func Setup(cmd *cobra.Command) (*mgr.Manager, error) {
	level, _ := cmd.Flags().GetString("log")
	if level == "" {
		level = "INFO"
	}

	LoggerFromFlag(level, nil)

	orb, err := setup()
	if err != nil {
		return nil, err
	}

	if orb.Ready() != true {
		_, err = orb.Initialize()
		if err != nil {
			return nil, err
		}
	}

	return orb, nil
}

func SetupInteractive(cmd *cobra.Command) (*mgr.Manager, error) {
	LoggerFromFlag(cmd.Flags().GetString("log"))

	orb, err := setup()
	if err != nil {
		return nil, err
	}

	if orb.Ready() != true {
		setupRoot := ui.NewSetupRoot()
		err := setupRoot.Run()
		if err != nil {
			return nil, err
		}

		if setupRoot.Confirm {
			_, err := orb.Initialize()
			if err != nil {
				return nil, err
			}
		}
	}

	return orb, nil
}

func setup() (*mgr.Manager, error) {
	if !sys.IsPrivilegedUser() {
		if !sys.SudoSessionActive() {
			fmt.Println("PEL Manager must run as a privileged user")
		}

		err := sys.InvokeSelfWithSudo()
		if err != nil {
			return nil, err
		}
	}

	return mgr.New(slog.New(Logger), vals.ManagedRoot, vals.TreeConfig)
}
