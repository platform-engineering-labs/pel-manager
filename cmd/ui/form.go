package ui

import (
	"fmt"
	"maps"
	"slices"

	"charm.land/huh/v2"
	"charm.land/lipgloss/v2"
	"github.com/platform-engineering-labs/orbital/opm/candidate"
	"github.com/platform-engineering-labs/orbital/opm/records"
	"github.com/platform-engineering-labs/pel-mananager/vals"
	"github.com/platform-engineering-labs/pelx/theme"
)

const (
	Cancel  Operation = "cancel"
	Install Operation = "install"
	Update  Operation = "update"
	Remove  Operation = "remove"

	Updatable   State = "updatable"
	Installable State = "installable"
	Removable   State = "removable"
)

type Operation = string
type State = string
type SetupRoot struct {
	Confirm bool

	form *huh.Form
}

func NewSetupRoot() *SetupRoot {
	sr := &SetupRoot{}
	sr.form = huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title(fmt.Sprintf("Setup software root at: %s", vals.ManagedRoot)).
				Affirmative("Yes").
				Negative("No").
				Value(&sr.Confirm),
		),
	).WithTheme(&theme.FormTheme{})

	return sr
}

func (sr *SetupRoot) Run() error {
	return sr.form.Run()
}

type Manager struct {
	Available map[string]*records.Status

	Selection string
	Operation Operation

	form *huh.Form
}

func NewManager(available map[string]*records.Status) *Manager {
	m := &Manager{Available: available, Operation: ""}
	var options []huh.Option[string]
	keys := slices.Collect(maps.Keys(available))
	slices.Sort(keys)

	for _, pkg := range keys {
		options = append(options, huh.NewOption(pkg, pkg))
	}

	m.form = huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select a tool").
				Options(options...).
				Value(&m.Selection),

			huh.NewNote().
				TitleFunc(func() string {
					hasUpdate, _ := m.Available[m.Selection].HasUpdate()
					if m.Available[m.Selection].Status == candidate.Available {
						return "status: absent • candidate: available"
					} else {
						if hasUpdate {
							return "status: installed • updates: available"
						} else {
							return "status: installed • updates: none"
						}
					}

				}, &m.Selection),
		).Title(
			lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFFDF5")).
				BorderStyle(lipgloss.DoubleBorder()).
				BorderForeground(lipgloss.Color("#FFFDF5")).
				BorderBottom(true).
				Render("▒▒░ PEL MANAGER"),
		),

		huh.NewGroup(
			NewMultiButton[Operation]().
				TitleFunc(func() string {
					return fmt.Sprintf("Manage: %s", m.Selection)
				}, &m.Selection).
				DescriptionFunc(func() string {
					return fmt.Sprintf("\nversion: %s\n", m.Available[m.Selection].Available[0].Version.Short())
				}, &m.Selection).
				OptionsFunc(func() []ButtonOption[Operation] {
					switch m.State() {
					case Installable:
						return []ButtonOption[Operation]{
							NewButtonOption("Install", Install),
							NewButtonOption("Cancel", Cancel),
						}
					case Removable:
						return []ButtonOption[Operation]{
							NewButtonOption("Cancel", Cancel),
							NewButtonOption("Remove", Remove),
						}
					case Updatable:
						return []ButtonOption[Operation]{
							NewButtonOption("Update", Update),
							NewButtonOption("Remove", Remove),
							NewButtonOption("Cancel", Cancel),
						}
					}
					return nil
				}, &m.Selection).
				Value(&m.Operation),
		),
	).WithTheme(&theme.FormTheme{})

	return m
}

func (m Manager) Run() error {
	return m.form.Run()
}

func (m Manager) State() State {
	if m.Available[m.Selection].Status == candidate.Available {
		return Installable
	} else {
		hasUpdate, _ := m.Available[m.Selection].HasUpdate()
		if hasUpdate {
			return Updatable
		} else {
			return Removable
		}
	}
}
