// Package calendar is the suite's shared month calendar: a selected day plus the
// navigation and rendering of the month it falls in. It owns the fiddly parts —
// the date arithmetic, the Mon–Sun grid layout, the bordered cells, and the
// selected/today highlight rules — so the suite's calendar tool (calendarium)
// and any date view (pensum's calendar screen, an album birthday view) share one
// implementation. The tool keeps the day's detail list (e.g. via browse.Cursor)
// and supplies a per-day marker.
//
// Spine-free: depends on bubbletea, lipgloss, and the porticus chrome only.
package calendar

import (
	"fmt"
	"strings"
	"time"

	"github.com/LinusNyman/porticus"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var weekdayHeads = []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"}

// Grid is a month calendar centred on a selected day.
type Grid struct {
	sel time.Time
}

// New builds a Grid on sel (a zero sel starts on today).
func New(sel time.Time) Grid {
	if sel.IsZero() {
		sel = time.Now()
	}
	return Grid{sel: sel}
}

// Selected is the highlighted day.
func (g Grid) Selected() time.Time { return g.sel }

// SetSelected moves the selection (ignored for the zero time).
func (g *Grid) SetSelected(t time.Time) {
	if !t.IsZero() {
		g.sel = t
	}
}

// Move applies a calendar navigation key and reports whether the selection
// changed: h/← and l/→ step a day, k/↑ and j/↓ a week, `[` and `]` a month
// (clamping the day-of-month), `t` jumps to today. Other keys are ignored
// (handled=false) so the caller can route them (tab to the day list, etc.).
func (g *Grid) Move(msg tea.KeyMsg) bool {
	prev := g.sel
	switch msg.String() {
	case "left", "h":
		g.sel = g.sel.AddDate(0, 0, -1)
	case "right", "l":
		g.sel = g.sel.AddDate(0, 0, 1)
	case "up", "k":
		g.sel = g.sel.AddDate(0, 0, -7)
	case "down", "j":
		g.sel = g.sel.AddDate(0, 0, 7)
	case "[":
		g.sel = AddMonthClamped(g.sel, -1)
	case "]":
		g.sel = AddMonthClamped(g.sel, 1)
	case "t":
		g.sel = time.Now()
	default:
		return false
	}
	return !g.sel.Equal(prev)
}

// View renders the month grid: the "January 2006" title, the Mon–Sun weekday
// header, then the bordered day cells (7 per week row, only the weeks the month
// spans). cells are width/7 wide. marker returns a short plain badge for a day
// (e.g. a count "3●"), or "" for none, and may be nil; the badge is styled in the
// suite count colour and right-aligned in the cell. The selected day wears the
// accent border, today the aegean-blue border — and the selection wins when they
// coincide.
func (g Grid) View(s porticus.Styles, width int, marker func(day time.Time) string) string {
	now := time.Now()
	sel := g.sel
	if sel.IsZero() {
		sel = now
	}
	loc := sel.Location()
	curMonth := sel.Month()

	// First grid cell = the Monday on or before the 1st of the month.
	first := time.Date(sel.Year(), curMonth, 1, 0, 0, 0, 0, loc)
	offset := (int(first.Weekday()) + 6) % 7 // Mon=0 … Sun=6
	gridStart := first.AddDate(0, 0, -offset)
	numWeeks := (offset + DaysInMonth(sel.Year(), curMonth) + 6) / 7

	cw := width / 7
	if cw < 6 {
		cw = 6
	}
	inner := cw - 2

	todayStr := now.Format("2006-01-02")
	selStr := sel.Format("2006-01-02")

	heads := make([]string, len(weekdayHeads))
	for i, wd := range weekdayHeads {
		heads[i] = s.Dim.Render(porticus.Center(wd, cw))
	}

	weekRows := make([]string, 0, numWeeks)
	d := gridStart
	for w := 0; w < numWeeks; w++ {
		cells := make([]string, 7)
		for c := 0; c < 7; c++ {
			mk := ""
			if marker != nil {
				mk = marker(d)
			}
			cells[c] = cell(s, d, inner, curMonth, todayStr, selStr, mk)
			d = d.AddDate(0, 0, 1)
		}
		weekRows = append(weekRows, lipgloss.JoinHorizontal(lipgloss.Top, cells...))
	}

	lines := []string{s.Title.Render(sel.Format("January 2006")), strings.Join(heads, "")}
	lines = append(lines, weekRows...)
	return strings.Join(lines, "\n")
}

// cell renders one day as a bordered box: day number left, marker badge right.
func cell(s porticus.Styles, d time.Time, inner int, curMonth time.Month, todayStr, selStr, marker string) string {
	if inner < 3 {
		inner = 3
	}
	day := fmt.Sprintf("%d", d.Day())
	numStyle := s.Name
	if d.Month() != curMonth {
		numStyle = s.Dim // dim days spilling in from adjacent months
	}
	gap := inner - lipgloss.Width(day) - lipgloss.Width(marker)
	if gap < 1 {
		gap = 1
	}
	content := numStyle.Render(day) + strings.Repeat(" ", gap) + s.Count.Render(marker)

	style := lipgloss.NewStyle().Width(inner).Height(1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(porticus.ColDivider))
	ds := d.Format("2006-01-02")
	// The focused day takes precedence over today: when today is selected it
	// wears the accent focus border, not the "today" border.
	switch {
	case ds == selStr:
		style = style.BorderForeground(lipgloss.Color(s.Accent)).Bold(true)
	case ds == todayStr:
		style = style.BorderForeground(lipgloss.Color(porticus.ColToday)).Bold(true)
	}
	return style.Render(content)
}

// AddMonthClamped shifts t by delta months, clamping the day-of-month to the
// target month's length (Jan 31 → Feb 28/29) so navigation never rolls over.
func AddMonthClamped(t time.Time, delta int) time.Time {
	first := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location()).AddDate(0, delta, 0)
	d := t.Day()
	if last := DaysInMonth(first.Year(), first.Month()); d > last {
		d = last
	}
	return time.Date(first.Year(), first.Month(), d, 0, 0, 0, 0, t.Location())
}

// DaysInMonth returns the number of days in the given month (day 0 of the next
// month is the last day of this one).
func DaysInMonth(year int, mo time.Month) int {
	return time.Date(year, mo+1, 0, 0, 0, 0, 0, time.UTC).Day()
}
