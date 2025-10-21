package models

import "midi-viewer/internal/midi"

// AppState represents the current state of the application
type AppState int

// Application state constants
const (
	StateDeviceSelection AppState = iota
	StateEventViewer
	StateFilterModal
)

// Filter represents the current filtering configuration
type Filter struct {
	HiddenChannels     map[uint8]bool   // channels to hide (empty = show all)
	HiddenMessageTypes map[string]bool  // message types to hide (empty = show all)
	HiddenColumns      map[string]bool  // columns to hide (empty = show all)
	ShowMusicalNotes   bool             // show musical note names (C4) instead of numbers (60)
}

// NewFilter creates a new empty filter (showing all events)
func NewFilter() Filter {
	return Filter{
		HiddenChannels:     make(map[uint8]bool),
		HiddenMessageTypes: make(map[string]bool),
		HiddenColumns:      make(map[string]bool),
		ShowMusicalNotes:   true, // Default to musical notes
	}
}

// ShouldShow returns true if an event should be displayed given the current filter
func (f Filter) ShouldShow(event midi.Event) bool {
	// Hide if channel is in hidden list
	if f.HiddenChannels[event.Channel] {
		return false
	}

	// Hide if message type is in hidden list
	if f.HiddenMessageTypes[event.MessageType] {
		return false
	}

	return true
}

// ToggleChannel toggles a channel's visibility
func (f *Filter) ToggleChannel(ch uint8) {
	if f.HiddenChannels[ch] {
		delete(f.HiddenChannels, ch)
	} else {
		f.HiddenChannels[ch] = true
	}
}

// ToggleMessageType toggles a message type's visibility
func (f *Filter) ToggleMessageType(msgType string) {
	if f.HiddenMessageTypes[msgType] {
		delete(f.HiddenMessageTypes, msgType)
	} else {
		f.HiddenMessageTypes[msgType] = true
	}
}

// IsChannelVisible returns true if a channel is currently visible
func (f Filter) IsChannelVisible(ch uint8) bool {
	return !f.HiddenChannels[ch]
}

// IsMessageTypeVisible returns true if a message type is currently visible
func (f Filter) IsMessageTypeVisible(msgType string) bool {
	return !f.HiddenMessageTypes[msgType]
}

// ToggleColumn toggles a column's visibility
func (f *Filter) ToggleColumn(col string) {
	if f.HiddenColumns[col] {
		delete(f.HiddenColumns, col)
	} else {
		f.HiddenColumns[col] = true
	}
}

// IsColumnVisible returns true if a column is currently visible
func (f Filter) IsColumnVisible(col string) bool {
	return !f.HiddenColumns[col]
}
