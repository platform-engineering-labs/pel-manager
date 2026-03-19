/*
Keeping this for later to do a full blown UI
*/
package ui

import (
	"fmt"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"
	"charm.land/huh/v2"
)

const (
	PackageList View = "packageList"
	Operations  View = "operations"
)

type View = string

// keyMap defines a set of keybindings. To work for help it must satisfy
// key.Map. It could also very easily be a map[string]key.Binding.
type keyMap struct {
	Quit key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Quit},
	}
}

var keys = keyMap{
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

type Model struct {
	Table     table.Model
	Current   View
	Selection string
	Help      help.Model
	Dialog    *huh.Form
	DialogOp  bool
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Run() error {
	m.Table = table.New(
		table.WithColumns([]table.Column{{Title: "AVAILABLE TOOLS", Width: 50}}),
		table.WithRows([]table.Row{
			{
				"formae",
			},
			{
				"ops",
			},
			{"pkl"},
		}),
		table.WithFocused(true),
		table.WithHeight(8),
		table.WithWidth(50),
	)

	m.Table.SetStyles(Styles)
	m.Help = help.New()
	m.Dialog = huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title(m.Selection).
				Affirmative("Install/Update").
				Negative("Remove").
				Value(&m.DialogOp),
		),
	)

	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Println("could not start program:", err)
		return err
	}

	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "esc":
			if m.Table.Focused() {
				m.Table.Blur()
			} else {
				m.Table.Focus()
			}
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			m.Current = Operations
			m.Selection = m.Table.SelectedRow()[0]
			return m, nil
		}
	}
	m.Table, cmd = m.Table.Update(msg)
	return m, cmd
}

func (m Model) View() tea.View {
	switch m.Current {
	default:
		return tea.NewView(TableStyle.Render(m.Table.View()) + "\n  " + m.Table.HelpView() + "  " + m.Help.View(keys) + "\n")
	case Operations:
		return tea.NewView(m.Dialog.View() + "\n")
	}
}
