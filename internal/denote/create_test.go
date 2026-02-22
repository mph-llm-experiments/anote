package denote

import (
	"testing"
)

func TestBuildIdeaFilename(t *testing.T) {
	filename := BuildIdeaFilename("01TESTID0000000000000000AB", "My Great Idea")
	want := "01TESTID0000000000000000AB--my-great-idea__idea.md"
	if filename != want {
		t.Errorf("BuildIdeaFilename: got %q, want %q", filename, want)
	}
}
