package porticus

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)

// taglineWidth caps how wide the one-sentence tagline is allowed to set, so on a
// broad terminal it reads as a centred subtitle line rather than stretching the
// full screen width. It still wraps narrower when the terminal is.
const taglineWidth = 52

// TitlePage renders a tool's cover screen: the sigil, the tool name in large
// heavy-block capitals (BigText) in the accent, the one-sentence tagline (when
// set), a hedera-and-rule ornament, then the version (when set) and the author —
// the whole composition centred in a block of exactly width×height so it sits as
// its own full-screen view (bound to "+", the suite's cover/meta screen, alongside
// "?" help). It carries no left-header bar; unlike Page/HelpPage a cover is a clean
// full-bleed composition. The name drops to spaced capitals when the banner is
// wider than the terminal; the tagline wraps to fit.
//
// The classical reading comes from the frame — the sigil crown, the ❧ over a ══
// rule, the spaced-caps author — not from the glyphs, which stay the suite's bold
// block letters. Version and tagline are shown here rather than in the help header.
func (s Styles) TitlePage(t Theme, width, height int) string {
	var lines []string
	add := func(styled string) { lines = append(lines, Center(styled, width)) }

	add(s.Title.Render(t.Sigil))
	add("")

	banner := BigText(t.Name)
	bannerLines := strings.Split(banner, "\n")
	if lipgloss.Width(bannerLines[0]) <= width {
		for _, ln := range bannerLines {
			add(s.Title.Render(ln))
		}
	} else {
		add(s.Title.Render(SpacedCaps(t.Name)))
	}

	add("")
	if t.Tagline != "" {
		measure := width
		if measure > taglineWidth {
			measure = taglineWidth
		}
		for _, ln := range strings.Split(ansi.Wrap(t.Tagline, measure, ""), "\n") {
			add(s.Desc.Render(ln))
		}
		add("")
	}
	add(s.Divider.Render("══════") + s.Title.Render("  ❧  ") + s.Divider.Render("══════"))
	add("")

	if t.Version != "" {
		add(s.Dim.Render(t.Version))
	}
	add(s.Dim.Render(SpacedCaps(Author)))

	block := strings.Join(lines, "\n")
	return lipgloss.Place(width, height, lipgloss.Left, lipgloss.Center, block)
}
