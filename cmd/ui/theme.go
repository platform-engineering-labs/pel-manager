package ui

import (
	"image/color"

	"charm.land/fang/v2"
	"charm.land/huh/v2"
	"charm.land/huh/v2/spinner"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/exp/charmtone"
)

type FormTheme struct{}

func (f *FormTheme) Theme(isDark bool) *huh.Styles {
	lightDark := lipgloss.LightDark(isDark)

	cream := lightDark(lipgloss.Color("#FFFDF5"), lipgloss.Color("#FFFDF5"))
	green := lightDark(lipgloss.Color("#02BA84"), lipgloss.Color("#02BF87"))
	orange := lipgloss.Color("#FF5F15")

	theme := huh.ThemeCharm(isDark)

	theme.Focused.SelectedOption = theme.Focused.SelectedOption.Foreground(orange)
	theme.Blurred.NoteTitle = theme.Focused.NoteTitle.Foreground(green).Bold(false).MarginBottom(1)
	theme.Focused.FocusedButton = theme.Focused.FocusedButton.Foreground(cream).Background(orange)

	return theme
}

type SpinnerTheme struct{}

func (s *SpinnerTheme) Theme(isDark bool) *spinner.Styles {
	orange := lipgloss.Color("#FF5F15")

	theme := spinner.ThemeDefault(isDark)
	theme.Spinner = theme.Spinner.Foreground(orange)

	return theme
}

func FangTheme(darkFunc lipgloss.LightDarkFunc) fang.ColorScheme {
	orange := lipgloss.Color("#FF5F15")

	theme := fang.DefaultColorScheme(darkFunc)
	theme.ErrorHeader = [2]color.Color{
		charmtone.Butter,
		orange,
	}

	theme.Command = darkFunc(orange, orange)

	return theme
}
