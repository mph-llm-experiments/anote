package denote

import (
	"fmt"
	"os"
	"path/filepath"
)

// Scanner finds and loads Denote files from a directory.
type Scanner struct {
	BaseDir string
}

// NewScanner creates a new scanner for the given directory.
func NewScanner(dir string) *Scanner {
	return &Scanner{BaseDir: dir}
}

// FindIdeas finds and parses all idea files in the directory.
func (s *Scanner) FindIdeas() ([]*Idea, error) {
	pattern := filepath.Join(s.BaseDir, "*__idea*.md")
	files, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to glob idea files: %w", err)
	}

	var ideas []*Idea
	for _, file := range files {
		idea, err := ParseIdeaFile(file)
		if err != nil {
			continue
		}
		ideas = append(ideas, idea)
	}

	return ideas, nil
}

// FindAllIdeaFiles returns basic file info for all idea files.
func (s *Scanner) FindAllIdeaFiles() ([]File, error) {
	pattern := filepath.Join(s.BaseDir, "*__idea*.md")
	paths, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to glob idea files: %w", err)
	}

	parser := NewParser()
	var allFiles []File
	for _, path := range paths {
		file, err := parser.ParseFilename(filepath.Base(path))
		if err != nil {
			continue
		}
		file.Path = path
		if info, err := os.Stat(path); err == nil {
			file.ModTime = info.ModTime()
		}
		allFiles = append(allFiles, *file)
	}

	return allFiles, nil
}
