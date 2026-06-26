package keys_test

import (
	"testing"

	"github.com/LinusNyman/porticus"
	"github.com/LinusNyman/porticus/keys"
)

func TestViewMapping(t *testing.T) {
	m := keys.Default()
	cases := map[string]int{
		"1": 1, "2": 2, "3": 3, "9": 9,
		"0": 0, "a": 0, "": 0, "10": 0, "tab": 0, "enter": 0,
	}
	for k, want := range cases {
		if got := m.View(k); got != want {
			t.Errorf("View(%q) = %d, want %d", k, got, want)
		}
	}
}

func TestHelpGroupsStandardOrder(t *testing.T) {
	m := keys.Default()
	extra := porticus.HelpGroup{Title: "Todos", Rows: [][2]string{{"a", "add todo"}}}
	groups := m.HelpGroups([]string{"list", "agenda", "calendar"}, extra)

	if len(groups) != 3 {
		t.Fatalf("got %d groups, want 3 (Navigate, View, Todos)", len(groups))
	}
	if groups[0].Title != "Navigate" {
		t.Errorf("first group = %q, want Navigate", groups[0].Title)
	}
	if groups[1].Title != "View" {
		t.Errorf("second group = %q, want View", groups[1].Title)
	}
	if groups[2].Title != "Todos" {
		t.Errorf("third group = %q, want the appended Todos", groups[2].Title)
	}
}

func TestHelpGroupsViewKeysContiguousFromOne(t *testing.T) {
	groups := keys.Default().HelpGroups([]string{"list", "agenda", "calendar"})
	view := groups[1]
	// The cover (+:title) always leads, then the numbered views from 1.
	want := [][2]string{{"+", "title"}, {"1", "list"}, {"2", "agenda"}, {"3", "calendar"}}
	if len(view.Rows) != len(want) {
		t.Fatalf("got %d view rows, want %d", len(view.Rows), len(want))
	}
	for i, w := range want {
		if view.Rows[i] != w {
			t.Errorf("view row %d = %v, want %v", i, view.Rows[i], w)
		}
	}
}

func TestHelpGroupsNoViews(t *testing.T) {
	groups := keys.Default().HelpGroups(nil)
	// Even with no numbered views, the cover screen gives a View group of one row.
	if len(groups) != 2 || groups[1].Title != "View" {
		t.Fatalf("want Navigate + View groups, got %d", len(groups))
	}
	if rows := groups[1].Rows; len(rows) != 1 || rows[0] != [2]string{"+", "title"} {
		t.Fatalf("View group should hold only the +:title row, got %v", groups[1].Rows)
	}
}
