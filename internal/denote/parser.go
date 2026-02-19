package denote

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

var denotePattern = regexp.MustCompile(`^(\d{8}T\d{6})-{1,2}([^_]+)(?:__(.+))?\.md$`)

// Parser handles parsing of Denote files.
type Parser struct{}

// NewParser creates a new parser.
func NewParser() *Parser {
	return &Parser{}
}

// ParseFilename extracts Denote components from a filename.
func (p *Parser) ParseFilename(filename string) (*File, error) {
	base := filepath.Base(filename)
	matches := denotePattern.FindStringSubmatch(base)
	if len(matches) < 3 {
		return nil, fmt.Errorf("not a valid denote filename: %s", base)
	}

	file := &File{
		ID:    matches[1],
		Slug:  matches[2],
		Title: titleFromSlug(matches[2]),
		Tags:  []string{},
		Path:  filename,
	}

	if len(matches) > 3 && matches[3] != "" {
		file.Tags = strings.Split(matches[3], "_")
	}

	return file, nil
}

// ParseIdeaFile reads and parses an idea file.
func ParseIdeaFile(path string) (*Idea, error) {
	p := NewParser()
	file, err := p.ParseFilename(path)
	if err != nil {
		return nil, err
	}

	if !contains(file.Tags, "idea") {
		return nil, fmt.Errorf("not an idea file: %s", path)
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("failed to stat file: %w", err)
	}

	idea := &Idea{
		File:    *file,
		ModTime: info.ModTime(),
		Content: string(content),
	}

	if fm, err := parseFrontmatter(content); err == nil {
		idea.IdeaMetadata = *fm
	}

	if idea.State == "" {
		idea.State = StateSeed
	}

	if idea.IdeaMetadata.Title != "" {
		idea.File.Title = idea.IdeaMetadata.Title
	}

	return idea, nil
}

func parseFrontmatter(content []byte) (*IdeaMetadata, error) {
	contentStr := string(content)
	if !strings.HasPrefix(contentStr, "---\n") {
		return nil, fmt.Errorf("no frontmatter found")
	}

	lines := strings.Split(contentStr, "\n")
	endLine := -1
	for i, line := range lines {
		if i == 0 {
			continue
		}
		if line == "---" {
			endLine = i
			break
		}
	}

	if endLine == -1 {
		return nil, fmt.Errorf("frontmatter not properly closed")
	}

	fmStr := strings.Join(lines[1:endLine], "\n")

	var meta IdeaMetadata
	if err := yaml.Unmarshal([]byte(fmStr), &meta); err != nil {
		return nil, fmt.Errorf("failed to parse frontmatter: %w", err)
	}

	// Backward compat: read old "related" and "project" field names
	var legacy ideaMetadataLegacy
	if err := yaml.Unmarshal([]byte(fmStr), &legacy); err == nil {
		if len(legacy.Related) > 0 && len(meta.RelatedIdeas) == 0 {
			meta.RelatedIdeas = legacy.Related
		}
		if len(legacy.Project) > 0 && len(meta.RelatedTasks) == 0 {
			meta.RelatedTasks = legacy.Project
		}
	}

	meta.EnsureRelationSlices()

	return &meta, nil
}

func titleFromSlug(slug string) string {
	return strings.ReplaceAll(slug, "-", " ")
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
