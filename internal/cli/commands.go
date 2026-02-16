package cli

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

// Command represents a CLI command.
type Command struct {
	Name        string
	Usage       string
	Description string
	Flags       *flag.FlagSet
	Run         func(cmd *Command, args []string) error
	Subcommands []*Command
}

// Execute runs the command, dispatching to subcommands if appropriate.
func (c *Command) Execute(args []string) error {
	if len(c.Subcommands) > 0 && len(args) > 0 && !strings.HasPrefix(args[0], "-") {
		for _, sub := range c.Subcommands {
			if sub.Name == args[0] {
				return sub.Execute(args[1:])
			}
		}
	}

	if c.Flags != nil {
		if err := c.Flags.Parse(args); err != nil {
			return err
		}
		args = c.Flags.Args()
	}

	if c.Run != nil {
		return c.Run(c, args)
	}

	c.PrintUsage()
	return nil
}

// PrintUsage prints command usage to stderr.
func (c *Command) PrintUsage() {
	fmt.Fprintf(os.Stderr, "Usage: %s\n\n", c.Usage)
	if c.Description != "" {
		fmt.Fprintf(os.Stderr, "%s\n\n", c.Description)
	}

	if len(c.Subcommands) > 0 {
		fmt.Fprintf(os.Stderr, "Commands:\n")
		maxLen := 0
		for _, sub := range c.Subcommands {
			if len(sub.Name) > maxLen {
				maxLen = len(sub.Name)
			}
		}
		for _, sub := range c.Subcommands {
			desc := sub.Description
			if idx := strings.Index(desc, "\n"); idx >= 0 {
				desc = desc[:idx]
			}
			fmt.Fprintf(os.Stderr, "  %-*s  %s\n", maxLen+2, sub.Name, desc)
		}
		fmt.Fprintln(os.Stderr)
	}

	if c.Flags != nil {
		fmt.Fprintf(os.Stderr, "Flags:\n")
		c.Flags.PrintDefaults()
	}
}

// GlobalFlags holds global CLI flags.
type GlobalFlags struct {
	Config  string
	Dir     string
	NoColor bool
	JSON    bool
	Quiet   bool
}

var globalFlags GlobalFlags

// GetGlobalFlags returns the current global flags.
func GetGlobalFlags() *GlobalFlags {
	return &globalFlags
}

// ParseGlobalFlags extracts global flags from args, returning remaining args.
func ParseGlobalFlags(args []string) ([]string, error) {
	var remaining []string
	i := 0
	for i < len(args) {
		arg := args[i]

		// Flags that take a value
		if (arg == "--config" || arg == "--dir") && i+1 < len(args) {
			switch arg {
			case "--config":
				globalFlags.Config = args[i+1]
			case "--dir":
				globalFlags.Dir = args[i+1]
			}
			i += 2
			continue
		}

		// Boolean flags
		switch arg {
		case "--no-color":
			globalFlags.NoColor = true
			i++
			continue
		case "--json":
			globalFlags.JSON = true
			i++
			continue
		case "--quiet", "-q":
			globalFlags.Quiet = true
			i++
			continue
		}

		// --flag=value syntax
		if strings.HasPrefix(arg, "--config=") {
			globalFlags.Config = strings.TrimPrefix(arg, "--config=")
			i++
			continue
		}
		if strings.HasPrefix(arg, "--dir=") {
			globalFlags.Dir = strings.TrimPrefix(arg, "--dir=")
			i++
			continue
		}

		remaining = append(remaining, arg)
		i++
	}

	return remaining, nil
}
