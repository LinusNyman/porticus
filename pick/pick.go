// Package pick is the suite's shared search/suggest overlay: a generic,
// type-parameterised picker — a query input over a live-filtered, scrollable,
// selectable results list — rendered as a page through the porticus chrome.
//
// It unifies the two overlays every tree tool re-implements: a free-text search
// (filter a flat index, jump to a hit) and a ranked suggest (score candidates,
// cap the list, jump to one). Both are the same component; the tool supplies the
// data via callbacks (Filter, Render) and keeps ownership of what a selection
// does and of any extra per-row keys.
//
// pick is presentation + interaction only: it depends on bubbletea, the bubbles
// text input, the porticus chrome, and the canonical keymap — never on the data
// spine. The item type T may hold a *tree.Node, but pick neither knows nor cares.
package pick

import (
	"fmt"
	"strings"

	"github.com/LinusNyman/porticus"
	"github.com/LinusNyman/porticus/keys"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// Opts configures a Picker without coupling it to any tool type.
type Opts[T any] struct {
	// Label is the word shown after the ❧ in the header where a node name
	// normally sits, e.g. "search" or "suggest".
	Label string
	// Placeholder is the text-input placeholder shown while the query is empty.
	Placeholder string
	// Limit caps how many results are listed (suggest uses 8); 0 = no cap.
	Limit int
	// Filter recomputes the results for a query. Called live on each keystroke,
	// so it should match an in-memory snapshot, not walk the disk.
	Filter func(query string) []T
	// Render returns one result row, styled by the tool (code, text, badges);
	// pick truncates it and applies the selection highlight.
	Render func(it T, width int) string
	// Empty returns the message shown when there are no results, given the
	// current query. Optional — a sensible default is used when nil.
	Empty func(query string) string
}

// Picker is a search/suggest overlay over items of type T. Embed one per
// overlay in the tool's Model; drive it with Open, Update, and View.
type Picker[T any] struct {
	opts    Opts[T]
	input   textinput.Model
	query   string
	results []T
	cursor  int
	typing  bool // true while the query input has focus
}

// New builds a Picker. Call Open to start a fresh query.
func New[T any](o Opts[T]) Picker[T] {
	ti := textinput.New()
	ti.Prompt = "› "
	ti.Placeholder = o.Placeholder
	return Picker[T]{opts: o, input: ti}
}

// Open clears the query and focuses the input for a fresh search. Return its
// command from the tool's Update so the input cursor blinks.
func (p *Picker[T]) Open() tea.Cmd {
	p.query = ""
	p.results = nil
	p.cursor = 0
	p.input.SetValue("")
	p.typing = true
	return p.input.Focus()
}

// Update routes a key. While typing it edits the live query (re-filtering on
// each keystroke) until enter commits to browsing; while browsing it moves the
// cursor (j/k/g/G), refines on `/`, and selects on l/right/enter. It returns the
// chosen item on select, a command to run (input blink), and handled=false for
// keys it does not consume — esc and quit included — so the tool can close the
// overlay or add its own keys (e.g. edit/delete on a hit).
func (p *Picker[T]) Update(msg tea.KeyMsg, km keys.Map) (selected *T, cmd tea.Cmd, handled bool) {
	if p.typing {
		switch msg.String() {
		case "esc":
			return nil, nil, false // let the caller close the overlay
		case "enter":
			p.typing = false
			p.input.Blur()
			return nil, nil, true
		}
		p.input, cmd = p.input.Update(msg)
		if q := p.input.Value(); q != p.query {
			p.query = q
			p.results = p.opts.Filter(q)
			p.clamp()
		}
		return nil, cmd, true
	}

	switch {
	case key.Matches(msg, km.Up):
		p.cursor--
	case key.Matches(msg, km.Down):
		p.cursor++
	case key.Matches(msg, km.Top):
		p.cursor = 0
	case key.Matches(msg, km.Bottom):
		p.cursor = p.limit() - 1
	case key.Matches(msg, km.Search):
		return nil, p.refine(), true
	case key.Matches(msg, km.Expand): // l / right / enter — select
		if it, ok := p.Selected(); ok {
			return &it, nil, true
		}
		return nil, nil, true
	default:
		return nil, nil, false
	}
	p.clamp()
	return nil, nil, true
}

// HandleMouse moves the selection on a wheel scroll while browsing results
// (ignored while typing the query) and reports whether it moved.
func (p *Picker[T]) HandleMouse(msg tea.MouseMsg) bool {
	if p.typing {
		return false
	}
	prev := p.cursor
	switch msg.Button {
	case tea.MouseButtonWheelUp:
		p.cursor--
	case tea.MouseButtonWheelDown:
		p.cursor++
	default:
		return false
	}
	p.clamp()
	return p.cursor != prev
}

// refine reopens the query input prefilled with the last query.
func (p *Picker[T]) refine() tea.Cmd {
	p.input.SetValue(p.query)
	p.input.CursorEnd()
	p.typing = true
	return p.input.Focus()
}

// Refresh re-runs the filter for the current query — call it after a mutation or
// a show/hide toggle so the results reflect new data. Keeps the cursor in range.
func (p *Picker[T]) Refresh() {
	p.results = p.opts.Filter(p.query)
	p.clamp()
}

// Selected returns the highlighted result, or ok=false when the list is empty.
func (p Picker[T]) Selected() (T, bool) {
	var zero T
	if p.cursor < 0 || p.cursor >= p.limit() {
		return zero, false
	}
	return p.results[p.cursor], true
}

// Query is the current query string.
func (p Picker[T]) Query() string { return p.query }

func (p Picker[T]) limit() int {
	n := len(p.results)
	if p.opts.Limit > 0 && n > p.opts.Limit {
		return p.opts.Limit
	}
	return n
}

func (p *Picker[T]) clamp() {
	lim := p.limit()
	if p.cursor >= lim {
		p.cursor = lim - 1
	}
	if p.cursor < 0 {
		p.cursor = 0
	}
}

func (p Picker[T]) emptyMsg() string {
	if p.opts.Empty != nil {
		return p.opts.Empty(p.query)
	}
	if p.query == "" {
		return "type to search"
	}
	return fmt.Sprintf("no matches for %q", p.query)
}

// View renders the overlay as a page: the porticus left header (sigil + tool
// name + ❧ + label) over the ══ rule, the query line (the live input while
// typing, else a dim echo), then the windowed results with the selection
// highlighted and a scroll hint. Returns exactly height lines; the tool draws
// the footer hint bar beneath as usual.
func (p Picker[T]) View(s porticus.Styles, t porticus.Theme, width, height int) string {
	lines := strings.Split(s.LeftHeader(t, p.opts.Label, width), "\n")

	switch {
	case p.typing:
		lines = append(lines, "  "+p.input.View())
	case p.query != "":
		lines = append(lines, s.Dim.Render("  query: ")+p.query)
	}

	lim := p.limit()
	if lim == 0 {
		lines = append(lines, s.Dim.Render("  "+p.emptyMsg()))
	} else {
		maxRows := height - len(lines)
		if maxRows < 1 {
			maxRows = 1
		}
		if lim > maxRows && maxRows > 1 {
			maxRows-- // reserve a row for the scroll hint
		}
		start := 0
		if p.cursor >= maxRows {
			start = p.cursor - maxRows + 1
		}
		end := start + maxRows
		if end > lim {
			end = lim
		}
		for i := start; i < end; i++ {
			line := porticus.Truncate(p.opts.Render(p.results[i], width), width)
			if i == p.cursor {
				line = s.SelFocus.Render(porticus.PadTo(line, width))
			}
			lines = append(lines, line)
		}
		if hint := s.ScrollHint(start, lim-end, width); hint != "" {
			lines = append(lines, hint)
		}
	}

	for len(lines) < height {
		lines = append(lines, "")
	}
	return strings.Join(lines, "\n")
}
