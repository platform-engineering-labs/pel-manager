package cli

import (
	"errors"
	"log/slog"

	"github.com/platform-engineering-labs/orbital/mgr"
	"github.com/platform-engineering-labs/pel-mananager/cmd/ui"
	"github.com/platform-engineering-labs/pel-mananager/vals"
	"github.com/spf13/cobra"
)

func Setup(cmd *cobra.Command) (*mgr.Manager, error) {
	channel, _ := cmd.Flags().GetString("channel")
	root, _ := cmd.Flags().GetString("install-path")
	yes, _ := cmd.Flags().GetBool("yes")

	LoggerFromCmd(cmd)

	cfg := vals.TreeConfig
	cfg.Repositories[0].Uri.Fragment = channel

	orb, err := mgr.New(slog.New(Logger), root, cfg)
	if err != nil {
		return nil, err
	}

	if orb.Ready() != true {
		setupRoot := ui.NewSetupRoot(yes)

		if !yes {
			err := setupRoot.Run()
			if err != nil {
				return nil, err
			}
		}

		if setupRoot.Confirm {
			_, err := orb.Initialize()
			if err != nil {
				return nil, err
			}
		} else {
			return nil, errors.New("cancelled")
		}
	}

	return orb, nil
}
