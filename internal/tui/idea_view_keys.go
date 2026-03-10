package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mph-llm-experiments/anote/internal/denote"
)

// handleIdeaViewKey handles keypresses in the idea detail view (ModeIdeaView).
func (m Model) handleIdeaViewKey(key string) (tea.Model, tea.Cmd) {
	switch key {
	case "q", "esc":
		if m.editingField != "" {
			// Cancel inline edit without leaving idea view
			m.editingField = ""
			m.editBuf.Clear()
		} else {
			m.mode = ModeNormal
			m.viewingIdea = nil
		}

	case "s":
		// Open kind-aware state picker
		if m.viewingIdea != nil {
			m.menuOptions = m.kindsConfig.ValidStatesFor(m.viewingIdea.Kind)
			if m.menuOptions == nil {
				m.menuOptions = []string{}
			}
			m.menuCursor = 0
			m.menuField = FieldState
			m.mode = ModeStateMenu
		}

	case "p":
		// Open purpose picker using the menu
		if m.viewingIdea != nil {
			// Build options list: empty string (unattached) + all purpose IDs
			m.menuOptions = make([]string, 0, len(m.purposes)+1)
			m.menuOptions = append(m.menuOptions, "") // unattached
			for _, p := range m.purposes {
				m.menuOptions = append(m.menuOptions, p.ID)
			}
			m.menuCursor = 0
			m.menuField = FieldPurpose
			m.mode = ModeStateMenu // reuse state menu mode for now
		}

	case "m":
		// Cycle maturity (no menu needed — just three values)
		if m.viewingIdea != nil {
			switch m.viewingIdea.Maturity {
			case "":
				m.viewingIdea.Maturity = denote.MaturityCrawl
			case denote.MaturityCrawl:
				m.viewingIdea.Maturity = denote.MaturityWalk
			case denote.MaturityWalk:
				m.viewingIdea.Maturity = denote.MaturityRun
			default:
				m.viewingIdea.Maturity = ""
			}
			if err := persistIdeaFrontmatter(m.viewingIdea); err != nil {
				m.statusMsg = "error saving maturity: " + err.Error()
			} else {
				if fresh, err := refreshIdea(m.viewingIdea); err == nil {
					m.viewingIdea = fresh
				}
				_ = m.loadIdeas()
			}
		}

	case "K":
		// Cycle kind to next in config
		if m.viewingIdea != nil {
			kinds := m.kindsConfig.AllKinds()
			for i, k := range kinds {
				if k == m.viewingIdea.Kind {
					if i+1 < len(kinds) {
						m.viewingIdea.Kind = kinds[i+1]
					} else {
						m.viewingIdea.Kind = kinds[0]
					}
					break
				}
			}
			if err := persistIdeaFrontmatter(m.viewingIdea); err != nil {
				m.statusMsg = "error saving kind: " + err.Error()
			} else {
				if fresh, err := refreshIdea(m.viewingIdea); err == nil {
					m.viewingIdea = fresh
				}
				_ = m.loadIdeas()
			}
		}

	case "t":
		if m.viewingIdea != nil {
			m.editBuf.SetValue(joinTags(m.viewingIdea.Tags))
			m.mode = ModeTagsEdit
		}

	case "l":
		m.logBuf.SetValue("")
		m.mode = ModeLogEntry

	case "r":
		if m.viewingIdea != nil {
			m.editingField = FieldTitle
			m.editBuf.SetValue(m.viewingIdea.Title)
		}

	case "enter":
		// Confirm inline rename if active
		if m.editingField == FieldTitle && m.viewingIdea != nil {
			newTitle := m.editBuf.Value()
			if newTitle != "" {
				m.viewingIdea.Title = newTitle
				if err := persistIdeaFrontmatter(m.viewingIdea); err != nil {
					m.statusMsg = "error saving title: " + err.Error()
				} else {
					if fresh, err := refreshIdea(m.viewingIdea); err == nil {
						m.viewingIdea = fresh
					}
					_ = m.loadIdeas()
				}
			}
			m.editingField = ""
			m.editBuf.Clear()
		}

	case "backspace":
		if m.editingField != "" {
			m.editBuf.Backspace()
		}

	case "x":
		m.mode = ModeConfirmDelete

	case "E":
		if m.viewingIdea != nil && m.viewingIdea.FilePath != "" {
			return m, openInEditor(m.viewingIdea.FilePath)
		}

	default:
		if m.editingField != "" && len(key) == 1 {
			m.editBuf.Insert(rune(key[0]))
		}
	}
	return m, nil
}

// joinTags joins a tag slice into a space-separated string.
func joinTags(tags []string) string {
	result := ""
	for i, t := range tags {
		if i > 0 {
			result += " "
		}
		result += t
	}
	return result
}
