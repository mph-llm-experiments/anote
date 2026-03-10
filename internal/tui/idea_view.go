package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	acoreui "github.com/mph-llm-experiments/acore/tui"
	"github.com/mph-llm-experiments/anote/internal/denote"
)

// viewIdeaDetail renders the full idea detail view and any active overlay modes.
func (m Model) viewIdeaDetail() string {
	if m.viewingIdea == nil {
		return "No idea selected\n\nq: back"
	}
	idea := m.viewingIdea
	var sb strings.Builder

	// Title (highlighted)
	sb.WriteString(acoreui.SelectedStyle.Render(idea.Title))
	sb.WriteString("\n\n")

	// Metadata fields
	sb.WriteString(m.renderMetaField("kind", KindSymbol(idea.Kind)+" "+idea.Kind, FieldKind))
	sb.WriteString(m.renderMetaField("state",
		StateSymbol(idea.State)+" "+denote.DisplayState(idea.State, idea.Kind), FieldState))
	if idea.Maturity != "" {
		sb.WriteString(m.renderMetaField("maturity", idea.Maturity, FieldMaturity))
	} else {
		sb.WriteString(m.renderMetaField("maturity", "—", FieldMaturity))
	}
	if idea.PurposeName != "" {
		sb.WriteString(m.renderMetaField("purpose", idea.PurposeName, FieldPurpose))
	} else if idea.PurposeID != "" {
		sb.WriteString(m.renderMetaField("purpose", idea.PurposeID, FieldPurpose))
	} else {
		sb.WriteString(m.renderMetaField("purpose", "—", FieldPurpose))
	}
	if len(idea.Tags) > 0 {
		sb.WriteString(m.renderMetaField("tags", strings.Join(idea.Tags, " "), FieldTags))
	} else {
		sb.WriteString(m.renderMetaField("tags", "—", FieldTags))
	}
	sb.WriteString("\n")

	// Body / content
	if idea.Content != "" {
		sb.WriteString(acoreui.MutedStyle.Render("───"))
		sb.WriteString("\n")
		sb.WriteString(idea.Content)
		sb.WriteString("\n")
	}

	// Overlay modes that appear on top of the detail view
	switch m.mode {
	case ModeStateMenu:
		sb.WriteString("\n")
		sb.WriteString(m.viewStateMenu())
	case ModeCompliancePrompt:
		sb.WriteString("\n")
		sb.WriteString(m.viewCompliancePrompt())
	case ModeLogEntry:
		sb.WriteString("\n")
		sb.WriteString(acoreui.HeaderStyle.Render("Add log entry:"))
		sb.WriteString("\n")
		sb.WriteString(m.logBuf.Render())
		sb.WriteString("\n")
		sb.WriteString(acoreui.MutedStyle.Render("enter: save  esc: cancel"))
	case ModeTagsEdit:
		sb.WriteString("\n")
		sb.WriteString(acoreui.HeaderStyle.Render("Edit tags (space-separated):"))
		sb.WriteString("\n")
		sb.WriteString(m.editBuf.Render())
		sb.WriteString("\n")
		sb.WriteString(acoreui.MutedStyle.Render("enter: save  esc: cancel"))
	}

	// Footer
	sb.WriteString("\n")
	sb.WriteString(acoreui.MutedStyle.Render(
		"s:state  p:purpose  m:maturity  K:kind  t:tags  l:log  r:rename  x:delete  E:editor  q:back",
	))

	return sb.String()
}

// renderMetaField renders a label: value metadata line.
// editingField comparison is for future inline edit highlighting.
func (m Model) renderMetaField(label, value, _ string) string {
	labelStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(acoreui.ColorHeader))
	return fmt.Sprintf("%s %s\n",
		labelStyle.Render(fmt.Sprintf("%-10s:", label)),
		acoreui.BodyStyle.Render(value),
	)
}

// viewStateMenu renders the kind-aware state picker overlay.
func (m Model) viewStateMenu() string {
	if m.viewingIdea == nil {
		return ""
	}
	var sb strings.Builder
	sb.WriteString(acoreui.HeaderStyle.Render("Choose state:"))
	sb.WriteString("\n")
	for i, s := range m.menuOptions {
		display := denote.DisplayState(s, m.viewingIdea.Kind)
		if i == m.menuCursor {
			sb.WriteString(acoreui.SelectedStyle.Render("  → " + display))
		} else {
			sb.WriteString(acoreui.MutedStyle.Render("    " + display))
		}
		sb.WriteString("\n")
	}
	sb.WriteString(acoreui.MutedStyle.Render("j/k: move  enter: select  esc: cancel"))
	return sb.String()
}

// viewCompliancePrompt renders the compliance fix overlay.
func (m Model) viewCompliancePrompt() string {
	if m.viewingIdea == nil {
		return ""
	}
	var sb strings.Builder
	sb.WriteString(acoreui.WarningStyle.Render(fmt.Sprintf(
		"State %q is not valid for kind %q. Choose a valid state:",
		m.viewingIdea.State, m.viewingIdea.Kind,
	)))
	sb.WriteString("\n")
	for i, s := range m.complianceOptions {
		display := denote.DisplayState(s, m.viewingIdea.Kind)
		if i == m.complianceCursor {
			sb.WriteString(acoreui.SelectedStyle.Render("  → " + display))
		} else {
			sb.WriteString(acoreui.MutedStyle.Render("    " + display))
		}
		sb.WriteString("\n")
	}
	sb.WriteString(acoreui.MutedStyle.Render("j/k: move  enter: fix  esc: back to list"))
	return sb.String()
}

// viewCreate renders the new idea creation form.
// (Stub for Task 10 — create_view.go will provide the full implementation.)
func (m Model) viewCreate() string {
	fieldLabel := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(acoreui.ColorHeader))

	switch m.createField {
	case 0:
		return acoreui.HeaderStyle.Render("New idea") + "\n\n" +
			fieldLabel.Render("Title: ") + m.editBuf.Render() + "\n" +
			acoreui.MutedStyle.Render("enter: next  esc: cancel")
	case 1:
		return acoreui.HeaderStyle.Render("New idea") + "\n\n" +
			fieldLabel.Render("Title: ") + acoreui.BodyStyle.Render(m.createTitle) + "\n" +
			fieldLabel.Render("Kind:  ") + m.editBuf.Render() + "\n" +
			acoreui.MutedStyle.Render("tab: cycle  enter: next  esc: cancel")
	default:
		return acoreui.HeaderStyle.Render("New idea") + "\n\n" +
			fieldLabel.Render("Title: ") + acoreui.BodyStyle.Render(m.createTitle) + "\n" +
			fieldLabel.Render("Kind:  ") + acoreui.BodyStyle.Render(m.createKind) + "\n" +
			fieldLabel.Render("Tags:  ") + m.editBuf.Render() + "\n" +
			acoreui.MutedStyle.Render("space-separated  enter: save  esc: cancel")
	}
}
