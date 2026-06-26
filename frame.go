package porticus

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// NarrowWidth is the terminal width below which the two-pane layout stacks
// vertically instead (suite standard §2).
const NarrowWidth = 80

// PaneWidths returns the left (tree/list) and right (working-area) pane widths
// for a given total terminal width, per the suite layout spec (§2): the left
// pane is width*2/5 clamped to [28,48]; the right takes the remainder after the
// 3-column divider, with a floor of 20.
func PaneWidths(total int) (left, right int) {
	left = total * 2 / 5
	if left < 28 {
		left = 28
	}
	if left > 48 {
		left = 48
	}
	right = total - left - 3 // 3 cols: space · ║ · space
	if right < 20 {
		right = 20
	}
	return left, right
}

// vDivider renders the vertical pane divider: ` ║ ` repeated for height rows so
// columns align across every row (suite standard §6).
func (s Styles) vDivider(height int) string {
	row := " " + s.Divider.Render("║") + " "
	rows := make([]string, height)
	for i := range rows {
		rows[i] = row
	}
	return strings.Join(rows, "\n")
}

// HRule renders a full-width single-rule line of the given width in the divider
// colour (the stacked-layout separator, U+2500).
func (s Styles) HRule(width int) string {
	return s.Divider.Render(strings.Repeat("─", width))
}

// TwoPane assembles the standard two-column body (terminal width ≥ NarrowWidth):
// a left pane and a right pane joined by the ` ║ ` divider. The caller supplies
// each pane's content via a render callback that is handed the pane's computed
// width and the shared body height, so the tool keeps full control of what fills
// each pane while porticus owns the geometry and the divider.
func (s Styles) TwoPane(width, height int, left, right func(w, h int) string) string {
	lw, rw := PaneWidths(width)
	l := left(lw, height)
	r := right(rw, height)
	return lipgloss.JoinHorizontal(lipgloss.Top, l, s.vDivider(height), r)
}

// Stacked assembles the narrow-fallback layout (terminal width < NarrowWidth):
// the top pane, a full-width ─ rule, then the bottom pane. The top pane gets half
// the height (min 4 rows), the bottom the rest (min 3) — see suite standard §2.
func (s Styles) Stacked(width, height int, top, bottom func(w, h int) string) string {
	topH := height / 2
	botH := height - topH - 1
	if topH < 4 {
		topH = 4
	}
	if botH < 3 {
		botH = 3
	}
	return top(width, topH) + "\n" + s.HRule(width) + "\n" + bottom(width, botH)
}

// LeftHeader renders the left pane's two-line header: the identity title
// (sigil + spaced-caps tool name + ❧ + node name) over a full-width ══ double
// rule (suite standard §3.1, §3.3). nodeName may be "" for the root.
func (s Styles) LeftHeader(t Theme, nodeName string, width int) string {
	line := t.Sigil + "  " + SpacedCaps(t.Name)
	if nodeName != "" {
		line += "  ❧  " + nodeName
	}
	title := Truncate(s.Title.Render(line), width)
	rule := s.Divider.Render(strings.Repeat("═", width))
	return title + "\n" + rule
}

// RightHeader renders the right pane's two-line header: a spaced-caps label in
// stone (no sigil) over a single ── rule. When the pane is focused the rule is
// drawn in the accent colour (suite standard §6).
func (s Styles) RightHeader(label string, focused bool, width int) string {
	return s.Dim.Render(SpacedCaps(label)) + "\n" + s.rightRule(focused, width)
}

// RightHeaderCount is RightHeader with an item-count suffix after the label: the
// count at the selected node in parens (verdigris) and the whole-forest total
// after a middle dot (dim), e.g. "T O D O S  (4) · 137". The counts show only
// when total > 0; both numbers come from the tool, which knows what an "item"
// is. (The right title carrying a count is a deliberate exception to the guide's
// "right title: single label" rule — see §3.2 / the §9 decision log.)
func (s Styles) RightHeaderCount(label string, focused bool, node, total, width int) string {
	title := s.Dim.Render(SpacedCaps(label))
	if total > 0 {
		title += "  " + s.Count.Render(fmt.Sprintf("(%d)", node)) +
			s.Dim.Render(fmt.Sprintf(" · %d", total))
	}
	return Truncate(title, width) + "\n" + s.rightRule(focused, width)
}

// rightRule renders the right pane's single ── underline: weathered stone when
// blurred, the accent colour when the pane is focused (suite standard §6).
func (s Styles) rightRule(focused bool, width int) string {
	ruleStyle := s.Dim
	if focused {
		ruleStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(s.Accent))
	}
	return ruleStyle.Render(strings.Repeat("─", width))
}

// ScrollHint renders a dim one-line indicator of how many rows are hidden above
// and/or below the visible window, e.g. "↑3  ↓5". Empty when nothing is hidden.
func (s Styles) ScrollHint(above, below, width int) string {
	if above <= 0 && below <= 0 {
		return ""
	}
	var parts []string
	if above > 0 {
		parts = append(parts, fmt.Sprintf("↑%d", above))
	}
	if below > 0 {
		parts = append(parts, fmt.Sprintf("↓%d", below))
	}
	return Truncate(s.Dim.Render("  "+strings.Join(parts, "  ")), width)
}

// NodeCount renders the canonical tree-row count badge: the node's own item
// count in parens (verdigris) then its descendants' total after a plus (dim),
// e.g. "(4) +12". Either part is dropped when zero; both zero yields "". This is
// the shared form every tool's browse.Opts.Annotate should return — browse
// right-aligns it and inserts the gap before it, so return no leading padding
// (suite standard §6, guide §3.3).
func (s Styles) NodeCount(own, desc int) string {
	if own == 0 && desc == 0 {
		return ""
	}
	var parts []string
	if own > 0 {
		parts = append(parts, s.Count.Render(fmt.Sprintf("(%d)", own)))
	}
	if desc > 0 {
		parts = append(parts, s.Dim.Render(fmt.Sprintf("+%d", desc)))
	}
	return strings.Join(parts, " ")
}

// Checkbox is the list status glyph shared by every tool: a dim empty ☐, or a
// laurel-green ☑ when done. Both are one cell wide so columns stay aligned.
func (s Styles) Checkbox(done bool) string {
	if done {
		return s.Completed.Render("☑")
	}
	return s.Dim.Render("☐")
}
