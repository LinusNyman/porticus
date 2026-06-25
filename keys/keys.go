// Package keys is the pantheon suite's canonical key grammar — the single
// source of truth for which keys do what across every tool, so navigation,
// pane switching, and view selection feel identical suite-wide.
//
// A tool builds a Map (start from Default) and matches incoming key messages
// against its bindings via the bubbles key.Matches helper. The same Map also
// generates the standard help groups (HelpGroups), so the help screen can never
// drift from the bindings it documents. keys depends on bubbletea/key and the
// porticus root (for HelpGroup) but never on the data spine, so non-tree tools
// can adopt the grammar without pulling the tree code.
package keys

import (
	"strconv"

	"github.com/LinusNyman/porticus"
	"github.com/charmbracelet/bubbles/key"
)

// Map is the suite key grammar. Field bindings carry both their keys and a
// help string; the help text is curated to read well in the help plaque rather
// than auto-derived.
type Map struct {
	// Navigation — identical in the tree and every list.
	Up       key.Binding // k / up   — move selection up
	Down     key.Binding // j / down — move selection down
	Top      key.Binding // g / home — jump to first
	Bottom   key.Binding // G / end  — jump to last
	HalfUp   key.Binding // ctrl+u / pgup   — half-screen up
	HalfDown key.Binding // ctrl+d / pgdown — half-screen down

	// Tree interaction.
	Collapse key.Binding // h / left            — collapse node / go to parent
	Expand   key.Binding // l / right / enter   — expand node / cross to working pane

	// Pane focus.
	NextPane key.Binding // tab
	PrevPane key.Binding // shift+tab

	// Common item actions (verbs are tool-worded in help; keys are fixed).
	Add    key.Binding // a
	Edit   key.Binding // e
	Delete key.Binding // x

	// Find / housekeeping.
	Search key.Binding // /
	Goto   key.Binding // f
	Reload key.Binding // r
	Filter key.Binding // t

	// App.
	Help key.Binding // ?
	Quit key.Binding // q / ctrl+c

	// Title — the cover / title screen, bound to 0 across every tool, the way
	// the numbered views are bound to 1…9.
	Title key.Binding

	// Views — the number keys 1..9 select views, bound contiguously from 1
	// upward; which number maps to which view is each tool's choice (see View
	// and HelpGroups).
	Views key.Binding
}

// Default returns the canonical suite key map. Tools should use this as-is so
// the grammar stays identical; override a field only with a logged reason.
func Default() Map {
	return Map{
		Up:       key.NewBinding(key.WithKeys("k", "up"), key.WithHelp("k", "up")),
		Down:     key.NewBinding(key.WithKeys("j", "down"), key.WithHelp("j", "down")),
		Top:      key.NewBinding(key.WithKeys("g", "home"), key.WithHelp("g", "top")),
		Bottom:   key.NewBinding(key.WithKeys("G", "end"), key.WithHelp("G", "bottom")),
		HalfUp:   key.NewBinding(key.WithKeys("ctrl+u", "pgup"), key.WithHelp("ctrl+u", "half up")),
		HalfDown: key.NewBinding(key.WithKeys("ctrl+d", "pgdown"), key.WithHelp("ctrl+d", "half down")),

		Collapse: key.NewBinding(key.WithKeys("h", "left"), key.WithHelp("h", "collapse / parent")),
		Expand:   key.NewBinding(key.WithKeys("l", "right", "enter"), key.WithHelp("l", "expand / open")),

		NextPane: key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "switch pane")),
		PrevPane: key.NewBinding(key.WithKeys("shift+tab"), key.WithHelp("shift+tab", "switch pane")),

		Add:    key.NewBinding(key.WithKeys("a"), key.WithHelp("a", "add")),
		Edit:   key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "edit")),
		Delete: key.NewBinding(key.WithKeys("x"), key.WithHelp("x", "delete")),

		Search: key.NewBinding(key.WithKeys("/"), key.WithHelp("/", "search")),
		Goto:   key.NewBinding(key.WithKeys("f"), key.WithHelp("f", "go to by code")),
		Reload: key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "reload")),
		Filter: key.NewBinding(key.WithKeys("t"), key.WithHelp("t", "filter")),

		Help: key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "help")),
		Quit: key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q", "quit")),

		Title: key.NewBinding(key.WithKeys("0"), key.WithHelp("0", "title")),
		Views: key.NewBinding(
			key.WithKeys("1", "2", "3", "4", "5", "6", "7", "8", "9"),
			key.WithHelp("1…9", "select view"),
		),
	}
}

// View returns the 1-based view index a digit key selects (1 for "1" … 9 for
// "9"), or 0 if k is not a view key. Pair it with key.Matches(msg, m.Views) to
// gate, then m.View(msg.String()) to read the index.
func (m Map) View(k string) int {
	if len(k) == 1 && k[0] >= '1' && k[0] <= '9' {
		return int(k[0] - '0')
	}
	return 0
}

// HelpGroups builds the standard, suite-wide help groups so every tool's help
// screen opens with the same Navigate and View sections in the same order. The
// View section always leads with the cover screen (0:title) and is followed by
// the tool's ordered viewLabels (rendered 1:first 2:second …), enforcing the
// contiguous-from-1 convention. extra groups (the tool's own actions, worded in
// its domain) are appended after, then passed to Styles.HelpPage for rendering.
func (m Map) HelpGroups(viewLabels []string, extra ...porticus.HelpGroup) []porticus.HelpGroup {
	groups := []porticus.HelpGroup{
		{Title: "Navigate", Rows: [][2]string{
			{"j / k", "move down / up"},
			{"g / G", "top / bottom"},
			{"tab", "switch pane"},
			{"h / l", "collapse / expand"},
			{"enter", "open / jump"},
		}},
	}
	// The cover screen (0) exists in every tool, so the View group is always
	// present and leads with it; the numbered views follow (1:first 2:second …).
	rows := [][2]string{{"0", "title"}}
	for i, lbl := range viewLabels {
		rows = append(rows, [2]string{strconv.Itoa(i + 1), lbl})
	}
	groups = append(groups, porticus.HelpGroup{Title: "View", Rows: rows})
	groups = append(groups, extra...)
	return groups
}
