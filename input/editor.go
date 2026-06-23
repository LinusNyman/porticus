package input

import (
	"github.com/LinusNyman/porticus"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// MaxRows caps how tall the free-text editor grows before it scrolls instead.
const MaxRows = 6

// Editor is the suite's free-text input for add/edit modes: a soft-wrapping
// textarea, stripped of decorative chrome, that grows with its content up to
// MaxRows then scrolls, with a `› ` prompt on the first row and a matching blank
// indent on wrap rows so a long entry reads as one hanging-indented paragraph.
// Enter commits (no newline is inserted); esc cancels.
type Editor struct {
	ta textarea.Model
}

// NewEditor builds an editor configured to the suite convention. Call SetWidth
// on resize and Open to begin entry.
func NewEditor() Editor {
	ta := textarea.New()
	ta.CharLimit = 280
	ta.ShowLineNumbers = false
	ta.MaxHeight = MaxRows
	// Strip the textarea's decorative chrome so it reads like a plain line that
	// happens to wrap (no line-number gutter, cursor-line tint, or prompt tint).
	plain := lipgloss.NewStyle()
	ta.FocusedStyle.Base, ta.BlurredStyle.Base = plain, plain
	ta.FocusedStyle.CursorLine, ta.BlurredStyle.CursorLine = plain, plain
	ta.FocusedStyle.Prompt, ta.BlurredStyle.Prompt = plain, plain
	ta.SetPromptFunc(2, func(line int) string {
		if line == 0 {
			return "› "
		}
		return "  "
	})
	ta.SetWidth(porticus.NarrowWidth)
	ta.SetHeight(1)
	return Editor{ta: ta}
}

// SetWidth sizes the editor (the footer width); call it on a window-size change.
func (e *Editor) SetWidth(w int) {
	if w > 0 {
		e.ta.SetWidth(w)
		e.sync()
	}
}

// Open focuses the editor with placeholder shown while empty and an initial
// value (empty for add, the existing text for edit), cursor at the end. Return
// its command so the cursor blinks.
func (e *Editor) Open(placeholder, value string) tea.Cmd {
	e.ta.Placeholder = placeholder
	e.ta.SetValue(value)
	e.ta.CursorEnd()
	cmd := e.ta.Focus()
	e.sync()
	return cmd
}

// Update feeds a key: esc cancels, enter commits (no newline), anything else
// edits and re-sizes the editor to its content. On Cancelled the editor blurs;
// on Committed read Value() and then either Clear() (sticky add) or stop.
func (e *Editor) Update(msg tea.KeyMsg) (Action, tea.Cmd) {
	switch msg.String() {
	case "esc":
		e.ta.Blur()
		return Cancelled, nil
	case "enter":
		return Committed, nil
	}
	var cmd tea.Cmd
	e.ta, cmd = e.ta.Update(msg)
	e.sync()
	return Editing, cmd
}

// Clear resets to an empty value while staying focused — for a sticky add that
// keeps the editor open after each commit. Return its command so the cursor
// keeps blinking.
func (e *Editor) Clear() tea.Cmd {
	e.ta.SetValue("")
	e.sync()
	return e.ta.Focus()
}

// sync resizes the textarea to its content height (min 1, capped by MaxHeight).
func (e *Editor) sync() {
	h := e.ta.LineInfo().Height
	if h < 1 {
		h = 1
	}
	e.ta.SetHeight(h)
}

// Value is the current text.
func (e Editor) Value() string { return e.ta.Value() }

// View renders the editor (place it in the footer area).
func (e Editor) View() string { return e.ta.View() }
