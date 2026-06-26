package porticus

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// windowLines windows a slice of pre-rendered lines to at most max rows around a
// scroll offset, returning the visible lines and how many are hidden above and
// below. When everything fits, scroll is ignored and above/below are 0. It is
// the shared scroll-offset windowing behind Page and the help plaque (the
// cursor-tracking variant is browse.Cursor.Window).
func windowLines(lines []string, max, scroll int) (visible []string, above, below int) {
	if max < 1 {
		max = 1
	}
	if len(lines) <= max {
		return lines, 0, 0
	}
	start := scroll
	if last := len(lines) - max; start > last {
		start = last
	}
	if start < 0 {
		start = 0
	}
	return lines[start : start+max], start, len(lines) - start - max
}

// PageRows reports how many body rows Styles.Page shows for a page of the given
// height holding total content lines: height-2 (the header line and its ══ rule),
// less one row reserved for the scroll hint when the content overflows. Never
// below 1. It is the single source of the page row budget — Page windows by it,
// and pager.Pager clamps its scroll offset to it so the two can't drift.
func PageRows(height, total int) int {
	avail := height - 2 // the header line and its ══ rule
	if avail < 1 {
		avail = 1
	}
	if total > avail && avail > 1 {
		avail-- // reserve a row for the scroll hint
	}
	return avail
}

// Page renders a scrollable read-only page through the standard chrome: the left
// header (sigil + spaced-caps tool name + ❧ + label) over the ══ rule, then body
// windowed to the remaining height and scrolled by scroll (the caller owns the
// offset, as with HelpPage), a scroll hint when it overflows, filled to exactly
// width×height so the footer stays pinned. Use it for any custom read-only
// screen — a markdown preview, a stats page; insights.InsightsPage is a thin
// wrapper over it, and pager.Pager wraps it with scroll state.
func (s Styles) Page(t Theme, label, body string, width, height, scroll int) string {
	lines := strings.Split(body, "\n")
	win := PageRows(height, len(lines))
	var above, below int
	if len(lines) > win {
		lines, above, below = windowLines(lines, win, scroll)
	}
	parts := []string{s.LeftHeader(t, label, width), strings.Join(lines, "\n")}
	if hint := s.ScrollHint(above, below, width); hint != "" {
		parts = append(parts, hint)
	}
	return lipgloss.Place(width, height, lipgloss.Left, lipgloss.Top, strings.Join(parts, "\n"))
}
