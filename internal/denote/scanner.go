package denote

import (
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
	sc := &acore.Scanner{Dir: s.BaseDir}
	paths, err := sc.FindByType(TypeIdea)
	if err != nil {
		return nil, err
	}

	var ideas []*Idea
	for _, path := range paths {
		idea, err := ParseIdeaFile(path)
		if err != nil {
			continue
		}
		ideas = append(ideas, idea)
	}

	return ideas, nil
}
