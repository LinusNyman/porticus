package insights

// Direction of a metric's movement relative to its recent average.
type Direction int

const (
	Flat Direction = iota
	Up
	Down
)

// Arrow is the glyph for a direction: ↑ up, ↓ down, → flat.
func (d Direction) Arrow() string {
	switch d {
	case Up:
		return "↑"
	case Down:
		return "↓"
	default:
		return "→"
	}
}

// Trend holds computed statistics for a single field over a data window.
type Trend struct {
	Today     *float64 // nil if today has no value
	Avg7d     *float64 // nil if fewer than 2 non-null values in the 7-day window
	Avg30d    *float64 // nil if fewer than 2 non-null values in the window
	Min       float64
	Max       float64
	Direction Direction
}

// Compute derives a Trend from a series of DayValues (oldest→newest). Skipped
// (nil) days are excluded from averages and do not count against N.
func Compute(values []DayValue) Trend {
	if len(values) == 0 {
		return Trend{}
	}

	// Today is the last element.
	last := values[len(values)-1]
	var today *float64
	if last.Value != nil {
		v := *last.Value
		today = &v
	}

	nonNil := make([]float64, 0, len(values))
	for _, dv := range values {
		if dv.Value != nil {
			nonNil = append(nonNil, *dv.Value)
		}
	}

	var minVal, maxVal float64
	for i, v := range nonNil {
		if i == 0 || v < minVal {
			minVal = v
		}
		if i == 0 || v > maxVal {
			maxVal = v
		}
	}

	avg := func(slice []float64) *float64 {
		if len(slice) < 2 {
			return nil
		}
		var sum float64
		for _, v := range slice {
			sum += v
		}
		a := sum / float64(len(slice))
		return &a
	}

	// 7-day window: collect non-nil values from the last 7 entries.
	window7 := make([]float64, 0, 7)
	start7 := len(values) - 7
	if start7 < 0 {
		start7 = 0
	}
	for _, dv := range values[start7:] {
		if dv.Value != nil {
			window7 = append(window7, *dv.Value)
		}
	}

	avg7d := avg(window7)
	avg30d := avg(nonNil)

	dir := directionFor(today, avg7d)

	return Trend{
		Today:     today,
		Avg7d:     avg7d,
		Avg30d:    avg30d,
		Min:       minVal,
		Max:       maxVal,
		Direction: dir,
	}
}

// directionFor computes the direction of today relative to avg7d: >10% above →
// Up, >10% below → Down, within 10% → Flat. For small ranges (avg < 5, e.g.
// 1–5 ratings) it uses an absolute threshold of 0.5 instead.
func directionFor(today, avg7d *float64) Direction {
	if today == nil || avg7d == nil {
		return Flat
	}
	t, a := *today, *avg7d
	if a == 0 {
		if t > 0 {
			return Up
		}
		return Flat
	}
	// Absolute 0.5 threshold when values are small (e.g. 1–5 ratings).
	if a < 5 {
		diff := t - a
		if diff > 0.5 {
			return Up
		}
		if diff < -0.5 {
			return Down
		}
		return Flat
	}
	// 10% relative threshold for larger counts.
	ratio := t / a
	if ratio > 1.10 {
		return Up
	}
	if ratio < 0.90 {
		return Down
	}
	return Flat
}
