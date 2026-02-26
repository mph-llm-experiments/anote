package cli

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/mph-llm-experiments/acore"
	"github.com/mph-llm-experiments/anote/internal/config"
	"github.com/mph-llm-experiments/anote/internal/denote"
	"github.com/mph-llm-experiments/anote/internal/idea"
)

func ideaNewCommand(cfg *config.Config) *Command {
	cmd := &Command{
		Name:        "new",
		Usage:       "anote new [--tag TAG]... [--kind KIND] <title>",
		Description: "Create a new idea",
	}

	cmd.Run = func(c *Command, args []string) error {
		// Manual flag parsing to allow: new "title" --tag X or new --tag X "title"
		var tags []string
		var titleParts []string
		var kind string
		for idx := 0; idx < len(args); idx++ {
			if args[idx] == "--tag" && idx+1 < len(args) {
				tags = append(tags, strings.TrimSpace(args[idx+1]))
				idx++
			} else if args[idx] == "--kind" && idx+1 < len(args) {
				kind = strings.TrimSpace(args[idx+1])
				idx++
			} else if !strings.HasPrefix(args[idx], "-") {
				titleParts = append(titleParts, args[idx])
			}
		}

		if len(titleParts) == 0 {
			return fmt.Errorf("title required: anote new \"My idea title\"")
		}

		if kind != "" && !denote.IsValidKind(kind) {
			return fmt.Errorf("invalid kind %q: use aspiration, belief, plan, note, or fact", kind)
		}

		title := strings.Join(titleParts, " ")

		created, err := idea.CreateIdea(cfg.IdeasDirectory, title, tags, kind)
		if err != nil {
			return err
		}

		if globalFlags.JSON {
			data, _ := json.MarshalIndent(created, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		if !globalFlags.Quiet {
			fmt.Printf("Created idea #%d: %q (%s)\n", created.IndexID, created.Title, created.FilePath)
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
		all        bool
		state      string
		maturity   string
		tag        string
		kindFilter string
		plannedFor string
	)

	cmd := &Command{
		Name:        "list",
		Usage:       "anote list [--state STATE] [--maturity LEVEL] [--kind KIND] [--tag TAG] [--planned-for DATE] [-a]",
		Description: "List ideas",
		Flags:       flag.NewFlagSet("list", flag.ContinueOnError),
	}

	cmd.Flags.BoolVar(&all, "a", false, "Show all ideas including terminal states")
	cmd.Flags.BoolVar(&all, "all", false, "Show all ideas including terminal states")
	cmd.Flags.StringVar(&state, "state", "", "Filter by state (accepts display labels like considering)")
	cmd.Flags.StringVar(&maturity, "maturity", "", "Filter by maturity")
	cmd.Flags.StringVar(&tag, "tag", "", "Filter by tag")
	cmd.Flags.StringVar(&kindFilter, "kind", "", "Filter by kind (aspiration, belief, plan, note, or fact)")
	cmd.Flags.StringVar(&plannedFor, "planned-for", "", "Filter by planned_for date (today, YYYY-MM-DD, or any)")

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

		// Resolve display label to canonical state for filtering
		filterState := state
		filterStateKind := ""
		if state != "" {
			filterState, filterStateKind = denote.ResolveDisplayState(state)
			if !denote.IsValidState(filterState) {
				return fmt.Errorf("invalid state %q", state)
			}
		}

		// Filter
		var filtered []*denote.Idea
		for _, i := range ideas {
			effectiveKind := i.Kind
			if effectiveKind == "" {
				effectiveKind = denote.KindAspiration
			}

			// Default: exclude terminal states unless -a or specific --state
			if !all && filterState == "" && isTerminalState(i.State) {
				continue
			}

			if filterState != "" && i.State != filterState {
				continue
			}

			// If user typed a kind-specific label, also filter by the implied kind
			if filterStateKind != "" && effectiveKind != filterStateKind {
				continue
			}

			if kindFilter != "" && effectiveKind != kindFilter {
				continue
			}

			if maturity != "" && i.Maturity != maturity {
				continue
			}

			if tag != "" && !i.HasTag(tag) {
				continue
			}

			if plannedFor != "" {
				switch strings.ToLower(plannedFor) {
				case "any":
					if i.PlannedFor == "" {
						continue
					}
				case "today":
					if i.PlannedFor != time.Now().Format("2006-01-02") {
						continue
					}
				default:
					if i.PlannedFor != plannedFor {
						continue
					}
				}
			}

			filtered = append(filtered, i)
		}

		// JSON output — use kind-specific display labels
		if globalFlags.JSON {
			type jsonIdea struct {
				denote.Idea
			}
			var output []jsonIdea
			for _, i := range filtered {
				ek := i.Kind
				if ek == "" {
					ek = denote.KindAspiration
				}
				ji := jsonIdea{Idea: *i}
				ji.State = denote.DisplayState(i.State, ek)
				ji.Kind = ek
				output = append(output, ji)
			}
			data, err := json.MarshalIndent(output, "", "  ")
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
		fmt.Printf("%-5s %-4s %-14s %-8s %-38s %s\n", "#", "KIND", "STATE", "MATURITY", "TITLE", "TAGS")
		fmt.Printf("%-5s %-4s %-14s %-8s %-38s %s\n", "---", "----", "-----", "--------", strings.Repeat("-", 38), "----")

		for _, i := range filtered {
			effectiveKind := i.Kind
			if effectiveKind == "" {
				effectiveKind = denote.KindAspiration
			}
			displayState := denote.DisplayState(i.State, effectiveKind)

			kindShort := "A"
			switch effectiveKind {
			case denote.KindBelief:
				kindShort = "B"
			case denote.KindPlan:
				kindShort = "P"
			case denote.KindNote:
				kindShort = "N"
			case denote.KindFact:
				kindShort = "F"
			}

			mat := i.Maturity
			if mat == "" {
				mat = "-"
			}
			if denote.IsSimpleKind(effectiveKind) {
				mat = "-"
			}

			tags := strings.Join(i.Tags, ", ")

			title := i.Title
			if len(title) > 38 {
				title = title[:35] + "..."
			}

			fmt.Printf("%-5d %-4s %-14s %-8s %-38s %s\n", i.IndexID, kindShort, displayState, mat, title, tags)
		}

		return nil
	}

	return cmd
}

var denoteIDPattern = regexp.MustCompile(`^\d{8}T\d{6}$`)

// lookupIdea finds an idea by index_id (number) or entity ID (ULID or legacy Denote ID).
func lookupIdea(dir string, ref string) (*denote.Idea, error) {
	// Try as Denote timestamp ID
	if denoteIDPattern.MatchString(ref) {
		return idea.FindIdeaByEntityID(dir, ref)
	}

	// Try as integer index_id
	id, err := strconv.Atoi(ref)
	if err != nil {
		// Must be a ULID
		return idea.FindIdeaByEntityID(dir, ref)
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

		effectiveKind := i.Kind
		if effectiveKind == "" {
			effectiveKind = denote.KindAspiration
		}
		displayState := denote.DisplayState(i.State, effectiveKind)

		// JSON output
		if globalFlags.JSON {
			type jsonIdea struct {
				denote.Idea
				Content string `json:"content,omitempty"`
			}
			ji := jsonIdea{
				Idea:    *i,
				Content: extractIdeaContent(i.Content),
			}
			ji.State = displayState
			ji.Kind = effectiveKind
			data, err := json.MarshalIndent(ji, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal JSON: %w", err)
			}
			fmt.Println(string(data))
			return nil
		}

		// Formatted output
		fmt.Printf("Idea #%d: %s\n", i.IndexID, i.Title)
		fmt.Printf("ID:         %s\n", i.ID)
		fmt.Printf("Kind:       %s\n", effectiveKind)
		fmt.Printf("State:      %s\n", displayState)
		if i.Maturity != "" {
			fmt.Printf("Maturity:   %s\n", i.Maturity)
		}
		if len(i.Tags) > 0 {
			fmt.Printf("Tags:       %s\n", strings.Join(i.Tags, ", "))
		}
		if i.Created != "" {
			fmt.Printf("Created:    %s\n", i.Created)
		}
		if i.Modified != "" {
			fmt.Printf("Modified:   %s\n", i.Modified)
		}
		if i.PlannedFor != "" {
			fmt.Printf("Planned:    %s\n", i.PlannedFor)
		}
		if i.RejectedReason != "" {
			fmt.Printf("Rejected:   %s\n", i.RejectedReason)
		}

		// Related ideas — resolve titles
		if len(i.RelatedIdeas) > 0 {
			fmt.Printf("Related ideas:\n")
			scanner := denote.NewScanner(cfg.IdeasDirectory)
			allIdeas, _ := scanner.FindIdeas()
			idMap := make(map[string]string)
			for _, a := range allIdeas {
				idMap[a.ID] = a.Title
			}
			for _, relID := range i.RelatedIdeas {
				title, ok := idMap[relID]
				if ok {
					fmt.Printf("  - %s (%s)\n", title, relID)
				} else {
					fmt.Printf("  - %s\n", relID)
				}
			}
		}

		// Related tasks
		if len(i.RelatedTasks) > 0 {
			fmt.Printf("Related tasks:\n")
			for _, taskID := range i.RelatedTasks {
				fmt.Printf("  - %s\n", taskID)
			}
		}

		// Related people
		if len(i.RelatedPeople) > 0 {
			fmt.Printf("Related people:\n")
			for _, personID := range i.RelatedPeople {
				fmt.Printf("  - %s\n", personID)
			}
		}

		fmt.Printf("File:       %s\n", i.FilePath)

		// Show content (everything after frontmatter)
		content := extractIdeaContent(i.Content)
		if content != "" {
			fmt.Printf("\n%s", content)
		}

		return nil
	}

	return cmd
}

func ideaUpdateCommand(cfg *config.Config) *Command {
	cmd := &Command{
		Name:        "update",
		Usage:       "anote update <id> [--title TITLE] [--state STATE] [--maturity LEVEL] [--kind KIND] [--plan-for DATE]",
		Description: "Update idea title, state, maturity, or kind",
	}

	cmd.Run = func(c *Command, args []string) error {
		// Manual flag parsing to allow: update <id> --state X or update --state X <id>
		var state, maturity, kind, title, body, idRef, planFor string
		var addPerson, removePerson, addTask, removeTask, addIdea, removeIdea string
		for idx := 0; idx < len(args); idx++ {
			switch args[idx] {
			case "--title":
				if idx+1 < len(args) {
					title = args[idx+1]
					idx++
				}
			case "--body":
				if idx+1 < len(args) {
					body = args[idx+1]
					idx++
				}
			case "--state":
				if idx+1 < len(args) {
					state = args[idx+1]
					idx++
				}
			case "--maturity":
				if idx+1 < len(args) {
					maturity = args[idx+1]
					idx++
				}
			case "--kind":
				if idx+1 < len(args) {
					kind = args[idx+1]
					idx++
				}
			case "--plan-for":
				if idx+1 < len(args) {
					planFor = args[idx+1]
					idx++
				}
			case "--add-person":
				if idx+1 < len(args) {
					addPerson = args[idx+1]
					idx++
				}
			case "--remove-person":
				if idx+1 < len(args) {
					removePerson = args[idx+1]
					idx++
				}
			case "--add-task":
				if idx+1 < len(args) {
					addTask = args[idx+1]
					idx++
				}
			case "--remove-task":
				if idx+1 < len(args) {
					removeTask = args[idx+1]
					idx++
				}
			case "--add-idea":
				if idx+1 < len(args) {
					addIdea = args[idx+1]
					idx++
				}
			case "--remove-idea":
				if idx+1 < len(args) {
					removeIdea = args[idx+1]
					idx++
				}
			default:
				if !strings.HasPrefix(args[idx], "-") && idRef == "" {
					idRef = args[idx]
				}
			}
		}

		if idRef == "" {
			return fmt.Errorf("idea ID required: anote update <id> --state STATE")
		}

		hasRelationUpdate := addPerson != "" || removePerson != "" || addTask != "" || removeTask != "" || addIdea != "" || removeIdea != ""
		if state == "" && maturity == "" && kind == "" && title == "" && body == "" && planFor == "" && !hasRelationUpdate {
			return fmt.Errorf("nothing to update: provide --title, --body, --state, --maturity, --kind, --plan-for, or relationship flags")
		}

		i, err := lookupIdea(cfg.IdeasDirectory, idRef)
		if err != nil {
			return err
		}

		if title != "" {
			i.Title = title
		}

		// Determine effective kind (consider --kind flag if provided)
		effectiveKind := i.Kind
		if kind != "" {
			effectiveKind = kind
		}
		if effectiveKind == "" {
			effectiveKind = denote.KindAspiration
		}

		// Resolve display label to canonical state
		if state != "" {
			resolved, _ := denote.ResolveDisplayState(state)
			if resolved == denote.StateRejected {
				return fmt.Errorf("use 'anote reject <id> \"reason\"' to reject an idea")
			}
			state = resolved
			if !denote.IsValidState(state) {
				return fmt.Errorf("invalid state %q", state)
			}
			if denote.IsSimpleKind(effectiveKind) && state != denote.StateActive && state != denote.StateArchived {
				return fmt.Errorf("note and fact kinds only support active and archived states")
			}
			i.State = state
		}

		// Validate and set kind
		if kind != "" {
			if !denote.IsValidKind(kind) {
				return fmt.Errorf("invalid kind %q: use aspiration, belief, plan, note, or fact", kind)
			}
			i.Kind = kind
		}

		// Validate and set maturity
		if maturity != "" {
			if denote.IsSimpleKind(effectiveKind) {
				return fmt.Errorf("note and fact kinds do not use maturity")
			}
			if !denote.IsValidMaturity(maturity) {
				return fmt.Errorf("invalid maturity %q: use crawl, walk, or run", maturity)
			}
			i.Maturity = maturity
		}

		// Resolve and set planned_for
		if planFor != "" {
			if strings.ToLower(planFor) == "none" {
				i.PlannedFor = ""
			} else {
				parsed, err := acore.ParseNaturalDate(planFor)
				if err != nil {
					return fmt.Errorf("invalid --plan-for date: %v", err)
				}
				i.PlannedFor = parsed
			}
		}

		// Apply cross-app relationship updates
		if addPerson != "" {
			acore.AddRelation(&i.RelatedPeople, addPerson)
			acore.SyncRelation(i.Type, i.ID, addPerson)
		}
		if removePerson != "" {
			acore.RemoveRelation(&i.RelatedPeople, removePerson)
			acore.UnsyncRelation(i.Type, i.ID, removePerson)
		}
		if addTask != "" {
			acore.AddRelation(&i.RelatedTasks, addTask)
			acore.SyncRelation(i.Type, i.ID, addTask)
		}
		if removeTask != "" {
			acore.RemoveRelation(&i.RelatedTasks, removeTask)
			acore.UnsyncRelation(i.Type, i.ID, removeTask)
		}
		if addIdea != "" {
			acore.AddRelation(&i.RelatedIdeas, addIdea)
			acore.SyncRelation(i.Type, i.ID, addIdea)
		}
		if removeIdea != "" {
			acore.RemoveRelation(&i.RelatedIdeas, removeIdea)
			acore.UnsyncRelation(i.Type, i.ID, removeIdea)
		}

		i.Modified = time.Now().Format(time.RFC3339)

		if body != "" {
			// Replace description (content before ## Log section), preserve log
			existingContent := extractIdeaContent(i.Content)
			newContent := replaceDescription(existingContent, body)
			if err := denote.WriteIdeaFile(i.FilePath, i, newContent); err != nil {
				return fmt.Errorf("failed to update idea: %w", err)
			}
		} else {
			if err := denote.UpdateIdeaFrontmatter(i.FilePath, i); err != nil {
				return fmt.Errorf("failed to update idea: %w", err)
			}
		}

		if globalFlags.JSON {
			reloaded, err := denote.ParseIdeaFile(i.FilePath)
			if err != nil {
				reloaded = i
			}
			data, _ := json.MarshalIndent(reloaded, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		if !globalFlags.Quiet {
			fmt.Printf("Updated idea #%d: %q", i.IndexID, i.Title)
			if kind != "" {
				fmt.Printf(" [kind: %s]", kind)
			}
			if state != "" {
				fmt.Printf(" [state: %s]", denote.DisplayState(state, effectiveKind))
			}
			if maturity != "" {
				fmt.Printf(" [maturity: %s]", maturity)
			}
			fmt.Println()

			// Encourage project link when aspiration goes active
			if state == denote.StateActive && len(i.RelatedTasks) == 0 && effectiveKind == denote.KindAspiration {
				fmt.Println("Hint: Consider linking an atask project with 'anote project <id> <project-id>'")
			}
		}

		return nil
	}

	return cmd
}

func ideaDeleteCommand(cfg *config.Config) *Command {
	cmd := &Command{
		Name:        "delete",
		Usage:       "anote delete <id> [--confirm]",
		Description: "Delete an idea file",
	}

	cmd.Run = func(c *Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("usage: anote delete <id> [--confirm]")
		}

		confirm := false
		idRef := ""
		for _, arg := range args {
			if arg == "--confirm" {
				confirm = true
			} else if idRef == "" {
				idRef = arg
			}
		}
		if idRef == "" {
			return fmt.Errorf("usage: anote delete <id> [--confirm]")
		}

		i, err := lookupIdea(cfg.IdeasDirectory, idRef)
		if err != nil {
			return err
		}

		if !confirm {
			return fmt.Errorf("use --confirm to delete idea '%s' (%s)", i.Title, i.FilePath)
		}

		if err := os.Remove(i.FilePath); err != nil {
			return fmt.Errorf("failed to delete idea: %w", err)
		}

		if globalFlags.JSON {
			result := map[string]interface{}{
				"deleted":  true,
				"index_id": i.IndexID,
				"title":    i.Title,
				"file":     i.FilePath,
			}
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		if !globalFlags.Quiet {
			fmt.Printf("Deleted idea #%d: %s\n", i.IndexID, i.Title)
		}
		return nil
	}

	return cmd
}

func ideaRejectCommand(cfg *config.Config) *Command {
	cmd := &Command{
		Name:        "reject",
		Usage:       "anote reject <id> <reason>",
		Description: "Reject an idea (reason required)",
	}

	cmd.Run = func(c *Command, args []string) error {
		if len(args) < 2 {
			return fmt.Errorf("usage: anote reject <id> \"reason for rejection\"")
		}

		i, err := lookupIdea(cfg.IdeasDirectory, args[0])
		if err != nil {
			return err
		}

		effectiveKind := i.Kind
		if effectiveKind == "" {
			effectiveKind = denote.KindAspiration
		}
		if denote.IsSimpleKind(effectiveKind) {
			return fmt.Errorf("note and fact kinds cannot be rejected; use 'anote update %s --state archived' instead", args[0])
		}

		reason := strings.Join(args[1:], " ")
		if strings.TrimSpace(reason) == "" {
			return fmt.Errorf("rejection reason cannot be empty")
		}

		i.State = denote.StateRejected
		i.RejectedReason = reason
		i.Modified = time.Now().Format(time.RFC3339)

		if err := denote.UpdateIdeaFrontmatter(i.FilePath, i); err != nil {
			return fmt.Errorf("failed to reject idea: %w", err)
		}

		if !globalFlags.Quiet {
			fmt.Printf("Rejected idea #%d: %q — %s\n", i.IndexID, i.Title, reason)
		}

		return nil
	}

	return cmd
}

func ideaTagCommand(cfg *config.Config) *Command {
	cmd := &Command{
		Name:        "tag",
		Usage:       "anote tag <id> <tag> [--remove]",
		Description: "Add or remove a tag",
	}

	cmd.Run = func(c *Command, args []string) error {
		// Parse: tag <id> <tagname> [--remove] or tag <id> --remove <tagname>
		var idRef, tagName string
		remove := false
		for _, arg := range args {
			if arg == "--remove" {
				remove = true
			} else if !strings.HasPrefix(arg, "-") {
				if idRef == "" {
					idRef = arg
				} else if tagName == "" {
					tagName = arg
				}
			}
		}

		if idRef == "" || tagName == "" {
			return fmt.Errorf("usage: anote tag <id> <tag> [--remove]")
		}

		i, err := lookupIdea(cfg.IdeasDirectory, idRef)
		if err != nil {
			return err
		}

		if remove {
			// Remove from frontmatter tags
			var newTags []string
			for _, t := range i.Tags {
				if t != tagName {
					newTags = append(newTags, t)
				}
			}
			i.Tags = newTags

			if !globalFlags.Quiet {
				fmt.Printf("Removed tag %q from idea #%d\n", tagName, i.IndexID)
			}
		} else {
			// Add to frontmatter tags (skip duplicates)
			found := false
			for _, t := range i.Tags {
				if t == tagName {
					found = true
					break
				}
			}
			if !found {
				i.Tags = append(i.Tags, tagName)
			}

			if !globalFlags.Quiet {
				fmt.Printf("Added tag %q to idea #%d\n", tagName, i.IndexID)
			}
		}

		i.Modified = time.Now().Format(time.RFC3339)
		if err := denote.UpdateIdeaFrontmatter(i.FilePath, i); err != nil {
			return fmt.Errorf("failed to update frontmatter: %w", err)
		}

		return nil
	}

	return cmd
}

func ideaLinkCommand(cfg *config.Config) *Command {
	cmd := &Command{
		Name:        "link",
		Usage:       "anote link <id1> <id2>",
		Description: "Link two related ideas (bidirectional)",
	}

	cmd.Run = func(c *Command, args []string) error {
		if len(args) < 2 {
			return fmt.Errorf("usage: anote link <id1> <id2>")
		}

		idea1, err := lookupIdea(cfg.IdeasDirectory, args[0])
		if err != nil {
			return fmt.Errorf("first idea: %w", err)
		}

		idea2, err := lookupIdea(cfg.IdeasDirectory, args[1])
		if err != nil {
			return fmt.Errorf("second idea: %w", err)
		}

		now := time.Now().Format(time.RFC3339)

		// Add idea2's ID to idea1's related ideas (skip duplicates)
		if !containsStr(idea1.RelatedIdeas, idea2.ID) {
			acore.AddRelation(&idea1.RelatedIdeas, idea2.ID)
			idea1.Modified = now
			if err := denote.UpdateIdeaFrontmatter(idea1.FilePath, idea1); err != nil {
				return fmt.Errorf("failed to update idea #%d: %w", idea1.IndexID, err)
			}
		}

		// Add idea1's ID to idea2's related ideas (skip duplicates)
		if !containsStr(idea2.RelatedIdeas, idea1.ID) {
			acore.AddRelation(&idea2.RelatedIdeas, idea1.ID)
			idea2.Modified = now
			if err := denote.UpdateIdeaFrontmatter(idea2.FilePath, idea2); err != nil {
				return fmt.Errorf("failed to update idea #%d: %w", idea2.IndexID, err)
			}
		}

		if !globalFlags.Quiet {
			fmt.Printf("Linked idea #%d ↔ idea #%d\n", idea1.IndexID, idea2.IndexID)
		}

		return nil
	}

	return cmd
}

func ideaProjectCommand(cfg *config.Config) *Command {
	cmd := &Command{
		Name:        "project",
		Usage:       "anote project <id> <project-id>",
		Description: "Link idea to an atask project",
	}

	cmd.Run = func(c *Command, args []string) error {
		if len(args) < 2 {
			return fmt.Errorf("usage: anote project <id> <project-id>")
		}

		i, err := lookupIdea(cfg.IdeasDirectory, args[0])
		if err != nil {
			return err
		}

		projectID := args[1]

		// Skip duplicates
		if containsStr(i.RelatedTasks, projectID) {
			if !globalFlags.Quiet {
				fmt.Printf("Idea #%d is already linked to project %s\n", i.IndexID, projectID)
			}
			return nil
		}

		acore.AddRelation(&i.RelatedTasks, projectID)
		i.Modified = time.Now().Format(time.RFC3339)

		if err := denote.UpdateIdeaFrontmatter(i.FilePath, i); err != nil {
			return fmt.Errorf("failed to update idea: %w", err)
		}

		if !globalFlags.Quiet {
			fmt.Printf("Linked idea #%d to project %s\n", i.IndexID, projectID)
		}

		return nil
	}

	return cmd
}

func containsStr(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func ideaLogCommand(cfg *config.Config) *Command {
	cmd := &Command{
		Name:        "log",
		Usage:       "anote log <id> \"message\"",
		Description: "Add a timestamped log entry to an idea",
	}

	cmd.Run = func(c *Command, args []string) error {
		if len(args) < 2 {
			return fmt.Errorf("usage: anote log <id> \"message\"")
		}

		idRef := args[0]
		message := strings.Join(args[1:], " ")

		i, err := lookupIdea(cfg.IdeasDirectory, idRef)
		if err != nil {
			return err
		}

		existingContent := extractIdeaContent(i.Content)
		newContent := addLogEntry(existingContent, message)

		i.Modified = time.Now().Format(time.RFC3339)
		if err := denote.WriteIdeaFile(i.FilePath, i, newContent); err != nil {
			return fmt.Errorf("failed to write idea: %w", err)
		}

		if !globalFlags.Quiet {
			fmt.Printf("Logged to idea #%d: %s\n", i.IndexID, i.Title)
		}

		return nil
	}

	return cmd
}

// replaceDescription replaces the description portion of content (before ## Log),
// preserving the log section.
func replaceDescription(content, newDesc string) string {
	logIdx := strings.Index(content, "\n## Log\n")
	if logIdx == -1 {
		logIdx = strings.Index(content, "## Log\n")
		if logIdx == 0 {
			// Content starts with ## Log, prepend description
			return newDesc + "\n\n" + content
		}
		// No log section, just replace everything
		return newDesc + "\n"
	}
	return newDesc + "\n" + content[logIdx:]
}

// addLogEntry appends a timestamped entry to the ## Log section.
func addLogEntry(content, message string) string {
	now := time.Now().Format("2006-01-02")
	entry := fmt.Sprintf("- **%s** %s", now, message)

	logIdx := strings.Index(content, "\n## Log\n")
	if logIdx != -1 {
		// Insert after the ## Log header
		insertAt := logIdx + len("\n## Log\n")
		return content[:insertAt] + entry + "\n" + content[insertAt:]
	}

	// Check if content starts with ## Log
	if strings.HasPrefix(content, "## Log\n") {
		insertAt := len("## Log\n")
		return content[:insertAt] + entry + "\n" + content[insertAt:]
	}

	// No log section yet, append one
	trimmed := strings.TrimRight(content, "\n")
	if trimmed == "" {
		return "## Log\n" + entry + "\n"
	}
	return trimmed + "\n\n## Log\n" + entry + "\n"
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
