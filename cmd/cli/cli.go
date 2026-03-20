package cli

import (
	"context"
	"fmt"
	"time"

	"charm.land/huh/v2/spinner"
	"github.com/platform-engineering-labs/pel-mananager/cmd/ui"
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
		orb, err := SetupInteractive(cmd)
		if err != nil {
			return err
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
