package cli

import (
	"github.com/spf13/cobra"
)

var Update = &cobra.Command{
	Use:   "update [name...]",
	Short: "update package(s)",

	RunE: func(cmd *cobra.Command, args []string) error {
		orb, err := Setup(cmd)
		if err != nil {
			return err
		}

		err = orb.Refresh()
		if err != nil {
			return err
		}

		return orb.Update(cmd.Flags().Args()...)
	},
}
