package porticus

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// HelpGroup is one labelled section of the help page: a heading over key/action
// rows. Rows are {keys, action} pairs, e.g. {"j / k", "move down / up"}.
type HelpGroup struct {
	Title string
	Rows  [][2]string
}

// helpKeyW is the column width the key field is padded to so actions align.
const helpKeyW = 6

// HelpPage renders the suite help screen as a page of the app (suite standard
// §6): the unchanged tool title with "help" where the node name normally sits,
// over a carved-stone plaque (a double-ruled border echoing the ══ rule) holding
// the binding groups in one or two columns. scroll is the first body row to show;
// the page windows the groups to the available height and shows a scroll hint in
// the header. Returns a block exactly width×height so the footer stays pinned.
func (s Styles) HelpPage(t Theme, groups []HelpGroup, width, height, scroll int) string {
	blocks := make([]string, len(groups))
	for i, g := range groups {
		blocks[i] = s.renderHelpGroup(g)
	}
	// The plaque chrome (border + horizontal padding) costs 6 columns; lay the
	// groups within whatever's left so it never spills past the terminal edge.
	body := s.helpColumns(blocks, width-6)

	// One row for the heading and the plaque's top/bottom rules take height the
	// body can't have; the top rule doubles as the heading rule. The scroll hint
	// rides in the header (not a reserved body row), so the body gets the full
	// avail rows. windowLines is the shared scroll-offset windowing (see page.go).
	avail := height - 1 - 2
	bodyLines, above, below := windowLines(strings.Split(body, "\n"), avail, scroll)
	// Draw the plaque full width so its top border spans the heading row and
	// reads as the rule beneath it (Width sets content+padding; the border adds 2).
	plaque := s.HelpPlaque.Width(width - 2).Render(strings.Join(bodyLines, "\n"))
	page := s.helpHeader(t, s.ScrollHint(above, below, width), width) + "\n" + plaque
	return lipgloss.Place(width, height, lipgloss.Left, lipgloss.Top, page)
}

// helpHeader renders the help heading: the tool title with "help" standing in
// for the node name, plus any scroll indicator. A single line at every width —
// the plaque's top border drawn beneath is the rule under it.
func (s Styles) helpHeader(t Theme, hint string, width int) string {
	title := s.Title.Render(t.Sigil + "  " + SpacedCaps(t.Name) + "  ❧  help")
	if hint != "" {
		title += "  " + hint
	}
	return Truncate(title, width)
}

// renderHelpGroup formats one section: a spaced-caps heading over "keys action"
// rows, keys padded to helpKeyW so actions align.
//
// Heading colour: the tool accent (s.Title). pensum is the suite's ground-truth
// help screen, and it renders headings in the accent; the earlier guide §6
// drift toward colHeading (terracotta) is resolved in pensum's favour here, so
// every tool's help looks like pensum's (guide §6/§9 updated 2026-06-23). Keep
// this as s.Title — do not switch to s.Heading.
func (s Styles) renderHelpGroup(g HelpGroup) string {
	lines := []string{s.Title.Render(SpacedCaps(g.Title))}
	for _, r := range g.Rows {
		lines = append(lines, "  "+PadTo(s.Code.Render(r[0]), helpKeyW)+"  "+s.Dim.Render(r[1]))
	}
	return strings.Join(lines, "\n")
}

// helpColumns lays the group blocks in two columns when they fit innerWidth,
// else stacks them in one so the plaque stays inside the terminal.
func (s Styles) helpColumns(blocks []string, innerWidth int) string {
	mid := (len(blocks) + 1) / 2
	left := lipgloss.JoinVertical(lipgloss.Left, joinWithGaps(blocks[:mid])...)
	right := lipgloss.JoinVertical(lipgloss.Left, joinWithGaps(blocks[mid:])...)
	two := lipgloss.JoinHorizontal(lipgloss.Top, left, "    ", right)
	if innerWidth <= 0 || lipgloss.Width(two) <= innerWidth {
		return two
	}
	return lipgloss.JoinVertical(lipgloss.Left, joinWithGaps(blocks)...)
}

// joinWithGaps interleaves blank lines between blocks so groups breathe.
func joinWithGaps(blocks []string) []string {
	out := make([]string, 0, len(blocks)*2)
	for i, b := range blocks {
		if i > 0 {
			out = append(out, "")
		}
		out = append(out, b)
	}
	return out
}
