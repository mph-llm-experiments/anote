package denote

import (
	"testing"
)

func TestIsValidState(t *testing.T) {
	valid := []string{"seed", "draft", "active", "iterating", "implemented", "archived", "rejected", "dropped"}
	for _, s := range valid {
		if !IsValidState(s) {
			t.Errorf("expected %q to be valid state", s)
		}
	}

	invalid := []string{"", "open", "done", "pending", "SEED", "Active"}
	for _, s := range invalid {
		if IsValidState(s) {
			t.Errorf("expected %q to be invalid state", s)
		}
	}
}

func TestIsValidMaturity(t *testing.T) {
	valid := []string{"crawl", "walk", "run"}
	for _, m := range valid {
		if !IsValidMaturity(m) {
			t.Errorf("expected %q to be valid maturity", m)
		}
	}

	invalid := []string{"", "sprint", "CRAWL", "Walk"}
	for _, m := range invalid {
		if IsValidMaturity(m) {
			t.Errorf("expected %q to be invalid maturity", m)
		}
	}
}

func TestValidateStateTransition(t *testing.T) {
	tests := []struct {
		from    string
		to      string
		wantErr bool
	}{
		// Valid transitions
		{"seed", "draft", false},
		{"draft", "active", false},
		{"active", "iterating", false},
		{"active", "implemented", false},
		{"active", "archived", false},
		{"active", "rejected", false},
		{"active", "dropped", false},
		{"iterating", "active", false},
		{"iterating", "implemented", false},
		{"iterating", "archived", false},
		{"iterating", "rejected", false},
		{"iterating", "dropped", false},
		{"archived", "active", false},

		// Invalid transitions
		{"seed", "active", true},      // must go through draft
		{"seed", "implemented", true},  // can't skip ahead
		{"draft", "iterating", true},   // must go through active
		{"draft", "seed", true},        // can't go backward
		{"implemented", "active", true}, // terminal state
		{"rejected", "active", true},    // terminal state
		{"dropped", "active", true},     // terminal state

		// Invalid state values
		{"bogus", "active", true},
		{"active", "bogus", true},
	}

	for _, tt := range tests {
		err := ValidateStateTransition(tt.from, tt.to)
		if (err != nil) != tt.wantErr {
			t.Errorf("ValidateStateTransition(%q, %q): got err=%v, wantErr=%v", tt.from, tt.to, err, tt.wantErr)
		}
	}
}

func TestValidateIdea(t *testing.T) {
	// Valid: rejected with reason
	idea := &Idea{}
	idea.State = StateRejected
	idea.RejectedReason = "Too expensive"
	if err := ValidateIdea(idea); err != nil {
		t.Errorf("expected no error for rejected with reason, got: %v", err)
	}

	// Invalid: rejected without reason
	idea2 := &Idea{}
	idea2.State = StateRejected
	idea2.RejectedReason = ""
	if err := ValidateIdea(idea2); err == nil {
		t.Error("expected error for rejected without reason")
	}

	// Invalid: rejected with whitespace-only reason
	idea3 := &Idea{}
	idea3.State = StateRejected
	idea3.RejectedReason = "   "
	if err := ValidateIdea(idea3); err == nil {
		t.Error("expected error for rejected with whitespace-only reason")
	}

	// Valid: active with no maturity (maturity is optional)
	idea4 := &Idea{}
	idea4.State = StateActive
	if err := ValidateIdea(idea4); err != nil {
		t.Errorf("expected no error for active without maturity, got: %v", err)
	}

	// Invalid state
	idea5 := &Idea{}
	idea5.State = "bogus"
	if err := ValidateIdea(idea5); err == nil {
		t.Error("expected error for invalid state")
	}

	// Invalid maturity
	idea6 := &Idea{}
	idea6.State = StateActive
	idea6.Maturity = "sprint"
	if err := ValidateIdea(idea6); err == nil {
		t.Error("expected error for invalid maturity")
	}
}

func TestHasTag(t *testing.T) {
	f := &File{Tags: []string{"idea", "coaching", "leadership"}}

	if !f.HasTag("idea") {
		t.Error("expected HasTag(idea) to be true")
	}
	if !f.HasTag("coaching") {
		t.Error("expected HasTag(coaching) to be true")
	}
	if f.HasTag("missing") {
		t.Error("expected HasTag(missing) to be false")
	}
}

func TestIsIdea(t *testing.T) {
	f1 := &File{Tags: []string{"idea", "coaching"}}
	if !f1.IsIdea() {
		t.Error("expected IsIdea() to be true when idea tag present")
	}

	f2 := &File{Tags: []string{"task", "work"}}
	if f2.IsIdea() {
		t.Error("expected IsIdea() to be false when idea tag absent")
	}

	f3 := &File{Tags: []string{}}
	if f3.IsIdea() {
		t.Error("expected IsIdea() to be false for empty tags")
	}
}
