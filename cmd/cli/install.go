package cli

import "github.com/spf13/cobra"

var Install = &cobra.Command{
	Use:   "install [name...] | [name@version...]",
	Short: "install package(s)",

	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}
