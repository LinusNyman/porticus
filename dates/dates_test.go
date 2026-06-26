package dates_test

import (
	"testing"
	"time"

	"github.com/LinusNyman/porticus/dates"
)

func TestRelativeDate(t *testing.T) {
	now := time.Date(2026, 6, 23, 12, 0, 0, 0, time.Local) // a Tuesday
	iso := func(d int) string {
		return now.AddDate(0, 0, d).Format("2006-01-02")
	}
	cases := []struct {
		offset int
		want   string
	}{
		{0, "today"},
		{1, "tmrw"},
		{-1, "yest"},
		{3, now.AddDate(0, 0, 3).Weekday().String()[:3]}, // within the week → weekday name
		{10, "in 10d"},
		{-10, "10d ago"},
		{60, "in 2mo"},
		{-60, "2mo ago"},
		{400, "in 1y"},
	}
	for _, c := range cases {
		if got := dates.RelativeDate(iso(c.offset), now); got != c.want {
			t.Errorf("RelativeDate(offset %d) = %q, want %q", c.offset, got, c.want)
		}
	}
}

func TestRelativeDateDST(t *testing.T) {
	// Stockholm springs forward on the last Sunday of March (2026-03-29 → a 23h
	// local day) and falls back on the last Sunday of October (2026-10-25 → 25h).
	// A truncating day-count drifts by one across the short day; rounding holds.
	loc, err := time.LoadLocation("Europe/Stockholm")
	if err != nil {
		t.Skip("no tzdata for Europe/Stockholm")
	}
	defer func(orig *time.Location) { time.Local = orig }(time.Local)
	time.Local = loc

	cases := []struct {
		nowDay string // "now" is noon on this local day
		offset int
		want   string
	}{
		{"2026-03-29", 1, "tmrw"},  // spring-forward day was the bug: read "today"
		{"2026-03-30", -1, "yest"}, // the morning after, looking back over it
		{"2026-10-25", 1, "tmrw"},  // fall-back (25h) stays correct too
		{"2026-06-01", 1, "tmrw"},  // control: no transition
	}
	for _, c := range cases {
		now, perr := time.ParseInLocation("2006-01-02", c.nowDay, time.Local)
		if perr != nil {
			t.Fatal(perr)
		}
		now = now.Add(12 * time.Hour)
		target := now.AddDate(0, 0, c.offset).Format("2006-01-02")
		if got := dates.RelativeDate(target, now); got != c.want {
			t.Errorf("from %s noon, %s = %q, want %q", c.nowDay, target, got, c.want)
		}
	}
}

func TestRelativeDatePassThrough(t *testing.T) {
	now := time.Now()
	if got := dates.RelativeDate("", now); got != "" {
		t.Errorf("empty should pass through, got %q", got)
	}
	if got := dates.RelativeDate("not-a-date", now); got != "not-a-date" {
		t.Errorf("unparseable should pass through, got %q", got)
	}
}

func TestFarOff(t *testing.T) {
	now := time.Date(2026, 6, 23, 0, 0, 0, 0, time.Local)
	near := now.AddDate(0, 0, 10).Format("2006-01-02")
	far := now.AddDate(0, 0, 60).Format("2006-01-02")
	if dates.FarOff(near, now) {
		t.Error("10 days out should not be far off")
	}
	if !dates.FarOff(far, now) {
		t.Error("60 days out should be far off")
	}
}

func TestCoarseSpan(t *testing.T) {
	if n, u := dates.CoarseSpan(60); n != 2 || u != "mo" {
		t.Errorf("CoarseSpan(60) = %d%s, want 2mo", n, u)
	}
	if n, u := dates.CoarseSpan(-400); n != 1 || u != "y" {
		t.Errorf("CoarseSpan(-400) = %d%s, want 1y", n, u)
	}
}
