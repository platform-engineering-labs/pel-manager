package ui

// multibutton.go
//
// A custom huh.Field that renders N options as horizontal buttons,
// modelled on huh.Confirm but generic over any comparable type T.
//
// Navigation : ← / h  and  → / l
// Confirm    : enter or space
// Back       : shift+tab

// multibutton.go
//
// A custom huh.Field that renders N options as horizontal buttons,
// modelled on huh.Confirm but generic over any comparable type T.
//
// Navigation : ← / h  and  → / l
// Confirm    : enter or space
// Back       : shift+tab

import (
	"fmt"
	"io"
	"reflect"
	"strings"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/huh/v2"
	"charm.land/lipgloss/v2"
)

// ── Option type ───────────────────────────────────────────────────────────────

// ButtonOption is a label + value pair for MultiButton.
type ButtonOption[T comparable] struct {
	Label string
	Val   T
}

func NewButtonOption[T comparable](label string, val T) ButtonOption[T] {
	return ButtonOption[T]{Label: label, Val: val}
}

// ── MultiButton ───────────────────────────────────────────────────────────────

// MultiButton is a huh.Field that renders N horizontal button options and
// stores the selected value in a variable of type T.
type MultiButton[T comparable] struct {
	fieldKey    string
	title       string
	description string

	// Static options (set via Options).
	options []ButtonOption[T]

	// Dynamic options (set via OptionsFunc).
	optionsFn   func() []ButtonOption[T]
	optionsBind any
	optionsSnap string

	// Dynamic title (set via TitleFunc).
	titleFn   func() string
	titleBind any
	titleSnap string

	// Dynamic description (set via DescriptionFunc).
	descFn   func() string
	descBind any
	descSnap string

	cursor  int
	value   *T
	focused bool
	theme   huh.Theme
	width   int

	validate func(T) error
	err      error
}

func NewMultiButton[T comparable]() *MultiButton[T] {
	return &MultiButton[T]{}
}

// ── Builder API ───────────────────────────────────────────────────────────────

func (m *MultiButton[T]) Key(k string) *MultiButton[T] {
	m.fieldKey = k
	return m
}

func (m *MultiButton[T]) Title(t string) *MultiButton[T] {
	m.title = t
	return m
}

func (m *MultiButton[T]) Description(d string) *MultiButton[T] {
	m.description = d
	return m
}

// Options sets a static slice of options, clearing any previously registered
// OptionsFunc.
func (m *MultiButton[T]) Options(opts ...ButtonOption[T]) *MultiButton[T] {
	m.optionsFn = nil
	m.optionsBind = nil
	m.setOptions(opts)
	return m
}

// OptionsFunc registers a dynamic options producer. f is re-evaluated whenever
// the value pointed-to by bindings changes. Pass nil to run f exactly once.
func (m *MultiButton[T]) OptionsFunc(f func() []ButtonOption[T], bindings any) *MultiButton[T] {
	m.optionsFn = f
	m.optionsBind = bindings
	m.optionsSnap = snapshot(bindings)
	m.setOptions(f())
	return m
}

// TitleFunc sets a dynamic title. f is re-evaluated whenever bindings changes.
// Pass nil to compute the title once.
//
//	NewMultiButton[string]().
//	    TitleFunc(func() string { return "Deploy to " + env }, &env)
func (m *MultiButton[T]) TitleFunc(f func() string, bindings any) *MultiButton[T] {
	m.titleFn = f
	m.titleBind = bindings
	m.titleSnap = snapshot(bindings)
	m.title = m.titleFn()
	return m
}

// DescriptionFunc sets a dynamic description. f is re-evaluated whenever
// bindings changes. Pass nil to compute the description once.
func (m *MultiButton[T]) DescriptionFunc(f func() string, bindings any) *MultiButton[T] {
	m.descFn = f
	m.descBind = bindings
	m.descSnap = snapshot(bindings)
	m.description = f()
	return m
}

func (m *MultiButton[T]) Value(v *T) *MultiButton[T] {
	m.value = v
	m.syncCursor()
	return m
}

func (m *MultiButton[T]) Validate(fn func(T) error) *MultiButton[T] {
	m.validate = fn
	return m
}

// ── Internal helpers ──────────────────────────────────────────────────────────

// setOptions replaces the current option slice and re-syncs the cursor so an
// existing *value pointer keeps pointing at the right button.
func (m *MultiButton[T]) setOptions(opts []ButtonOption[T]) {
	m.options = opts
	m.syncCursor()
}

// syncCursor moves the cursor to the option matching *m.value (if any).
func (m *MultiButton[T]) syncCursor() {
	if m.value == nil || len(m.options) == 0 {
		return
	}
	for i, o := range m.options {
		if o.Val == *m.value {
			m.cursor = i
			return
		}
	}
	// value not in new option set → clamp cursor
	if m.cursor >= len(m.options) {
		m.cursor = 0
	}
}

// snapshot returns a string of the *dereferenced* value of v so that pointer
// identity (which never changes) is not mistaken for value equality.
// Supports single pointers and slices of pointers via reflect.
func snapshot(v any) string {
	if v == nil {
		return "<nil>"
	}
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr && !rv.IsNil() {
		return fmt.Sprintf("%v", rv.Elem().Interface())
	}
	return fmt.Sprintf("%v", v)
}

// recomputeIfChanged re-evaluates any registered Func when its binding has
// changed since the last snapshot. Returns true if anything was refreshed.
func (m *MultiButton[T]) recomputeIfChanged() bool {
	changed := false
	if m.optionsFn != nil {
		if cur := snapshot(m.optionsBind); cur != m.optionsSnap {
			m.optionsSnap = cur
			m.setOptions(m.optionsFn())
			changed = true
		}
	}
	if m.titleFn != nil {
		if cur := snapshot(m.titleBind); cur != m.titleSnap {
			m.titleSnap = cur
			m.title = m.titleFn()
			changed = true
		}
	}
	if m.descFn != nil {
		if cur := snapshot(m.descBind); cur != m.descSnap {
			m.descSnap = cur
			m.description = m.descFn()
			changed = true
		}
	}
	return changed
}

// ── huh.Field / tea.Model interface ──────────────────────────────────────────

func (m *MultiButton[T]) Init() tea.Cmd {
	// Evaluate all Funcs once on init so nothing is blank on first render,
	// even when the Func methods were called before the form was run.
	if m.optionsFn != nil && len(m.options) == 0 {
		m.setOptions(m.optionsFn())
	}
	if m.titleFn != nil && m.title == "" {
		m.title = m.titleFn()
	}
	if m.descFn != nil && m.description == "" {
		m.description = m.descFn()
	}
	return nil
}

func (m *MultiButton[T]) Error() error {
	return m.err
}

func (m *MultiButton[T]) Update(msg tea.Msg) (huh.Model, tea.Cmd) {
	// Always check for binding changes, even when not focused, so that options
	// are fresh when this field comes into view (mirrors huh's behaviour).
	m.recomputeIfChanged()

	if !m.focused {
		return m, nil
	}
	km, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}
	switch km.String() {
	case "left", "h":
		if m.cursor > 0 {
			m.cursor--
		}
	case "right", "l":
		if m.cursor < len(m.options)-1 {
			m.cursor++
		}
	case "enter", " ":
		return m.submit()
	case "shift+tab":
		return m, huh.PrevField
	}
	return m, nil
}

func (m *MultiButton[T]) submit() (huh.Model, tea.Cmd) {
	if len(m.options) == 0 {
		return m, nil
	}
	val := m.options[m.cursor].Val
	if m.validate != nil {
		if err := m.validate(val); err != nil {
			m.err = err
			return m, nil
		}
	}
	m.err = nil
	if m.value != nil {
		*m.value = val
	}
	return m, huh.NextField
}

func (m *MultiButton[T]) View() string {
	st := m.activeStyles()
	var b strings.Builder

	if m.title != "" {
		b.WriteString(st.Title.Render(m.title))
		b.WriteByte('\n')
	}
	if m.description != "" {
		b.WriteString(st.Description.Render(m.description))
		b.WriteByte('\n')
	}

	// Render all buttons side-by-side.
	btns := make([]string, len(m.options))
	for i, o := range m.options {
		if i == m.cursor {
			btns[i] = st.FocusedButton.Render(o.Label)
		} else {
			btns[i] = st.BlurredButton.Render(o.Label)
		}
	}
	b.WriteString(lipgloss.JoinHorizontal(lipgloss.Left, btns...))

	if m.err != nil {
		b.WriteByte('\n')
		b.WriteString(st.ErrorMessage.Render(m.err.Error()))
	}

	return st.Base.Render(b.String())
}

func (m *MultiButton[T]) activeStyles() huh.FieldStyles {
	var s *huh.Styles
	if m.theme != nil {
		s = m.theme.Theme(false)
	} else {
		s = huh.ThemeCharm(false)
	}
	if m.focused {
		return s.Focused
	}
	return s.Blurred
}

func (m *MultiButton[T]) Focus() tea.Cmd {
	m.focused = true
	m.recomputeIfChanged()
	return nil
}
func (m *MultiButton[T]) Blur() tea.Cmd { m.focused = false; return nil }

func (m *MultiButton[T]) Run() error {
	return huh.NewForm(huh.NewGroup(m)).Run()
}

func (m *MultiButton[T]) RunAccessible(w io.Writer, r io.Reader) error {
	if m.optionsFn != nil {
		m.setOptions(m.optionsFn())
	}
	if m.title != "" {
		fmt.Fprintln(w, m.title)
	}
	for i, o := range m.options {
		fmt.Fprintf(w, "  %d) %s\n", i+1, o.Label)
	}
	for {
		fmt.Fprint(w, "Choice [1]: ")
		var raw string
		fmt.Fscan(r, &raw) //nolint:errcheck
		if raw == "" {
			raw = "1"
		}
		var idx int
		if _, err := fmt.Sscan(raw, &idx); err != nil || idx < 1 || idx > len(m.options) {
			fmt.Fprintf(w, "Enter a number between 1 and %d.\n", len(m.options))
			continue
		}
		val := m.options[idx-1].Val
		if m.validate != nil {
			if err := m.validate(val); err != nil {
				fmt.Fprintln(w, "Error:", err)
				continue
			}
		}
		if m.value != nil {
			*m.value = val
		}
		return nil
	}
}

func (m *MultiButton[T]) Skip() bool { return false }
func (m *MultiButton[T]) Zoom() bool { return false }

func (m *MultiButton[T]) KeyBinds() []key.Binding {
	return []key.Binding{
		key.NewBinding(key.WithKeys("left", "h"), key.WithHelp("←/h", "prev")),
		key.NewBinding(key.WithKeys("right", "l"), key.WithHelp("→/l", "next")),
		key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "select")),
	}
}

func (m *MultiButton[T]) WithTheme(t huh.Theme) huh.Field            { m.theme = t; return m }
func (m *MultiButton[T]) WithKeyMap(_ *huh.KeyMap) huh.Field         { return m }
func (m *MultiButton[T]) WithWidth(w int) huh.Field                  { m.width = w; return m }
func (m *MultiButton[T]) WithHeight(_ int) huh.Field                 { return m }
func (m *MultiButton[T]) WithPosition(_ huh.FieldPosition) huh.Field { return m }

func (m *MultiButton[T]) GetKey() string { return m.fieldKey }
func (m *MultiButton[T]) GetValue() any {
	if m.value == nil || len(m.options) == 0 {
		return nil
	}
	return *m.value
}
