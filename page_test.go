package porticus

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
)

func TestWindowLines(t *testing.T) {
	lines := []string{"0", "1", "2", "3", "4", "5"}

	// Fits: returns all, nothing hidden, scroll ignored.
	vis, above, below := windowLines(lines, 10, 3)
	if len(vis) != 6 || above != 0 || below != 0 {
		t.Errorf("fitting window = %v (%d,%d), want all/0/0", vis, above, below)
	}

	// Overflow at top.
	vis, above, below = windowLines(lines, 2, 0)
	if strings.Join(vis, "") != "01" || above != 0 || below != 4 {
		t.Errorf("top window = %v (%d,%d), want 01/0/4", vis, above, below)
	}

	// Scroll clamps to the last page.
	vis, above, below = windowLines(lines, 2, 100)
	if strings.Join(vis, "") != "45" || above != 4 || below != 0 {
		t.Errorf("clamped window = %v (%d,%d), want 45/4/0", vis, above, below)
	}
}

// The version is shown on the title screen (TestTitlePageShowsAuthorAndVersion),
// not in the help header — confirm it has left the header.
func TestHelpPageOmitsVersion(t *testing.T) {
	s := NewStyles("#f5a623")
	groups := []HelpGroup{{Title: "Navigate", Rows: [][2]string{{"j", "down"}}}}
	album := Theme{Name: "album", Sigil: "❦", Accent: "#f5a623"}

	if strings.Contains(s.HelpPage(album.WithVersion("v1.2.3"), groups, 80, 12, 0), "v1.2.3") {
		t.Error("help header should no longer show the version (it lives on the title page)")
	}
}

func TestPageGeometry(t *testing.T) {
	s := NewStyles("#e06474")
	tm := Theme{Name: "pensum", Sigil: "✎", Accent: "#e06474"}
	body := strings.Repeat("row\n", 50)
	out := s.Page(tm, "preview", body, 60, 12, 0)

	lines := strings.Split(out, "\n")
	if len(lines) != 12 {
		t.Fatalf("Page height = %d lines, want 12", len(lines))
	}
	for i, ln := range lines {
		if w := lipgloss.Width(ln); w > 60 {
			t.Errorf("line %d width %d exceeds 60", i, w)
		}
	}
	// A 50-row body in a 12-row page overflows, so a scroll hint must show.
	if !strings.Contains(out, "↓") {
		t.Error("expected a downward scroll hint on an overflowing page")
	}
}
