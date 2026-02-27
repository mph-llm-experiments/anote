package idea

import (
	"fmt"
	"path/filepath"

	"github.com/mph-llm-experiments/acore"
	"github.com/mph-llm-experiments/anote/internal/denote"
)

// CreateIdea creates a new idea file with YAML frontmatter.
func CreateIdea(dir, title string, tags []string, kind string, body string) (*denote.Idea, error) {
	if kind == "" {
		kind = denote.KindAspiration
	}

	counter, err := denote.NewIDCounter(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to get ID counter: %w", err)
	}

	indexID, err := counter.Next()
	if err != nil {
		return nil, fmt.Errorf("failed to get next index ID: %w", err)
	}

	id := acore.NewID()
	now := acore.Now()

	// Build filename: {ulid}--{slug}__idea.md
	filename := denote.BuildIdeaFilename(id, title)
	path := filepath.Join(dir, filename)

	// Ensure "idea" is in the tags array for frontmatter
	allTags := make([]string, 0, len(tags)+1)
	allTags = append(allTags, "idea")
	for _, tag := range tags {
		if tag != "idea" {
			allTags = append(allTags, tag)
		}
	}

	idea := &denote.Idea{}
	idea.ID = id
	idea.Title = title
	idea.IndexID = indexID
	idea.Type = denote.TypeIdea
	idea.Tags = allTags
	idea.Created = now
	idea.Modified = now
	idea.Kind = kind
	if denote.IsSimpleKind(kind) {
		idea.State = denote.StateActive
	} else {
		idea.State = denote.StateSeed
	}
	idea.FilePath = path

	content := ""
	if body != "" {
		content = body + "\n"
	}
	if err := denote.WriteIdeaFile(path, idea, content); err != nil {
		return nil, fmt.Errorf("failed to write idea file: %w", err)
	}

	// Parse back to get consistent state (ModTime, etc.)
	return denote.ParseIdeaFile(path)
}
