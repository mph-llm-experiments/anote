package denote

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/mph-llm-experiments/acore"
)

var (
	// Legacy Denote filename pattern for backward compatibility during migration
	legacyDenotePattern = regexp.MustCompile(`^(\d{8}T\d{6})-{1,2}([^_]+)(?:__(.+))?\.md$`)
)

// ParseIdeaFile reads and parses an idea file using acore.
func ParseIdeaFile(path string) (*Idea, error) {
	var idea Idea
	content, err := acore.ReadFile(path, &idea)
	if err != nil {
		return nil, fmt.Errorf("failed to parse idea file: %w", err)
	}
	idea.Content = content
	idea.FilePath = path

	// Get file modification time
	if info, err := os.Stat(path); err == nil {
		idea.ModTime = info.ModTime()
	}

	// If ID not in frontmatter, extract from filename (legacy Denote files)
	if idea.ID == "" {
		base := filepath.Base(path)
		if m := legacyDenotePattern.FindStringSubmatch(base); len(m) > 1 {
			idea.ID = m[1]
		}
	}

	// Set defaults per spec
	if idea.State == "" {
		idea.State = StateSeed
	}
	if idea.Type == "" {
		idea.Type = TypeIdea
	}

	// Ensure relation slices for JSON output
	idea.EnsureSlices()

	return &idea, nil
}
