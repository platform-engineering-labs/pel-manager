package cli

import (
	"fmt"

	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/table"
	"github.com/spf13/cobra"
)

var List = &cobra.Command{
	Use:   "list",
	Short: "list installed packages",

	RunE: func(cmd *cobra.Command, args []string) error {
		orb, err := Setup(cmd)
		if err != nil {
			return err
		}

		packages, err := orb.List()
		if err != nil {
			return err
		}

		if len(packages) == 0 {
			return nil
		}

		tbl := table.New()
		tbl.Headers("Name", "Version", "Summary", "Status")

		for _, pkg := range packages {
			status := ""

			if pkg.Frozen {
				status = lipgloss.NewStyle().Foreground(lipgloss.Color("12")).Render("*")
			}

			tbl.Row(pkg.Name, pkg.Version.String(), pkg.Summary, status)
		}

		fmt.Println(tbl.Render())

		return nil
	},
}
