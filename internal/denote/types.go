package denote

import (
	"fmt"
	"strings"
	"time"
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
)

// Type constant.
const TypeIdea = "idea"

// File represents the basic Denote file structure.
type File struct {
	ID      string    `json:"denote_id"`
	Title   string    `json:"-"`
	Slug    string    `json:"slug,omitempty"`
	Tags    []string  `json:"filename_tags,omitempty"`
	Path    string    `json:"path,omitempty"`
	ModTime time.Time `json:"-"`
}

// IsIdea checks if the file is an idea based on filename tags.
func (f *File) IsIdea() bool {
	return f.HasTag("idea")
}

// HasTag checks if the file has a specific tag.
func (f *File) HasTag(tag string) bool {
	for _, t := range f.Tags {
		if t == tag {
			return true
		}
	}
	return false
}

// MatchesSearch checks if the file matches a search query.
func (f *File) MatchesSearch(query string) bool {
	query = strings.ToLower(query)
	if fuzzyMatch(strings.ToLower(f.Title), query) {
		return true
	}
	if fuzzyMatch(strings.ToLower(f.Slug), query) {
		return true
	}
	for _, tag := range f.Tags {
		if fuzzyMatch(strings.ToLower(tag), query) {
			return true
		}
	}
	return false
}

// MatchesTag checks if the file has a tag matching the query.
func (f *File) MatchesTag(query string) bool {
	query = strings.ToLower(query)
	for _, tag := range f.Tags {
		if fuzzyMatch(strings.ToLower(tag), query) {
			return true
		}
	}
	return false
}

func fuzzyMatch(text, pattern string) bool {
	if pattern == "" {
		return true
	}
	patternIdx := 0
	for _, ch := range text {
		if patternIdx < len(pattern) && ch == rune(pattern[patternIdx]) {
			patternIdx++
		}
	}
	return patternIdx == len(pattern)
}

// IdeaMetadata represents idea-specific YAML frontmatter.
type IdeaMetadata struct {
	Title          string   `yaml:"title" json:"title"`
	IndexID        int      `yaml:"index_id" json:"index_id"`
	Type           string   `yaml:"type,omitempty" json:"type,omitempty"`
	Kind           string   `yaml:"kind,omitempty" json:"kind,omitempty"`
	State          string   `yaml:"state,omitempty" json:"state,omitempty"`
	Maturity       string   `yaml:"maturity,omitempty" json:"maturity,omitempty"`
	Tags           []string `yaml:"tags,omitempty" json:"tags,omitempty"`
	Related        []string `yaml:"related,omitempty" json:"related,omitempty"`
	Project        []string `yaml:"project,omitempty" json:"project,omitempty"`
	RejectedReason string   `yaml:"rejected_reason,omitempty" json:"rejected_reason,omitempty"`
	Created        string   `yaml:"created,omitempty" json:"created,omitempty"`
	Modified       string   `yaml:"modified,omitempty" json:"modified,omitempty"`
}

// Idea combines File info with IdeaMetadata and content.
type Idea struct {
	File
	IdeaMetadata
	Content string    `json:"-"`
	ModTime time.Time `json:"modified_at"`
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
	case KindAspiration, KindBelief:
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
	for canonical, label := range displayLabels[KindBelief] {
		if label == display {
			return canonical, KindBelief
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
		return fmt.Errorf("invalid kind: %s (use aspiration or belief)", idea.Kind)
	}

	if idea.Maturity != "" && !IsValidMaturity(idea.Maturity) {
		return fmt.Errorf("invalid maturity: %s", idea.Maturity)
	}

	if idea.State == StateRejected && strings.TrimSpace(idea.RejectedReason) == "" {
		return fmt.Errorf("rejected state requires a rejected_reason")
	}

	return nil
}
