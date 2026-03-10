package tui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	acoreui "github.com/mph-llm-experiments/acore/tui"
	"github.com/mph-llm-experiments/anote/internal/denote"
)

// KindSymbol returns the single-character symbol for a kind.
func KindSymbol(kind string) string {
	switch kind {
	case denote.KindAspiration:
		return SymbolAspiration
	case denote.KindBelief:
		return SymbolBelief
	case denote.KindPlan:
		return SymbolPlan
	case denote.KindNote:
		return SymbolNote
	case denote.KindFact:
		return SymbolFact
	case denote.KindPurpose:
		return SymbolPurpose
	}
	return "?"
}

// StateSymbol returns the single-character symbol for a state.
func StateSymbol(state string) string {
	switch state {
	case denote.StateSeed:
		return SymbolSeed
	case denote.StateDraft:
		return SymbolDraft
	case denote.StateActive:
		return SymbolActive
	case denote.StateIterating:
		return SymbolIterating
	case denote.StateImplemented:
		return SymbolImplemented
	case denote.StateArchived:
		return SymbolArchived
	case denote.StateRejected:
		return SymbolRejected
	case denote.StateDropped:
		return SymbolDropped
	}
	return "?"
}

// stateColor returns a lipgloss color string for a canonical state.
func stateColor(state string) string {
	switch state {
	case denote.StateActive, denote.StateIterating:
		return acoreui.ColorSuccess
	case denote.StateImplemented:
		return acoreui.ColorMuted
	case denote.StateRejected, denote.StateDropped:
		return acoreui.ColorError
	case denote.StateSeed, denote.StateDraft:
		return acoreui.ColorWarning
	case denote.StateArchived:
		return acoreui.ColorMuted
	}
	return acoreui.ColorBody
}

// FormatListRow renders a single idea as a list row string.
// selected applies the selected highlight color to the title.
func FormatListRow(title, kind, state, maturity, purposeName string, selected bool, width int) string {
	kindSym := lipgloss.NewStyle().
		Foreground(lipgloss.Color(acoreui.ColorMuted)).
		Render(KindSymbol(kind))

	stateSym := lipgloss.NewStyle().
		Foreground(lipgloss.Color(stateColor(state))).
		Render(StateSymbol(state))

	purposeStr := ""
	if purposeName != "" {
		purposeStr = lipgloss.NewStyle().
			Foreground(lipgloss.Color(acoreui.ColorMuted)).
			Render(fmt.Sprintf(" [%s]", purposeName))
	}

	maturityStr := ""
	if maturity != "" {
		maturityStr = lipgloss.NewStyle().
			Foreground(lipgloss.Color(acoreui.ColorMuted)).
			Render(fmt.Sprintf(" %s", maturity))
	}

	titleStyle := acoreui.BodyStyle
	if selected {
		titleStyle = acoreui.SelectedStyle
	}

	return fmt.Sprintf("%s %s %s%s%s",
		kindSym, stateSym, titleStyle.Render(title), purposeStr, maturityStr,
	)
}
