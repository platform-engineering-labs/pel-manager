package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var Remove = &cobra.Command{
	Use:   "remove [name...]",
	Short: "remove package(s)",

	RunE: func(cmd *cobra.Command, args []string) error {
		orb, err := Setup(cmd)
		if err != nil {
			return err
		}

		if cmd.Flags().NArg() == 0 {
			return fmt.Errorf("remove: must specify at least one package")
		}

		return orb.Remove(cmd.Flags().Args()...)
	},
}
