package denote

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mph-llm-experiments/acore"
)

func TestWriteIdeaFile(t *testing.T) {
	dir := t.TempDir()

	idea := &Idea{}
	idea.ID = "01TESTID0000000000000000EF"
	idea.Title = "Test Idea"
	idea.IndexID = 1
	idea.Type = TypeIdea
	idea.Tags = []string{"idea", "coaching"}
	idea.Created = "2026-02-16T10:30:45Z"
	idea.Modified = "2026-02-16T10:30:45Z"
	idea.State = StateSeed
	idea.Kind = KindAspiration

	path := filepath.Join(dir, "01TESTID0000000000000000EF--test-idea__idea.md")

	if err := WriteIdeaFile(path, idea, ""); err != nil {
		t.Fatalf("WriteIdeaFile: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}

	content := string(data)
	if !strings.HasPrefix(content, "---\n") {
		t.Error("file should start with ---")
	}
	if !strings.Contains(content, "title: Test Idea") {
		t.Error("file should contain title")
	}
	if !strings.Contains(content, "index_id: 1") {
		t.Error("file should contain index_id")
	}
	if !strings.Contains(content, "state: seed") {
		t.Error("file should contain state")
	}
}

func TestFrontmatterRoundTrip(t *testing.T) {
	dir := t.TempDir()

	original := &Idea{}
	original.ID = "01TESTID0000000000000000GH"
	original.Title = "Round Trip Test"
	original.IndexID = 42
	original.Type = TypeIdea
	original.Tags = []string{"idea", "testing", "roundtrip"}
	original.Created = "2026-02-16T10:30:45Z"
	original.Modified = "2026-02-16T11:15:22Z"
	original.State = StateActive
	original.Maturity = MaturityCrawl
	original.Kind = KindAspiration
	original.RelatedIdeas = []string{"20260301T091500"}
	original.RelatedTasks = []string{"20260215T140000"}

	path := filepath.Join(dir, "01TESTID0000000000000000GH--round-trip-test__idea.md")
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
	if idea.Title != original.Title {
		t.Errorf("Title: got %q, want %q", idea.Title, original.Title)
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
	if len(idea.Tags) != len(original.Tags) {
		t.Errorf("Tags: got %v, want %v", idea.Tags, original.Tags)
	}
	if len(idea.RelatedIdeas) != 1 || idea.RelatedIdeas[0] != "20260301T091500" {
		t.Errorf("RelatedIdeas: got %v", idea.RelatedIdeas)
	}
	if len(idea.RelatedTasks) != 1 || idea.RelatedTasks[0] != "20260215T140000" {
		t.Errorf("RelatedTasks: got %v", idea.RelatedTasks)
	}
	if idea.Created != original.Created {
		t.Errorf("Created: got %q, want %q", idea.Created, original.Created)
	}
	if idea.Modified != original.Modified {
		t.Errorf("Modified: got %q, want %q", idea.Modified, original.Modified)
	}
}

func TestUpdateIdeaFrontmatter_PreservesContent(t *testing.T) {
	dir := t.TempDir()

	original := &Idea{}
	original.ID = "01TESTID0000000000000000IJ"
	original.Title = "Preserve Content"
	original.IndexID = 1
	original.Type = TypeIdea
	original.Tags = []string{"idea"}
	original.Created = "2026-02-16T10:30:45Z"
	original.Modified = "2026-02-16T10:30:45Z"
	original.State = StateSeed

	bodyContent := "## My Idea\n\nThis content should survive frontmatter updates.\n"
	path := filepath.Join(dir, "01TESTID0000000000000000IJ--preserve-content__idea.md")

	if err := WriteIdeaFile(path, original, bodyContent); err != nil {
		t.Fatalf("WriteIdeaFile: %v", err)
	}

	// Update frontmatter
	updated := &Idea{}
	updated.ID = original.ID
	updated.Title = original.Title
	updated.IndexID = original.IndexID
	updated.Type = original.Type
	updated.Tags = original.Tags
	updated.Created = original.Created
	updated.Modified = "2026-02-16T12:00:00Z"
	updated.State = StateDraft

	if err := UpdateIdeaFrontmatter(path, updated); err != nil {
		t.Fatalf("UpdateIdeaFrontmatter: %v", err)
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

func TestWriteIdeaFile_OmitsEmptyFields(t *testing.T) {
	dir := t.TempDir()

	idea := &Idea{}
	idea.ID = "01TESTID0000000000000000KL"
	idea.Title = "Minimal"
	idea.IndexID = 1
	idea.Type = TypeIdea
	idea.State = StateSeed
	idea.Created = acore.Now()
	idea.Modified = idea.Created

	path := filepath.Join(dir, "01TESTID0000000000000000KL--minimal__idea.md")

	if err := WriteIdeaFile(path, idea, ""); err != nil {
		t.Fatalf("WriteIdeaFile: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}

	content := string(data)
	if strings.Contains(content, "maturity") {
		t.Error("empty maturity should be omitted")
	}
	if strings.Contains(content, "rejected_reason") {
		t.Error("empty rejected_reason should be omitted")
	}
	if strings.Contains(content, "related_ideas") {
		t.Error("empty related_ideas should be omitted")
	}
	if strings.Contains(content, "related_tasks") {
		t.Error("empty related_tasks should be omitted")
	}
	if strings.Contains(content, "related_people") {
		t.Error("empty related_people should be omitted")
	}
}
