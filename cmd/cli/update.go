package cli

import "github.com/spf13/cobra"

var Update = &cobra.Command{
	Use:   "update [name...]",
	Short: "update package(s)",

	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}
