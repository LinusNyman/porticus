# porticus

The pantheon suite's shared TUI façade — the colonnade every tool presents in
common. It holds what is identical across the suite so it lives in one place
instead of being copy-pasted into each tool.

The **root package** (`porticus`) is the dependency-light chrome — lipgloss
only, no bubbletea, no data spine:

- **Palette & styles** — the fixed base palette and the lipgloss style set
  (`NewStyles(accent)`); the accent colour is the only per-tool input.
- **Identity** — `Theme{Name, Sigil, Accent}` and the canonical `Tools` table.
- **Layout** — `TwoPane` / `Stacked`, pane widths, the `║` divider, `══`/`──`
  pane headers.
- **Footer** — the `Hints` bar (groups joined by `❧`, hints by `·`, wrapped).
- **Help** — the carved-stone `HelpPage` plaque.
- **Text helpers** — `SpacedCaps`, `PadTo`, `Truncate`, `WrapRows`, `ScrollHint`.

**Sub-packages** add the shared interactive layer that the chrome alone can't
cover (each opt-in, imported only where needed):

- **`porticus/keys`** — the canonical key grammar (`Map`, `Default()`), the
  view-key helper (`View`), and the help-group generator (`HelpGroups`), so
  navigation, pane switching, and view selection are identical suite-wide and
  the help screen can't drift from the bindings. Depends on bubbletea, not the
  spine.
- **`porticus/browse`** — the stateful two-pane tree-browser scaffold a tool
  embeds: `TreePane` (the left-pane node tree), `Cursor` (list selection +
  windowing), and `Status` (the auto-clearing status line). Depends on the data
  spine (`pantheon/tree`) for the node type, kept here in its own package.
- **`porticus/insights`** — shared charts (`Sparkline`, `Heatmap`,
  `CalendarHeatmap`, `Compute`/`Trend`) and a scrollable `InsightsPage`, for the
  graph/stats screens (album demographics, speculum habit stats). Pure render.
- **`porticus/pick`** — the search/suggest overlay: a generic
  `Picker[T]` (query input + live-filtered, scrollable, selectable results) the
  tool drives with `Open`/`Update`/`View`, supplying `Filter`/`Render` and what a
  selection does. One component for both free-text search and ranked suggest (set
  `Limit`). Depends on bubbletea + the bubbles text input, not the spine.
- **`porticus/input`** — text-entry chrome: `Editor` (the soft-wrapping,
  content-sized add/edit textarea with the hanging-indent prompt), `Field` (the
  single-line input for dates/codes/names), and `Confirm` (the y/n prompt). All
  share enter-commits / esc-cancels and report an `Action`/`Answer` the tool acts
  on. Spine-free.
- **`porticus/calendar`** — the month calendar: a `Grid` (selected day +
  day/week/month navigation + per-day marker callback + selected/today
  highlight), the shared base for calendarium and any date view. The day-detail
  list is the tool's (via `browse.Cursor`). Spine-free.
- **`porticus/dates`** — the relative-date vocabulary: `RelativeDate` /
  `FarOff` / `CoarseSpan` render an ISO date as `today` / `tmrw` / `in 3d` /
  `in 2mo`, so dates read identically suite-wide. Pure.

The root also has **`Styles.Page`** — the shared scrollable read-only page chrome
(header + windowed body + scroll hint), behind `HelpPage`/`InsightsPage` and
reusable as a `Pager` for custom screens (a markdown preview, a stats page) — and
the `browse`/`pick` cursors take wheel scrolling (`HandleMouse`) and reordering
(`Cursor.Reorder`).

## Layering

```
              bubbletea + lipgloss
                       │
   pantheon ───────────┤   porticus (root: chrome, presentation only)
   (data spine)        │       ├── keys      (no spine)
        │              │       ├── insights  (no spine)
        │              │       ├── pick      (no spine)
        │              │       ├── input     (no spine)
        │              │       ├── calendar  (no spine)
        │              │       ├── dates     (no spine)
        │              ├───────┴── browse    (uses pantheon/tree)
        └──────────────┴───────────────┐
                                    each tool
```

The **root** package is presentation only — no spine import, no tool-specific
types. A tool renders its own rows (a todo, a contact, a habit) and hands the
strings to porticus's layout helpers, supplying its `Theme`. The only package
that touches the spine is `browse`, which needs the `tree.Node` type for the
shared tree pane; tools that aren't tree-shaped (e.g. speculum) use the root,
`keys`, and `insights` without it.

The authoritative design rationale is the suite TUI design guide
(`aoapp_tui_design.md`); this module is its executable form. When the two drift,
fix it here once and record the decision in the guide's §9 log.

## Use

```go
theme := porticus.Tools["album"]      // or a Theme literal
st := theme.Styles()                  // == porticus.NewStyles(theme.Accent)

body := st.TwoPane(width, height,
    func(w, h int) string { return renderListPane(st, w, h) },
    func(w, h int) string { return renderDetailPane(st, w, h) },
)
footer := st.Hints([][]string{{"j/k:move", "tab:pane"}, {"?:help", "q:quit"}}, width)
```

## Build

```
GOWORK=off go build ./... && GOWORK=off go test ./...
```

The suite `go.work` makes local edits here visible to every tool immediately; on
release, tag the module and bump each tool's pin.
