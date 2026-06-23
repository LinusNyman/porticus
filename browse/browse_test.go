package browse

import (
	"testing"

	"github.com/LinusNyman/pantheon/tree"
	"github.com/LinusNyman/porticus/keys"
	tea "github.com/charmbracelet/bubbletea"
)

// keyMsg builds a rune key message whose String() is s (e.g. "j", "l"), which
// is what key.Matches compares against the bindings.
func keyMsg(s string) tea.KeyMsg {
	return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(s)}
}

func defaultMap() keys.Map { return keys.Default() }

func wheel(b tea.MouseButton) tea.MouseMsg { return tea.MouseMsg{Button: b} }

// fixture builds a small forest:
//
//	a            (children a1, a2)
//	  a1         (leaf)
//	  a2         (children a2x)
//	    a2x      (leaf)
//	b            (leaf)
func fixture() *tree.Tree {
	a := &tree.Node{Code: "a", Name: "alpha"}
	a1 := &tree.Node{Code: "a1", Name: "alpha_one", Parent: a}
	a2 := &tree.Node{Code: "a2", Name: "alpha_two", Parent: a}
	a2x := &tree.Node{Code: "a2x", Name: "alpha_two_x", Parent: a2}
	a.Children = []*tree.Node{a1, a2}
	a2.Children = []*tree.Node{a2x}
	b := &tree.Node{Code: "b", Name: "beta"}
	return &tree.Tree{Roots: []*tree.Node{a, b}}
}

func codes(p TreePane) []string {
	out := make([]string, len(p.rows))
	for i, r := range p.rows {
		out[i] = r.Node.Code
	}
	return out
}

func eq(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestRebuildCollapsedShowsRootsOnly(t *testing.T) {
	p := New(fixture(), Opts{})
	if got := codes(p); !eq(got, []string{"a", "b"}) {
		t.Errorf("collapsed rows = %v, want [a b]", got)
	}
	if n := p.Selected(); n == nil || n.Code != "a" {
		t.Errorf("initial selection = %v, want a", n)
	}
}

func TestExpandRevealsChildren(t *testing.T) {
	p := New(fixture(), Opts{})
	// Expand the selected root a (l/right/enter).
	nav := p.Update(keyMsg("l"), defaultMap(), 1)
	if !nav.Handled || nav.CrossRight {
		t.Fatalf("expanding a node with children: nav = %+v", nav)
	}
	if got := codes(p); !eq(got, []string{"a", "a1", "a2", "b"}) {
		t.Errorf("after expand a, rows = %v, want [a a1 a2 b]", got)
	}
	// Selection is preserved on a by code across the rebuild.
	if n := p.Selected(); n == nil || n.Code != "a" {
		t.Errorf("selection after expand = %v, want a", n)
	}
}

func TestExpandLeafCrossesRight(t *testing.T) {
	p := New(fixture(), Opts{})
	p.cursor = 1 // b, a leaf root
	nav := p.Update(keyMsg("l"), defaultMap(), 1)
	if !nav.Handled || !nav.CrossRight {
		t.Errorf("expand on leaf b should signal CrossRight, got %+v", nav)
	}
}

func TestDownMovesAndReportsMoved(t *testing.T) {
	p := New(fixture(), Opts{})
	nav := p.Update(keyMsg("j"), defaultMap(), 1)
	if !nav.Handled || !nav.Moved {
		t.Fatalf("j on first row should move, got %+v", nav)
	}
	if n := p.Selected(); n == nil || n.Code != "b" {
		t.Errorf("after j, selection = %v, want b", n)
	}
}

func TestGoToRevealsDeepNode(t *testing.T) {
	f := fixture()
	a2x := f.Roots[0].Children[1].Children[0]
	p := New(f, Opts{})
	p.GoTo(a2x)
	if got := codes(p); !eq(got, []string{"a", "a1", "a2", "a2x", "b"}) {
		t.Errorf("after GoTo(a2x), rows = %v, want fully expanded a subtree", got)
	}
	if n := p.Selected(); n == nil || n.Code != "a2x" {
		t.Errorf("GoTo should select a2x, got %v", n)
	}
}

func TestOnlyWithItemsFilter(t *testing.T) {
	// Only a2/a2x carry items; the filter should drop a1 and b but keep the
	// ancestors of a kept node.
	keep := map[string]bool{"a": true, "a2": true, "a2x": true}
	p := New(fixture(), Opts{HasItems: func(n *tree.Node) bool { return keep[n.Code] }})
	p.expanded["a"] = true
	p.expanded["a2"] = true
	p.SetOnlyWithItems(true)
	if got := codes(p); !eq(got, []string{"a", "a2", "a2x"}) {
		t.Errorf("filtered rows = %v, want [a a2 a2x]", got)
	}
}

func TestCursorWindow(t *testing.T) {
	c := Cursor{Index: 0, Len: 10}
	start, end, above, below := c.Window(4)
	if start != 0 || end != 4 || above != 0 || below != 6 {
		t.Errorf("top window = (%d,%d,%d,%d), want (0,4,0,6)", start, end, above, below)
	}
	c.Index = 9 // bottom
	start, end, above, below = c.Window(4)
	if start != 6 || end != 10 || above != 6 || below != 0 {
		t.Errorf("bottom window = (%d,%d,%d,%d), want (6,10,6,0)", start, end, above, below)
	}
}

func TestCursorMouseWheel(t *testing.T) {
	c := Cursor{Index: 2, Len: 5}
	if !c.HandleMouse(wheel(tea.MouseButtonWheelDown)) || c.Index != 3 {
		t.Errorf("wheel down should move to 3, got %d", c.Index)
	}
	if !c.HandleMouse(wheel(tea.MouseButtonWheelUp)) || c.Index != 2 {
		t.Errorf("wheel up should move to 2, got %d", c.Index)
	}
	c.Index = 0
	if c.HandleMouse(wheel(tea.MouseButtonWheelUp)) {
		t.Error("wheel up at top should not move")
	}
}

func TestCursorReorder(t *testing.T) {
	items := []string{"a", "b", "c"}
	c := Cursor{Index: 0, Len: 3}
	swap := func(from, to int) bool {
		items[from], items[to] = items[to], items[from]
		return true
	}
	if !c.Reorder(+1, swap) || c.Index != 1 || items[0] != "b" {
		t.Errorf("reorder +1 failed: items=%v cursor=%d", items, c.Index)
	}
	c.Index = 0
	if c.Reorder(-1, swap) {
		t.Error("reorder up at the top should be a no-op")
	}
}

func TestTreePaneMouseWheel(t *testing.T) {
	p := New(fixture(), Opts{})
	if !p.HandleMouse(wheel(tea.MouseButtonWheelDown)) || p.Selected().Code != "b" {
		t.Errorf("wheel down should select b, got %v", p.Selected())
	}
}

func TestStatusGenerationGuard(t *testing.T) {
	var s Status
	s.Set("first")  // gen 1
	s.Set("second") // gen 2
	// A stale clear tick (gen 1) is consumed but must not wipe the newer message.
	if !s.Handle(clearMsg{gen: 1}) {
		t.Error("Handle should consume a clearMsg")
	}
	if s.Msg != "second" {
		t.Errorf("stale clear wiped current message, Msg = %q", s.Msg)
	}
	// The matching clear (gen 2) wipes it.
	s.Handle(clearMsg{gen: 2})
	if s.Msg != "" {
		t.Errorf("matching clear should wipe message, Msg = %q", s.Msg)
	}
}
