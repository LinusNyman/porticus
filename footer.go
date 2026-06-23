package porticus

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Hints renders the footer hint bar from groups of "key:action" hints, wrapped
// to width so no hint is pushed off-screen (suite standard §6). Hints within a
// group are separated by " · " (dim); groups by " ❧ " (the accent hedera). A
// single "key:action" hint is never split across lines. width<=0 falls back to
// a sane default.
func (s Styles) Hints(groups [][]string, width int) string {
	if width <= 0 {
		width = NarrowWidth
	}
	hedera := " " + s.Title.Render("❧") + " "
	hederaW := lipgloss.Width(hedera)
	const dot = " · "
	dotW := len(dot)

	type token struct {
		s string
		w int
	}
	var toks []token
	for gi, g := range groups {
		if gi > 0 {
			toks = append(toks, token{hedera, hederaW})
		}
		for hi, h := range g {
			if hi > 0 {
				toks = append(toks, token{dot, dotW})
			}
			r := s.Dim.Render(h)
			toks = append(toks, token{r, lipgloss.Width(r)})
		}
	}

	var lines []string
	cur, curW := "", 0
	for _, t := range toks {
		switch {
		case cur == "":
			cur, curW = t.s, t.w
		case curW+t.w <= width:
			cur += t.s
			curW += t.w
		default:
			lines = append(lines, cur)
			cur, curW = t.s, t.w
		}
	}
	if cur != "" {
		lines = append(lines, cur)
	}
	return strings.Join(lines, "\n")
}
