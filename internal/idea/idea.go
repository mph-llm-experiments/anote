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

// FindIdeaByDenoteID finds an idea by its Denote timestamp ID.
func FindIdeaByDenoteID(dir string, denoteID string) (*denote.Idea, error) {
	scanner := denote.NewScanner(dir)
	ideas, err := scanner.FindIdeas()
	if err != nil {
		return nil, err
	}

	for _, idea := range ideas {
		if idea.ID == denoteID {
			return idea, nil
		}
	}

	return nil, fmt.Errorf("idea with Denote ID %s not found", denoteID)
}
