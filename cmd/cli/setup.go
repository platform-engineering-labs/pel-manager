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
	level, _ := cmd.Flags().GetString("log")
	if level == "" {
		level = "INFO"
	}

	LoggerFromFlag(level, nil)

	orb, err := setup(root, channel)
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
	channel, _ := cmd.Flags().GetString("channel")
	root, _ := cmd.Flags().GetString("install-path")
	LoggerFromFlag(cmd.Flags().GetString("log"))

	orb, err := setup(root, channel)
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
		} else {
			return nil, errors.New("cancelled")
		}
	}

	return orb, nil
}

func setup(root string, channel string) (*mgr.Manager, error) {
	cfg := vals.TreeConfig
	cfg.Repositories[0].Uri.Fragment = channel

	return mgr.New(slog.New(Logger), root, cfg)
}
