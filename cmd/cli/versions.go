package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var Versions = &cobra.Command{
	Use:   "versions [name]",
	Short: "show versions for [name]",

	RunE: func(cmd *cobra.Command, args []string) error {
		orb, err := Setup(cmd)
		if err != nil {
			return err
		}

		if cmd.Flags().NArg() == 0 {
			return fmt.Errorf("versions: must specify a package")
		}

		err = orb.Refresh()
		if err != nil {
			return err
		}

		available, err := orb.Available()
		if err != nil {
			return err
		}

		if _, ok := available[cmd.Flags().Arg(0)]; !ok {
			return fmt.Errorf("versions: no versions available for %s", cmd.Flags().Arg(0))
		}

		fmt.Printf("Status: %s\n", available[cmd.Flags().Arg(0)].Status)
		for _, pkg := range available[cmd.Flags().Arg(0)].Available {
			fmt.Printf("%s  %s\n", pkg.Name, pkg.Version.String())
		}

		return nil
	},
}
