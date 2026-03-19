package ui

import (
	"charm.land/bubbles/v2/table"
	"charm.land/lipgloss/v2"
	"github.com/platform-engineering-labs/pel-mananager/fmx"
)

var DefaultTable = table.DefaultStyles()

var TableStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.DoubleBorder()).
	BorderForeground(lipgloss.Yellow)

var Styles table.Styles = table.Styles{
	Header: DefaultTable.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Yellow).
		BorderBottom(true).
		Bold(true).
		Width(50).
		Align(lipgloss.Center),

	Selected: DefaultTable.Selected.
		Foreground(lipgloss.Color("229")).
		Bold(true).
		Align(lipgloss.Center).
		Transform(func(s string) string {
			return fmx.Insert(s, "❯❯ ")
		}),

	Cell: DefaultTable.Cell.
		Width(50).
		Align(lipgloss.Center, lipgloss.Center).
		MarginTop(1),
}
