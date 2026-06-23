// Package input is the suite's shared text-entry chrome: the free-text add/edit
// editor (a soft-wrapping, content-sized textarea), the single-line field (for
// dates, codes, names…), and the y/n confirm prompt. Each tool re-implemented
// these; centralising them keeps entry behaviour — the hanging-indent prompt,
// the grow-then-scroll height, enter-commits/esc-cancels — identical everywhere.
//
// Presentation + interaction only: depends on bubbletea, the bubbles text
// widgets, and the porticus chrome — never on the data spine. The tool owns what
// a commit does; these own the editing and the rendering.
package input

import (
	"github.com/LinusNyman/porticus"
	tea "github.com/charmbracelet/bubbletea"
)

// Action is the outcome of feeding a key to an Editor or Field.
type Action int

const (
	Editing   Action = iota // still editing; the widget consumed the key
	Committed               // enter pressed — read Value() and act on it
	Cancelled               // esc pressed — close the input
)

// Answer is the outcome of feeding a key to a Confirm prompt.
type Answer int

const (
	Pending Answer = iota // no decision yet
	Yes                   // y / Y / enter
	No                    // n / N / esc
)

// Confirm is a y/n confirmation line. It carries only the question; Update is a
// pure decision over the key, so a tool can hold one per pending action (delete,
// clean, …) and render it in the footer.
type Confirm struct {
	Question string // e.g. "delete this todo?"
}

// Update maps a key to an Answer: y/Y/enter confirm, n/N/esc cancel, anything
// else leaves it Pending.
func (c Confirm) Update(msg tea.KeyMsg) Answer {
	switch msg.String() {
	case "y", "Y", "enter":
		return Yes
	case "n", "N", "esc":
		return No
	}
	return Pending
}

// View renders the prompt in the suite error colour (a confirmation is a
// stop-and-think), truncated to width.
func (c Confirm) View(s porticus.Styles, width int) string {
	return porticus.Truncate(s.Err.Render(c.Question+"  — y/enter to confirm, esc to cancel"), width)
}
