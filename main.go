package main

import (
	"fmt"
	"os"

	"github.com/mph-llm-experiments/anote/internal/cli"
	"github.com/mph-llm-experiments/anote/internal/config"
)

var version = "0.2.0"
var specVersion = "0.2.0"

func main() {
	for _, arg := range os.Args[1:] {
		if arg == "--version" || arg == "-version" {
			fmt.Printf("anote v%s (spec v%s)\n", version, specVersion)
			os.Exit(0)
		}
	}

	cfg, err := config.Load("")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	if err := cli.Run(cfg, os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
