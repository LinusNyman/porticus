// Package porticus is the pantheon suite's shared TUI façade: the chrome every
// tool presents in common — palette, styles, two-pane layout, title bars, the
// footer hint bar, and the help plaque. It is the colonnade the whole suite
// shares; only the accent colour, sigil and name vary per tool (see Theme).
//
// porticus is presentation only: it has no dependency on the data spine
// (github.com/LinusNyman/pantheon) and no tool-specific types. Tools render
// their own rows (a todo, a contact, a habit) and hand the rendered strings to
// the layout helpers here. The canonical design rationale lives in the suite
// TUI design guide (aoapp_tui_design.md); this package is its executable form.
package porticus

import "github.com/charmbracelet/lipgloss"

// Shared base palette (suite standard §5). Fixed across every tool — do not
// vary these per tool. Only the accent (Theme.Accent) changes. Kept as exported
// constants so a tool can reach for a raw colour when a Style isn't enough.
const (
	ColCode      = "#4fa6e0" // aegean blue: node codes
	ColName      = "#f2ead0" // marble ivory: primary text, selection fg
	ColOwnCount  = "#3fc79a" // verdigris bronze: own-item counts
	ColDim       = "#9a917b" // weathered stone: dim text, hints, inactive
	ColDue       = "#f5a623" // ochre: due dates, soft urgency
	ColRecur     = "#bd6ad8" // Tyrian purple: recurrence specs
	ColOverdue   = "#e8492c" // Pompeian red: errors, overdue
	ColCompleted = "#aac63f" // laurel green: done marks, success
	ColHeading   = "#ee7f44" // terracotta: section headings
	ColDesc      = "#d4c79e" // aged parchment: descriptions
	ColSelFocus  = "#4a4231" // dark umber: focused selection bg
	ColSelBlur   = "#332d20" // deep umber: blurred selection bg
	ColSelFg     = "#f2ead0" // marble ivory: selection text
	ColDivider   = "#4a4231" // umber: pane divider, rules, plaque border
	ColPaneOff   = "#9a917b" // weathered stone: unfocused pane title
	ColToday     = "#4fa6e0" // aegean blue: the "today" calendar day, distinct from the accent focus
)

// Styles bundles every shared lipgloss style, built once per tool by NewStyles.
// The accent colour is the only per-tool input; every other style uses a fixed
// suite colour so look stays identical across tools. Field names mirror the
// suite standard §7 variable names so style code reads the same everywhere.
type Styles struct {
	Accent string // raw accent hex (Theme.Accent), for ad-hoc styling

	Title     lipgloss.Style // bold accent — left pane title, footer hederas
	Code      lipgloss.Style
	Name      lipgloss.Style
	Count     lipgloss.Style
	Dim       lipgloss.Style
	Done      lipgloss.Style // dim + strikethrough
	Due       lipgloss.Style
	Recur     lipgloss.Style
	Overdue   lipgloss.Style
	Completed lipgloss.Style
	Heading   lipgloss.Style // bold terracotta
	Desc      lipgloss.Style // parchment italic
	Err       lipgloss.Style // overdue colour — error text
	OK        lipgloss.Style // laurel green — status text

	SelFocus lipgloss.Style // focused selection: umber bg, ivory fg, bold
	SelBlur  lipgloss.Style // blurred selection: deep umber bg, ivory fg

	Divider    lipgloss.Style // ║ divider and ══ / ── rules
	PaneOn     lipgloss.Style // bold accent — focused pane title
	PaneOff    lipgloss.Style // dim — unfocused pane title
	HelpPlaque lipgloss.Style // double border in umber — the help tablet
	Today      lipgloss.Style // aegean blue — the "today" calendar day
}

// NewStyles builds the suite style set for one tool. accent is the tool's
// identity colour (Theme.Accent, suite standard §4); everything else is fixed.
func NewStyles(accent string) Styles {
	fg := func(hex string) lipgloss.Style {
		return lipgloss.NewStyle().Foreground(lipgloss.Color(hex))
	}
	return Styles{
		Accent:    accent,
		Title:     lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(accent)),
		Code:      fg(ColCode),
		Name:      fg(ColName),
		Count:     fg(ColOwnCount),
		Dim:       fg(ColDim),
		Done:      fg(ColDim).Strikethrough(true),
		Due:       fg(ColDue),
		Recur:     fg(ColRecur),
		Overdue:   fg(ColOverdue),
		Completed: fg(ColCompleted),
		Heading:   lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(ColHeading)),
		Desc:      fg(ColDesc).Italic(true),
		Err:       fg(ColOverdue),
		OK:        fg(ColCompleted),
		SelFocus:  lipgloss.NewStyle().Background(lipgloss.Color(ColSelFocus)).Foreground(lipgloss.Color(ColSelFg)).Bold(true),
		SelBlur:   lipgloss.NewStyle().Background(lipgloss.Color(ColSelBlur)).Foreground(lipgloss.Color(ColSelFg)),
		Divider:   fg(ColDivider),
		PaneOn:    lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(accent)),
		PaneOff:   fg(ColPaneOff),
		HelpPlaque: lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(lipgloss.Color(ColDivider)).
			Padding(0, 2),
		Today: fg(ColToday),
	}
}
