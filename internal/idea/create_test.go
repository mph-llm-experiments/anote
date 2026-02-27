package idea

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mph-llm-experiments/anote/internal/denote"
)

func TestCreateIdea_Basic(t *testing.T) {
	dir := t.TempDir()

	idea, err := CreateIdea(dir,"My test idea", nil, "", "")
	if err != nil {
		t.Fatalf("CreateIdea: %v", err)
	}

	if idea.Title != "My test idea" {
		t.Errorf("Title: got %q, want %q", idea.Title, "My test idea")
	}

	if idea.IndexID != 1 {
		t.Errorf("IndexID: got %d, want 1", idea.IndexID)
	}

	if idea.State != denote.StateSeed {
		t.Errorf("State: got %q, want %q", idea.State, denote.StateSeed)
	}

	if idea.Type != denote.TypeIdea {
		t.Errorf("Type: got %q, want %q", idea.Type, denote.TypeIdea)
	}

	if idea.Created == "" {
		t.Error("Created should not be empty")
	}

	if idea.ID == "" {
		t.Error("ID (ULID) should not be empty")
	}

	// File should exist
	if _, err := os.Stat(idea.FilePath); os.IsNotExist(err) {
		t.Errorf("file should exist: %s", idea.FilePath)
	}
}

func TestCreateIdea_IdeaTagAlwaysInFilename(t *testing.T) {
	dir := t.TempDir()

	idea, err := CreateIdea(dir,"Tag test", nil, "", "")
	if err != nil {
		t.Fatalf("CreateIdea: %v", err)
	}

	filename := filepath.Base(idea.FilePath)
	if !strings.Contains(filename, "__idea") {
		t.Errorf("filename should contain __idea: %s", filename)
	}
}

func TestCreateIdea_WithTags(t *testing.T) {
	dir := t.TempDir()

	idea, err := CreateIdea(dir,"Tagged idea", []string{"coaching", "leadership"}, "", "")
	if err != nil {
		t.Fatalf("CreateIdea: %v", err)
	}

	// Tags should be in frontmatter (with "idea" first)
	if len(idea.Tags) != 3 {
		t.Errorf("frontmatter tags: got %v, want [idea, coaching, leadership]", idea.Tags)
	}
	if idea.Tags[0] != "idea" {
		t.Errorf("first tag should be 'idea', got %q", idea.Tags[0])
	}
}

func TestCreateIdea_IdeaTagNotDuplicated(t *testing.T) {
	dir := t.TempDir()

	// Pass "idea" as a user tag — should not appear twice
	idea, err := CreateIdea(dir,"Dedup test", []string{"idea", "coaching"}, "", "")
	if err != nil {
		t.Fatalf("CreateIdea: %v", err)
	}

	ideaCount := 0
	for _, tag := range idea.Tags {
		if tag == "idea" {
			ideaCount++
		}
	}
	if ideaCount != 1 {
		t.Errorf("idea tag should appear exactly once, got %d times in %v", ideaCount, idea.Tags)
	}
}

func TestCreateIdea_SequentialIDs(t *testing.T) {
	dir := t.TempDir()

	idea1, err := CreateIdea(dir,"First", nil, "", "")
	if err != nil {
		t.Fatalf("CreateIdea 1: %v", err)
	}

	idea2, err := CreateIdea(dir,"Second", nil, "", "")
	if err != nil {
		t.Fatalf("CreateIdea 2: %v", err)
	}

	idea3, err := CreateIdea(dir,"Third", nil, "", "")
	if err != nil {
		t.Fatalf("CreateIdea 3: %v", err)
	}

	if idea1.IndexID != 1 {
		t.Errorf("idea1 IndexID: got %d, want 1", idea1.IndexID)
	}
	if idea2.IndexID != 2 {
		t.Errorf("idea2 IndexID: got %d, want 2", idea2.IndexID)
	}
	if idea3.IndexID != 3 {
		t.Errorf("idea3 IndexID: got %d, want 3", idea3.IndexID)
	}
}

func TestCreateIdea_WithKind(t *testing.T) {
	dir := t.TempDir()

	idea, err := CreateIdea(dir,"Chaos is unsustainable", nil, "belief", "")
	if err != nil {
		t.Fatalf("CreateIdea: %v", err)
	}

	if idea.Kind != denote.KindBelief {
		t.Errorf("Kind: got %q, want %q", idea.Kind, denote.KindBelief)
	}

	if idea.State != denote.StateSeed {
		t.Errorf("State: got %q, want %q", idea.State, denote.StateSeed)
	}
}

func TestCreateIdea_DefaultKind(t *testing.T) {
	dir := t.TempDir()

	idea, err := CreateIdea(dir,"Build a widget", nil, "", "")
	if err != nil {
		t.Fatalf("CreateIdea: %v", err)
	}

	if idea.Kind != denote.KindAspiration {
		t.Errorf("Kind: got %q, want %q", idea.Kind, denote.KindAspiration)
	}
}

func TestCreateIdea_ULIDFilenameFormat(t *testing.T) {
	dir := t.TempDir()

	idea, err := CreateIdea(dir,"My Great Idea", nil, "", "")
	if err != nil {
		t.Fatalf("CreateIdea: %v", err)
	}

	filename := filepath.Base(idea.FilePath)

	// Should end with __idea.md
	if !strings.HasSuffix(filename, "__idea.md") {
		t.Errorf("filename should end with __idea.md: %s", filename)
	}

	// Should contain the slug
	if !strings.Contains(filename, "my-great-idea") {
		t.Errorf("filename should contain slug: %s", filename)
	}

	// ID should be 26 chars (ULID)
	if len(idea.ID) != 26 {
		t.Errorf("ID should be 26 chars (ULID): got %d chars (%s)", len(idea.ID), idea.ID)
	}
}
