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

// The cover shows the one-sentence tagline when set, hides it when empty, and a
// long tagline wraps without spilling past the terminal edge (even when narrow).
func TestTitlePageTagline(t *testing.T) {
	s := NewStyles("#e06474")
	tm := Theme{Name: "pensum", Sigil: "✎", Accent: "#e06474"}

	tag := "the work you set yourself, kept in order"
	if out := s.TitlePage(tm.WithTagline(tag), 100, 30); !strings.Contains(out, tag) {
		t.Errorf("title page should show the tagline when set")
	}
	if strings.Contains(s.TitlePage(tm, 100, 30), tag) {
		t.Error("no tagline set → none shown")
	}

	// A tagline wider than the terminal must wrap, never overflow, at any width.
	long := "a long one-sentence description that comfortably exceeds a narrow terminal width and must wrap"
	for _, w := range []int{40, 60, 100} {
		out := s.TitlePage(tm.WithTagline(long), w, 30)
		for i, ln := range strings.Split(out, "\n") {
			if lw := lipgloss.Width(ln); lw > w {
				t.Errorf("width %d: tagline line %d width %d exceeds %d", w, i, lw, w)
			}
		}
	}
}

// The big-letter banner appears at a comfortable width and degrades to spaced
// capitals when the terminal is too narrow for it.
func TestTitlePageNarrowFallback(t *testing.T) {
	s := NewStyles("#f5a623")
	tm := Theme{Name: "fasti", Sigil: "⊕", Accent: "#f5a623"}

	wide := s.TitlePage(tm, 120, 30)
	if !strings.Contains(wide, "█") {
		t.Error("wide title page should render the heavy-block banner")
	}
	narrow := s.TitlePage(tm, 30, 30)
	if strings.Contains(narrow, "█") {
		t.Error("narrow title page should fall back to spaced caps, not the banner")
	}
	if !strings.Contains(narrow, SpacedCaps("fasti")) {
		t.Error("narrow fallback should show the spaced-caps name")
	}
	for i, ln := range strings.Split(narrow, "\n") {
		if w := lipgloss.Width(ln); w > 30 {
			t.Errorf("narrow line %d width %d exceeds 30", i, w)
		}
	}
}
