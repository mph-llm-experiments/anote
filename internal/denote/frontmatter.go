package denote

import (
	"github.com/mph-llm-experiments/acore"
)

// WriteIdeaFile writes a complete idea file (frontmatter + content).
func WriteIdeaFile(path string, idea *Idea, content string) error {
	return acore.WriteFile(path, idea, content)
}

// UpdateIdeaFrontmatter replaces the frontmatter in an existing file, preserving content below.
func UpdateIdeaFrontmatter(path string, idea *Idea) error {
	return acore.UpdateFrontmatter(path, idea)
}
