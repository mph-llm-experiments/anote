package denote

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWriteFrontmatter(t *testing.T) {
	meta := &IdeaMetadata{
		Title:   "Test Idea",
		IndexID: 1,
		Type:    TypeIdea,
		State:   StateSeed,
		Tags:    []string{"coaching"},
		Created: "2026-02-16T10:30:45Z",
	}

	fm, err := WriteFrontmatter(meta)
	if err != nil {
		t.Fatalf("WriteFrontmatter: %v", err)
	}

	if !strings.HasPrefix(fm, "---\n") {
		t.Error("frontmatter should start with ---")
	}
	if !strings.HasSuffix(fm, "---\n") {
		t.Error("frontmatter should end with ---")
	}
	if !strings.Contains(fm, "title: Test Idea") {
		t.Error("frontmatter should contain title")
	}
	if !strings.Contains(fm, "index_id: 1") {
		t.Error("frontmatter should contain index_id")
	}
	if !strings.Contains(fm, "state: seed") {
		t.Error("frontmatter should contain state")
	}
}

func TestFrontmatterRoundTrip(t *testing.T) {
	dir := t.TempDir()

	original := &IdeaMetadata{
		Title:    "Round Trip Test",
		IndexID:  42,
		Type:     TypeIdea,
		State:    StateActive,
		Maturity: MaturityCrawl,
		Tags:     []string{"testing", "roundtrip"},
		Related:  []string{"20260301T091500"},
		Project:  []string{"20260215T140000"},
		Created:  "2026-02-16T10:30:45Z",
		Modified: "2026-02-16T11:15:22Z",
	}

	filename := "20260216T103045--round-trip-test__idea_testing.md"
	path := filepath.Join(dir, filename)

	bodyContent := "## The Idea\n\nSome content here.\n"

	if err := WriteIdeaFile(path, original, bodyContent); err != nil {
		t.Fatalf("WriteIdeaFile: %v", err)
	}

	// Parse it back
	idea, err := ParseIdeaFile(path)
	if err != nil {
		t.Fatalf("ParseIdeaFile: %v", err)
	}

	// Verify all fields round-tripped
	if idea.IdeaMetadata.Title != original.Title {
		t.Errorf("Title: got %q, want %q", idea.IdeaMetadata.Title, original.Title)
	}
	if idea.IndexID != original.IndexID {
		t.Errorf("IndexID: got %d, want %d", idea.IndexID, original.IndexID)
	}
	if idea.Type != original.Type {
		t.Errorf("Type: got %q, want %q", idea.Type, original.Type)
	}
	if idea.State != original.State {
		t.Errorf("State: got %q, want %q", idea.State, original.State)
	}
	if idea.Maturity != original.Maturity {
		t.Errorf("Maturity: got %q, want %q", idea.Maturity, original.Maturity)
	}
	if len(idea.IdeaMetadata.Tags) != len(original.Tags) {
		t.Errorf("Tags: got %v, want %v", idea.IdeaMetadata.Tags, original.Tags)
	}
	if len(idea.Related) != 1 || idea.Related[0] != "20260301T091500" {
		t.Errorf("Related: got %v", idea.Related)
	}
	if len(idea.Project) != 1 || idea.Project[0] != "20260215T140000" {
		t.Errorf("Project: got %v", idea.Project)
	}
	if idea.Created != original.Created {
		t.Errorf("Created: got %q, want %q", idea.Created, original.Created)
	}
	if idea.Modified != original.Modified {
		t.Errorf("Modified: got %q, want %q", idea.Modified, original.Modified)
	}
}

func TestUpdateFrontmatter_PreservesContent(t *testing.T) {
	dir := t.TempDir()

	original := &IdeaMetadata{
		Title:   "Preserve Content",
		IndexID: 1,
		Type:    TypeIdea,
		State:   StateSeed,
		Created: "2026-02-16T10:30:45Z",
	}

	bodyContent := "## My Idea\n\nThis content should survive frontmatter updates.\n"
	filename := "20260216T103045--preserve-content__idea.md"
	path := filepath.Join(dir, filename)

	if err := WriteIdeaFile(path, original, bodyContent); err != nil {
		t.Fatalf("WriteIdeaFile: %v", err)
	}

	// Update frontmatter
	updated := &IdeaMetadata{
		Title:    "Preserve Content",
		IndexID:  1,
		Type:     TypeIdea,
		State:    StateDraft,
		Modified: "2026-02-16T12:00:00Z",
		Created:  "2026-02-16T10:30:45Z",
	}

	if err := UpdateFrontmatter(path, updated); err != nil {
		t.Fatalf("UpdateFrontmatter: %v", err)
	}

	// Read back
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}

	content := string(data)
	if !strings.Contains(content, "state: draft") {
		t.Error("updated frontmatter should contain new state")
	}
	if !strings.Contains(content, "This content should survive") {
		t.Error("body content should be preserved after frontmatter update")
	}
}

func TestWriteFrontmatter_OmitsEmptyFields(t *testing.T) {
	meta := &IdeaMetadata{
		Title:   "Minimal",
		IndexID: 1,
		State:   StateSeed,
	}

	fm, err := WriteFrontmatter(meta)
	if err != nil {
		t.Fatalf("WriteFrontmatter: %v", err)
	}

	if strings.Contains(fm, "maturity") {
		t.Error("empty maturity should be omitted")
	}
	if strings.Contains(fm, "rejected_reason") {
		t.Error("empty rejected_reason should be omitted")
	}
	if strings.Contains(fm, "related") {
		t.Error("empty related should be omitted")
	}
	if strings.Contains(fm, "project") {
		t.Error("empty project should be omitted")
	}
}
