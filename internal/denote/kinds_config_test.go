package denote_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mph-llm-experiments/anote/internal/denote"
)

func TestDefaultKindsConfig_ContainsExpectedKinds(t *testing.T) {
	cfg := denote.DefaultKindsConfig()
	expected := []string{"aspiration", "belief", "plan", "note", "fact", "purpose"}
	for _, kind := range expected {
		if _, ok := cfg.Kinds[kind]; !ok {
			t.Errorf("expected default config to contain kind %q", kind)
		}
	}
}

func TestKindsConfig_WriteAndLoad(t *testing.T) {
	dir := t.TempDir()
	cfg := denote.DefaultKindsConfig()

	if err := cfg.WriteToDir(dir); err != nil {
		t.Fatalf("WriteToDir failed: %v", err)
	}

	loaded, err := denote.LoadKindsConfig(dir)
	if err != nil {
		t.Fatalf("LoadKindsConfig failed: %v", err)
	}

	if len(loaded.Kinds) != len(cfg.Kinds) {
		t.Errorf("expected %d kinds, got %d", len(cfg.Kinds), len(loaded.Kinds))
	}
}

func TestKindsConfig_LoadMissingWritesDefault(t *testing.T) {
	dir := t.TempDir()
	cfg, err := denote.LoadKindsConfig(dir)
	if err != nil {
		t.Fatalf("expected LoadKindsConfig to succeed and write defaults, got: %v", err)
	}
	if cfg == nil {
		t.Fatal("expected non-nil config")
	}
	if _, err := os.Stat(filepath.Join(dir, "kinds.json")); err != nil {
		t.Error("expected kinds.json to be written on first load")
	}
}

func TestKindsConfig_IsCompliant(t *testing.T) {
	cfg := denote.DefaultKindsConfig()

	tests := []struct {
		kind, state string
		want        bool
	}{
		{"aspiration", "seed", true},
		{"aspiration", "active", true},
		{"aspiration", "considering", false}, // belief display label, not canonical
		{"belief", "seed", true},
		{"belief", "active", true},
		{"note", "active", true},
		{"note", "iterating", false},
		{"purpose", "active", true},
		{"purpose", "seed", false},
		{"unknown_kind", "active", false},
	}

	for _, tt := range tests {
		got := cfg.IsCompliant(tt.kind, tt.state)
		if got != tt.want {
			t.Errorf("IsCompliant(%q, %q) = %v, want %v", tt.kind, tt.state, got, tt.want)
		}
	}
}

func TestKindsConfig_ValidStatesFor(t *testing.T) {
	cfg := denote.DefaultKindsConfig()
	states := cfg.ValidStatesFor("belief")
	if len(states) == 0 {
		t.Error("expected non-empty states for belief")
	}
	found := false
	for _, s := range states {
		if s == "seed" {
			found = true
		}
	}
	if !found {
		t.Error("expected 'seed' in valid states for belief")
	}
}

func TestKindsConfig_DefaultStateFor(t *testing.T) {
	cfg := denote.DefaultKindsConfig()
	if cfg.DefaultStateFor("aspiration") != "seed" {
		t.Error("expected default state 'seed' for aspiration")
	}
	if cfg.DefaultStateFor("note") != "active" {
		t.Error("expected default state 'active' for note")
	}
	if cfg.DefaultStateFor("unknown") != "seed" {
		t.Error("expected fallback default state 'seed' for unknown kind")
	}
}

func TestKindsConfig_PurposeRequired(t *testing.T) {
	cfg := denote.DefaultKindsConfig()
	if !cfg.PurposeRequired("aspiration") {
		t.Error("expected purpose_required=true for aspiration")
	}
	if cfg.PurposeRequired("note") {
		t.Error("expected purpose_required=false for note")
	}
	if cfg.PurposeRequired("purpose") {
		t.Error("expected purpose_required=false for purpose")
	}
}

func TestKindsConfig_AllKinds(t *testing.T) {
	cfg := denote.DefaultKindsConfig()
	kinds := cfg.AllKinds()
	if len(kinds) < 6 {
		t.Errorf("expected at least 6 kinds, got %d", len(kinds))
	}
}

func TestKindsConfig_KindExists(t *testing.T) {
	cfg := denote.DefaultKindsConfig()
	if !cfg.KindExists("aspiration") {
		t.Error("expected aspiration to exist")
	}
	if cfg.KindExists("bogus") {
		t.Error("expected bogus to not exist")
	}
}
