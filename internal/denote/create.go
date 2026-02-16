package denote

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// titleToSlug converts a title to a Denote-compatible kebab-case slug.
func titleToSlug(title string) string {
	slug := strings.ToLower(title)

	slug = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			return r
		}
		return '-'
	}, slug)

	for strings.Contains(slug, "--") {
		slug = strings.ReplaceAll(slug, "--", "-")
	}

	slug = strings.Trim(slug, "-")

	return slug
}

// TitleToSlug is the exported version for use by other packages.
func TitleToSlug(title string) string {
	return titleToSlug(title)
}

// BuildDenoteFilename constructs a Denote filename from components.
func BuildDenoteFilename(id, slug string, tags []string) string {
	tagString := ""
	if len(tags) > 0 {
		tagString = "__" + strings.Join(tags, "_")
	}
	return fmt.Sprintf("%s--%s%s.md", id, slug, tagString)
}

// RenameFileForTags renames a Denote file to reflect updated tags.
// Returns the new file path.
func RenameFileForTags(oldPath string, newTags []string) (string, error) {
	p := NewParser()
	file, err := p.ParseFilename(oldPath)
	if err != nil {
		return "", err
	}

	newFilename := BuildDenoteFilename(file.ID, file.Slug, newTags)
	dir := filepath.Dir(oldPath)
	newPath := filepath.Join(dir, newFilename)

	if oldPath == newPath {
		return oldPath, nil
	}

	if err := os.Rename(oldPath, newPath); err != nil {
		return "", fmt.Errorf("failed to rename file: %w", err)
	}

	return newPath, nil
}
