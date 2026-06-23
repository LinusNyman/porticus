package porticus

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
)

func TestSpacedCaps(t *testing.T) {
	if got := SpacedCaps("pensum"); got != "P E N S U M" {
		t.Errorf("SpacedCaps = %q, want %q", got, "P E N S U M")
	}
}

func TestPadTo(t *testing.T) {
	if got := PadTo("ab", 5); got != "ab   " {
		t.Errorf("PadTo short = %q, want %q", got, "ab   ")
	}
	if got := PadTo("abcde", 3); got != "abcde" { // never truncates
		t.Errorf("PadTo over = %q, want unchanged", got)
	}
}

func TestTruncate(t *testing.T) {
	if got := Truncate("hello", 10); got != "hello" {
		t.Errorf("Truncate fits = %q, want unchanged", got)
	}
	if w := lipgloss.Width(Truncate("hello world", 5)); w > 5 {
		t.Errorf("Truncate width = %d, want <= 5", w)
	}
}

func TestNewStylesAccentOnly(t *testing.T) {
	const accent = "#e06474"
	st := NewStyles(accent)
	if st.Accent != accent {
		t.Errorf("Accent = %q, want %q", st.Accent, accent)
	}
	// The accent flows into the identity styles…
	if got := st.Title.GetForeground(); got != lipgloss.Color(accent) {
		t.Errorf("Title fg = %v, want accent %v", got, accent)
	}
	// …but shared semantic colours stay fixed regardless of accent.
	if got := st.Code.GetForeground(); got != lipgloss.Color(ColCode) {
		t.Errorf("Code fg = %v, want shared %v", got, ColCode)
	}
}

func TestPaneWidthsClamp(t *testing.T) {
	if lw, _ := PaneWidths(50); lw != 28 { // 50*2/5=20 -> floor 28
		t.Errorf("PaneWidths(50) left = %d, want 28", lw)
	}
	if lw, _ := PaneWidths(200); lw != 48 { // 200*2/5=80 -> cap 48
		t.Errorf("PaneWidths(200) left = %d, want 48", lw)
	}
}

func TestTwoPaneGeometry(t *testing.T) {
	st := NewStyles("#e06474")
	rect := func(ch string) func(w, h int) string {
		return func(w, h int) string {
			rows := make([]string, h)
			for i := range rows {
				rows[i] = strings.Repeat(ch, w)
			}
			return strings.Join(rows, "\n")
		}
	}
	const w, h = 120, 10
	out := st.TwoPane(w, h, rect("L"), rect("R"))
	lines := strings.Split(out, "\n")
	if len(lines) != h {
		t.Fatalf("TwoPane height = %d lines, want %d", len(lines), h)
	}
	lw, rw := PaneWidths(w)
	want := lw + 3 + rw // 3-col divider
	for i, line := range lines {
		if got := lipgloss.Width(line); got != want {
			t.Errorf("line %d width = %d, want %d", i, got, want)
		}
	}
}

func TestHintsFitWidth(t *testing.T) {
	st := NewStyles("#e06474")
	groups := [][]string{{"j/k:move", "tab:pane"}, {"a:add", "d:done"}, {"?:help", "q:quit"}}
	const width = 40
	out := st.Hints(groups, width)
	for _, line := range strings.Split(out, "\n") {
		if w := lipgloss.Width(line); w > width {
			t.Errorf("hint line %q width = %d, want <= %d", line, w, width)
		}
	}
	if !strings.Contains(out, "j/k:move") || !strings.Contains(out, "q:quit") {
		t.Errorf("hints missing expected text:\n%s", out)
	}
}

func TestLeftHeaderTwoLines(t *testing.T) {
	st := NewStyles("#e06474")
	hdr := st.LeftHeader(Tools["pensum"], "inbox", 40)
	lines := strings.Split(hdr, "\n")
	if len(lines) != 2 {
		t.Fatalf("LeftHeader = %d lines, want 2 (title + rule)", len(lines))
	}
	if !strings.Contains(lines[0], "P E N S U M") || !strings.Contains(lines[0], "inbox") {
		t.Errorf("title line missing name/node: %q", lines[0])
	}
}
