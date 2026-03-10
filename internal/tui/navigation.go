package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

// handleNavKey delegates a keypress to the NavigationHandler.
// Returns true if the key was a navigation key and the cursor was updated.
func (m *Model) handleNavKey(key string) bool {
	switch key {
	case "j", "down", "k", "up", "g", "G", "ctrl+d", "ctrl+u":
		m.nav.HandleKey(key)
		return true
	}
	return false
}

// handleWindowSize handles terminal resize events.
func handleWindowSize(m Model, msg tea.WindowSizeMsg) Model {
	m.width = msg.Width
	m.height = msg.Height
	return m
}
