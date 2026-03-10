package tui

import (
	"os"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// Update implements tea.Model.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m = handleWindowSize(m, msg)
		return m, nil
	case tea.KeyMsg:
		return m.handleKey(msg.String())
	}
	return m, nil
}

func (m Model) handleKey(key string) (tea.Model, tea.Cmd) {
	switch m.mode {
	case ModeIdeaView:
		return m.handleIdeaViewKey(key)
	case ModeSearch:
		return m.handleSearchKey(key)
	case ModeStateMenu:
		return m.handleMenuKey(key)
	case ModeCompliancePrompt:
		return m.handleComplianceKey(key)
	case ModeCreate:
		return m.handleCreateKey(key)
	case ModeLogEntry:
		return m.handleLogEntryKey(key)
	case ModeTagsEdit:
		return m.handleTagsEditKey(key)
	case ModeHelp:
		m.mode = ModeNormal
		return m, nil
	default:
		return m.handleNormalKey(key)
	}
}

func (m Model) handleNormalKey(key string) (tea.Model, tea.Cmd) {
	if m.handleNavKey(key) {
		return m, nil
	}

	switch key {
	case "q", "ctrl+c":
		return m, tea.Quit

	case "?":
		m.mode = ModeHelp
		return m, nil

	case "/":
		m.mode = ModeSearch
		m.editBuf.SetValue(m.searchQuery)
		return m, nil

	case "enter":
		if len(m.filtered) == 0 {
			return m, nil
		}
		idea := m.filtered[m.nav.Cursor()]
		m.viewingIdea = &idea
		// Compliance check: open compliance prompt if state is invalid for kind
		if !m.kindsConfig.IsCompliant(idea.Kind, idea.State) || !m.kindsConfig.KindExists(idea.Kind) {
			m.complianceOptions = m.kindsConfig.ValidStatesFor(idea.Kind)
			if m.complianceOptions == nil {
				m.complianceOptions = []string{}
			}
			m.complianceCursor = 0
			m.mode = ModeCompliancePrompt
		} else {
			m.mode = ModeIdeaView
		}
		return m, nil

	case "c":
		m.mode = ModeCreate
		m.createTitle = ""
		m.createKind = "note"
		m.createTags = ""
		m.createField = 0
		m.editBuf.SetValue("")
		return m, nil

	case "K":
		// Cycle kind filter through available kinds, then clear
		kinds := m.kindsConfig.AllKinds()
		if len(kinds) == 0 {
			return m, nil
		}
		if m.kindFilter == "" {
			m.kindFilter = kinds[0]
		} else {
			found := false
			for i, k := range kinds {
				if k == m.kindFilter {
					if i+1 < len(kinds) {
						m.kindFilter = kinds[i+1]
					} else {
						m.kindFilter = ""
					}
					found = true
					break
				}
			}
			if !found {
				m.kindFilter = ""
			}
		}
		m.applyFilters()
		return m, nil

	case "P":
		// Cycle purpose filter through available purposes, then clear
		if len(m.purposes) == 0 {
			return m, nil
		}
		if m.purposeFilter == "" {
			m.purposeFilter = m.purposes[0].ID
		} else {
			found := false
			for i, p := range m.purposes {
				if p.ID == m.purposeFilter {
					if i+1 < len(m.purposes) {
						m.purposeFilter = m.purposes[i+1].ID
					} else {
						m.purposeFilter = ""
					}
					found = true
					break
				}
			}
			if !found {
				m.purposeFilter = ""
			}
		}
		m.applyFilters()
		return m, nil

	case "esc":
		// Clear active filters on esc
		if m.kindFilter != "" || m.stateFilter != "" || m.purposeFilter != "" || m.searchQuery != "" {
			m.kindFilter = ""
			m.stateFilter = ""
			m.purposeFilter = ""
			m.searchQuery = ""
			m.applyFilters()
		}
		return m, nil
	}

	return m, nil
}

func (m Model) handleSearchKey(key string) (tea.Model, tea.Cmd) {
	switch key {
	case "enter", "esc":
		m.searchQuery = m.editBuf.Value()
		m.applyFilters()
		m.mode = ModeNormal
	case "backspace":
		m.editBuf.Backspace()
		m.searchQuery = m.editBuf.Value()
		m.applyFilters()
	default:
		if len(key) == 1 {
			m.editBuf.Insert(rune(key[0]))
			m.searchQuery = m.editBuf.Value()
			m.applyFilters()
		}
	}
	return m, nil
}

// handleMenuKey handles j/k/enter/esc in the state picker menu (ModeStateMenu).
func (m Model) handleMenuKey(key string) (tea.Model, tea.Cmd) {
	switch key {
	case "j", "down":
		if m.menuCursor < len(m.menuOptions)-1 {
			m.menuCursor++
		}
	case "k", "up":
		if m.menuCursor > 0 {
			m.menuCursor--
		}
	case "enter":
		if m.viewingIdea != nil && m.menuCursor < len(m.menuOptions) {
			selected := m.menuOptions[m.menuCursor]
			switch m.menuField {
			case FieldState:
				m.viewingIdea.State = selected
				if err := persistIdeaFrontmatter(m.viewingIdea); err != nil {
					m.statusMsg = "error saving state: " + err.Error()
				} else {
					if fresh, err := refreshIdea(m.viewingIdea); err == nil {
						m.viewingIdea = fresh
					}
					_ = m.loadIdeas()
				}
			case FieldPurpose:
				if selected == "" {
					m.viewingIdea.PurposeID = ""
					m.viewingIdea.PurposeName = ""
				} else {
					m.viewingIdea.PurposeID = selected
					m.viewingIdea.PurposeName = m.purposeNameFor(selected)
				}
				if err := persistIdeaFrontmatter(m.viewingIdea); err != nil {
					m.statusMsg = "error saving purpose: " + err.Error()
				} else {
					if fresh, err := refreshIdea(m.viewingIdea); err == nil {
						m.viewingIdea = fresh
					}
					_ = m.loadIdeas()
				}
			}
		}
		m.mode = ModeIdeaView
	case "esc", "q":
		m.mode = ModeIdeaView
	}
	return m, nil
}

func (m Model) handleComplianceKey(key string) (tea.Model, tea.Cmd) {
	switch key {
	case "j", "down":
		if m.complianceCursor < len(m.complianceOptions)-1 {
			m.complianceCursor++
		}
	case "k", "up":
		if m.complianceCursor > 0 {
			m.complianceCursor--
		}
	case "enter":
		if m.viewingIdea != nil && m.complianceCursor < len(m.complianceOptions) {
			m.viewingIdea.State = m.complianceOptions[m.complianceCursor]
			if err := persistIdeaFrontmatter(m.viewingIdea); err != nil {
				m.statusMsg = "error saving state: " + err.Error()
			} else {
				if fresh, err := refreshIdea(m.viewingIdea); err == nil {
					m.viewingIdea = fresh
				}
				_ = m.loadIdeas()
			}
		}
		m.mode = ModeIdeaView
	case "esc", "q":
		m.mode = ModeNormal
		m.viewingIdea = nil
	}
	return m, nil
}

func (m Model) handleCreateKey(key string) (tea.Model, tea.Cmd) {
	switch key {
	case "enter":
		if m.createField == 0 {
			m.createTitle = m.editBuf.Value()
			if m.createTitle == "" {
				return m, nil // don't advance without a title
			}
			m.createField = 1
			m.editBuf.SetValue(m.createKind)
		} else if m.createField == 1 {
			val := m.editBuf.Value()
			if val != "" {
				m.createKind = val
			}
			m.createField = 2
			m.editBuf.SetValue("")
		} else {
			m.createTags = m.editBuf.Value()
			// Parse space-separated tags
			var tags []string
			for _, t := range strings.Fields(m.createTags) {
				if t != "" {
					tags = append(tags, t)
				}
			}
			created, err := createIdea(m.cfg, m.createTitle, m.createKind, tags)
			if err != nil {
				m.statusMsg = "error creating idea: " + err.Error()
			} else {
				m.statusMsg = "Created: " + created.Title
				_ = m.loadIdeas()
			}
			m.mode = ModeNormal
			m.editBuf.Clear()
		}
	case "tab":
		// Cycle through available kinds when on kind field
		if m.createField == 1 {
			kinds := m.kindsConfig.AllKinds()
			current := m.editBuf.Value()
			for i, k := range kinds {
				if k == current {
					if i+1 < len(kinds) {
						m.editBuf.SetValue(kinds[i+1])
					} else {
						m.editBuf.SetValue(kinds[0])
					}
					break
				}
			}
		}
	case "esc":
		m.mode = ModeNormal
		m.editBuf.Clear()
	case "backspace":
		m.editBuf.Backspace()
	default:
		if len(key) == 1 {
			m.editBuf.Insert(rune(key[0]))
		}
	}
	return m, nil
}

func (m Model) handleLogEntryKey(key string) (tea.Model, tea.Cmd) {
	switch key {
	case "enter":
		entry := m.logBuf.Value()
		if entry != "" && m.viewingIdea != nil {
			if err := persistLogEntry(m.viewingIdea, entry); err != nil {
				m.statusMsg = "error saving log: " + err.Error()
			} else {
				if fresh, err := refreshIdea(m.viewingIdea); err == nil {
					m.viewingIdea = fresh
				}
				_ = m.loadIdeas()
			}
		}
		m.logBuf.Clear()
		m.mode = ModeIdeaView
	case "esc":
		m.logBuf.Clear()
		m.mode = ModeIdeaView
	case "backspace":
		m.logBuf.Backspace()
	default:
		if len(key) == 1 {
			m.logBuf.Insert(rune(key[0]))
		}
	}
	return m, nil
}

func (m Model) handleTagsEditKey(key string) (tea.Model, tea.Cmd) {
	switch key {
	case "enter":
		if m.viewingIdea != nil {
			// Parse space-separated tags and update frontmatter
			var newTags []string
			for _, t := range strings.Fields(m.editBuf.Value()) {
				if t != "" {
					newTags = append(newTags, t)
				}
			}
			m.viewingIdea.Tags = newTags
			if err := persistIdeaFrontmatter(m.viewingIdea); err != nil {
				m.statusMsg = "error saving tags: " + err.Error()
			} else {
				if fresh, err := refreshIdea(m.viewingIdea); err == nil {
					m.viewingIdea = fresh
				}
				_ = m.loadIdeas()
			}
		}
		m.mode = ModeIdeaView
	case "esc":
		m.mode = ModeIdeaView
	case "backspace":
		m.editBuf.Backspace()
	default:
		if len(key) == 1 {
			m.editBuf.Insert(rune(key[0]))
		}
	}
	return m, nil
}

// openInEditor opens the file at path in $EDITOR, suspending the TUI.
func openInEditor(filePath string) tea.Cmd {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vi"
	}
	c := exec.Command(editor, filePath)
	return tea.ExecProcess(c, nil)
}
