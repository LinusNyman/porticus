// Package dates is the suite's shared relative-date vocabulary: it renders an
// ISO date (YYYY-MM-DD) the way every tool should show it at a glance — "today",
// "tmrw", "in 3d", "in 2mo" — so due dates, agendas, calendars, and birthdays
// read identically across the suite. Pure presentation: no spine, no bubbletea.
package dates

import (
	"fmt"
	"math"
	"time"
)

// daysUntil returns the signed number of calendar days from now to the ISO date
// (YYYY-MM-DD), and whether it parsed. Both endpoints are taken at local midnight
// and the elapsed span is *rounded* to the nearest day: a DST transition makes a
// local day 23 or 25 hours long, so truncating would shift the count by one (e.g.
// label tomorrow "today" on the spring-forward day). Rounding absorbs the skew.
func daysUntil(iso string, now time.Time) (int, bool) {
	d, err := time.Parse("2006-01-02", iso)
	if err != nil {
		return 0, false
	}
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	day := time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, time.Local)
	return int(math.Round(day.Sub(today).Hours() / 24)), true
}

// RelativeDate renders an ISO date (YYYY-MM-DD) relative to now, for at-a-glance
// scanning: "today", "tmrw", "yest", a weekday name within the coming week,
// "in Nd" / "Nd ago" within a few weeks, and coarser "in Nmo" / "in Ny" labels
// (with the symmetric "ago" forms) further out. An empty or unparseable value is
// returned unchanged so callers can pass raw fields through safely.
func RelativeDate(iso string, now time.Time) string {
	days, ok := daysUntil(iso, now)
	if !ok {
		return iso
	}

	switch {
	case days == 0:
		return "today"
	case days == 1:
		return "tmrw"
	case days == -1:
		return "yest"
	case days > 1 && days <= 6:
		// Within the coming week, the weekday name reads fastest. AddDate walks
		// whole calendar days, so the weekday is correct across a DST transition.
		return now.AddDate(0, 0, days).Weekday().String()[:3]
	case days >= 2 && days <= 27:
		return fmt.Sprintf("in %dd", days)
	case days < -1 && days >= -27:
		return fmt.Sprintf("%dd ago", -days)
	}
	// Beyond a month, coarsen to whole months ("mo"), then years ("y"), so the
	// label stays short; callers that need the exact day show the ISO date too.
	n, unit := CoarseSpan(days)
	if days > 0 {
		return fmt.Sprintf("in %d%s", n, unit)
	}
	return fmt.Sprintf("%d%s ago", n, unit)
}

// FarOff reports whether a date lands outside the day-granular window — one that
// RelativeDate renders as a coarse month/year span. Callers with room for the
// exact day (e.g. a detail pane) append the ISO date for these.
func FarOff(iso string, now time.Time) bool {
	days, ok := daysUntil(iso, now)
	if !ok {
		return false
	}
	return days > 27 || days < -27
}

// CoarseSpan rounds a (signed) day count to a whole number of months under a
// year, or years beyond, returning the magnitude and its unit suffix.
func CoarseSpan(days int) (int, string) {
	if days < 0 {
		days = -days
	}
	if days < 365 {
		mo := (days + 15) / 30 // round to the nearest month
		if mo < 1 {
			mo = 1
		}
		return mo, "mo"
	}
	return (days + 182) / 365, "y" // round to the nearest year
}
