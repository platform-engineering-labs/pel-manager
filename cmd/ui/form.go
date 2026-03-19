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
)

const (
	Cancel  Operation = "cancel"
	Install Operation = "install"
	Update  Operation = "update"
	Remove  Operation = "remove"
)

type Operation = string
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
				Affirmative("Yes!").
				Negative("No.").
				Value(&sr.Confirm),
		),
	)

	return sr
}

func (sr *SetupRoot) Run() error {
	return sr.form.Run()
}

type Manager struct {
	Available map[string]*records.Status

	Selection      string
	Operation      bool
	OperationLabel string

	form *huh.Form
}

func NewManager(available map[string]*records.Status) *Manager {
	m := &Manager{Available: available, Operation: true}
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
			huh.NewConfirm().
				TitleFunc(func() string {
					return fmt.Sprintf("Manage: %s", m.Selection)
				}, &m.Selection).
				Affirmative("Cancel").
				Negative("Remove").
				Value(&m.Operation),
		).WithHideFunc(func() bool {
			hasUpdate, _ := m.Available[m.Selection].HasUpdate()
			if m.Available[m.Selection].Status == candidate.Available {
				return true
			} else {
				return hasUpdate == true
			}
		}),
		huh.NewGroup(
			huh.NewConfirm().
				TitleFunc(func() string {
					return fmt.Sprintf("Manage: %s", m.Selection)
				}, &m.Selection).
				DescriptionFunc(func() string {
					return fmt.Sprintf("\nversion: %s", m.Available[m.Selection].Available[0].Version.Short())
				}, &m.Selection).
				Affirmative("Update").
				Negative("Remove").
				Value(&m.Operation),
		).WithHideFunc(func() bool {
			hasUpdate, _ := m.Available[m.Selection].HasUpdate()
			if m.Available[m.Selection].Status == candidate.Available {
				return true
			} else {
				return hasUpdate == false
			}
		}),
		huh.NewGroup(
			huh.NewConfirm().
				TitleFunc(func() string {
					return fmt.Sprintf("Manage: %s", m.Selection)
				}, &m.Selection).
				DescriptionFunc(func() string {
					return fmt.Sprintf("\nversion: %s", m.Available[m.Selection].Available[0].Version.Short())
				}, &m.Selection).
				Affirmative("Install").
				Negative("Cancel").
				Value(&m.Operation),
		).WithHideFunc(func() bool {
			return m.Available[m.Selection].Status != candidate.Available
		}),
	).WithTheme(&FormTheme{})

	return m
}

func (m Manager) Request() Operation {
	installed := m.Available[m.Selection].Status == candidate.Frozen || m.Available[m.Selection].Status == candidate.Installed
	hasUpdate, _ := m.Available[m.Selection].HasUpdate()

	switch m.Operation {
	case true:
		if hasUpdate == false {
			return Cancel
		}

		if installed {
			return Update
		} else {
			return Install
		}
	default:
		if installed {
			return Remove
		} else {
			return Cancel
		}
	}
}

func (m Manager) Run() error {
	return m.form.Run()
}
