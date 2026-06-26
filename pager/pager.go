// Package pager is the pantheon suite's shared scroll state for read-only
// screens: a Pager owns the scroll offset over a block of text and renders it
// through porticus.Styles.Page, handling the suite's navigation keys and the
// mouse wheel. It is the line-scroll analogue of browse.Cursor (which owns item
// selection) — drop one into any tool that shows a markdown preview, a changelog,
// or a long answer, and scrolling behaves identically everywhere.
//
// Presentation + interaction only: depends on bubbletea, porticus/keys, and the
// porticus chrome — never on the data spine.
package pager

import (
	"strings"

	"github.com/LinusNyman/porticus"
	"github.com/LinusNyman/porticus/keys"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

// Pager holds the scroll offset over a body of text. The zero value is an empty,
// top-of-page pager ready for SetContent. It remembers the last viewport height
// it was given (by Handle or View) so the mouse wheel — which carries no height —
// can still clamp to the bottom.
type Pager struct {
	body   string
	scroll int
	height int // last height seen, so HandleMouse can clamp the wheel
}

// SetContent replaces the text and clamps the offset to the new content, so a
// shorter body can't leave the page scrolled past its end. Use Top to jump back
// to the first row explicitly.
func (p *Pager) SetContent(body string) {
	p.body = body
	p.clamp()
}

// Top resets the scroll to the first row.
func (p *Pager) Top() { p.scroll = 0 }

// Handle applies a navigation key and reports whether the offset moved: j/k step
// a line, ctrl+d/ctrl+u a half page, g/G jump to top/bottom. height is the page
// height so the half-page step and the bottom stop match what View renders. Keys
// it doesn't recognise leave the pager unchanged and return false, so the caller
// can route them.
func (p *Pager) Handle(msg tea.KeyMsg, km keys.Map, height int) bool {
	p.height = height
	prev := p.scroll
	switch {
	case key.Matches(msg, km.Down):
		p.scroll++
	case key.Matches(msg, km.Up):
		p.scroll--
	case key.Matches(msg, km.HalfDown):
		p.scroll += keys.PageStep(height)
	case key.Matches(msg, km.HalfUp):
		p.scroll -= keys.PageStep(height)
	case key.Matches(msg, km.Top):
		p.scroll = 0
	case key.Matches(msg, km.Bottom):
		p.scroll = p.maxScroll()
	default:
		return false
	}
	p.clamp()
	return p.scroll != prev
}

// HandleMouse scrolls one line per wheel notch and reports whether it moved.
// Other mouse events return false. It clamps against the last height the pager
// saw (from Handle or View), mirroring browse.Cursor.HandleMouse.
func (p *Pager) HandleMouse(msg tea.MouseMsg) bool {
	prev := p.scroll
	switch msg.Button {
	case tea.MouseButtonWheelUp:
		p.scroll--
	case tea.MouseButtonWheelDown:
		p.scroll++
	default:
		return false
	}
	p.clamp()
	return p.scroll != prev
}

// View renders the page at the current offset through the standard read-only
// chrome (the header label over the ══ rule, the windowed body, a scroll hint
// when it overflows), filled to exactly width×height so the footer stays pinned.
// It records height so a later wheel event clamps correctly.
func (p *Pager) View(s porticus.Styles, t porticus.Theme, label string, width, height int) string {
	p.height = height
	p.clamp()
	return s.Page(t, label, p.body, width, height, p.scroll)
}

// Scroll is the current offset (the first body row shown).
func (p *Pager) Scroll() int { return p.scroll }

// maxScroll is the largest offset that still shows content: the total line count
// minus the rows Page renders at the current height (0 when everything fits). It
// shares porticus.PageRows with Page so the clamp can't drift from the render.
func (p *Pager) maxScroll() int {
	total := strings.Count(p.body, "\n") + 1
	if m := total - porticus.PageRows(p.height, total); m > 0 {
		return m
	}
	return 0
}

func (p *Pager) clamp() {
	if max := p.maxScroll(); p.scroll > max {
		p.scroll = max
	}
	if p.scroll < 0 {
		p.scroll = 0
	}
}
