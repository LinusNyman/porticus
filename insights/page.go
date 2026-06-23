package insights

import (
	"strings"

	"github.com/LinusNyman/porticus"
)

// InsightsPage renders a full-screen insights view as a page of the app: the
// tool's left header (sigil + spaced-caps name + ❧ + label) over the ══ rule,
// then the caller's pre-rendered blocks stacked with breathing room, windowed to
// height and scrolled by scroll (the caller owns the offset). A scroll hint
// appears when the body overflows; the result is exactly width×height.
//
// Each block is a self-contained rendered section (a heading, a Sparkline, a
// CalendarHeatmap, a table…); blocks are joined with a blank line between them.
// It is a thin wrapper over porticus.Styles.Page — the shared scrollable-page
// chrome — so insights screens, help, and any custom Pager stay consistent.
func InsightsPage(s porticus.Styles, t porticus.Theme, label string, blocks []string, width, height, scroll int) string {
	return s.Page(t, label, strings.Join(blocks, "\n\n"), width, height, scroll)
}
