// Package browse is the pantheon suite's shared two-pane tree-browser scaffold:
// the stateful left-pane node tree (TreePane) and the list selection cursor
// (Cursor). Each is embedded by a tool's own bubbletea Model — porticus owns the
// navigation grammar, geometry, and rendering so the tree/working screen behaves
// identically across tools, while the tool keeps its own data and working-pane
// content. (The transient status line lives in porticus/status, spine-free.)
//
// browse depends on the data spine (github.com/LinusNyman/pantheon/tree) for
// the node type and on bubbletea for messages/commands. Per the suite rule,
// this spine-coupled, interactive layer lives in its own package, separate from
// the dependency-light porticus chrome.
package browse

import (
	"fmt"
	"strings"

	"github.com/LinusNyman/pantheon/tree"
	"github.com/LinusNyman/porticus"
	"github.com/LinusNyman/porticus/keys"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Row is one visible tree line: a node at an indentation depth.
type Row struct {
	Node  *tree.Node
	Depth int
}

// Opts configures a TreePane's tool-specific behaviour without coupling it to
// any tool type.
type Opts struct {
	// Annotate returns the trailing badge for a node row — e.g. a count like
	// "(4) +12" (see porticus.Styles.NodeCount) or an overdue mark, already
	// styled. The badge is right-aligned at the pane's right edge: a long node
	// name is truncated before the badge is, so the count is always visible.
	// porticus inserts the gap before it, so return the badge content alone with
	// no leading padding. May be nil for a bare tree.
	Annotate func(n *tree.Node) string
	// HasItems reports whether a node survives the "only nodes with items"
	// filter (true when the node or any descendant has items). Consulted only
	// while the filter is on; nil means every node passes.
	HasItems func(n *tree.Node) bool
}

// TreePane is the shared left-pane node tree. It owns the forest, the expansion
// state, the selection cursor, and the scroll window, and renders through the
// porticus chrome so every tool's tree looks and navigates identically. The
// tool supplies per-node badges and the items filter via Opts; the working-pane
// content (todos, contacts, …) stays in the tool.
type TreePane struct {
	forest   *tree.Tree
	opts     Opts
	expanded map[string]bool
	rows     []Row
	cursor   int
	onlyWith bool
}

// New builds a TreePane over forest and lays out its rows (roots collapsed).
func New(forest *tree.Tree, o Opts) TreePane {
	p := TreePane{forest: forest, opts: o, expanded: map[string]bool{}}
	p.Rebuild()
	return p
}

// Nav reports what a key press did so the tool can react: refresh its working
// pane when the selection Moved, or move focus into it on CrossRight.
type Nav struct {
	Handled    bool // the key was a tree key and was consumed
	Moved      bool // the selected node changed -> refresh the working pane
	CrossRight bool // expand pressed on a leaf/expanded node -> focus the working pane
}

// Update applies a navigation/interaction key. pageStep is the half-screen jump
// (see PageStep): j/k/g/G and half-page move the cursor; expand (l/right/enter)
// opens a collapsed node or, on a leaf/already-
// expanded node, signals CrossRight; collapse (h/left) closes an expanded node
// or jumps to the parent.
func (p *TreePane) Update(msg tea.KeyMsg, km keys.Map, pageStep int) Nav {
	prev := p.cursor
	switch {
	case key.Matches(msg, km.Down):
		p.cursor++
	case key.Matches(msg, km.Up):
		p.cursor--
	case key.Matches(msg, km.Top):
		p.cursor = 0
	case key.Matches(msg, km.Bottom):
		p.cursor = len(p.rows) - 1
	case key.Matches(msg, km.HalfDown):
		p.cursor += pageStep
	case key.Matches(msg, km.HalfUp):
		p.cursor -= pageStep
	case key.Matches(msg, km.Expand):
		n := p.Selected()
		if n == nil {
			return Nav{}
		}
		if !p.expanded[n.Code] && len(n.Children) > 0 {
			p.expanded[n.Code] = true
			p.Rebuild()
			return Nav{Handled: true}
		}
		// A leaf or an already-expanded node: cross into the working pane.
		return Nav{Handled: true, CrossRight: true}
	case key.Matches(msg, km.Collapse):
		n := p.Selected()
		if n == nil {
			return Nav{}
		}
		if p.expanded[n.Code] && len(n.Children) > 0 {
			p.expanded[n.Code] = false
			p.Rebuild()
			return Nav{Handled: true}
		}
		if n.Parent != nil {
			for i, r := range p.rows {
				if r.Node == n.Parent {
					p.cursor = i
					break
				}
			}
		}
		p.clamp()
		return Nav{Handled: true, Moved: p.cursor != prev}
	default:
		return Nav{}
	}
	p.clamp()
	return Nav{Handled: true, Moved: p.cursor != prev}
}

// HandleMouse moves the tree selection on a wheel scroll (up/down by one) and
// reports whether it moved — so the wheel scrolls the tree like j/k. The caller
// refreshes the working pane when it returns true, as after a key move.
func (p *TreePane) HandleMouse(msg tea.MouseMsg) bool {
	prev := p.cursor
	switch msg.Button {
	case tea.MouseButtonWheelUp:
		p.cursor--
	case tea.MouseButtonWheelDown:
		p.cursor++
	default:
		return false
	}
	p.clamp()
	return p.cursor != prev
}

// Rebuild lays out the visible rows from the forest, honouring the expansion
// state and the items filter, with no disk reads. It preserves the selection by
// node code across the rebuild, falling back to a clamped index.
func (p *TreePane) Rebuild() {
	prevCode := ""
	if n := p.Selected(); n != nil {
		prevCode = n.Code
	}
	p.rows = p.rows[:0]
	for _, r := range p.forest.Roots {
		p.appendTree(r, 0)
	}
	if prevCode != "" {
		for i, r := range p.rows {
			if r.Node.Code == prevCode {
				p.cursor = i
				break
			}
		}
	}
	p.clamp()
}

func (p *TreePane) appendTree(n *tree.Node, depth int) {
	if p.onlyWith && p.opts.HasItems != nil && !p.opts.HasItems(n) {
		return
	}
	p.rows = append(p.rows, Row{Node: n, Depth: depth})
	if !p.expanded[n.Code] {
		return
	}
	for _, c := range n.Children {
		p.appendTree(c, depth+1)
	}
}

func (p *TreePane) clamp() {
	if p.cursor >= len(p.rows) {
		p.cursor = len(p.rows) - 1
	}
	if p.cursor < 0 {
		p.cursor = 0
	}
}

// Selected returns the node under the cursor, or nil when the tree is empty.
func (p TreePane) Selected() *tree.Node {
	if p.cursor < 0 || p.cursor >= len(p.rows) {
		return nil
	}
	return p.rows[p.cursor].Node
}

// GoTo reveals n: expands its ancestors, rebuilds, and places the cursor on it.
// Shared by go-to-by-code and any "reveal this node" caller (search/suggest).
func (p *TreePane) GoTo(n *tree.Node) {
	for a := n.Parent; a != nil; a = a.Parent {
		p.expanded[a.Code] = true
	}
	p.Rebuild()
	for i, r := range p.rows {
		if r.Node == n {
			p.cursor = i
			break
		}
	}
	p.clamp()
}

// SetForest swaps in a freshly scanned forest (after a reload) and rebuilds,
// preserving the selection by code where possible.
func (p *TreePane) SetForest(forest *tree.Tree) {
	p.forest = forest
	p.Rebuild()
}

// SetOnlyWithItems toggles the "only nodes with items" filter and rebuilds.
func (p *TreePane) SetOnlyWithItems(on bool) {
	p.onlyWith = on
	p.Rebuild()
}

// OnlyWithItems reports the current filter state.
func (p TreePane) OnlyWithItems() bool { return p.onlyWith }

// View renders the pane: the porticus left header (sigil + tool name + selected
// node) over the ══ rule, the windowed rows with the selection highlighted
// (umber when focused, deep umber when blurred), and a scroll hint. It returns
// exactly height lines so the two-pane layout stays aligned.
func (p TreePane) View(s porticus.Styles, t porticus.Theme, width, height int, focused bool) string {
	nodeName := ""
	if n := p.Selected(); n != nil {
		nodeName = n.Display()
	}
	// LeftHeader returns the title line and the ══ rule as two lines.
	lines := strings.Split(s.LeftHeader(t, nodeName, width), "\n")

	maxRows := height - 2
	if maxRows < 1 {
		maxRows = 1
	}
	// Reserve a row for the scroll hint when the list overflows.
	if len(p.rows) > maxRows && maxRows > 1 {
		maxRows--
	}
	start, end, above, below := Cursor{Index: p.cursor, Len: len(p.rows)}.Window(maxRows)

	for i := start; i < end; i++ {
		line := p.renderRow(s, i, width)
		if i == p.cursor {
			line = porticus.PadTo(line, width)
			if focused {
				line = s.SelFocus.Render(line)
			} else {
				line = s.SelBlur.Render(line)
			}
		}
		lines = append(lines, line)
	}
	if hint := s.ScrollHint(above, below, width); hint != "" {
		lines = append(lines, hint)
	}
	for len(lines) < height {
		lines = append(lines, "")
	}
	return strings.Join(lines, "\n")
}

func (p TreePane) renderRow(s porticus.Styles, i, width int) string {
	r := p.rows[i]
	indent := strings.Repeat("  ", r.Depth)
	arrow := " "
	if len(r.Node.Children) > 0 {
		if p.expanded[r.Node.Code] {
			arrow = "▾"
		} else {
			arrow = "▸"
		}
	}
	left := fmt.Sprintf("%s%s %s  %s",
		indent, arrow,
		s.Code.Render(r.Node.Code),
		s.Name.Render(r.Node.Display()))

	badge := ""
	if p.opts.Annotate != nil {
		badge = p.opts.Annotate(r.Node)
	}
	if badge == "" {
		return porticus.Truncate(left, width)
	}
	// Right-align the badge: the count is the salient datum, so reserve it at
	// the pane's right edge and truncate the (variable-length) name rather than
	// let a long name push the badge off the line (the bug this fixed). Keep a
	// one-cell gap before it.
	avail := width - lipgloss.Width(badge) - 1
	if avail < 1 {
		// Too narrow to show both; the badge wins.
		return porticus.Truncate(badge, width)
	}
	return porticus.PadTo(porticus.Truncate(left, avail), avail) + " " + badge
}
