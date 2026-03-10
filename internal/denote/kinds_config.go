package denote

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
)

const kindsConfigFilename = "kinds.json"

// KindEntry defines the valid states and behavior for a single kind.
type KindEntry struct {
	States          []string `json:"states"`
	Terminal        []string `json:"terminal"`
	Default         string   `json:"default"`
	PurposeRequired bool     `json:"purpose_required"`
}

// KindsConfig is the full configuration loaded from kinds.json.
type KindsConfig struct {
	Kinds map[string]KindEntry `json:"kinds"`
}

// DefaultKindsConfig returns the built-in default configuration.
// Uses canonical state values from types.go.
func DefaultKindsConfig() *KindsConfig {
	return &KindsConfig{
		Kinds: map[string]KindEntry{
			KindAspiration: {
				States:          []string{StateSeed, StateDraft, StateActive, StateIterating, StateImplemented, StateDropped},
				Terminal:        []string{StateImplemented, StateDropped},
				Default:         StateSeed,
				PurposeRequired: true,
			},
			KindBelief: {
				States:          []string{StateSeed, StateActive, StateIterating, StateImplemented, StateRejected},
				Terminal:        []string{StateImplemented, StateRejected},
				Default:         StateSeed,
				PurposeRequired: true,
			},
			KindPlan: {
				States:          []string{StateSeed, StateDraft, StateActive, StateIterating, StateImplemented, StateDropped},
				Terminal:        []string{StateImplemented, StateDropped},
				Default:         StateSeed,
				PurposeRequired: true,
			},
			KindNote: {
				States:          []string{StateActive, StateArchived},
				Terminal:        []string{StateArchived},
				Default:         StateActive,
				PurposeRequired: false,
			},
			KindFact: {
				States:          []string{StateActive, StateArchived},
				Terminal:        []string{StateArchived},
				Default:         StateActive,
				PurposeRequired: false,
			},
			KindPurpose: {
				States:          []string{StateActive, StateArchived},
				Terminal:        []string{StateArchived},
				Default:         StateActive,
				PurposeRequired: false,
			},
		},
	}
}

// LoadKindsConfig reads kinds.json from dir. If missing, writes defaults and returns them.
func LoadKindsConfig(dir string) (*KindsConfig, error) {
	path := filepath.Join(dir, kindsConfigFilename)
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		cfg := DefaultKindsConfig()
		if writeErr := cfg.WriteToDir(dir); writeErr != nil {
			return nil, writeErr
		}
		return cfg, nil
	}
	if err != nil {
		return nil, err
	}
	var cfg KindsConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// WriteToDir writes the config to kinds.json in dir.
func (kc *KindsConfig) WriteToDir(dir string) error {
	data, err := json.MarshalIndent(kc, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, kindsConfigFilename), data, 0644)
}

// IsCompliant returns true if state is a valid canonical state for kind.
func (kc *KindsConfig) IsCompliant(kind, state string) bool {
	entry, ok := kc.Kinds[kind]
	if !ok {
		return false
	}
	for _, s := range entry.States {
		if s == state {
			return true
		}
	}
	return false
}

// ValidStatesFor returns the valid canonical states for a kind.
// Returns nil if the kind is not in the config.
func (kc *KindsConfig) ValidStatesFor(kind string) []string {
	if entry, ok := kc.Kinds[kind]; ok {
		return entry.States
	}
	return nil
}

// DefaultStateFor returns the default state for a kind.
// Falls back to StateSeed if the kind is unknown.
func (kc *KindsConfig) DefaultStateFor(kind string) string {
	if entry, ok := kc.Kinds[kind]; ok && entry.Default != "" {
		return entry.Default
	}
	return StateSeed
}

// PurposeRequired returns true if ideas of this kind require a purpose.
func (kc *KindsConfig) PurposeRequired(kind string) bool {
	if entry, ok := kc.Kinds[kind]; ok {
		return entry.PurposeRequired
	}
	return false
}

// KindExists returns true if kind is in the config.
func (kc *KindsConfig) KindExists(kind string) bool {
	_, ok := kc.Kinds[kind]
	return ok
}

// AllKinds returns all kind names in the config, sorted alphabetically.
func (kc *KindsConfig) AllKinds() []string {
	kinds := make([]string, 0, len(kc.Kinds))
	for k := range kc.Kinds {
		kinds = append(kinds, k)
	}
	sort.Strings(kinds)
	return kinds
}
