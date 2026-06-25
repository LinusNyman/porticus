package porticus

import "strings"

// bigHeight is the row count of every glyph in bigFont; all glyphs are this tall
// so BigText can concatenate them column-by-column and have the rows line up.
const bigHeight = 6

// bigFont is the suite's large-letter alphabet — the "ANSI Shadow" heavy-block
// capitals, the same family speculum hand-draws for its splash. Every glyph is
// bigHeight rows and internally rectangular (all its rows share one width), so a
// name renders by joining glyphs left-to-right per row. Lower-case input is
// upper-cased by BigText; anything outside A–Z (and space) falls back to a blank
// glyph of space's width so an unknown rune never breaks the banner.
var bigFont = map[rune][]string{
	' ': {
		"    ",
		"    ",
		"    ",
		"    ",
		"    ",
		"    ",
	},
	'A': {
		" █████╗ ",
		"██╔══██╗",
		"███████║",
		"██╔══██║",
		"██║  ██║",
		"╚═╝  ╚═╝",
	},
	'B': {
		"██████╗ ",
		"██╔══██╗",
		"██████╔╝",
		"██╔══██╗",
		"██████╔╝",
		"╚═════╝ ",
	},
	'C': {
		" ██████╗",
		"██╔════╝",
		"██║     ",
		"██║     ",
		"╚██████╗",
		" ╚═════╝",
	},
	'D': {
		"██████╗ ",
		"██╔══██╗",
		"██║  ██║",
		"██║  ██║",
		"██████╔╝",
		"╚═════╝ ",
	},
	'E': {
		"███████╗",
		"██╔════╝",
		"█████╗  ",
		"██╔══╝  ",
		"███████╗",
		"╚══════╝",
	},
	'F': {
		"███████╗",
		"██╔════╝",
		"█████╗  ",
		"██╔══╝  ",
		"██║     ",
		"╚═╝     ",
	},
	'G': {
		" ██████╗ ",
		"██╔════╝ ",
		"██║  ███╗",
		"██║   ██║",
		"╚██████╔╝",
		" ╚═════╝ ",
	},
	'H': {
		"██╗  ██╗",
		"██║  ██║",
		"███████║",
		"██╔══██║",
		"██║  ██║",
		"╚═╝  ╚═╝",
	},
	'I': {
		"██╗",
		"██║",
		"██║",
		"██║",
		"██║",
		"╚═╝",
	},
	'J': {
		"     ██╗",
		"     ██║",
		"     ██║",
		"██   ██║",
		"╚█████╔╝",
		" ╚════╝ ",
	},
	'K': {
		"██╗  ██╗",
		"██║ ██╔╝",
		"█████╔╝ ",
		"██╔═██╗ ",
		"██║  ██╗",
		"╚═╝  ╚═╝",
	},
	'L': {
		"██╗     ",
		"██║     ",
		"██║     ",
		"██║     ",
		"███████╗",
		"╚══════╝",
	},
	'M': {
		"███╗   ███╗",
		"████╗ ████║",
		"██╔████╔██║",
		"██║╚██╔╝██║",
		"██║ ╚═╝ ██║",
		"╚═╝     ╚═╝",
	},
	'N': {
		"███╗   ██╗",
		"████╗  ██║",
		"██╔██╗ ██║",
		"██║╚██╗██║",
		"██║ ╚████║",
		"╚═╝  ╚═══╝",
	},
	'O': {
		" ██████╗ ",
		"██╔═══██╗",
		"██║   ██║",
		"██║   ██║",
		"╚██████╔╝",
		" ╚═════╝ ",
	},
	'P': {
		"██████╗ ",
		"██╔══██╗",
		"██████╔╝",
		"██╔═══╝ ",
		"██║     ",
		"╚═╝     ",
	},
	'Q': {
		" ██████╗ ",
		"██╔═══██╗",
		"██║   ██║",
		"██║▄▄ ██║",
		"╚██████╔╝",
		" ╚══▀▀═╝ ",
	},
	'R': {
		"██████╗ ",
		"██╔══██╗",
		"██████╔╝",
		"██╔══██╗",
		"██║  ██║",
		"╚═╝  ╚═╝",
	},
	'S': {
		"███████╗",
		"██╔════╝",
		"███████╗",
		"╚════██║",
		"███████║",
		"╚══════╝",
	},
	'T': {
		"████████╗",
		"╚══██╔══╝",
		"   ██║   ",
		"   ██║   ",
		"   ██║   ",
		"   ╚═╝   ",
	},
	'U': {
		"██╗   ██╗",
		"██║   ██║",
		"██║   ██║",
		"██║   ██║",
		"╚██████╔╝",
		" ╚═════╝ ",
	},
	'V': {
		"██╗   ██╗",
		"██║   ██║",
		"██║   ██║",
		"╚██╗ ██╔╝",
		" ╚████╔╝ ",
		"  ╚═══╝  ",
	},
	'W': {
		"██╗    ██╗",
		"██║    ██║",
		"██║ █╗ ██║",
		"██║███╗██║",
		"╚███╔███╔╝",
		" ╚══╝╚══╝ ",
	},
	'X': {
		"██╗  ██╗",
		"╚██╗██╔╝",
		" ╚███╔╝ ",
		" ██╔██╗ ",
		"██╔╝ ██╗",
		"╚═╝  ╚═╝",
	},
	'Y': {
		"██╗   ██╗",
		"╚██╗ ██╔╝",
		" ╚████╔╝ ",
		"  ╚██╔╝  ",
		"   ██║   ",
		"   ╚═╝   ",
	},
	'Z': {
		"███████╗",
		"╚══███╔╝",
		"  ███╔╝ ",
		" ███╔╝  ",
		"███████╗",
		"╚══════╝",
	},
}

// BigText renders s as a bigHeight-row banner in the suite's heavy-block capitals
// (see bigFont), the "large characters" of a tool's title screen. It upper-cases
// s, looks up each rune's glyph (a blank space-width glyph for anything unmapped),
// and joins the glyphs column-by-column so the rows align. The returned string has
// exactly bigHeight lines, each the same display width. Empty input yields
// bigHeight empty lines.
func BigText(s string) string {
	rows := make([]strings.Builder, bigHeight)
	for _, r := range strings.ToUpper(s) {
		g, ok := bigFont[r]
		if !ok {
			g = bigFont[' ']
		}
		for i := range bigHeight {
			rows[i].WriteString(g[i])
		}
	}
	lines := make([]string, bigHeight)
	for i := range rows {
		lines[i] = rows[i].String()
	}
	return strings.Join(lines, "\n")
}
