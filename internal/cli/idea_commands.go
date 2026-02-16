package cli

import (
	"encoding/json"
	"flag"
	"fmt"
	"regexp"
	"sort"
	"strconv"
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

var denoteIDPattern = regexp.MustCompile(`^\d{8}T\d{6}$`)

// lookupIdea finds an idea by index_id (number) or Denote ID (timestamp).
func lookupIdea(dir string, ref string) (*denote.Idea, error) {
	if denoteIDPattern.MatchString(ref) {
		return idea.FindIdeaByDenoteID(dir, ref)
	}

	id, err := strconv.Atoi(ref)
	if err != nil {
		return nil, fmt.Errorf("invalid idea reference %q: use a number or Denote ID", ref)
	}

	return idea.FindIdeaByID(dir, id)
}

func ideaShowCommand(cfg *config.Config) *Command {
	cmd := &Command{
		Name:        "show",
		Usage:       "anote show <id>",
		Description: "Show idea details",
	}

	cmd.Run = func(c *Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("idea ID required: anote show <id>")
		}

		i, err := lookupIdea(cfg.IdeasDirectory, args[0])
		if err != nil {
			return err
		}

		// JSON output
		if globalFlags.JSON {
			data, err := json.MarshalIndent(i, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal JSON: %w", err)
			}
			fmt.Println(string(data))
			return nil
		}

		// Formatted output
		fmt.Printf("Idea #%d: %s\n", i.IndexID, i.IdeaMetadata.Title)
		fmt.Printf("Denote ID:  %s\n", i.File.ID)
		fmt.Printf("State:      %s\n", i.State)
		if i.Maturity != "" {
			fmt.Printf("Maturity:   %s\n", i.Maturity)
		}
		if len(i.IdeaMetadata.Tags) > 0 {
			fmt.Printf("Tags:       %s\n", strings.Join(i.IdeaMetadata.Tags, ", "))
		}
		if i.Created != "" {
			fmt.Printf("Created:    %s\n", i.Created)
		}
		if i.Modified != "" {
			fmt.Printf("Modified:   %s\n", i.Modified)
		}
		if i.RejectedReason != "" {
			fmt.Printf("Rejected:   %s\n", i.RejectedReason)
		}

		// Related ideas — resolve titles
		if len(i.Related) > 0 {
			fmt.Printf("Related:\n")
			scanner := denote.NewScanner(cfg.IdeasDirectory)
			allIdeas, _ := scanner.FindIdeas()
			idMap := make(map[string]string)
			for _, a := range allIdeas {
				idMap[a.File.ID] = a.IdeaMetadata.Title
			}
			for _, relID := range i.Related {
				title, ok := idMap[relID]
				if ok {
					fmt.Printf("  - %s (%s)\n", title, relID)
				} else {
					fmt.Printf("  - %s\n", relID)
				}
			}
		}

		// Linked projects
		if len(i.Project) > 0 {
			fmt.Printf("Projects:\n")
			for _, projID := range i.Project {
				fmt.Printf("  - %s\n", projID)
			}
		}

		fmt.Printf("File:       %s\n", i.File.Path)

		// Show content (everything after frontmatter)
		content := extractIdeaContent(i.Content)
		if content != "" {
			fmt.Printf("\n%s", content)
		}

		return nil
	}

	return cmd
}

// extractIdeaContent extracts the body content after YAML frontmatter.
func extractIdeaContent(fullContent string) string {
	if !strings.HasPrefix(fullContent, "---\n") {
		return fullContent
	}

	lines := strings.Split(fullContent, "\n")
	for idx, line := range lines {
		if idx == 0 {
			continue
		}
		if line == "---" {
			rest := strings.Join(lines[idx+1:], "\n")
			return strings.TrimPrefix(rest, "\n")
		}
	}

	return ""
}
