package denote

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseFilename(t *testing.T) {
	p := NewParser()

	tests := []struct {
		name     string
		filename string
		wantID   string
		wantSlug string
		wantTags []string
		wantErr  bool
	}{
		{
			name:     "full filename with double dash",
			filename: "20260216T103045--coaching-practice__idea_coaching_leadership.md",
			wantID:   "20260216T103045",
			wantSlug: "coaching-practice",
			wantTags: []string{"idea", "coaching", "leadership"},
		},
		{
			name:     "single dash (backward compat)",
			filename: "20260216T103045-coaching-practice__idea_coaching.md",
			wantID:   "20260216T103045",
			wantSlug: "coaching-practice",
			wantTags: []string{"idea", "coaching"},
		},
		{
			name:     "no tags",
			filename: "20260216T103045--my-idea.md",
			wantID:   "20260216T103045",
			wantSlug: "my-idea",
			wantTags: []string{},
		},
		{
			name:     "single tag",
			filename: "20260216T103045--my-idea__idea.md",
			wantID:   "20260216T103045",
			wantSlug: "my-idea",
			wantTags: []string{"idea"},
		},
		{
			name:     "with path prefix",
			filename: "/home/user/ideas/20260216T103045--my-idea__idea.md",
			wantID:   "20260216T103045",
			wantSlug: "my-idea",
			wantTags: []string{"idea"},
		},
		{
			name:     "not a denote file",
			filename: "random-file.md",
			wantErr:  true,
		},
		{
			name:     "wrong extension",
			filename: "20260216T103045--my-idea__idea.txt",
			wantErr:  true,
		},
		{
			name:     "invalid timestamp",
			filename: "2026021--my-idea__idea.md",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file, err := p.ParseFilename(tt.filename)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ParseFilename(%q): err=%v, wantErr=%v", tt.filename, err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}

			if file.ID != tt.wantID {
				t.Errorf("ID: got %q, want %q", file.ID, tt.wantID)
			}
			if file.Slug != tt.wantSlug {
				t.Errorf("Slug: got %q, want %q", file.Slug, tt.wantSlug)
			}
			if len(file.Tags) != len(tt.wantTags) {
				t.Errorf("Tags: got %v, want %v", file.Tags, tt.wantTags)
			} else {
				for i, tag := range file.Tags {
					if tag != tt.wantTags[i] {
						t.Errorf("Tags[%d]: got %q, want %q", i, tag, tt.wantTags[i])
					}
				}
			}
		})
	}
}

func TestParseIdeaFile(t *testing.T) {
	dir := t.TempDir()

	content := `---
title: Coaching Practice
index_id: 1
type: idea
state: active
maturity: crawl
tags: [coaching, leadership]
related:
  - "20260301T091500"
project:
  - "20260215T140000"
created: "2026-02-16T10:30:45Z"
modified: "2026-02-16T11:15:22Z"
---

## The Idea

This is about coaching non-traditional managers.
`

	filename := "20260216T103045--coaching-practice__idea_coaching_leadership.md"
	path := filepath.Join(dir, filename)
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	idea, err := ParseIdeaFile(path)
	if err != nil {
		t.Fatalf("ParseIdeaFile: %v", err)
	}

	// Check File fields
	if idea.ID != "20260216T103045" {
		t.Errorf("ID: got %q, want %q", idea.ID, "20260216T103045")
	}

	// Check metadata from frontmatter
	if idea.IdeaMetadata.Title != "Coaching Practice" {
		t.Errorf("Title: got %q, want %q", idea.IdeaMetadata.Title, "Coaching Practice")
	}
	if idea.IndexID != 1 {
		t.Errorf("IndexID: got %d, want %d", idea.IndexID, 1)
	}
	if idea.State != "active" {
		t.Errorf("State: got %q, want %q", idea.State, "active")
	}
	if idea.Maturity != "crawl" {
		t.Errorf("Maturity: got %q, want %q", idea.Maturity, "crawl")
	}
	if len(idea.IdeaMetadata.Tags) != 2 {
		t.Errorf("Tags: got %v, want 2 tags", idea.IdeaMetadata.Tags)
	}
	if len(idea.Related) != 1 || idea.Related[0] != "20260301T091500" {
		t.Errorf("Related: got %v, want [20260301T091500]", idea.Related)
	}
	if len(idea.Project) != 1 || idea.Project[0] != "20260215T140000" {
		t.Errorf("Project: got %v, want [20260215T140000]", idea.Project)
	}
	if idea.Created != "2026-02-16T10:30:45Z" {
		t.Errorf("Created: got %q", idea.Created)
	}
}

func TestParseIdeaFile_DefaultState(t *testing.T) {
	dir := t.TempDir()

	content := `---
title: Quick Thought
index_id: 2
type: idea
---

Just a seed.
`
	filename := "20260216T120000--quick-thought__idea.md"
	path := filepath.Join(dir, filename)
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	idea, err := ParseIdeaFile(path)
	if err != nil {
		t.Fatalf("ParseIdeaFile: %v", err)
	}

	if idea.State != StateSeed {
		t.Errorf("default state: got %q, want %q", idea.State, StateSeed)
	}
}

func TestParseIdeaFile_NotAnIdea(t *testing.T) {
	dir := t.TempDir()

	content := "---\ntitle: A Task\n---\n"
	filename := "20260216T120000--a-task__task.md"
	path := filepath.Join(dir, filename)
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	_, err := ParseIdeaFile(path)
	if err == nil {
		t.Error("expected error when parsing non-idea file")
	}
}

func TestTitleFromSlug(t *testing.T) {
	tests := []struct {
		slug string
		want string
	}{
		{"coaching-practice", "coaching practice"},
		{"single", "single"},
		{"a-b-c", "a b c"},
	}

	for _, tt := range tests {
		got := titleFromSlug(tt.slug)
		if got != tt.want {
			t.Errorf("titleFromSlug(%q): got %q, want %q", tt.slug, got, tt.want)
		}
	}
}
