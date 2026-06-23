// Package browse is the pantheon suite's shared two-pane tree-browser scaffold:
// the stateful left-pane node tree (TreePane), the list selection cursor
// (Cursor), and the auto-clearing status line (Status). Each is embedded by a
// tool's own bubbletea Model — porticus owns the navigation grammar, geometry,
// and rendering so the tree/working screen behaves identically across tools,
// while the tool keeps its own data and working-pane content.
//
// browse depends on the data spine (github.com/LinusNyman/pantheon/tree) for
// the node type and on bubbletea for messages/commands. Per the suite rule,
// this spine-coupled, interactive layer lives in its own package, separate from
// the dependency-light porticus chrome.
package browse

import (
	"time"

	"github.com/LinusNyman/porticus"
	tea "github.com/charmbracelet/bubbletea"
)

// statusTimeout is how long a status/error line lingers before auto-clearing.
const statusTimeout = 4 * time.Second

// clearMsg asks a Status to clear itself if it is still the generation that
// scheduled the tick — a superseded message must not wipe a newer one.
type clearMsg struct{ gen int }

// Status is the transient status/error line shared by every tool: set a message
// (or error), and it auto-clears after a few seconds. The generation guard
// means a stale clear tick from an earlier message never wipes a current one.
type Status struct {
	Msg string // success/info text, rendered in laurel green
	Err string // error text, rendered in Pompeian red
	gen int    // bumped on each Set/SetErr; tags the clear tick
}

// Set shows msg (clearing any error) and returns the command that schedules its
// auto-clear. Return the command from the tool's Update so the tick runs.
func (s *Status) Set(msg string) tea.Cmd {
	s.Msg, s.Err = msg, ""
	return s.tick()
}

// SetErr shows err (clearing any status) and returns its auto-clear command.
func (s *Status) SetErr(err string) tea.Cmd {
	s.Err, s.Msg = err, ""
	return s.tick()
}

func (s *Status) tick() tea.Cmd {
	s.gen++
	g := s.gen
	return tea.Tick(statusTimeout, func(time.Time) tea.Msg { return clearMsg{gen: g} })
}

// Handle clears the line when msg is the matching (current-generation) clear
// tick and reports whether it consumed the message. Call it first in the tool's
// Update so the auto-clear works: if it returns true, return early.
func (s *Status) Handle(msg tea.Msg) bool {
	cm, ok := msg.(clearMsg)
	if !ok {
		return false
	}
	if cm.gen == s.gen {
		s.Msg, s.Err = "", ""
	}
	return true
}

// View renders the line (error takes precedence over status), truncated to
// width. Empty when there is nothing to show.
func (s Status) View(st porticus.Styles, width int) string {
	switch {
	case s.Err != "":
		return porticus.Truncate(st.Err.Render(s.Err), width)
	case s.Msg != "":
		return porticus.Truncate(st.OK.Render(s.Msg), width)
	}
	return ""
}
