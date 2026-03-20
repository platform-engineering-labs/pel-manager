package cli

import "github.com/spf13/cobra"

var Versions = &cobra.Command{
	Use:   "versions [name]",
	Short: "show versions for [name]",

	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}
