package cli

import (
	"flag"
	"fmt"
	"strings"

	"github.com/mph-llm-experiments/anote/internal/config"
	"github.com/mph-llm-experiments/anote/internal/idea"
)

// tagList is a flag.Value that collects repeated --tag flags.
type tagList []string

func (t *tagList) String() string { return strings.Join(*t, ", ") }
func (t *tagList) Set(value string) error {
	*t = append(*t, strings.TrimSpace(value))
	return nil
}

func ideaNewCommand(cfg *config.Config) *Command {
	var tags tagList

	cmd := &Command{
		Name:        "new",
		Usage:       "anote new [--tag TAG]... <title>",
		Description: "Create a new idea",
		Flags:       flag.NewFlagSet("new", flag.ContinueOnError),
	}

	cmd.Flags.Var(&tags, "tag", "Add a tag (can be repeated)")

	cmd.Run = func(c *Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("title required: anote new \"My idea title\"")
		}

		title := strings.Join(args, " ")

		created, err := idea.CreateIdea(cfg.IdeasDirectory, title, []string(tags))
		if err != nil {
			return err
		}

		if !globalFlags.Quiet {
			fmt.Printf("Created idea #%d: %q (%s)\n", created.IndexID, created.IdeaMetadata.Title, created.File.Path)
		}

		return nil
	}

	return cmd
}
