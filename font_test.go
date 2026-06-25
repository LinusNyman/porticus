package porticus

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
)

// Every glyph must be bigHeight rows of equal width, or BigText can't align rows.
func TestBigFontGlyphsRectangular(t *testing.T) {
	for r, g := range bigFont {
		if len(g) != bigHeight {
			t.Errorf("glyph %q has %d rows, want %d", r, len(g), bigHeight)
			continue
		}
		w := lipgloss.Width(g[0])
		for i, row := range g {
			if got := lipgloss.Width(row); got != w {
				t.Errorf("glyph %q row %d width %d, want %d", r, i, got, w)
			}
		}
	}
}

// BigText renders every A–Z (the letters tool names are built from) as a banner of
// exactly bigHeight rows, all the same display width.
func TestBigTextEveryLetter(t *testing.T) {
	for r := 'a'; r <= 'z'; r++ {
		out := BigText(string(r))
		lines := strings.Split(out, "\n")
		if len(lines) != bigHeight {
			t.Fatalf("BigText(%q) = %d rows, want %d", r, len(lines), bigHeight)
		}
		w := lipgloss.Width(lines[0])
		for i, ln := range lines {
			if got := lipgloss.Width(ln); got != w {
				t.Errorf("BigText(%q) row %d width %d, want %d", r, i, got, w)
			}
		}
		if w == 0 {
			t.Errorf("BigText(%q) rendered empty", r)
		}
	}
}

// A multi-letter name keeps the rows aligned, and lower-case input is upper-cased.
func TestBigTextNameAligned(t *testing.T) {
	out := BigText("pensum")
	lines := strings.Split(out, "\n")
	if len(lines) != bigHeight {
		t.Fatalf("rows = %d, want %d", len(lines), bigHeight)
	}
	w := lipgloss.Width(lines[0])
	for i, ln := range lines {
		if got := lipgloss.Width(ln); got != w {
			t.Errorf("row %d width %d, want %d (rows must align)", i, got, w)
		}
	}
	// "pensum" and "PENSUM" must render identically (input is upper-cased).
	if BigText("pensum") != BigText("PENSUM") {
		t.Error("BigText should upper-case its input")
	}
}

// An unmapped rune falls back to the blank glyph rather than panicking.
func TestBigTextUnknownRune(t *testing.T) {
	out := BigText("@") // not in the font
	lines := strings.Split(out, "\n")
	if len(lines) != bigHeight {
		t.Fatalf("rows = %d, want %d", len(lines), bigHeight)
	}
	if strings.TrimSpace(strings.Join(lines, "")) != "" {
		t.Errorf("unknown rune should render blank, got %q", out)
	}
}
