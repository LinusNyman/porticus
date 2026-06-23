package calendar_test

import (
	"strings"
	"testing"
	"time"

	"github.com/LinusNyman/porticus"
	"github.com/LinusNyman/porticus/calendar"
	tea "github.com/charmbracelet/bubbletea"
)

func key(s string) tea.KeyMsg { return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(s)} }

func TestDaysInMonth(t *testing.T) {
	cases := []struct {
		y int
		m time.Month
		d int
	}{
		{2026, time.February, 28},
		{2024, time.February, 29}, // leap
		{2026, time.April, 30},
		{2026, time.January, 31},
	}
	for _, c := range cases {
		if got := calendar.DaysInMonth(c.y, c.m); got != c.d {
			t.Errorf("DaysInMonth(%d, %v) = %d, want %d", c.y, c.m, got, c.d)
		}
	}
}

func TestAddMonthClamped(t *testing.T) {
	jan31 := time.Date(2026, 1, 31, 0, 0, 0, 0, time.UTC)
	got := calendar.AddMonthClamped(jan31, 1)
	if got.Month() != time.February || got.Day() != 28 {
		t.Errorf("Jan 31 + 1mo = %v, want 2026-02-28", got.Format("2006-01-02"))
	}
}

func TestMove(t *testing.T) {
	start := time.Date(2026, 6, 15, 0, 0, 0, 0, time.UTC)
	g := calendar.New(start)

	if !g.Move(key("l")) || g.Selected().Day() != 16 {
		t.Errorf("l should step +1 day, got %v", g.Selected())
	}
	if !g.Move(key("j")) || g.Selected().Day() != 23 {
		t.Errorf("j should step +7 days, got %v", g.Selected())
	}
	if !g.Move(key("]")) || g.Selected().Month() != time.July {
		t.Errorf("] should advance a month, got %v", g.Selected())
	}
	// An unrelated key is not handled.
	if g.Move(key("z")) {
		t.Error("z should not be handled by the grid")
	}
}

func TestViewLayout(t *testing.T) {
	g := calendar.New(time.Date(2026, 6, 15, 0, 0, 0, 0, time.UTC))
	s := porticus.NewStyles("#e06474")
	out := g.View(s, 70, func(d time.Time) string {
		if d.Day() == 15 {
			return "2●"
		}
		return ""
	})
	if !strings.Contains(out, "June 2026") {
		t.Errorf("month title missing from grid:\n%s", out)
	}
	for _, wd := range []string{"Mon", "Tue", "Sun"} {
		if !strings.Contains(out, wd) {
			t.Errorf("weekday header %q missing", wd)
		}
	}
	if !strings.Contains(out, "2●") {
		t.Error("per-day marker not rendered")
	}
}
