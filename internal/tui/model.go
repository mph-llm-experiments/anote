package tui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	acoreui "github.com/mph-llm-experiments/acore/tui"
	"github.com/mph-llm-experiments/anote/internal/config"
	"github.com/mph-llm-experiments/anote/internal/denote"
)

// Model is the root Bubble Tea model for the anote TUI.
type Model struct {
	cfg         *config.Config
	kindsConfig *denote.KindsConfig

	// List state
	ideas    []denote.Idea
	filtered []denote.Idea
	nav      *acoreui.NavigationHandler

	// UI dimensions
	width  int
	height int

	// Current mode
	mode Mode

	// Filters
	kindFilter    string
	stateFilter   string
	purposeFilter string
	searchQuery   string

	// Sort
	sortBy      string
	reverseSort bool

	// Idea detail view
	viewingIdea  *denote.Idea
	editingField string
	editBuf      *acoreui.EditBuffer

	// State menu (kind-aware) and purpose picker
	menuOptions []string
	menuCursor  int
	menuField   string // which field the menu is editing

	// Create mode
	createTitle string
	createKind  string
	createTags  string
	createField int

	// Purposes (cached list of purpose-kind ideas)
	purposes []denote.Idea

	// Compliance prompt
	complianceOptions []string
	complianceCursor  int

	// Log entry
	logBuf *acoreui.EditBuffer

	// Status message (shown in footer)
	statusMsg string
}

// NewModel initializes the TUI model by loading ideas and config.
func NewModel(cfg *config.Config) (*Model, error) {
	kindsConfig, err := denote.LoadKindsConfig(cfg.IdeasDirectory)
	if err != nil {
		return nil, err
	}

	m := &Model{
		cfg:         cfg,
		kindsConfig: kindsConfig,
		nav:         acoreui.NewNavigationHandler(0, false),
		editBuf:     acoreui.NewEditBuffer(""),
		logBuf:      acoreui.NewEditBuffer(""),
		sortBy:      "modified",
	}

	if err := m.loadIdeas(); err != nil {
		return nil, err
	}

	return m, nil
}

// loadIdeas scans the data directory and populates m.ideas and m.filtered.
func (m *Model) loadIdeas() error {
	scanner := denote.NewScanner(m.cfg.IdeasDirectory)
	ptrs, err := scanner.FindIdeas()
	if err != nil {
		return err
	}

	m.ideas = make([]denote.Idea, 0, len(ptrs))
	for _, p := range ptrs {
		m.ideas = append(m.ideas, *p)
	}

	// Cache purposes (ideas where Kind == KindPurpose)
	m.purposes = make([]denote.Idea, 0)
	for _, idea := range m.ideas {
		if idea.Kind == denote.KindPurpose {
			m.purposes = append(m.purposes, idea)
		}
	}

	m.applyFilters()
	return nil
}

// applyFilters filters m.ideas into m.filtered based on active filters.
func (m *Model) applyFilters() {
	filtered := make([]denote.Idea, 0, len(m.ideas))
	for _, idea := range m.ideas {
		if m.kindFilter != "" && idea.Kind != m.kindFilter {
			continue
		}
		if m.stateFilter != "" && idea.State != m.stateFilter {
			continue
		}
		if m.purposeFilter != "" && idea.PurposeID != m.purposeFilter {
			continue
		}
		if m.searchQuery != "" && !containsFold(idea.Title, m.searchQuery) {
			continue
		}
		filtered = append(filtered, idea)
	}
	m.filtered = filtered
	m.nav.SetMax(len(m.filtered))
}

func containsFold(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}

// purposeNameFor looks up the purpose name by ID from the cached purposes list.
func (m *Model) purposeNameFor(purposeID string) string {
	if purposeID == "" {
		return ""
	}
	for _, p := range m.purposes {
		if p.ID == purposeID {
			return p.Title
		}
	}
	return ""
}

// Init implements tea.Model.
func (m Model) Init() tea.Cmd { return nil }

// Update implements tea.Model — stubbed here, filled in by keys.go.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

// View implements tea.Model — stubbed here, filled in by views.go.
func (m Model) View() string {
	return "anote (loading...)"
}
