package cli

import (
	"context"
	"fmt"
	"log/slog"

	"charm.land/huh/v2/spinner"
	"github.com/charmbracelet/log"
	"github.com/platform-engineering-labs/orbital/mgr"
	"github.com/platform-engineering-labs/pel-mananager/cmd/ui"
	"github.com/platform-engineering-labs/pel-mananager/sys"
	"github.com/platform-engineering-labs/pel-mananager/vals"
	"github.com/spf13/cobra"
)

func init() {
	Root.PersistentFlags().String("log", "", "log level: ERR, WARN, INFO, DEBUG, FATAL")
}

var Root = &cobra.Command{
	Use:   "pelmgr",
	Short: "pel manager - install/update/remove PEL tools",

	RunE: func(cmd *cobra.Command, args []string) error {
		logLevel, _ := cmd.Flags().GetString("log")

		if !sys.IsPrivilegedUser() {
			if !sys.SudoSessionActive() {
				fmt.Println("PEL Manager must run as a privileged user")
			}
			err := sys.InvokeSelfWithSudo(args...)
			if err != nil {
				return err
			}
		}

		switch logLevel {
		case "ERR":
			Logger.SetLevel(log.ErrorLevel)
		case "WARN":
			Logger.SetLevel(log.WarnLevel)
		case "INFO":
			Logger.SetLevel(log.InfoLevel)
		case "DEBUG":
			Logger.SetLevel(log.DebugLevel)
		default:
			Logger.SetLevel(log.FatalLevel)
		}

		orb, err := mgr.New(slog.New(Logger), vals.ManagedRoot, vals.TreeConfig)
		if err != nil {
			return err
		}

		if orb.Ready() != true {
			setupRoot := ui.NewSetupRoot()
			err := setupRoot.Run()
			if err != nil {
				return err
			}

			if setupRoot.Confirm {
				_, err := orb.Initialize()
				if err != nil {
					return err
				}
			}
		}

		err = spinner.New().
			Title("Loading catalog").
			ActionWithErr(func(ctx context.Context) error {
				return orb.Refresh()
			}).
			Run()
		if err != nil {
			return err
		}

		available, err := orb.Available()
		if err != nil {
			return err
		}

		manager := ui.NewManager(available)
		err = manager.Run()
		if err != nil {
			return err
		}

		reqOp := ""
		switch manager.Request() {
		case ui.Install:
			reqOp = "Installing"
		case ui.Update:
			reqOp = "Updating"
		case ui.Remove:
			reqOp = "Removing"
		default:
			return nil
		}

		err = spinner.New().
			Title(fmt.Sprintf("%s: %s", reqOp, manager.Selection)).
			ActionWithErr(func(ctx context.Context) error {
				switch manager.Request() {
				case ui.Install:
					return orb.Install(manager.Selection)
				case ui.Update:
					return orb.Update(manager.Selection)
				case ui.Remove:
					return orb.Remove(manager.Selection)
				}
				return nil
			}).
			Run()

		return err
	},
}
