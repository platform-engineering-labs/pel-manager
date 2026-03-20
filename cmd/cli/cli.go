package cli

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"charm.land/huh/v2/spinner"
	"github.com/platform-engineering-labs/orbital/mgr"
	"github.com/platform-engineering-labs/pel-mananager/cmd/ui"
	"github.com/platform-engineering-labs/pel-mananager/sys"
	"github.com/platform-engineering-labs/pel-mananager/vals"
	"github.com/spf13/cobra"
)

func init() {
	Root.AddCommand(Install)
	Root.AddCommand(List)
	Root.AddCommand(Remove)
	Root.AddCommand(Update)
	Root.AddCommand(Versions)

	Root.PersistentFlags().String("log", "", "log level: ERR | WARN | INFO | DEBUG | FATAL")
}

var Root = &cobra.Command{
	Use:   "pelmgr",
	Short: "pel manager - install/update/remove PEL tools",

	RunE: func(cmd *cobra.Command, args []string) error {
		level := LoggerFromFlag(cmd.Flags().GetString("log"))

		if !sys.IsPrivilegedUser() {
			if !sys.SudoSessionActive() {
				fmt.Println("PEL Manager must run as a privileged user")
			}

			var cmdArgs []string
			if level != "" {
				cmdArgs = append(cmdArgs, "--log", level)
			}

			err := sys.InvokeSelfWithSudo(cmdArgs...)
			if err != nil {
				return err
			}
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
			WithTheme(&ui.SpinnerTheme{}).
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

		return spinner.New().
			Title(fmt.Sprintf("%s: %s", reqOp, manager.Selection)).
			ActionWithErr(func(ctx context.Context) error {
				switch manager.Request() {
				case ui.Install:
					return orb.Install(manager.Selection)
				case ui.Update:
					time.Sleep(1 * time.Second)
					return orb.Update(manager.Selection)
				case ui.Remove:
					time.Sleep(1 * time.Second)
					return orb.Remove(manager.Selection)
				}
				return nil
			}).
			WithTheme(&ui.SpinnerTheme{}).
			Run()
	},
}
