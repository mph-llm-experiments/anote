package denote

import (
	"path/filepath"

	"github.com/mph-llm-experiments/acore"
)

// WriteIdeaFile writes a complete idea file (frontmatter + content).
func WriteIdeaFile(path string, idea *Idea, content string) error {
	store := acore.NewLocalStore(filepath.Dir(path))
	return acore.WriteFile(store, filepath.Base(path), idea, content)
}

// UpdateIdeaFrontmatter replaces the frontmatter in an existing file, preserving content below.
func UpdateIdeaFrontmatter(path string, idea *Idea) error {
	store := acore.NewLocalStore(filepath.Dir(path))
	return acore.UpdateFrontmatter(store, filepath.Base(path), idea)
}
