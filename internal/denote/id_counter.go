package denote

import (
	"path/filepath"

	"github.com/mph-llm-experiments/acore"
)

// NewIDCounter creates or loads the index counter for anote.
func NewIDCounter(dir string) (*acore.IndexCounter, error) {
	store := acore.NewLocalStore(dir)
	counter, err := acore.NewIndexCounter(store, "anote")
	if err != nil {
		return nil, err
	}

	// Initialize from existing files if the counter is at 1 (fresh)
	err = counter.InitFromFiles(TypeIdea, func(name string) (int, error) {
		idea, err := ParseIdeaFile(filepath.Join(dir, name))
		if err != nil {
			return 0, err
		}
		return idea.IndexID, nil
	})
	if err != nil {
		return nil, err
	}

	return counter, nil
}
