package cli

import "github.com/spf13/cobra"

var List = &cobra.Command{
	Use:   "list",
	Short: "list installed packages",

	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}
