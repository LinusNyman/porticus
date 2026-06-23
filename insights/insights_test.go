package insights

import (
	"strings"
	"testing"
	"time"
)

// ── helpers ──────────────────────────────────────────────────────────────────

func dayValues(floats ...float64) []DayValue {
	base := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	out := make([]DayValue, len(floats))
	for i, f := range floats {
		v := f
		out[i] = DayValue{Date: base.AddDate(0, 0, i), Value: &v}
	}
	return out
}

func dayValuesWithNils(vals ...*float64) []DayValue {
	base := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	out := make([]DayValue, len(vals))
	for i, v := range vals {
		out[i] = DayValue{Date: base.AddDate(0, 0, i), Value: v}
	}
	return out
}

func dayValuesFrom(start time.Time, vals ...*float64) []DayValue {
	out := make([]DayValue, len(vals))
	for i, v := range vals {
		out[i] = DayValue{Date: start.AddDate(0, 0, i), Value: v}
	}
	return out
}

func fp(f float64) *float64 { return &f }

// stripANSI removes ANSI escape sequences so plain runes can be counted.
func stripANSI(s string) string {
	var out strings.Builder
	inEsc := false
	for _, r := range s {
		if r == '\x1b' {
			inEsc = true
			continue
		}
		if inEsc {
			if r == 'm' {
				inEsc = false
			}
			continue
		}
		out.WriteRune(r)
	}
	return out.String()
}

var weekStart = time.Date(2026, 1, 5, 0, 0, 0, 0, time.UTC) // a Monday

// ── Compute ──────────────────────────────────────────────────────────────────

func TestComputeEmpty(t *testing.T) {
	tr := Compute(nil)
	if tr.Today != nil || tr.Avg7d != nil || tr.Direction != Flat {
		t.Errorf("empty Compute = %+v, want zero/Flat", tr)
	}
}

func TestComputeSingleValue(t *testing.T) {
	tr := Compute(dayValues(3.0))
	if tr.Today == nil || *tr.Today != 3.0 {
		t.Errorf("Today = %v, want 3.0", tr.Today)
	}
	if tr.Avg7d != nil {
		t.Errorf("single value should have no average, got %v", *tr.Avg7d)
	}
}

func TestComputeDirection(t *testing.T) {
	if d := Compute(dayValues(3, 4, 3, 4, 3, 4, 3, 4, 3, 7)).Direction; d != Up {
		t.Errorf("rising series Direction = %v, want Up", d)
	}
	if d := Compute(dayValues(100, 100, 100, 100, 100, 100, 100, 100, 100, 10)).Direction; d != Down {
		t.Errorf("falling series Direction = %v, want Down", d)
	}
	if d := Compute(dayValues(5, 5, 5, 5, 5, 5, 5, 5)).Direction; d != Flat {
		t.Errorf("steady series Direction = %v, want Flat", d)
	}
}

func TestComputeRatingThreshold(t *testing.T) {
	// Ratings (avg < 5): absolute 0.5 threshold.
	if d := Compute(dayValues(3, 3, 3, 3, 3, 3, 3, 3, 3, 3.4)).Direction; d != Flat {
		t.Errorf("0.4 diff should be Flat, got %v", d)
	}
	if d := Compute(dayValues(3, 3, 3, 3, 3, 3, 3, 3, 3, 3.6)).Direction; d != Up {
		t.Errorf("0.6 diff should be Up, got %v", d)
	}
}

func TestComputeAvg7dWindow(t *testing.T) {
	base := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	vals := make([]DayValue, 30)
	for i := 0; i < 23; i++ {
		v := 10.0
		vals[i] = DayValue{Date: base.AddDate(0, 0, i), Value: &v}
	}
	for i := 23; i < 30; i++ {
		v := 20.0
		vals[i] = DayValue{Date: base.AddDate(0, 0, i), Value: &v}
	}
	tr := Compute(vals)
	if tr.Avg7d == nil || *tr.Avg7d < 19.99 || *tr.Avg7d > 20.01 {
		t.Errorf("avg7d = %v, want 20 (last 7 days)", tr.Avg7d)
	}
	if tr.Avg30d == nil || *tr.Avg30d <= 10.0 || *tr.Avg30d >= 20.0 {
		t.Errorf("avg30d = %v, want between 10 and 20", tr.Avg30d)
	}
}

func TestComputeNilToday(t *testing.T) {
	tr := Compute(dayValuesWithNils(fp(3.0), fp(4.0), nil))
	if tr.Today != nil || tr.Direction != Flat {
		t.Errorf("nil today should give nil/Flat, got %+v", tr)
	}
}

// ── Sparkline ────────────────────────────────────────────────────────────────

func TestSparklineWidth(t *testing.T) {
	plain := stripANSI(Sparkline(dayValues(1, 2, 3, 4, 5), 10))
	if n := len([]rune(plain)); n != 10 {
		t.Errorf("sparkline width = %d runes, want 10", n)
	}
}

func TestSparklineAllNil(t *testing.T) {
	plain := stripANSI(Sparkline(dayValuesWithNils(nil, nil, nil), 5))
	if plain != strings.Repeat(" ", 5) {
		t.Errorf("all-nil sparkline = %q, want 5 spaces", plain)
	}
}

func TestSparklineMonotoneIncreasing(t *testing.T) {
	runes := []rune(stripANSI(Sparkline(dayValues(1, 2, 3, 4, 5, 6, 7, 8), 8)))
	for i := 0; i < len(runes)-1; i++ {
		if runes[i] > runes[i+1] {
			t.Errorf("rune %d (%c) > rune %d (%c); should be non-decreasing", i, runes[i], i+1, runes[i+1])
		}
	}
}

func TestSparklineConstant(t *testing.T) {
	runes := []rune(stripANSI(Sparkline(dayValues(5, 5, 5, 5), 4)))
	for i := 1; i < len(runes); i++ {
		if runes[i] != runes[0] {
			t.Errorf("constant values should all map to the same char, got %q", string(runes))
			break
		}
	}
}

// ── Heatmap ──────────────────────────────────────────────────────────────────

func TestHeatmapWidth(t *testing.T) {
	plain := stripANSI(Heatmap(dayValues(1, 0, 1, 1, 0), 10))
	if n := len([]rune(plain)); n != 10 {
		t.Errorf("heatmap width = %d runes, want 10", n)
	}
}

func TestHeatmapSymbols(t *testing.T) {
	plain := stripANSI(Heatmap(dayValuesWithNils(fp(1.0), fp(0.0), nil), 3))
	for _, sym := range []string{"▓", "░", "·"} {
		if !strings.Contains(plain, sym) {
			t.Errorf("heatmap %q missing %q", plain, sym)
		}
	}
}

// ── CalendarHeatmap ───────────────────────────────────────────────────────────

func TestCalendarHeatmapEmpty(t *testing.T) {
	if CalendarHeatmap(nil) != "" || CalendarHeatmap([]DayValue{}) != "" {
		t.Error("empty CalendarHeatmap should be the empty string")
	}
}

func TestCalendarHeatmapFullWeek(t *testing.T) {
	out := stripANSI(CalendarHeatmap(dayValuesFrom(weekStart, fp(1), fp(1), fp(1), fp(1), fp(1), fp(1), fp(1))))
	lines := strings.Split(out, "\n")
	if len(lines) != 4 {
		t.Fatalf("got %d lines, want 4 (header, marks, blank, hit rate)", len(lines))
	}
	if lines[0] != "  M T W T F S S" {
		t.Errorf("header = %q", lines[0])
	}
	if lines[1] != "  ✓ ✓ ✓ ✓ ✓ ✓ ✓" {
		t.Errorf("marks = %q", lines[1])
	}
	if lines[3] != "  Hit rate: 100%" {
		t.Errorf("hit rate = %q", lines[3])
	}
}

func TestCalendarHeatmapHitRateRoundsDown(t *testing.T) {
	out := stripANSI(CalendarHeatmap(dayValuesFrom(weekStart, fp(1), fp(1), fp(1), fp(0), fp(0), fp(0), fp(0))))
	if !strings.Contains(out, "Hit rate: 42%") {
		t.Errorf("3 of 7 should round down to 42%%, got %q", out)
	}
}

func TestCalendarHeatmapNilsExcluded(t *testing.T) {
	out := stripANSI(CalendarHeatmap(dayValuesFrom(weekStart, fp(1), nil, fp(1), nil, fp(0), nil, nil)))
	lines := strings.Split(out, "\n")
	if lines[1] != "  ✓ · ✓ · ✗ · ·" {
		t.Errorf("marks = %q", lines[1])
	}
	if !strings.Contains(out, "Hit rate: 66%") {
		t.Errorf("2 of 3 non-nil should be 66%%, got %q", out)
	}
}

func TestCalendarHeatmapStartsMidWeek(t *testing.T) {
	thu := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC) // Thursday
	out := stripANSI(CalendarHeatmap(dayValuesFrom(thu, fp(1), fp(0), fp(1))))
	lines := strings.Split(out, "\n")
	if lines[0] != "  M T W T F S" {
		t.Errorf("header = %q", lines[0])
	}
	if lines[1] != "        ✓ ✗ ✓" {
		t.Errorf("marks = %q, want leading padding for Mon-Wed", lines[1])
	}
}

func TestCalendarHeatmapWeekBoundary(t *testing.T) {
	vals := make([]*float64, 14)
	for i := range vals {
		vals[i] = fp(1)
	}
	out := stripANSI(CalendarHeatmap(dayValuesFrom(weekStart, vals...)))
	lines := strings.Split(out, "\n")
	if !strings.Contains(lines[0], "S  M") {
		t.Errorf("expected double-space week boundary in header %q", lines[0])
	}
	if !strings.Contains(lines[1], "✓  ✓") {
		t.Errorf("expected double-space week boundary in marks %q", lines[1])
	}
}
