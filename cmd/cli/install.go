package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var Install = &cobra.Command{
	Use:   "install [name...] | [name@version...]",
	Short: "install package(s)",

	RunE: func(cmd *cobra.Command, args []string) error {
		orb, err := Setup(cmd)
		if err != nil {
			return err
		}

		if cmd.Flags().NArg() == 0 {
			return fmt.Errorf("install: must specify at least one package")
		}

		err = orb.Refresh()
		if err != nil {
			return err
		}

		return orb.Install(cmd.Flags().Args()...)
	},
}
