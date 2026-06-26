package porticus

// Theme is a tool's whole visual identity: the three things that vary across the
// suite (suite standard §4). Everything else — palette, layout, key grammar — is
// shared. Build a tool's styles with NewStyles(theme.Accent).
type Theme struct {
	Name    string // tool name, e.g. "pensum"; rendered spaced-caps in titles
	Sigil   string // identity glyph, e.g. "✎"
	Accent  string // accent hex, e.g. "#e06474"
	Version string // optional app version shown on the title screen, e.g. "v1.2.3"; "" hides it
	Tagline string // optional one-sentence description shown on the title screen; "" hides it
}

// Author is the suite's author, shown on the title screen (TitlePage). It is the
// same across every tool, so it lives here rather than on Theme.
const Author = "Linus Nyman"

// Styles is a convenience for theme.Styles() == NewStyles(theme.Accent).
func (t Theme) Styles() Styles { return NewStyles(t.Accent) }

// WithVersion returns a copy of the theme with Version set, so a tool can stamp
// its build version onto its identity in one line, e.g.
// porticus.Tools["album"].WithVersion(version). The Tools table itself carries no
// version (it's runtime info).
func (t Theme) WithVersion(v string) Theme {
	t.Version = v
	return t
}

// WithTagline returns a copy of the theme with Tagline set — a tool's one-sentence
// description for the title screen, e.g. porticus.Tools["pensum"].WithTagline(
// "the work you set yourself, kept in order"). The Tools table leaves it empty;
// each tool supplies its own copy, the same way it supplies its build Version.
func (t Theme) WithTagline(s string) Theme {
	t.Tagline = s
	return t
}

// Tools is the canonical per-tool identity table (suite standard §4), kept in
// one place so a tool can pull its identity by name rather than hard-coding the
// hex and glyph. A tool may also construct its own Theme literal.
var Tools = map[string]Theme{
	"pantheon":    {Name: "pantheon", Sigil: "✦", Accent: "#c8c0b0"},
	"pensum":      {Name: "pensum", Sigil: "✎", Accent: "#e06474"},
	"tabella":     {Name: "tabella", Sigil: "≡", Accent: "#4fa6e0"},
	"decreta":     {Name: "decreta", Sigil: "⚖", Accent: "#c45f9c"},
	"speculum":    {Name: "speculum", Sigil: "○", Accent: "#3fc79a"},
	"studium":     {Name: "studium", Sigil: "⊙", Accent: "#a78bfa"},
	"atrium":      {Name: "atrium", Sigil: "◈", Accent: "#ee7f44"},
	"album":       {Name: "album", Sigil: "❦", Accent: "#f5a623"},
	"fasti": {Name: "fasti", Sigil: "⊕", Accent: "#5b7dc8"},
}
