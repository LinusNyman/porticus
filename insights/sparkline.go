package insights

import "strings"

var sparkChars = []rune{'▁', '▂', '▃', '▄', '▅', '▆', '▇', '█'}

// Sparkline renders a horizontal bar chart for daily values (oldest→newest),
// exactly width runes wide. Nil days render as spaces so positions are kept;
// when fewer than width values are given the line is right-aligned in the
// window. The whole line is tinted by the trend direction (Compute): laurel
// green rising, ochre falling, stone flat.
func Sparkline(values []DayValue, width int) string {
	if width <= 0 {
		width = 30
	}

	// Collect non-nil values to find the range.
	var nonNil []float64
	for _, dv := range values {
		if dv.Value != nil {
			nonNil = append(nonNil, *dv.Value)
		}
	}

	// All nil → blank line of the right width.
	if len(nonNil) == 0 {
		return strings.Repeat(" ", width)
	}

	var minV, maxV float64
	for i, v := range nonNil {
		if i == 0 || v < minV {
			minV = v
		}
		if i == 0 || v > maxV {
			maxV = v
		}
	}

	scale := func(v float64) rune {
		if maxV == minV {
			return sparkChars[len(sparkChars)/2]
		}
		idx := int((v - minV) / (maxV - minV) * float64(len(sparkChars)-1))
		if idx < 0 {
			idx = 0
		}
		if idx >= len(sparkChars) {
			idx = len(sparkChars) - 1
		}
		return sparkChars[idx]
	}

	// Build a rune slice of exactly width entries (values right-aligned).
	runes := make([]rune, width)
	for i := range runes {
		runes[i] = ' '
	}
	offset := width - len(values)
	if offset < 0 {
		offset = 0
		values = values[len(values)-width:]
	}
	for i, dv := range values {
		if dv.Value == nil {
			runes[offset+i] = ' '
		} else {
			runes[offset+i] = scale(*dv.Value)
		}
	}

	s := string(runes)
	switch Compute(values).Direction {
	case Up:
		return chartUp.Render(s)
	case Down:
		return chartDown.Render(s)
	default:
		return chartFlat.Render(s)
	}
}
