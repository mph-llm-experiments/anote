package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	acoreui "github.com/mph-llm-experiments/acore/tui"
)

// View implements tea.Model by dispatching to the appropriate renderer.
func (m Model) View() string {
	switch m.mode {
	case ModeIdeaView, ModeStateMenu, ModeTagsEdit, ModeLogEntry, ModeCompliancePrompt, ModeConfirmDelete:
		return m.viewIdeaDetail()
	case ModeHelp:
		return m.viewHelp()
	case ModeCreate, ModeCreateTags:
		return m.viewCreate()
	default:
		return m.viewList()
	}
}

func (m Model) viewList() string {
	var sb strings.Builder

	sb.WriteString(m.renderHeader())
	sb.WriteString("\n")

	visibleHeight := m.height - HeaderFooterHeight
	if visibleHeight < MinVisibleHeight {
		visibleHeight = MinVisibleHeight
	}

	cursor := m.nav.Cursor()
	start := 0
	if cursor >= visibleHeight {
		start = cursor - visibleHeight + 1
	}

	for i := start; i < len(m.filtered) && i < start+visibleHeight; i++ {
		idea := m.filtered[i]
		selected := i == cursor
		row := FormatListRow(
			idea.Title,
			idea.Kind,
			idea.State,
			idea.Maturity,
			idea.PurposeName,
			selected,
			m.width,
		)
		sb.WriteString(row)
		sb.WriteString("\n")
	}

	if len(m.filtered) == 0 {
		sb.WriteString(acoreui.MutedStyle.Render("  no notes"))
		sb.WriteString("\n")
	}

	sb.WriteString(m.renderFooter())
	return sb.String()
}

func (m Model) renderHeader() string {
	count := fmt.Sprintf("%d notes", len(m.filtered))
	if len(m.filtered) != len(m.ideas) {
		count = fmt.Sprintf("%d / %d", len(m.filtered), len(m.ideas))
	}
	filters := m.activeFilterSummary()
	title := "anote"
	if filters != "" {
		title = fmt.Sprintf("anote  %s", filters)
	}
	right := acoreui.MutedStyle.Render(count)
	left := acoreui.HeaderStyle.Render(title)
	space := m.width - lipgloss.Width(left) - lipgloss.Width(right)
	if space < 1 {
		space = 1
	}
	return left + strings.Repeat(" ", space) + right
}

func (m Model) activeFilterSummary() string {
	parts := []string{}
	if m.kindFilter != "" {
		parts = append(parts, "kind:"+m.kindFilter)
	}
	if m.stateFilter != "" {
		parts = append(parts, "state:"+m.stateFilter)
	}
	if m.purposeFilter != "" {
		purposeName := m.purposeNameFor(m.purposeFilter)
		if purposeName == "" {
			purposeName = m.purposeFilter
		}
		parts = append(parts, "purpose:"+purposeName)
	}
	if m.searchQuery != "" {
		parts = append(parts, "/"+m.searchQuery)
	}
	if len(parts) == 0 {
		return ""
	}
	return acoreui.WarningStyle.Render("[" + strings.Join(parts, " ") + "]")
}

func (m Model) renderFooter() string {
	hints := "c:new  /:search  enter:open  K:kind  P:purpose  ?:help  q:quit"
	if m.statusMsg != "" {
		hints = m.statusMsg
	}
	return acoreui.MutedStyle.Render(hints)
}

func (m Model) viewHelp() string {
	bindings := []acoreui.KeyBinding{
		{Key: "j / k", Desc: "move up/down"},
		{Key: "g / G", Desc: "top / bottom"},
		{Key: "ctrl+d/u", Desc: "page down/up"},
		{Key: "enter", Desc: "open idea"},
		{Key: "c", Desc: "create new idea"},
		{Key: "/", Desc: "search by title"},
		{Key: "K", Desc: "filter by kind (list) / change kind (idea view)"},
		{Key: "P", Desc: "filter by purpose"},
		{Key: "?", Desc: "this help"},
		{Key: "q / esc", Desc: "back / quit"},
		{Key: "— idea view —", Desc: ""},
		{Key: "s", Desc: "change state (kind-aware)"},
		{Key: "p", Desc: "change purpose"},
		{Key: "m", Desc: "change maturity"},
		{Key: "K", Desc: "change kind"},
		{Key: "t", Desc: "edit tags"},
		{Key: "l", Desc: "add log entry"},
		{Key: "r", Desc: "rename title"},
		{Key: "x", Desc: "delete note"},
		{Key: "E", Desc: "open in $EDITOR"},
	}
	return acoreui.RenderHelp(bindings, m.width)
}

