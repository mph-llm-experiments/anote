package denote

import (
	"testing"
)

func TestTitleToSlug(t *testing.T) {
	tests := []struct {
		title string
		want  string
	}{
		{"Coaching Practice", "coaching-practice"},
		{"Hello World", "hello-world"},
		{"My  Extra   Spaces", "my-extra-spaces"},
		{"Punctuation! & Symbols?", "punctuation-symbols"},
		{"already-slug", "already-slug"},
		{"  Leading/Trailing  ", "leading-trailing"},
		{"Numbers 123 Here", "numbers-123-here"},
		{"ALLCAPS", "allcaps"},
		{"", ""},
	}

	for _, tt := range tests {
		got := titleToSlug(tt.title)
		if got != tt.want {
			t.Errorf("titleToSlug(%q): got %q, want %q", tt.title, got, tt.want)
		}
	}
}

func TestBuildDenoteFilename(t *testing.T) {
	tests := []struct {
		id   string
		slug string
		tags []string
		want string
	}{
		{
			"20260216T103045", "coaching-practice", []string{"idea", "coaching"},
			"20260216T103045--coaching-practice__idea_coaching.md",
		},
		{
			"20260216T103045", "simple-idea", []string{"idea"},
			"20260216T103045--simple-idea__idea.md",
		},
		{
			"20260216T103045", "no-tags", nil,
			"20260216T103045--no-tags.md",
		},
		{
			"20260216T103045", "no-tags", []string{},
			"20260216T103045--no-tags.md",
		},
	}

	for _, tt := range tests {
		got := BuildDenoteFilename(tt.id, tt.slug, tt.tags)
		if got != tt.want {
			t.Errorf("BuildDenoteFilename(%q, %q, %v): got %q, want %q",
				tt.id, tt.slug, tt.tags, got, tt.want)
		}
	}
}
