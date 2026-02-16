package denote

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// WriteFrontmatter serializes IdeaMetadata to a YAML frontmatter string with delimiters.
func WriteFrontmatter(meta *IdeaMetadata) (string, error) {
	data, err := yaml.Marshal(meta)
	if err != nil {
		return "", fmt.Errorf("failed to marshal frontmatter: %w", err)
	}
	return "---\n" + string(data) + "---\n", nil
}

// WriteIdeaFile writes a complete idea file (frontmatter + content).
func WriteIdeaFile(path string, meta *IdeaMetadata, content string) error {
	fm, err := WriteFrontmatter(meta)
	if err != nil {
		return err
	}

	body := fm
	if content != "" {
		body += "\n" + content
	}

	if err := os.WriteFile(path, []byte(body), 0644); err != nil {
		return fmt.Errorf("failed to write idea file: %w", err)
	}

	return nil
}

// UpdateFrontmatter replaces the frontmatter in an existing file, preserving content below.
func UpdateFrontmatter(path string, meta *IdeaMetadata) error {
	existing, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	content := extractContent(string(existing))

	return WriteIdeaFile(path, meta, content)
}

// extractContent returns everything after the closing frontmatter delimiter.
func extractContent(fileContent string) string {
	if !strings.HasPrefix(fileContent, "---\n") {
		return fileContent
	}

	lines := strings.Split(fileContent, "\n")
	for i, line := range lines {
		if i == 0 {
			continue
		}
		if line == "---" {
			// Everything after the closing delimiter
			rest := strings.Join(lines[i+1:], "\n")
			// Trim leading newline that separates frontmatter from content
			rest = strings.TrimPrefix(rest, "\n")
			return rest
		}
	}

	return fileContent
}
