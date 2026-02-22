package denote

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewIDCounter_CreatesOnFirstUse(t *testing.T) {
	dir := t.TempDir()

	counter, err := NewIDCounter(dir)
	if err != nil {
		t.Fatalf("NewIDCounter: %v", err)
	}

	// First ID should be 1 (empty directory)
	id, err := counter.Next()
	if err != nil {
		t.Fatalf("Next: %v", err)
	}
	if id != 1 {
		t.Errorf("first ID: got %d, want 1", id)
	}

	// Counter file should exist after first use
	counterFile := filepath.Join(dir, ".anote-counter.json")
	if _, err := os.Stat(counterFile); os.IsNotExist(err) {
		t.Fatal("counter file should have been created after Next()")
	}
}

func TestNewIDCounter_Increments(t *testing.T) {
	dir := t.TempDir()

	counter, err := NewIDCounter(dir)
	if err != nil {
		t.Fatalf("NewIDCounter: %v", err)
	}

	id1, err := counter.Next()
	if err != nil {
		t.Fatalf("Next: %v", err)
	}
	if id1 != 1 {
		t.Errorf("first ID: got %d, want 1", id1)
	}

	id2, err := counter.Next()
	if err != nil {
		t.Fatalf("Next: %v", err)
	}
	if id2 != 2 {
		t.Errorf("second ID: got %d, want 2", id2)
	}

	id3, err := counter.Next()
	if err != nil {
		t.Fatalf("Next: %v", err)
	}
	if id3 != 3 {
		t.Errorf("third ID: got %d, want 3", id3)
	}
}

func TestNewIDCounter_RecoveryFromExistingFiles(t *testing.T) {
	dir := t.TempDir()

	// Create some idea files with index_ids to simulate recovery
	ideas := []struct {
		filename string
		content  string
	}{
		{
			"01TESTID0000000000000000A1--idea-one__idea.md",
			"---\nid: 01TESTID0000000000000000A1\ntitle: Idea One\nindex_id: 5\ntype: idea\ncreated: \"2026-02-16T10:30:45Z\"\nmodified: \"2026-02-16T10:30:45Z\"\n---\n",
		},
		{
			"01TESTID0000000000000000A2--idea-two__idea.md",
			"---\nid: 01TESTID0000000000000000A2\ntitle: Idea Two\nindex_id: 12\ntype: idea\ncreated: \"2026-02-16T10:30:45Z\"\nmodified: \"2026-02-16T10:30:45Z\"\n---\n",
		},
		{
			"01TESTID0000000000000000A3--idea-three__idea.md",
			"---\nid: 01TESTID0000000000000000A3\ntitle: Idea Three\nindex_id: 8\ntype: idea\ncreated: \"2026-02-16T10:30:45Z\"\nmodified: \"2026-02-16T10:30:45Z\"\n---\n",
		},
	}

	for _, idea := range ideas {
		path := filepath.Join(dir, idea.filename)
		if err := os.WriteFile(path, []byte(idea.content), 0644); err != nil {
			t.Fatalf("failed to write %s: %v", idea.filename, err)
		}
	}

	// No counter file exists — should recover from max index_id (12)
	counter, err := NewIDCounter(dir)
	if err != nil {
		t.Fatalf("NewIDCounter: %v", err)
	}

	id, err := counter.Next()
	if err != nil {
		t.Fatalf("Next: %v", err)
	}

	// Max existing is 12, so next should be 13
	if id != 13 {
		t.Errorf("recovered ID: got %d, want 13", id)
	}
}
