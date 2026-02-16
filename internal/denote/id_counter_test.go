package denote

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestIDCounter_CreatesOnFirstUse(t *testing.T) {
	ResetSingleton()
	defer ResetSingleton()

	dir := t.TempDir()

	counter, err := GetIDCounter(dir)
	if err != nil {
		t.Fatalf("GetIDCounter: %v", err)
	}

	// Counter file should exist
	counterFile := filepath.Join(dir, ".anote-counter.json")
	if _, err := os.Stat(counterFile); os.IsNotExist(err) {
		t.Fatal("counter file should have been created")
	}

	// First ID should be 1 (empty directory, maxID=0, next=1)
	if counter.NextIndexID != 1 {
		t.Errorf("NextIndexID: got %d, want 1", counter.NextIndexID)
	}

	if counter.SpecVersion != "0.1.0" {
		t.Errorf("SpecVersion: got %q, want %q", counter.SpecVersion, "0.1.0")
	}
}

func TestIDCounter_Increments(t *testing.T) {
	ResetSingleton()
	defer ResetSingleton()

	dir := t.TempDir()

	counter, err := GetIDCounter(dir)
	if err != nil {
		t.Fatalf("GetIDCounter: %v", err)
	}

	id1, err := counter.NextID()
	if err != nil {
		t.Fatalf("NextID: %v", err)
	}
	if id1 != 1 {
		t.Errorf("first ID: got %d, want 1", id1)
	}

	id2, err := counter.NextID()
	if err != nil {
		t.Fatalf("NextID: %v", err)
	}
	if id2 != 2 {
		t.Errorf("second ID: got %d, want 2", id2)
	}

	id3, err := counter.NextID()
	if err != nil {
		t.Fatalf("NextID: %v", err)
	}
	if id3 != 3 {
		t.Errorf("third ID: got %d, want 3", id3)
	}
}

func TestIDCounter_PersistsToDisk(t *testing.T) {
	ResetSingleton()
	defer ResetSingleton()

	dir := t.TempDir()

	counter, err := GetIDCounter(dir)
	if err != nil {
		t.Fatalf("GetIDCounter: %v", err)
	}

	// Get a few IDs
	counter.NextID()
	counter.NextID()
	counter.NextID()

	// Read the file directly
	counterFile := filepath.Join(dir, ".anote-counter.json")
	data, err := os.ReadFile(counterFile)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}

	var saved CounterData
	if err := json.Unmarshal(data, &saved); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}

	if saved.NextIndexID != 4 {
		t.Errorf("persisted NextIndexID: got %d, want 4", saved.NextIndexID)
	}
}

func TestIDCounter_LoadsExisting(t *testing.T) {
	ResetSingleton()
	defer ResetSingleton()

	dir := t.TempDir()

	// Pre-create counter file at ID 50
	counterFile := filepath.Join(dir, ".anote-counter.json")
	data, _ := json.MarshalIndent(CounterData{
		NextIndexID: 50,
		SpecVersion: "0.1.0",
	}, "", "  ")
	os.WriteFile(counterFile, data, 0644)

	counter, err := GetIDCounter(dir)
	if err != nil {
		t.Fatalf("GetIDCounter: %v", err)
	}

	id, err := counter.NextID()
	if err != nil {
		t.Fatalf("NextID: %v", err)
	}

	if id != 50 {
		t.Errorf("ID from pre-existing counter: got %d, want 50", id)
	}
}

func TestIDCounter_RecoveryFromExistingFiles(t *testing.T) {
	ResetSingleton()
	defer ResetSingleton()

	dir := t.TempDir()

	// Create some idea files with index_ids to simulate recovery
	ideas := []struct {
		filename string
		content  string
	}{
		{
			"20260216T103045--idea-one__idea.md",
			"---\ntitle: Idea One\nindex_id: 5\ntype: idea\n---\n",
		},
		{
			"20260216T110000--idea-two__idea.md",
			"---\ntitle: Idea Two\nindex_id: 12\ntype: idea\n---\n",
		},
		{
			"20260216T120000--idea-three__idea.md",
			"---\ntitle: Idea Three\nindex_id: 8\ntype: idea\n---\n",
		},
	}

	for _, idea := range ideas {
		path := filepath.Join(dir, idea.filename)
		if err := os.WriteFile(path, []byte(idea.content), 0644); err != nil {
			t.Fatalf("failed to write %s: %v", idea.filename, err)
		}
	}

	// No counter file exists — should recover from max index_id (12)
	counter, err := GetIDCounter(dir)
	if err != nil {
		t.Fatalf("GetIDCounter: %v", err)
	}

	id, err := counter.NextID()
	if err != nil {
		t.Fatalf("NextID: %v", err)
	}

	// Max existing is 12, so next should be 13
	if id != 13 {
		t.Errorf("recovered ID: got %d, want 13", id)
	}
}

func TestIDCounter_Singleton(t *testing.T) {
	ResetSingleton()
	defer ResetSingleton()

	dir := t.TempDir()

	c1, err := GetIDCounter(dir)
	if err != nil {
		t.Fatalf("GetIDCounter: %v", err)
	}

	c2, err := GetIDCounter(dir)
	if err != nil {
		t.Fatalf("GetIDCounter: %v", err)
	}

	if c1 != c2 {
		t.Error("GetIDCounter should return the same instance")
	}
}

func TestIDCounter_CounterFileFormat(t *testing.T) {
	ResetSingleton()
	defer ResetSingleton()

	dir := t.TempDir()

	_, err := GetIDCounter(dir)
	if err != nil {
		t.Fatalf("GetIDCounter: %v", err)
	}

	counterFile := filepath.Join(dir, ".anote-counter.json")
	data, err := os.ReadFile(counterFile)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}

	var saved CounterData
	if err := json.Unmarshal(data, &saved); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}

	if saved.SpecVersion != "0.1.0" {
		t.Errorf("SpecVersion: got %q, want %q", saved.SpecVersion, "0.1.0")
	}
}
