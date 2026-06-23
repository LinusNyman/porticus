package browse

import (
	"github.com/LinusNyman/porticus/keys"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

// Cursor is a selection index over a list of Len items with the suite's shared
// navigation (j/k, g/G, ctrl+d/ctrl+u) and viewport windowing. Embed one per
// working-pane list (todos, contacts, search results, …); keep Len in sync with
// the data before handling keys or rendering.
type Cursor struct {
	Index int
	Len   int
}

// PageStep is the half-screen jump for ctrl+d/ctrl+u, derived from the body
// height exactly as the reference impl: (height-3)/2, at least 1.
func PageStep(height int) int {
	step := (height - 3) / 2
	if step < 1 {
		step = 1
	}
	return step
}

// Handle applies a navigation key to the cursor and reports whether it moved.
// pageStep is the half-screen jump (see PageStep). Keys it doesn't recognise
// leave the cursor unchanged and return false, so the caller can handle them.
func (c *Cursor) Handle(msg tea.KeyMsg, km keys.Map, pageStep int) bool {
	prev := c.Index
	switch {
	case key.Matches(msg, km.Down):
		c.Index++
	case key.Matches(msg, km.Up):
		c.Index--
	case key.Matches(msg, km.Top):
		c.Index = 0
	case key.Matches(msg, km.Bottom):
		c.Index = c.Len - 1
	case key.Matches(msg, km.HalfDown):
		c.Index += pageStep
	case key.Matches(msg, km.HalfUp):
		c.Index -= pageStep
	default:
		return false
	}
	c.clamp()
	return c.Index != prev
}

// HandleMouse moves the cursor on a wheel scroll (up/down by one) and reports
// whether it moved. Other mouse events return false. Mirrors the suite
// convention that the wheel scrolls the focused list like j/k.
func (c *Cursor) HandleMouse(msg tea.MouseMsg) bool {
	prev := c.Index
	switch msg.Button {
	case tea.MouseButtonWheelUp:
		c.Index--
	case tea.MouseButtonWheelDown:
		c.Index++
	default:
		return false
	}
	c.clamp()
	return c.Index != prev
}

// Reorder moves the selected item by delta (−1 up, +1 down) by asking move to
// perform the data swap into the new slot, then follows it with the cursor.
// move(from, to) does the move and reports success. Returns whether the item
// moved — the shared "J/K reorder" behaviour (pensum todos, decreta principles).
func (c *Cursor) Reorder(delta int, move func(from, to int) bool) bool {
	if delta == 0 {
		return false
	}
	to := c.Index + delta
	if to < 0 || to >= c.Len {
		return false
	}
	if !move(c.Index, to) {
		return false
	}
	c.Index = to
	return true
}

func (c *Cursor) clamp() {
	if c.Index >= c.Len {
		c.Index = c.Len - 1
	}
	if c.Index < 0 {
		c.Index = 0
	}
}

// Window returns the visible row range [start,end) for a viewport of maxRows
// that keeps the cursor in view, plus the number of rows hidden above and below
// (feed those to porticus.Styles.ScrollHint). Matches the reference tree/list
// windowing so every scrolling pane behaves the same.
func (c Cursor) Window(maxRows int) (start, end, above, below int) {
	if maxRows < 1 {
		maxRows = 1
	}
	if c.Index >= maxRows {
		start = c.Index - maxRows + 1
	}
	end = start + maxRows
	if end > c.Len {
		end = c.Len
	}
	above = start
	below = c.Len - end
	return start, end, above, below
}
