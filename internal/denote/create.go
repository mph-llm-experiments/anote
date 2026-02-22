package denote

import (
	"github.com/mph-llm-experiments/acore"
)

// BuildIdeaFilename constructs an idea filename using acore format.
func BuildIdeaFilename(id, title string) string {
	return acore.BuildFilename(id, title, TypeIdea)
}
