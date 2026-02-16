package cli

import (
	"fmt"

	"github.com/mph-llm-experiments/anote/internal/config"
)

// Run executes the CLI with the given config and arguments.
func Run(cfg *config.Config, args []string) error {
	remaining, err := ParseGlobalFlags(args)
	if err != nil {
		return err
	}

	// Reload config if --config flag was provided
	if globalFlags.Config != "" {
		newCfg, err := config.Load(globalFlags.Config)
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}
		cfg = newCfg
	}

	// Override ideas directory if --dir flag was provided
	if globalFlags.Dir != "" {
		cfg.IdeasDirectory = globalFlags.Dir
	}

	root := &Command{
		Name:  "anote",
		Usage: "anote <command> [options]",
		Description: `Agent-first idea management using Denote file conventions.

Commands:
  new        Create a new idea
  list       List ideas
  show       Show idea details
  update     Update idea state or maturity
  reject     Reject an idea (with reason)
  tag        Add or remove tags
  link       Link related ideas

Global Options:
  --config PATH  Use specific config file
  --dir PATH     Override ideas directory
  --json         Output in JSON format
  --no-color     Disable color output
  --quiet, -q    Minimal output`,
	}

	root.Subcommands = append(root.Subcommands,
		ideaNewCommand(cfg),
		ideaListCommand(cfg),
		ideaShowCommand(cfg),
	)

	if len(remaining) == 0 {
		root.PrintUsage()
		return nil
	}

	return root.Execute(remaining)
}
