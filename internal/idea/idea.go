package idea

import (
	"fmt"

	"github.com/mph-llm-experiments/anote/internal/denote"
)

// FindIdeaByID finds an idea by its sequential index ID.
func FindIdeaByID(dir string, id int) (*denote.Idea, error) {
	scanner := denote.NewScanner(dir)
	ideas, err := scanner.FindIdeas()
	if err != nil {
		return nil, err
	}

	for _, idea := range ideas {
		if idea.IndexID == id {
			return idea, nil
		}
	}

	return nil, fmt.Errorf("idea %d not found", id)
}

// FindIdeaByEntityID finds an idea by its entity ID (ULID or legacy Denote ID).
func FindIdeaByEntityID(dir string, entityID string) (*denote.Idea, error) {
	scanner := denote.NewScanner(dir)
	ideas, err := scanner.FindIdeas()
	if err != nil {
		return nil, err
	}

	for _, idea := range ideas {
		if idea.ID == entityID {
			return idea, nil
		}
	}

	return nil, fmt.Errorf("idea with ID %s not found", entityID)
}
