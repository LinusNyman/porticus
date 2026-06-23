package input

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// Field is the suite's single-line text input — for dates, codes, names,
// descriptions, and other short values — sharing the editor's enter-commits /
// esc-cancels grammar and the `› ` prompt. (Live-filtered search lives in
// porticus/pick; this is for committed single values.)
type Field struct {
	ti textinput.Model
}

// NewField builds a single-line field configured to the suite convention.
func NewField() Field {
	ti := textinput.New()
	ti.Prompt = "› "
	ti.CharLimit = 280
	return Field{ti: ti}
}

// Open focuses the field with placeholder shown while empty and an initial value
// (empty, or a prefill to edit), cursor at the end. Return its command so the
// cursor blinks.
func (f *Field) Open(placeholder, value string) tea.Cmd {
	f.ti.Placeholder = placeholder
	f.ti.SetValue(value)
	f.ti.CursorEnd()
	return f.ti.Focus()
}

// Update feeds a key: esc cancels, enter commits, anything else edits.
func (f *Field) Update(msg tea.KeyMsg) (Action, tea.Cmd) {
	switch msg.String() {
	case "esc":
		f.ti.Blur()
		return Cancelled, nil
	case "enter":
		f.ti.Blur()
		return Committed, nil
	}
	var cmd tea.Cmd
	f.ti, cmd = f.ti.Update(msg)
	return Editing, cmd
}

// Value is the current text.
func (f Field) Value() string { return f.ti.Value() }

// View renders the field (place it in the footer area).
func (f Field) View() string { return f.ti.View() }
