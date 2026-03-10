package tui

import acoreui "github.com/mph-llm-experiments/acore/tui"

// Re-export shared layout constants.
const (
	HeaderFooterHeight = acoreui.HeaderFooterHeight
	MinVisibleHeight   = acoreui.MinVisibleHeight
)

// Mode represents the current TUI state.
type Mode int

const (
	ModeNormal Mode = iota
	ModeSearch
	ModeHelp
	ModeIdeaView
	ModeCreate
	ModeCreateTags
	ModeKindFilter
	ModeStateFilter
	ModePurposeFilter
	ModeStateMenu
	ModeConfirmDelete
	ModeLogEntry
	ModeTagsEdit
	ModeCompliancePrompt
	ModeSortMenu
	ModeFilterMenu
)

// anote-specific status symbols.
const (
	SymbolSeed        = "○"
	SymbolActive      = "●"
	SymbolDraft       = "◐"
	SymbolIterating   = "↻"
	SymbolImplemented = "✓"
	SymbolArchived    = "⊘"
	SymbolRejected    = "⨯"
	SymbolDropped     = "⨯"
)

// Kind symbols for list display.
const (
	SymbolAspiration = "↑"
	SymbolBelief     = "◆"
	SymbolPlan       = "→"
	SymbolNote       = "·"
	SymbolFact       = "="
	SymbolPurpose    = "★"
)

// Editable field names.
const (
	FieldState    = "state"
	FieldPurpose  = "purpose"
	FieldMaturity = "maturity"
	FieldKind     = "kind"
	FieldTags     = "tags"
	FieldLog      = "log"
	FieldTitle    = "title"
)
