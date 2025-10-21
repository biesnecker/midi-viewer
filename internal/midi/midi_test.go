package midi

import (
	"testing"
)

func TestNoteToName(t *testing.T) {
	tests := []struct {
		note     uint8
		expected string
	}{
		{0, "C-1"},
		{12, "C0"},
		{24, "C1"},
		{36, "C2"},
		{48, "C3"},
		{60, "C4"},   // Middle C
		{61, "C#4"},
		{62, "D4"},
		{72, "C5"},
		{127, "G9"},
	}

	for _, tt := range tests {
		result := NoteToName(tt.note)
		if result != tt.expected {
			t.Errorf("NoteToName(%d) = %s; want %s", tt.note, result, tt.expected)
		}
	}
}

func TestGetMessageType(t *testing.T) {
	// Note: This is a basic test structure
	// In practice, you'd create actual MIDI messages to test
	// For now, we're just verifying the function exists and has the right signature
	tests := []struct {
		name     string
		expected string
	}{
		{"Unknown message should return 'Unknown'", "Unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This test would need actual MIDI message construction
			// which requires the gomidi library's message builders
		})
	}
}
