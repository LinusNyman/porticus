package porticus

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
)

// TitlePage fills exactly width×height (like Page) so the cover sits as its own
// full-screen view, and never spills past the terminal edge.
func TestTitlePageGeometry(t *testing.T) {
	s := NewStyles("#e06474")
	tm := Theme{Name: "pensum", Sigil: "✎", Accent: "#e06474"}
	out := s.TitlePage(tm, 100, 30)

	lines := strings.Split(out, "\n")
	if len(lines) != 30 {
		t.Fatalf("TitlePage height = %d lines, want 30", len(lines))
	}
	for i, ln := range lines {
		if w := lipgloss.Width(ln); w > 100 {
			t.Errorf("line %d width %d exceeds 100", i, w)
		}
	}
}

// The cover shows the author and, when set, the version — version having moved off
// the help header onto the title screen.
func TestTitlePageShowsAuthorAndVersion(t *testing.T) {
	s := NewStyles("#3fc79a")
	tm := Theme{Name: "speculum", Sigil: "○", Accent: "#3fc79a"}

	out := s.TitlePage(tm.WithVersion("v1.2.3"), 100, 30)
	if !strings.Contains(out, SpacedCaps(Author)) {
		t.Errorf("title page should show the author %q", SpacedCaps(Author))
	}
	if !strings.Contains(out, "v1.2.3") {
		t.Error("title page should show the version when set")
	}
	if strings.Contains(s.TitlePage(tm, 100, 30), "v1.2.3") {
		t.Error("no version set → none shown")
	}
}

// The big-letter banner appears at a comfortable width and degrades to spaced
// capitals when the terminal is too narrow for it.
func TestTitlePageNarrowFallback(t *testing.T) {
	s := NewStyles("#f5a623")
	tm := Theme{Name: "calendarium", Sigil: "⊕", Accent: "#f5a623"}

	wide := s.TitlePage(tm, 120, 30)
	if !strings.Contains(wide, "█") {
		t.Error("wide title page should render the heavy-block banner")
	}
	narrow := s.TitlePage(tm, 40, 30)
	if strings.Contains(narrow, "█") {
		t.Error("narrow title page should fall back to spaced caps, not the banner")
	}
	if !strings.Contains(narrow, SpacedCaps("calendarium")) {
		t.Error("narrow fallback should show the spaced-caps name")
	}
	for i, ln := range strings.Split(narrow, "\n") {
		if w := lipgloss.Width(ln); w > 40 {
			t.Errorf("narrow line %d width %d exceeds 40", i, w)
		}
	}
}
