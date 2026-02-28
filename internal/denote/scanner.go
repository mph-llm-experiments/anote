package denote

import (
	"path/filepath"

	"github.com/mph-llm-experiments/acore"
)

// Scanner finds and loads idea files from a directory.
type Scanner struct {
	BaseDir string
}

// NewScanner creates a new scanner for the given directory.
func NewScanner(dir string) *Scanner {
	return &Scanner{BaseDir: dir}
}

// FindIdeas finds and parses all idea files in the directory.
func (s *Scanner) FindIdeas() ([]*Idea, error) {
	sc := &acore.Scanner{Store: acore.NewLocalStore(s.BaseDir)}
	names, err := sc.FindByType(TypeIdea)
	if err != nil {
		return nil, err
	}

	var ideas []*Idea
	for _, name := range names {
		idea, err := ParseIdeaFile(filepath.Join(s.BaseDir, name))
		if err != nil {
			continue
		}
		ideas = append(ideas, idea)
	}

	return ideas, nil
}
