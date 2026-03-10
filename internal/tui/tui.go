package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mph-llm-experiments/anote/internal/config"
)

// Run launches the anote TUI. Returns when the user quits.
func Run(cfg *config.Config) error {
	m, err := NewModel(cfg)
	if err != nil {
		return fmt.Errorf("failed to initialize TUI: %w", err)
	}

	p := tea.NewProgram(m, tea.WithAltScreen())
	_, err = p.Run()
	return err
}
