package denote

import (
	"fmt"
	"strings"
	"time"

	"github.com/mph-llm-experiments/acore"
)

// State constants for the idea lifecycle.
const (
	StateSeed        = "seed"
	StateDraft       = "draft"
	StateActive      = "active"
	StateIterating   = "iterating"
	StateImplemented = "implemented"
	StateArchived    = "archived"
	StateRejected    = "rejected"
	StateDropped     = "dropped"
)

// Maturity constants (orthogonal to state).
const (
	MaturityCrawl = "crawl"
	MaturityWalk  = "walk"
	MaturityRun   = "run"
)

// Kind constants (orthogonal to state and maturity).
const (
	KindAspiration = "aspiration"
	KindBelief     = "belief"
	KindPlan       = "plan"
)

// Type constant.
const TypeIdea = "idea"

// IdeaMetadata holds domain-specific idea fields.
// Common fields (ID, Title, IndexID, Type, Tags, Created, Modified,
// RelatedPeople, RelatedTasks, RelatedIdeas) come from embedded acore.Entity.
type IdeaMetadata struct {
	Kind           string `yaml:"kind,omitempty" json:"kind,omitempty"`
	State          string `yaml:"state,omitempty" json:"state,omitempty"`
	Maturity       string `yaml:"maturity,omitempty" json:"maturity,omitempty"`
	RejectedReason string `yaml:"rejected_reason,omitempty" json:"rejected_reason,omitempty"`
}

// Idea combines acore.Entity with idea-specific metadata and content.
type Idea struct {
	acore.Entity `yaml:",inline"`
	IdeaMetadata `yaml:",inline"`
	Content      string    `yaml:"-" json:"-"`
	ModTime      time.Time `yaml:"-" json:"-"`
}

// IsValidState checks if a state value is valid.
func IsValidState(state string) bool {
	switch state {
	case StateSeed, StateDraft, StateActive, StateIterating,
		StateImplemented, StateArchived, StateRejected, StateDropped:
		return true
	}
	return false
}

// IsValidMaturity checks if a maturity value is valid.
func IsValidMaturity(maturity string) bool {
	switch maturity {
	case MaturityCrawl, MaturityWalk, MaturityRun:
		return true
	}
	return false
}

// allowedTransitions defines which state transitions are valid per the spec.
var allowedTransitions = map[string][]string{
	StateSeed:      {StateDraft},
	StateDraft:     {StateActive},
	StateActive:    {StateIterating, StateImplemented, StateArchived, StateRejected, StateDropped},
	StateIterating: {StateActive, StateImplemented, StateArchived, StateRejected, StateDropped},
	StateArchived:  {StateActive},
	// Terminal states with no outbound transitions:
	// StateImplemented, StateRejected, StateDropped
}

// ValidateStateTransition checks if a state transition is allowed.
func ValidateStateTransition(from, to string) error {
	if !IsValidState(from) {
		return fmt.Errorf("invalid source state: %s", from)
	}
	if !IsValidState(to) {
		return fmt.Errorf("invalid target state: %s", to)
	}

	allowed, ok := allowedTransitions[from]
	if !ok {
		return fmt.Errorf("no transitions allowed from terminal state %q", from)
	}

	for _, s := range allowed {
		if s == to {
			return nil
		}
	}

	return fmt.Errorf("transition from %q to %q is not allowed", from, to)
}

// IsValidKind checks if a kind value is valid.
func IsValidKind(kind string) bool {
	switch kind {
	case KindAspiration, KindBelief, KindPlan:
		return true
	}
	return false
}

// displayLabels maps canonical states to kind-specific display labels.
// Only states that differ between kinds are listed.
var displayLabels = map[string]map[string]string{
	KindBelief: {
		StateActive:      "considering",
		StateIterating:   "reconsidering",
		StateImplemented: "accepted",
	},
	KindPlan: {
		StateActive:      "committed",
		StateIterating:   "replanning",
		StateImplemented: "completed",
	},
}

// DisplayState returns the kind-specific display label for a canonical state.
func DisplayState(canonical, kind string) string {
	if kindMap, ok := displayLabels[kind]; ok {
		if label, ok := kindMap[canonical]; ok {
			return label
		}
	}
	return canonical
}

// ResolveDisplayState converts a display label (possibly kind-specific) back
// to the canonical state value. Returns the canonical state and the kind it
// belongs to (empty string if the label is shared or already canonical).
func ResolveDisplayState(display string) (canonical string, matchedKind string) {
	if IsValidState(display) {
		return display, ""
	}
	for kind, kindMap := range displayLabels {
		for canonical, label := range kindMap {
			if label == display {
				return canonical, kind
			}
		}
	}
	return display, ""
}

// ValidateIdea checks business rules for an idea.
func ValidateIdea(idea *Idea) error {
	if idea.State != "" && !IsValidState(idea.State) {
		return fmt.Errorf("invalid state: %s", idea.State)
	}

	if idea.Kind != "" && !IsValidKind(idea.Kind) {
		return fmt.Errorf("invalid kind: %s (use aspiration, belief, or plan)", idea.Kind)
	}

	if idea.Maturity != "" && !IsValidMaturity(idea.Maturity) {
		return fmt.Errorf("invalid maturity: %s", idea.Maturity)
	}

	if idea.State == StateRejected && strings.TrimSpace(idea.RejectedReason) == "" {
		return fmt.Errorf("rejected state requires a rejected_reason")
	}

	return nil
}
