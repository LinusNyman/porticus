// Package insights renders the suite's shared statistical chrome: compact
// sparklines, binary heatmaps, trend computation, and a scrollable insights
// page. It is the single charting source for tools with a graph/insights
// screen — album's contact demographics, speculum's habit stats — so those
// screens look consistent and improve together.
//
// Pure presentation: it depends only on lipgloss and the porticus root (for the
// suite palette and page chrome), never on bubbletea or the data spine. As with
// HelpPage, any scroll offset lives in the caller.
package insights

import (
	"time"

	"github.com/LinusNyman/porticus"
	"github.com/charmbracelet/lipgloss"
)

// DayValue is one day's optional numeric value, ordered oldest→newest by the
// callers. A nil Value marks a skipped or missing day, which the charts keep in
// position rather than collapsing. It is deliberately decoupled from any tool's
// storage type so this package stays presentation-only.
type DayValue struct {
	Date  time.Time
	Value *float64
}

// Suite-themed chart colours (suite standard §5): rising/present in laurel
// green, falling in ochre, missing or flat in weathered stone, and absent
// (false) in Pompeian red.
var (
	chartUp   = lipgloss.NewStyle().Foreground(lipgloss.Color(porticus.ColCompleted))
	chartDown = lipgloss.NewStyle().Foreground(lipgloss.Color(porticus.ColDue))
	chartFlat = lipgloss.NewStyle().Foreground(lipgloss.Color(porticus.ColDim))
	heatOn    = chartUp
	heatOff   = lipgloss.NewStyle().Foreground(lipgloss.Color(porticus.ColOverdue))
	heatNil   = chartFlat
)
