// Package status is the pantheon suite's shared transient status line: set a
// message — a success, a neutral notice, or an error — and it auto-clears after a
// few seconds. Every tool shows one after a mutation ("saved", "deleted", "no
// matches", "error: …"); centralising it keeps the timing, the generation guard,
// and the colour roles identical suite-wide.
//
// Presentation + interaction only: depends on bubbletea (for the auto-clear tick)
// and the porticus chrome — never on the data spine — so any tool can embed one,
// whether or not it uses the tree browser. (Lifted out of porticus/browse so a
// non-tree tool needn't pull the spine to show a status line.)
package status

import (
	"time"

	"github.com/LinusNyman/porticus"
	tea "github.com/charmbracelet/bubbletea"
)

// timeout is how long a status line lingers before auto-clearing.
const timeout = 4 * time.Second

// Kind is the colour role of a status message.
type Kind int

const (
	OK    Kind = iota // success, laurel green
	Info              // neutral notice, marble ivory
	Error             // error, Pompeian red
)

// clearMsg asks a Line to clear itself if it is still the generation that
// scheduled the tick — a superseded message must not wipe a newer one.
type clearMsg struct{ gen int }

// Line is the transient status line shared by every tool: set a message with a
// kind and it auto-clears after a few seconds. The generation guard means a stale
// clear tick from an earlier message never wipes a current one. The zero value is
// ready to use (an empty line that renders as nothing).
type Line struct {
	text string
	kind Kind
	gen  int // bumped on each set; tags the clear tick
}

// Set shows msg as a success (laurel green) and returns the command that
// schedules its auto-clear. Return the command from the tool's Update so the tick
// runs.
func (l *Line) Set(msg string) tea.Cmd { return l.flash(msg, OK) }

// SetInfo shows msg as a neutral notice (marble ivory) and returns its auto-clear
// command.
func (l *Line) SetInfo(msg string) tea.Cmd { return l.flash(msg, Info) }

// SetErr shows msg as an error (Pompeian red) and returns its auto-clear command.
func (l *Line) SetErr(msg string) tea.Cmd { return l.flash(msg, Error) }

// flash sets the text and kind, bumps the generation, and schedules the clear.
func (l *Line) flash(msg string, k Kind) tea.Cmd {
	l.text, l.kind = msg, k
	l.gen++
	g := l.gen
	return tea.Tick(timeout, func(time.Time) tea.Msg { return clearMsg{gen: g} })
}

// Handle clears the line when msg is the matching (current-generation) clear tick
// and reports whether it consumed the message. Call it first in the tool's Update
// so the auto-clear works: if it returns true, return early.
func (l *Line) Handle(msg tea.Msg) bool {
	cm, ok := msg.(clearMsg)
	if !ok {
		return false
	}
	if cm.gen == l.gen {
		l.text = ""
	}
	return true
}

// Text is the current message ("" when cleared).
func (l Line) Text() string { return l.text }

// Kind is the colour role of the current message.
func (l Line) Kind() Kind { return l.kind }

// View renders the line in its kind's colour, truncated to width. Empty when
// there is nothing to show.
func (l Line) View(s porticus.Styles, width int) string {
	if l.text == "" {
		return ""
	}
	style := s.Name // Info: marble ivory
	switch l.kind {
	case OK:
		style = s.OK
	case Error:
		style = s.Err
	}
	return porticus.Truncate(style.Render(l.text), width)
}
