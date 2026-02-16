package cli

import (
	"encoding/json"
	"flag"
	"fmt"
	"sort"
	"strings"

	"github.com/mph-llm-experiments/anote/internal/config"
	"github.com/mph-llm-experiments/anote/internal/denote"
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

// isTerminalState returns true for end-of-lifecycle states.
func isTerminalState(state string) bool {
	switch state {
	case denote.StateImplemented, denote.StateArchived, denote.StateRejected, denote.StateDropped:
		return true
	}
	return false
}

func ideaListCommand(cfg *config.Config) *Command {
	var (
		all      bool
		state    string
		maturity string
		tag      string
	)

	cmd := &Command{
		Name:        "list",
		Usage:       "anote list [--state STATE] [--maturity LEVEL] [--tag TAG] [-a]",
		Description: "List ideas",
		Flags:       flag.NewFlagSet("list", flag.ContinueOnError),
	}

	cmd.Flags.BoolVar(&all, "a", false, "Show all ideas including terminal states")
	cmd.Flags.BoolVar(&all, "all", false, "Show all ideas including terminal states")
	cmd.Flags.StringVar(&state, "state", "", "Filter by state")
	cmd.Flags.StringVar(&maturity, "maturity", "", "Filter by maturity")
	cmd.Flags.StringVar(&tag, "tag", "", "Filter by tag")

	cmd.Run = func(c *Command, args []string) error {
		scanner := denote.NewScanner(cfg.IdeasDirectory)
		ideas, err := scanner.FindIdeas()
		if err != nil {
			return fmt.Errorf("failed to scan ideas: %w", err)
		}

		// Sort by modification time, most recent first
		sort.Slice(ideas, func(i, j int) bool {
			return ideas[i].ModTime.After(ideas[j].ModTime)
		})

		// Filter
		var filtered []*denote.Idea
		for _, i := range ideas {
			// Default: exclude terminal states unless -a or specific --state
			if !all && state == "" && isTerminalState(i.State) {
				continue
			}

			if state != "" && i.State != state {
				continue
			}

			if maturity != "" && i.Maturity != maturity {
				continue
			}

			if tag != "" {
				found := false
				for _, t := range i.IdeaMetadata.Tags {
					if t == tag {
						found = true
						break
					}
				}
				// Also check filename tags
				if !found {
					for _, t := range i.File.Tags {
						if t == tag {
							found = true
							break
						}
					}
				}
				if !found {
					continue
				}
			}

			filtered = append(filtered, i)
		}

		// JSON output
		if globalFlags.JSON {
			data, err := json.MarshalIndent(filtered, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal JSON: %w", err)
			}
			fmt.Println(string(data))
			return nil
		}

		// Tabular output
		if len(filtered) == 0 {
			if !globalFlags.Quiet {
				fmt.Println("No ideas found.")
			}
			return nil
		}

		// Header
		fmt.Printf("%-5s %-12s %-8s %-40s %s\n", "#", "STATE", "MATURITY", "TITLE", "TAGS")
		fmt.Printf("%-5s %-12s %-8s %-40s %s\n", "---", "-----", "--------", strings.Repeat("-", 40), "----")

		for _, i := range filtered {
			mat := i.Maturity
			if mat == "" {
				mat = "-"
			}

			tags := strings.Join(i.IdeaMetadata.Tags, ", ")

			title := i.IdeaMetadata.Title
			if len(title) > 40 {
				title = title[:37] + "..."
			}

			fmt.Printf("%-5d %-12s %-8s %-40s %s\n", i.IndexID, i.State, mat, title, tags)
		}

		return nil
	}

	return cmd
}
