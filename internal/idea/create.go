package idea

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/mph-llm-experiments/anote/internal/denote"
)

// CreateIdea creates a new idea file with YAML frontmatter.
func CreateIdea(dir, title string, tags []string) (*denote.Idea, error) {
	counter, err := denote.GetIDCounter(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to get ID counter: %w", err)
	}

	indexID, err := counter.NextID()
	if err != nil {
		return nil, fmt.Errorf("failed to get next index ID: %w", err)
	}

	now := time.Now()
	denoteID := now.Format("20060102T150405")

	slug := denote.TitleToSlug(title)

	// Ensure "idea" tag is first in filename tags
	filenameTags := make([]string, 0, len(tags)+1)
	filenameTags = append(filenameTags, "idea")
	for _, tag := range tags {
		if tag != "idea" {
			filenameTags = append(filenameTags, tag)
		}
	}

	filename := denote.BuildDenoteFilename(denoteID, slug, filenameTags)
	path := filepath.Join(dir, filename)

	meta := &denote.IdeaMetadata{
		Title:   title,
		IndexID: indexID,
		Type:    denote.TypeIdea,
		State:   denote.StateSeed,
		Tags:    tags,
		Created: now.Format(time.RFC3339),
	}

	if err := denote.WriteIdeaFile(path, meta, ""); err != nil {
		return nil, fmt.Errorf("failed to write idea file: %w", err)
	}

	return denote.ParseIdeaFile(path)
}
