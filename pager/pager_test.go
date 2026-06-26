package pager

import (
	"strings"
	"testing"

	"github.com/LinusNyman/porticus"
	"github.com/LinusNyman/porticus/keys"
	tea "github.com/charmbracelet/bubbletea"
)

// keyMsg builds a rune key message (j/k/g/G); ctrlKey builds a control key.
func keyMsg(s string) tea.KeyMsg           { return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(s)} }
func ctrlKey(t tea.KeyType) tea.KeyMsg     { return tea.KeyMsg{Type: t} }
func wheel(b tea.MouseButton) tea.MouseMsg { return tea.MouseMsg{Button: b} }

func body(n int) string { return strings.TrimRight(strings.Repeat("row\n", n), "\n") }

// A 50-row body in a 12-row page: PageRows(12,50)=9, so maxScroll = 50-9 = 41.
const (
	tall    = 50
	height  = 12
	maxRoll = 41
)

func TestNavKeys(t *testing.T) {
	km := keys.Default()
	var p Pager
	p.SetContent(body(tall))

	// Down/Up step one line; Up at the top is a no-op.
	if !p.Handle(keyMsg("j"), km, height) || p.Scroll() != 1 {
		t.Fatalf("j should scroll to 1, got %d", p.Scroll())
	}
	if !p.Handle(keyMsg("k"), km, height) || p.Scroll() != 0 {
		t.Fatalf("k should scroll back to 0, got %d", p.Scroll())
	}
	if p.Handle(keyMsg("k"), km, height) {
		t.Error("k at the top should not move")
	}

	// Half page = PageStep(12) = (12-3)/2 = 4.
	if !p.Handle(ctrlKey(tea.KeyCtrlD), km, height) || p.Scroll() != 4 {
		t.Errorf("ctrl+d should scroll by a half page (4), got %d", p.Scroll())
	}

	// G jumps to the bottom (maxScroll); G again is a no-op; g returns to the top.
	if !p.Handle(keyMsg("G"), km, height) || p.Scroll() != maxRoll {
		t.Fatalf("G should land on maxScroll %d, got %d", maxRoll, p.Scroll())
	}
	if p.Handle(keyMsg("G"), km, height) {
		t.Error("G at the bottom should not move")
	}
	if !p.Handle(keyMsg("g"), km, height) || p.Scroll() != 0 {
		t.Errorf("g should return to the top, got %d", p.Scroll())
	}

	// An unrecognised key is left for the caller.
	if p.Handle(keyMsg("z"), km, height) {
		t.Error("an unrecognised key should not be consumed")
	}
}

func TestWheel(t *testing.T) {
	var p Pager
	p.SetContent(body(tall))
	// Give the pager a height (a render would do this) so the wheel can clamp.
	_ = p.View(porticus.NewStyles("#e06474"), porticus.Theme{Name: "pensum"}, "preview", 40, height)

	if !p.HandleMouse(wheel(tea.MouseButtonWheelDown)) || p.Scroll() != 1 {
		t.Fatalf("wheel down should scroll to 1, got %d", p.Scroll())
	}
	if !p.HandleMouse(wheel(tea.MouseButtonWheelUp)) || p.Scroll() != 0 {
		t.Fatalf("wheel up should scroll to 0, got %d", p.Scroll())
	}
	if p.HandleMouse(wheel(tea.MouseButtonWheelUp)) {
		t.Error("wheel up at the top should not move")
	}
}

func TestSetContentClamps(t *testing.T) {
	km := keys.Default()
	var p Pager
	p.SetContent(body(tall))
	p.Handle(keyMsg("G"), km, height) // scroll to the bottom of the tall body

	// Swapping in a body that fits the page must pull the offset back to 0.
	p.SetContent(body(3))
	if p.Scroll() != 0 {
		t.Errorf("SetContent with a fitting body should reset scroll to 0, got %d", p.Scroll())
	}
}

func TestViewMatchesPage(t *testing.T) {
	s := porticus.NewStyles("#e06474")
	tm := porticus.Theme{Name: "pensum", Sigil: "✎", Accent: "#e06474"}
	km := keys.Default()
	b := body(tall)

	var p Pager
	p.SetContent(b)
	p.Handle(ctrlKey(tea.KeyCtrlD), km, height) // scroll to 4

	got := p.View(s, tm, "preview", 60, height)
	want := s.Page(tm, "preview", b, 60, height, p.Scroll())
	if got != want {
		t.Error("Pager.View must equal Styles.Page at the same offset")
	}
	// Overflowing content shows a downward scroll hint.
	if !strings.Contains(got, "↓") {
		t.Error("an overflowing pager should render a scroll hint")
	}
}
