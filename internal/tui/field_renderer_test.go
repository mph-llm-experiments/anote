package tui_test

import (
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/mph-llm-experiments/anote/internal/tui"
	"github.com/muesli/termenv"
)

func TestKindSymbol_KnownKinds(t *testing.T) {
	tests := []struct{ kind, want string }{
		{"aspiration", tui.SymbolAspiration},
		{"belief", tui.SymbolBelief},
		{"plan", tui.SymbolPlan},
		{"note", tui.SymbolNote},
		{"fact", tui.SymbolFact},
		{"purpose", tui.SymbolPurpose},
		{"unknown", "?"},
	}
	for _, tt := range tests {
		got := tui.KindSymbol(tt.kind)
		if got != tt.want {
			t.Errorf("KindSymbol(%q) = %q, want %q", tt.kind, got, tt.want)
		}
	}
}

func TestStateSymbol_KnownStates(t *testing.T) {
	states := []string{"seed", "draft", "active", "iterating", "implemented", "archived", "rejected", "dropped"}
	for _, s := range states {
		got := tui.StateSymbol(s)
		if got == "" || got == "?" {
			t.Errorf("StateSymbol(%q) returned %q, expected a real symbol", s, got)
		}
	}
}

func TestStateSymbol_Unknown(t *testing.T) {
	got := tui.StateSymbol("bogus")
	if got != "?" {
		t.Errorf("expected '?' for unknown state, got %q", got)
	}
}

func TestFormatListRow_NotEmpty(t *testing.T) {
	row := tui.FormatListRow("My Idea", "aspiration", "active", "walk", "photography", false, 80)
	if row == "" {
		t.Error("expected non-empty row")
	}
}

func TestFormatListRow_Selected(t *testing.T) {
	// Force color output so lipgloss ANSI codes are emitted in a non-TTY test environment.
	lipgloss.SetColorProfile(termenv.ANSI256)
	t.Cleanup(func() { lipgloss.SetColorProfile(termenv.Ascii) })

	normal := tui.FormatListRow("Test", "note", "active", "", "", false, 80)
	selected := tui.FormatListRow("Test", "note", "active", "", "", true, 80)
	// Selected and normal should differ (different styling applied)
	if normal == selected {
		t.Error("expected selected and normal rows to render differently")
	}
}
