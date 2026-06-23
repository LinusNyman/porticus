package porticus

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)

// SpacedCaps converts a string to uppercase with a space between every letter,
// e.g. "pensum" → "P E N S U M". Used for tool and pane titles (suite standard
// §3.1). Works on multi-byte runes (åäö) since it splits on runes.
func SpacedCaps(s string) string {
	return strings.Join(strings.Split(strings.ToUpper(s), ""), " ")
}

// PadTo right-pads s with spaces to width display cells, accounting for ANSI
// styling (it measures with lipgloss.Width). It never truncates: a string wider
// than width is returned unchanged.
func PadTo(s string, width int) string {
	if width <= 0 {
		return s
	}
	w := lipgloss.Width(s)
	if w >= width {
		return s
	}
	return s + strings.Repeat(" ", width-w)
}

// Center pads s with spaces on both sides to center it within width display
// cells, accounting for ANSI styling. A string already at least width wide is
// returned unchanged; any odd remainder favours the right side.
func Center(s string, width int) string {
	pad := width - lipgloss.Width(s)
	if pad <= 0 {
		return s
	}
	left := pad / 2
	return strings.Repeat(" ", left) + s + strings.Repeat(" ", pad-left)
}

// Truncate visually truncates an ANSI-styled string to width cells, appending an
// ellipsis when it overflows. width<=0 yields "".
func Truncate(s string, width int) string {
	if width <= 0 {
		return ""
	}
	if lipgloss.Width(s) <= width {
		return s
	}
	if width == 1 {
		return "…"
	}
	return ansi.Truncate(s, width, "…")
}

// WrapRows wraps a list row's styled body so the whole row fits within width,
// returning one string per visual row. The first row carries the prefix (index,
// checkbox, columns…); continuation rows are indented by the prefix width so the
// leading column reads blank and the body lines up as one paragraph. Shared by
// every tool's list panes so wrapping behaves identically everywhere.
func WrapRows(prefix, body string, width int) []string {
	pw := lipgloss.Width(prefix)
	avail := width - pw
	if avail < 1 {
		// Too narrow for a hanging indent; truncate the joined line instead.
		return []string{Truncate(prefix+body, width)}
	}
	parts := strings.Split(ansi.Wrap(body, avail, ""), "\n")
	rows := make([]string, 0, len(parts))
	indent := strings.Repeat(" ", pw)
	for i, p := range parts {
		if i == 0 {
			rows = append(rows, prefix+p)
		} else {
			rows = append(rows, indent+p)
		}
	}
	return rows
}
