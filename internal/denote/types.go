package denote

import (
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
