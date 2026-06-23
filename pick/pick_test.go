package pick_test

import (
	"strings"
	"testing"

	"github.com/LinusNyman/porticus/keys"
	"github.com/LinusNyman/porticus/pick"
	tea "github.com/charmbracelet/bubbletea"
)

// items filters a fixed corpus by substring — stands in for a tool's index.
var corpus = []string{"apple", "apricot", "banana", "cherry"}

func contains(q string) []string {
	if q == "" {
		return nil
	}
	var out []string
	for _, w := range corpus {
		if strings.Contains(w, q) {
			out = append(out, w)
		}
	}
	return out
}

func newPicker(limit int) pick.Picker[string] {
	return pick.New(pick.Opts[string]{
		Label:  "search",
		Limit:  limit,
		Filter: contains,
		Render: func(it string, width int) string { return it },
	})
}

func rune1(s string) tea.KeyMsg        { return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(s)} }
func special(t tea.KeyType) tea.KeyMsg { return tea.KeyMsg{Type: t} }

func TestLiveFilterWhileTyping(t *testing.T) {
	p := newPicker(0)
	p.Open()
	km := keys.Default()
	// Type "ap" → matches apple, apricot.
	p.Update(rune1("a"), km)
	p.Update(rune1("p"), km)
	if p.Query() != "ap" {
		t.Fatalf("query = %q, want \"ap\"", p.Query())
	}
	if _, ok := p.Selected(); !ok {
		t.Fatal("expected a selected result after filtering")
	}
	if got, _ := p.Selected(); got != "apple" {
		t.Errorf("first result = %q, want apple", got)
	}
}

func TestEnterCommitsThenNavigateAndSelect(t *testing.T) {
	p := newPicker(0)
	p.Open()
	km := keys.Default()
	p.Update(rune1("a"), km) // matches apple, apricot, banana
	// Commit the query — leave typing mode.
	if _, _, handled := p.Update(special(tea.KeyEnter), km); !handled {
		t.Fatal("enter should be handled (commit query)")
	}
	// Now in browse mode: j moves down.
	if _, _, handled := p.Update(rune1("j"), km); !handled {
		t.Fatal("j should be handled in browse mode")
	}
	got, ok := p.Selected()
	if !ok || got != "apricot" {
		t.Errorf("after j, selection = %q (ok=%v), want apricot", got, ok)
	}
	// Select with enter → returns the item.
	sel, _, handled := p.Update(special(tea.KeyEnter), km)
	if !handled || sel == nil || *sel != "apricot" {
		t.Errorf("enter should select apricot, got %v handled=%v", sel, handled)
	}
}

func TestLimitCapsResultsAndCursor(t *testing.T) {
	p := newPicker(2) // suggest-style cap
	p.Open()
	km := keys.Default()
	p.Update(rune1("a"), km)            // apple, apricot, banana (3) → capped to 2
	p.Update(special(tea.KeyEnter), km) // browse mode
	p.Update(rune1("G"), km)            // jump to bottom (cursor → 1, not 2)
	if got, _ := p.Selected(); got != "apricot" {
		t.Errorf("with limit 2, bottom selection = %q, want apricot (index 1)", got)
	}
}

func TestEscIsNotHandledSoToolCanClose(t *testing.T) {
	p := newPicker(0)
	p.Open()
	km := keys.Default()
	// esc while typing is left to the caller (to close the overlay).
	if _, _, handled := p.Update(special(tea.KeyEsc), km); handled {
		t.Error("esc while typing should be unhandled so the tool can close")
	}
}

func TestMouseWheelBrowse(t *testing.T) {
	p := newPicker(0)
	p.Open()
	km := keys.Default()
	p.Update(rune1("a"), km)            // apple, apricot, banana
	p.Update(special(tea.KeyEnter), km) // browse mode
	if !p.HandleMouse(tea.MouseMsg{Button: tea.MouseButtonWheelDown}) {
		t.Fatal("wheel should move in browse mode")
	}
	if got, _ := p.Selected(); got != "apricot" {
		t.Errorf("after wheel down, selection = %q, want apricot", got)
	}
}

func TestMouseIgnoredWhileTyping(t *testing.T) {
	p := newPicker(0)
	p.Open() // typing mode
	if p.HandleMouse(tea.MouseMsg{Button: tea.MouseButtonWheelDown}) {
		t.Error("wheel should be ignored while typing")
	}
}

func TestRefineReopensQuery(t *testing.T) {
	p := newPicker(0)
	p.Open()
	km := keys.Default()
	p.Update(rune1("b"), km)            // banana
	p.Update(special(tea.KeyEnter), km) // browse
	// `/` refines: re-enter typing with the prior query preserved.
	if _, _, handled := p.Update(rune1("/"), km); !handled {
		t.Fatal("/ should be handled (refine)")
	}
	if p.Query() != "b" {
		t.Errorf("refine should keep the query, got %q", p.Query())
	}
}
