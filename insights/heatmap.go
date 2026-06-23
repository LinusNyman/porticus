package insights

import (
	"fmt"
	"strings"
	"time"
)

// Heatmap renders a compact binary calendar using ▓ (true), ░ (false), and ·
// (skipped/missing). Values are ordered oldest→newest and the result is exactly
// width characters wide, padded on the left with · when short.
func Heatmap(values []DayValue, width int) string {
	if width <= 0 {
		width = 30
	}
	if len(values) == 0 {
		return strings.Repeat("·", width)
	}

	// Use the last width values.
	if len(values) > width {
		values = values[len(values)-width:]
	}

	var sb strings.Builder
	for _, dv := range values {
		switch {
		case dv.Value == nil:
			sb.WriteString(heatNil.Render("·"))
		case *dv.Value >= 0.5:
			sb.WriteString(heatOn.Render("▓"))
		default:
			sb.WriteString(heatOff.Render("░"))
		}
	}

	// Pad left with neutral chars when shorter than width.
	if pad := width - len(values); pad > 0 {
		return heatNil.Render(strings.Repeat("·", pad)) + sb.String()
	}
	return sb.String()
}

// CalendarHeatmap renders a two-line week calendar: day-of-week headers above,
// ✓/✗/· marks below, grouped into Mon–Sun week columns separated by a double
// space, with a hit-rate line beneath. Returns a multi-line string, or "" when
// there are no values.
func CalendarHeatmap(values []DayValue) string {
	if len(values) == 0 {
		return ""
	}

	dayNames := []string{"M", "T", "W", "T", "F", "S", "S"}
	// isoWeekday: Mon=1 … Sun=7.
	isoWeekday := func(t time.Time) int {
		wd := int(t.Weekday())
		if wd == 0 {
			return 7
		}
		return wd
	}

	// Pad the beginning so the first day aligns to its weekday.
	startOffset := isoWeekday(values[0].Date) - 1 // 0=Mon … 6=Sun

	type cell struct {
		header string
		mark   string
	}
	cells := make([]cell, 0, startOffset+len(values))
	for i := 0; i < startOffset; i++ {
		cells = append(cells, cell{dayNames[i], " "})
	}
	for _, dv := range values {
		wd := isoWeekday(dv.Date) - 1
		var mark string
		switch {
		case dv.Value == nil:
			mark = heatNil.Render("·")
		case *dv.Value >= 0.5:
			mark = heatOn.Render("✓")
		default:
			mark = heatOff.Render("✗")
		}
		cells = append(cells, cell{dayNames[wd], mark})
	}

	var headerLine, markLine strings.Builder
	for i, c := range cells {
		weekday := i % 7
		if i > 0 && weekday == 0 {
			// Week boundary: a two-space separator.
			headerLine.WriteString("  ")
			markLine.WriteString("  ")
		} else if i > 0 {
			headerLine.WriteString(" ")
			markLine.WriteString(" ")
		}
		headerLine.WriteString(heatNil.Render(c.header))
		markLine.WriteString(c.mark)
	}

	// Hit-rate line (nil days excluded).
	var trueCount, total int
	for _, dv := range values {
		if dv.Value == nil {
			continue
		}
		total++
		if *dv.Value >= 0.5 {
			trueCount++
		}
	}
	hitRate := 0
	if total > 0 {
		hitRate = trueCount * 100 / total
	}
	hitLine := heatNil.Render(fmt.Sprintf("Hit rate: %d%%", hitRate))

	return "  " + headerLine.String() + "\n" +
		"  " + markLine.String() + "\n\n" +
		"  " + hitLine
}
