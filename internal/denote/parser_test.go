package denote

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseIdeaFile(t *testing.T) {
	dir := t.TempDir()

	// Use acore filename format: {ulid}--{slug}__{type}.md
	content := `---
id: 01TESTID0000000000000000AB
title: Coaching Practice
index_id: 1
type: idea
kind: aspiration
state: active
maturity: crawl
tags: [idea, coaching, leadership]
related_ideas:
  - "20260301T091500"
related_tasks:
  - "20260215T140000"
created: "2026-02-16T10:30:45Z"
modified: "2026-02-16T11:15:22Z"
---

## The Idea

This is about coaching non-traditional managers.
`

	filename := "01TESTID0000000000000000AB--coaching-practice__idea.md"
	path := filepath.Join(dir, filename)
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	idea, err := ParseIdeaFile(path)
	if err != nil {
		t.Fatalf("ParseIdeaFile: %v", err)
	}

	// Check Entity fields
	if idea.ID != "01TESTID0000000000000000AB" {
		t.Errorf("ID: got %q, want %q", idea.ID, "01TESTID0000000000000000AB")
	}
	if idea.Title != "Coaching Practice" {
		t.Errorf("Title: got %q, want %q", idea.Title, "Coaching Practice")
	}
	if idea.IndexID != 1 {
		t.Errorf("IndexID: got %d, want %d", idea.IndexID, 1)
	}

	// Check IdeaMetadata fields
	if idea.State != "active" {
		t.Errorf("State: got %q, want %q", idea.State, "active")
	}
	if idea.Maturity != "crawl" {
		t.Errorf("Maturity: got %q, want %q", idea.Maturity, "crawl")
	}
	if idea.Kind != "aspiration" {
		t.Errorf("Kind: got %q, want %q", idea.Kind, "aspiration")
	}

	// Check tags
	if len(idea.Tags) != 3 {
		t.Errorf("Tags: got %v, want 3 tags", idea.Tags)
	}

	// Check relations
	if len(idea.RelatedIdeas) != 1 || idea.RelatedIdeas[0] != "20260301T091500" {
		t.Errorf("RelatedIdeas: got %v, want [20260301T091500]", idea.RelatedIdeas)
	}
	if len(idea.RelatedTasks) != 1 || idea.RelatedTasks[0] != "20260215T140000" {
		t.Errorf("RelatedTasks: got %v, want [20260215T140000]", idea.RelatedTasks)
	}
	if idea.Created != "2026-02-16T10:30:45Z" {
		t.Errorf("Created: got %q", idea.Created)
	}
}

func TestParseIdeaFile_LegacyDenoteFilename(t *testing.T) {
	dir := t.TempDir()

	// Legacy Denote format file — ID comes from filename
	content := `---
title: Legacy Idea
index_id: 5
type: idea
state: draft
---

Old format.
`
	filename := "20260216T103045--legacy-idea__idea.md"
	path := filepath.Join(dir, filename)
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	idea, err := ParseIdeaFile(path)
	if err != nil {
		t.Fatalf("ParseIdeaFile: %v", err)
	}

	// ID should be extracted from filename
	if idea.ID != "20260216T103045" {
		t.Errorf("ID: got %q, want %q", idea.ID, "20260216T103045")
	}
	if idea.Title != "Legacy Idea" {
		t.Errorf("Title: got %q, want %q", idea.Title, "Legacy Idea")
	}
}

func TestParseIdeaFile_DefaultState(t *testing.T) {
	dir := t.TempDir()

	content := `---
id: 01TESTID0000000000000000CD
title: Quick Thought
index_id: 2
type: idea
---

Just a seed.
`
	filename := "01TESTID0000000000000000CD--quick-thought__idea.md"
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
