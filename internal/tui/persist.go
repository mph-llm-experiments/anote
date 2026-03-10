package tui

import (
	"strings"
	"time"

	"github.com/mph-llm-experiments/anote/internal/config"
	"github.com/mph-llm-experiments/anote/internal/denote"
	"github.com/mph-llm-experiments/anote/internal/idea"
)

// persistIdeaFrontmatter writes updated frontmatter for the idea, updating the
// Modified timestamp. The idea's FilePath must be set.
func persistIdeaFrontmatter(idea *denote.Idea) error {
	idea.Modified = time.Now().Format(time.RFC3339)
	return denote.UpdateIdeaFrontmatter(idea.FilePath, idea)
}

// persistLogEntry appends a timestamped log entry to the idea file.
// It re-reads the file first to extract the current body content, then writes
// the full file back with the log entry appended.
func persistLogEntry(i *denote.Idea, entry string) error {
	existingContent := extractContent(i.Content)
	newContent := appendLogEntry(existingContent, entry)
	i.Modified = time.Now().Format(time.RFC3339)
	return denote.WriteIdeaFile(i.FilePath, i, newContent)
}

// createIdea creates a new idea file and returns the parsed result.
func createIdea(cfg *config.Config, title, kind string, tags []string) (*denote.Idea, error) {
	return idea.CreateIdea(cfg.IdeasDirectory, title, tags, kind, "")
}

// refreshIdea re-reads the idea from disk, returning the fresh version.
func refreshIdea(i *denote.Idea) (*denote.Idea, error) {
	return denote.ParseIdeaFile(i.FilePath)
}

// extractContent extracts the body content after YAML frontmatter.
// Mirrors the CLI's extractIdeaContent function.
func extractContent(fullContent string) string {
	if !strings.HasPrefix(fullContent, "---\n") {
		return fullContent
	}
	lines := strings.Split(fullContent, "\n")
	for idx, line := range lines {
		if idx == 0 {
			continue
		}
		if line == "---" {
			rest := strings.Join(lines[idx+1:], "\n")
			return strings.TrimPrefix(rest, "\n")
		}
	}
	return ""
}

// appendLogEntry appends a timestamped entry to the ## Log section.
// Mirrors the CLI's addLogEntry function.
func appendLogEntry(content, message string) string {
	now := time.Now().Format("2006-01-02")
	entry := "- **" + now + "** " + message

	logIdx := strings.Index(content, "\n## Log\n")
	if logIdx != -1 {
		insertAt := logIdx + len("\n## Log\n")
		return content[:insertAt] + entry + "\n" + content[insertAt:]
	}

	if strings.HasPrefix(content, "## Log\n") {
		insertAt := len("## Log\n")
		return content[:insertAt] + entry + "\n" + content[insertAt:]
	}

	// No log section yet — append one.
	trimmed := strings.TrimRight(content, "\n")
	if trimmed == "" {
		return "## Log\n" + entry + "\n"
	}
	return trimmed + "\n\n## Log\n" + entry + "\n"
}
