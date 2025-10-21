package models

import (
	"testing"
	"time"

	"midi-viewer/internal/midi"
)

func TestNewFilter(t *testing.T) {
	filter := NewFilter()

	if len(filter.HiddenChannels) != 0 {
		t.Errorf("NewFilter() HiddenChannels = %v; want empty map", filter.HiddenChannels)
	}

	if len(filter.HiddenMessageTypes) != 0 {
		t.Errorf("NewFilter() HiddenMessageTypes = %v; want empty map", filter.HiddenMessageTypes)
	}

	if len(filter.HiddenColumns) != 0 {
		t.Errorf("NewFilter() HiddenColumns = %v; want empty map", filter.HiddenColumns)
	}
}

func TestFilterToggleChannel(t *testing.T) {
	filter := NewFilter()

	// Initially visible
	if !filter.IsChannelVisible(5) {
		t.Error("Channel 5 should be visible by default")
	}

	// Toggle to hide
	filter.ToggleChannel(5)
	if filter.IsChannelVisible(5) {
		t.Error("ToggleChannel(5) should hide channel 5")
	}

	// Toggle to show
	filter.ToggleChannel(5)
	if !filter.IsChannelVisible(5) {
		t.Error("ToggleChannel(5) second time should show channel 5")
	}
}

func TestFilterToggleMessageType(t *testing.T) {
	filter := NewFilter()

	// Initially visible
	if !filter.IsMessageTypeVisible("Note On") {
		t.Error("Message type 'Note On' should be visible by default")
	}

	// Toggle to hide
	filter.ToggleMessageType("Note On")
	if filter.IsMessageTypeVisible("Note On") {
		t.Error("ToggleMessageType('Note On') should hide 'Note On'")
	}

	// Toggle to show
	filter.ToggleMessageType("Note On")
	if !filter.IsMessageTypeVisible("Note On") {
		t.Error("ToggleMessageType('Note On') second time should show 'Note On'")
	}
}

func TestFilterShouldShow_NoFilters(t *testing.T) {
	filter := NewFilter()
	event := midi.Event{
		Timestamp:   time.Now(),
		Channel:     1,
		MessageType: "Note On",
		Data:        "Note: 60, Velocity: 100",
	}

	if !filter.ShouldShow(event) {
		t.Error("ShouldShow should return true when no filters are active")
	}
}

func TestFilterShouldShow_ChannelFilter(t *testing.T) {
	filter := NewFilter()
	filter.ToggleChannel(1) // Hide channel 1
	filter.ToggleChannel(2) // Hide channel 2

	event1 := midi.Event{
		Channel:     1,
		MessageType: "Note On",
	}

	event2 := midi.Event{
		Channel:     3,
		MessageType: "Note On",
	}

	if filter.ShouldShow(event1) {
		t.Error("ShouldShow should return false for hidden channel 1")
	}

	if !filter.ShouldShow(event2) {
		t.Error("ShouldShow should return true for visible channel 3")
	}
}

func TestFilterShouldShow_MessageTypeFilter(t *testing.T) {
	filter := NewFilter()
	filter.ToggleMessageType("Note On")  // Hide Note On
	filter.ToggleMessageType("Note Off") // Hide Note Off

	event1 := midi.Event{
		MessageType: "Note On",
		Channel:     1,
	}

	event2 := midi.Event{
		MessageType: "CC",
		Channel:     1,
	}

	if filter.ShouldShow(event1) {
		t.Error("ShouldShow should return false for hidden 'Note On'")
	}

	if !filter.ShouldShow(event2) {
		t.Error("ShouldShow should return true for visible 'CC'")
	}
}

func TestFilterShouldShow_CombinedFilters(t *testing.T) {
	filter := NewFilter()
	filter.ToggleChannel(1)           // Hide channel 1
	filter.ToggleMessageType("Note On") // Hide Note On

	event1 := midi.Event{
		Channel:     1,
		MessageType: "Note On",
	}

	event2 := midi.Event{
		Channel:     1,
		MessageType: "CC",
	}

	event3 := midi.Event{
		Channel:     2,
		MessageType: "Note On",
	}

	if filter.ShouldShow(event1) {
		t.Error("ShouldShow should return false for hidden channel 1 + hidden Note On")
	}

	if filter.ShouldShow(event2) {
		t.Error("ShouldShow should return false for hidden channel 1 + visible CC")
	}

	if filter.ShouldShow(event3) {
		t.Error("ShouldShow should return false for visible channel 2 + hidden Note On")
	}
}

func TestFilterIsChannelVisible(t *testing.T) {
	filter := NewFilter()
	filter.ToggleChannel(5) // Hide channel 5

	if filter.IsChannelVisible(5) {
		t.Error("IsChannelVisible(5) should return false for hidden channel")
	}

	if !filter.IsChannelVisible(3) {
		t.Error("IsChannelVisible(3) should return true for non-hidden channel")
	}
}

func TestFilterIsMessageTypeVisible(t *testing.T) {
	filter := NewFilter()
	filter.ToggleMessageType("Note On") // Hide Note On

	if filter.IsMessageTypeVisible("Note On") {
		t.Error("IsMessageTypeVisible('Note On') should return false for hidden type")
	}

	if !filter.IsMessageTypeVisible("CC") {
		t.Error("IsMessageTypeVisible('CC') should return true for non-hidden type")
	}
}

func TestFilterToggleColumn(t *testing.T) {
	filter := NewFilter()

	// Initially visible
	if !filter.IsColumnVisible("Note") {
		t.Error("Column 'Note' should be visible by default")
	}

	// Toggle to hide
	filter.ToggleColumn("Note")
	if filter.IsColumnVisible("Note") {
		t.Error("ToggleColumn('Note') should hide 'Note'")
	}

	// Toggle to show
	filter.ToggleColumn("Note")
	if !filter.IsColumnVisible("Note") {
		t.Error("ToggleColumn('Note') second time should show 'Note'")
	}
}

func TestFilterIsColumnVisible(t *testing.T) {
	filter := NewFilter()
	filter.ToggleColumn("Vel") // Hide Vel column

	if filter.IsColumnVisible("Vel") {
		t.Error("IsColumnVisible('Vel') should return false for hidden column")
	}

	if !filter.IsColumnVisible("Note") {
		t.Error("IsColumnVisible('Note') should return true for non-hidden column")
	}
}
