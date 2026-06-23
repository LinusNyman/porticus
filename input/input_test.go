package input_test

import (
	"testing"

	"github.com/LinusNyman/porticus/input"
	tea "github.com/charmbracelet/bubbletea"
)

func runes(s string) tea.KeyMsg        { return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(s)} }
func special(t tea.KeyType) tea.KeyMsg { return tea.KeyMsg{Type: t} }

func TestEditorTypeCommitClear(t *testing.T) {
	e := input.NewEditor()
	e.Open("add todo", "")
	for _, r := range []string{"b", "u", "y"} {
		if a, _ := e.Update(runes(r)); a != input.Editing {
			t.Fatalf("typing %q should be Editing, got %v", r, a)
		}
	}
	if e.Value() != "buy" {
		t.Fatalf("value = %q, want buy", e.Value())
	}
	// Enter commits without inserting a newline.
	if a, _ := e.Update(special(tea.KeyEnter)); a != input.Committed {
		t.Fatalf("enter should Commit, got %v", a)
	}
	if e.Value() != "buy" {
		t.Errorf("commit must not alter value, got %q", e.Value())
	}
	// Sticky add: clear keeps it ready for the next entry.
	e.Clear()
	if e.Value() != "" {
		t.Errorf("Clear should empty the editor, got %q", e.Value())
	}
}

func TestEditorEscCancels(t *testing.T) {
	e := input.NewEditor()
	e.Open("add todo", "draft")
	if a, _ := e.Update(special(tea.KeyEsc)); a != input.Cancelled {
		t.Errorf("esc should Cancel, got %v", a)
	}
}

func TestFieldPrefillEditCommit(t *testing.T) {
	f := input.NewField()
	f.Open("due", "2026-01-01")
	if f.Value() != "2026-01-01" {
		t.Fatalf("prefill value = %q", f.Value())
	}
	if a, _ := f.Update(runes("x")); a != input.Editing {
		t.Fatalf("typing should be Editing, got %v", a)
	}
	if a, _ := f.Update(special(tea.KeyEnter)); a != input.Committed {
		t.Errorf("enter should Commit, got %v", a)
	}
	if f.Value() != "2026-01-01x" {
		t.Errorf("value = %q, want 2026-01-01x", f.Value())
	}
}

func TestFieldEscCancels(t *testing.T) {
	f := input.NewField()
	f.Open("code", "")
	if a, _ := f.Update(special(tea.KeyEsc)); a != input.Cancelled {
		t.Errorf("esc should Cancel, got %v", a)
	}
}

func TestConfirmAnswers(t *testing.T) {
	c := input.Confirm{Question: "delete this?"}
	yes := []tea.KeyMsg{runes("y"), runes("Y"), special(tea.KeyEnter)}
	for _, k := range yes {
		if got := c.Update(k); got != input.Yes {
			t.Errorf("key %v should be Yes, got %v", k, got)
		}
	}
	no := []tea.KeyMsg{runes("n"), runes("N"), special(tea.KeyEsc)}
	for _, k := range no {
		if got := c.Update(k); got != input.No {
			t.Errorf("key %v should be No, got %v", k, got)
		}
	}
	if got := c.Update(runes("z")); got != input.Pending {
		t.Errorf("unrelated key should be Pending, got %v", got)
	}
}
