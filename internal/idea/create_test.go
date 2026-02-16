package idea

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mph-llm-experiments/anote/internal/denote"
)

func TestCreateIdea_Basic(t *testing.T) {
	denote.ResetSingleton()
	defer denote.ResetSingleton()

	dir := t.TempDir()

	idea, err := CreateIdea(dir, "My test idea", nil)
	if err != nil {
		t.Fatalf("CreateIdea: %v", err)
	}

	if idea.IdeaMetadata.Title != "My test idea" {
		t.Errorf("Title: got %q, want %q", idea.IdeaMetadata.Title, "My test idea")
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

	// File should exist
	if _, err := os.Stat(idea.File.Path); os.IsNotExist(err) {
		t.Errorf("file should exist: %s", idea.File.Path)
	}
}

func TestCreateIdea_IdeaTagAlwaysInFilename(t *testing.T) {
	denote.ResetSingleton()
	defer denote.ResetSingleton()

	dir := t.TempDir()

	idea, err := CreateIdea(dir, "Tag test", nil)
	if err != nil {
		t.Fatalf("CreateIdea: %v", err)
	}

	filename := filepath.Base(idea.File.Path)
	if !strings.Contains(filename, "__idea") {
		t.Errorf("filename should contain __idea: %s", filename)
	}
}

func TestCreateIdea_WithTags(t *testing.T) {
	denote.ResetSingleton()
	defer denote.ResetSingleton()

	dir := t.TempDir()

	idea, err := CreateIdea(dir, "Tagged idea", []string{"coaching", "leadership"})
	if err != nil {
		t.Fatalf("CreateIdea: %v", err)
	}

	// Tags should be in filename
	filename := filepath.Base(idea.File.Path)
	if !strings.Contains(filename, "__idea_coaching_leadership") {
		t.Errorf("filename should contain tags: %s", filename)
	}

	// Tags should be in frontmatter
	if len(idea.IdeaMetadata.Tags) != 2 {
		t.Errorf("frontmatter tags: got %v, want [coaching, leadership]", idea.IdeaMetadata.Tags)
	}
}

func TestCreateIdea_IdeaTagNotDuplicated(t *testing.T) {
	denote.ResetSingleton()
	defer denote.ResetSingleton()

	dir := t.TempDir()

	// Pass "idea" as a user tag — should not appear twice in filename
	idea, err := CreateIdea(dir, "Dedup test", []string{"idea", "coaching"})
	if err != nil {
		t.Fatalf("CreateIdea: %v", err)
	}

	filename := filepath.Base(idea.File.Path)
	// Should be __idea_coaching, not __idea_idea_coaching
	if strings.Contains(filename, "idea_idea") {
		t.Errorf("idea tag should not be duplicated in filename: %s", filename)
	}
}

func TestCreateIdea_SequentialIDs(t *testing.T) {
	denote.ResetSingleton()
	defer denote.ResetSingleton()

	dir := t.TempDir()

	idea1, err := CreateIdea(dir, "First", nil)
	if err != nil {
		t.Fatalf("CreateIdea 1: %v", err)
	}

	idea2, err := CreateIdea(dir, "Second", nil)
	if err != nil {
		t.Fatalf("CreateIdea 2: %v", err)
	}

	idea3, err := CreateIdea(dir, "Third", nil)
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

func TestCreateIdea_DenoteFilenameFormat(t *testing.T) {
	denote.ResetSingleton()
	defer denote.ResetSingleton()

	dir := t.TempDir()

	idea, err := CreateIdea(dir, "My Great Idea", nil)
	if err != nil {
		t.Fatalf("CreateIdea: %v", err)
	}

	filename := filepath.Base(idea.File.Path)

	// Should match Denote pattern
	p := denote.NewParser()
	file, err := p.ParseFilename(filename)
	if err != nil {
		t.Fatalf("filename should be valid Denote format: %v", err)
	}

	if file.Slug != "my-great-idea" {
		t.Errorf("slug: got %q, want %q", file.Slug, "my-great-idea")
	}

	if !file.IsIdea() {
		t.Error("file should be recognized as an idea")
	}
}
