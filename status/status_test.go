package status

import (
	"strings"
	"testing"

	"github.com/LinusNyman/porticus"
	tea "github.com/charmbracelet/bubbletea"
)

func TestGenerationGuard(t *testing.T) {
	var l Line
	l.Set("first")  // gen 1
	l.Set("second") // gen 2
	// A stale clear tick (gen 1) is consumed but must not wipe the newer message.
	if !l.Handle(clearMsg{gen: 1}) {
		t.Error("Handle should consume a clearMsg")
	}
	if l.Text() != "second" {
		t.Errorf("stale clear wiped current message, Text = %q", l.Text())
	}
	// The matching clear (gen 2) wipes it.
	l.Handle(clearMsg{gen: 2})
	if l.Text() != "" {
		t.Errorf("matching clear should wipe message, Text = %q", l.Text())
	}
}

func TestHandleIgnoresOtherMessages(t *testing.T) {
	var l Line
	l.Set("hi")
	if l.Handle(tea.KeyMsg{}) {
		t.Error("Handle should not consume a non-clear message")
	}
	if l.Text() != "hi" {
		t.Errorf("a non-clear message must not clear the line, Text = %q", l.Text())
	}
}

func TestKinds(t *testing.T) {
	s := porticus.NewStyles("#e06474")
	var l Line

	// The zero value renders as nothing.
	if got := l.View(s, 40); got != "" {
		t.Errorf("zero-value line should render empty, got %q", got)
	}

	// Each setter selects its kind, stamps the text, returns an auto-clear command
	// so the line eventually fades, and renders the message within a wide field.
	for _, c := range []struct {
		name string
		set  func() tea.Cmd
		text string
		want Kind
	}{
		{"Set", func() tea.Cmd { return l.Set("saved") }, "saved", OK},
		{"SetInfo", func() tea.Cmd { return l.SetInfo("no matches") }, "no matches", Info},
		{"SetErr", func() tea.Cmd { return l.SetErr("disk full") }, "disk full", Error},
	} {
		if c.set() == nil {
			t.Errorf("%s should return an auto-clear command", c.name)
		}
		if l.Kind() != c.want || l.Text() != c.text {
			t.Errorf("%s: kind=%v text=%q, want %v/%q", c.name, l.Kind(), l.Text(), c.want, c.text)
		}
		if out := l.View(s, 40); !strings.Contains(out, c.text) {
			t.Errorf("%s: View should contain the message, got %q", c.name, out)
		}
	}
}
