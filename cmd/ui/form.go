package ui

import (
	"fmt"

	"charm.land/huh/v2"
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

	for pkg, _ := range available {
		options = append(options, huh.NewOption(pkg, pkg))
	}

	m.form = huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select a tool").
				Options(options...).
				Value(&m.Selection),
		),
		huh.NewGroup(
			huh.NewConfirm().
				TitleFunc(func() string {
					return fmt.Sprintf("Manage: %s", m.Selection)
				}, &m.Selection).
				Affirmative("Update").
				Negative("Remove").
				Value(&m.Operation),
		).WithHideFunc(func() bool {
			return m.Available[m.Selection].Status != candidate.Frozen && m.Available[m.Selection].Status != candidate.Installed
		}),
		huh.NewGroup(
			huh.NewConfirm().
				TitleFunc(func() string {
					return fmt.Sprintf("Manage: %s", m.Selection)
				}, &m.Selection).
				Affirmative("Install").
				Negative("Cancel").
				Value(&m.Operation),
		).WithHideFunc(func() bool {
			return m.Available[m.Selection].Status == candidate.Frozen || m.Available[m.Selection].Status == candidate.Installed
		}),
	)

	return m
}

func (m Manager) Request() Operation {
	installed := m.Available[m.Selection].Status == candidate.Frozen || m.Available[m.Selection].Status == candidate.Installed

	switch m.Operation {
	case true:
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
