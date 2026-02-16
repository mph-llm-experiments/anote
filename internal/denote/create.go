package denote

import (
	"fmt"
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
