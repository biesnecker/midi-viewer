package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	gomidi "gitlab.com/gomidi/midi/v2"
	"midi-viewer/internal/midi"
	"midi-viewer/internal/models"
	"midi-viewer/internal/ui/theme"
)

type eventViewerKeyMap struct {
	Pause   key.Binding
	Options key.Binding
	Clear   key.Binding
	Back    key.Binding
	Quit    key.Binding
}

var eventViewerKeys = eventViewerKeyMap{
	Pause: key.NewBinding(
		key.WithKeys(" "),
		key.WithHelp("space", "pause/unpause"),
	),
	Options: key.NewBinding(
		key.WithKeys("o"),
		key.WithHelp("o", "options"),
	),
	Clear: key.NewBinding(
		key.WithKeys("c"),
		key.WithHelp("c", "clear events"),
	),
	Back: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back to devices"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

// EventViewer displays MIDI events in a scrolling list
type EventViewer struct {
	events       []midi.Event
	device       midi.Device
	theme        theme.Theme
	width        int
	height       int
	paused       bool
	filter       models.Filter
	maxEvents    int
	activeNotes  map[uint8]map[uint8]bool // channel -> note -> active
}

// NewEventViewer creates a new event viewer
func NewEventViewer(device midi.Device, t theme.Theme) EventViewer {
	return EventViewer{
		device:      device,
		theme:       t,
		events:      make([]midi.Event, 0),
		paused:      false,
		filter:      models.NewFilter(),
		maxEvents:   1000, // Keep last 1000 events
		activeNotes: make(map[uint8]map[uint8]bool),
	}
}

// GetFilter returns the current filter
func (e EventViewer) GetFilter() models.Filter {
	return e.filter
}

// Update handles messages
func (e EventViewer) Update(msg tea.Msg) (EventViewer, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		e.width = msg.Width
		e.height = msg.Height

	case MIDIEventMsg:
		if !e.paused {
			event := midi.ParseMessage(msg.Message)

			// Track note on/off for active notes display
			e.updateActiveNotes(event)

			if e.filter.ShouldShow(event) {
				e.events = append(e.events, event)
				// Keep only last maxEvents
				if len(e.events) > e.maxEvents {
					e.events = e.events[len(e.events)-e.maxEvents:]
				}
			}
		}

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, eventViewerKeys.Back):
			return e, func() tea.Msg {
				return BackToDeviceSelectionMsg{}
			}
		case key.Matches(msg, eventViewerKeys.Quit):
			return e, tea.Quit
		case key.Matches(msg, eventViewerKeys.Pause):
			e.paused = !e.paused
		case key.Matches(msg, eventViewerKeys.Clear):
			e.events = make([]midi.Event, 0)
		case key.Matches(msg, eventViewerKeys.Options):
			return e, func() tea.Msg {
				return OpenOptionsModalMsg{}
			}
		}

	case FilterUpdatedMsg:
		e.filter = msg.Filter
	}

	return e, nil
}

// View renders the event viewer
func (e EventViewer) View() string {
	headerStyle := lipgloss.NewStyle().
		Foreground(e.theme.Primary).
		Background(e.theme.Background).
		Bold(true).
		Padding(0, 1).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(e.theme.Border).
		BorderBottom(true)

	statusStyle := lipgloss.NewStyle().
		Foreground(e.theme.Muted).
		Background(e.theme.Background).
		Padding(0, 1)

	pausedStyle := lipgloss.NewStyle().
		Foreground(e.theme.Warning).
		Background(e.theme.Background).
		Bold(true).
		Padding(0, 1)

	helpStyle := lipgloss.NewStyle().
		Foreground(e.theme.Muted).
		Background(e.theme.Background).
		Padding(0, 1).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(e.theme.Border).
		BorderTop(true)

	var b strings.Builder

	// Header
	header := fmt.Sprintf("MIDI Monitor - %s", e.device.Name)
	if e.paused {
		status := pausedStyle.Render(" [PAUSED] ")
		header += status
	}
	if e.hasActiveFilters() {
		filterIndicator := statusStyle.Render(" [FILTERED] ")
		header += filterIndicator
	}
	b.WriteString(headerStyle.Width(e.width).Render(header))
	b.WriteString("\n")

	// Calculate column widths
	timeWidth := 12
	eventWidth := 16
	chanWidth := 5
	noteWidth := 8
	velWidth := 6
	ctrlWidth := 7
	valWidth := 6

	// Column header styles
	colHeaderStyle := lipgloss.NewStyle().
		Foreground(e.theme.Secondary).
		Bold(true).
		Underline(true)

	// Build column headers
	var headerRow strings.Builder
	headerRow.WriteString("  ") // Left padding
	if e.filter.IsColumnVisible("Time") {
		headerRow.WriteString(colHeaderStyle.Width(timeWidth).Render("Time"))
		headerRow.WriteString("  ")
	}
	if e.filter.IsColumnVisible("Chan") {
		headerRow.WriteString(colHeaderStyle.Width(chanWidth).Render("Chan"))
		headerRow.WriteString("  ")
	}
	if e.filter.IsColumnVisible("Event") {
		headerRow.WriteString(colHeaderStyle.Width(eventWidth).Render("Event"))
		headerRow.WriteString("  ")
	}
	if e.filter.IsColumnVisible("Note") {
		headerRow.WriteString(colHeaderStyle.Width(noteWidth).Render("Note"))
		headerRow.WriteString("  ")
	}
	if e.filter.IsColumnVisible("Vel") {
		headerRow.WriteString(colHeaderStyle.Width(velWidth).Render("Vel"))
		headerRow.WriteString("  ")
	}
	if e.filter.IsColumnVisible("Ctrl") {
		headerRow.WriteString(colHeaderStyle.Width(ctrlWidth).Render("Ctrl"))
		headerRow.WriteString("  ")
	}
	if e.filter.IsColumnVisible("Val") {
		headerRow.WriteString(colHeaderStyle.Width(valWidth).Render("Val"))
	}
	b.WriteString(headerRow.String())
	b.WriteString("\n")

	// Calculate available height for events (accounting for column header and active notes)
	activeNotesHeight := 3 // "Active Notes:" + notes line + blank line
	availableHeight := e.height - 5 - activeNotesHeight // header + column header + help + active notes + padding
	if availableHeight < 0 {
		availableHeight = 0
	}

	// Events (show most recent at top)
	eventCount := len(e.events)
	startIdx := 0
	if eventCount > availableHeight && availableHeight > 0 {
		startIdx = eventCount - availableHeight
	}

	visibleEvents := []midi.Event{}
	if availableHeight > 0 {
		visibleEvents = e.events[startIdx:]
	}

	// Column value styles
	timeColStyle := lipgloss.NewStyle().Foreground(e.theme.Muted).Width(timeWidth)
	eventColStyle := lipgloss.NewStyle().Foreground(e.theme.Secondary).Bold(true).Width(eventWidth)
	chanColStyle := lipgloss.NewStyle().Foreground(e.theme.Primary).Width(chanWidth)
	dataColStyle := lipgloss.NewStyle().Foreground(e.theme.Foreground)

	// Reverse the slice to show newest at top
	for i := len(visibleEvents) - 1; i >= 0; i-- {
		event := visibleEvents[i]
		var row strings.Builder
		row.WriteString("  ") // Left padding

		if e.filter.IsColumnVisible("Time") {
			row.WriteString(timeColStyle.Render(event.Timestamp.Format("15:04:05.000")))
			row.WriteString("  ")
		}

		if e.filter.IsColumnVisible("Chan") {
			chanVal := ""
			if event.MessageType != "Unknown" && event.MessageType != "SysEx" {
				chanVal = fmt.Sprintf("%d", event.Channel+1)
			}
			row.WriteString(chanColStyle.Render(chanVal))
			row.WriteString("  ")
		}

		if e.filter.IsColumnVisible("Event") {
			row.WriteString(eventColStyle.Render(event.MessageType))
			row.WriteString("  ")
		}

		// Extract note, velocity, controller, and value from event data
		note, vel, ctrl, val := e.parseEventData(event)

		if e.filter.IsColumnVisible("Note") {
			row.WriteString(dataColStyle.Width(noteWidth).Render(note))
			row.WriteString("  ")
		}

		if e.filter.IsColumnVisible("Vel") {
			row.WriteString(dataColStyle.Width(velWidth).Render(vel))
			row.WriteString("  ")
		}

		if e.filter.IsColumnVisible("Ctrl") {
			row.WriteString(dataColStyle.Width(ctrlWidth).Render(ctrl))
			row.WriteString("  ")
		}

		if e.filter.IsColumnVisible("Val") {
			row.WriteString(dataColStyle.Width(valWidth).Render(val))
		}

		b.WriteString(row.String())
		b.WriteString("\n")
	}

	// Pad remaining space
	for i := len(visibleEvents); i < availableHeight; i++ {
		b.WriteString("\n")
	}

	// Active notes section
	b.WriteString("\n")
	b.WriteString(e.renderActiveNotes())
	b.WriteString("\n")

	// Help
	helpText := "space: pause • o: options • c: clear • esc: devices • q: quit"
	b.WriteString(helpStyle.Width(e.width).Render(helpText))

	return b.String()
}

// renderActiveNotes renders the currently playing notes
func (e EventViewer) renderActiveNotes() string {
	labelStyle := lipgloss.NewStyle().
		Foreground(e.theme.Secondary).
		Bold(true)

	noteStyle := lipgloss.NewStyle().
		Foreground(e.theme.Success).
		Bold(true)

	var notes []string

	// Collect all active notes across all channels
	for ch := uint8(0); ch < 16; ch++ {
		if channelNotes, ok := e.activeNotes[ch]; ok && len(channelNotes) > 0 {
			for note := range channelNotes {
				var noteName string
				if e.filter.ShowMusicalNotes {
					noteName = midi.NoteToName(note)
				} else {
					noteName = fmt.Sprintf("%d", note)
				}
				notes = append(notes, fmt.Sprintf("Ch%d:%s", ch+1, noteName))
			}
		}
	}

	var result strings.Builder
	result.WriteString(labelStyle.Render("Active Notes: "))

	if len(notes) == 0 {
		result.WriteString(lipgloss.NewStyle().Foreground(e.theme.Muted).Render("(none)"))
	} else {
		result.WriteString(noteStyle.Render(strings.Join(notes, " ")))
	}

	return result.String()
}

func (e EventViewer) hasActiveFilters() bool {
	return len(e.filter.HiddenChannels) > 0 || len(e.filter.HiddenMessageTypes) > 0
}

// updateActiveNotes tracks which notes are currently playing
func (e *EventViewer) updateActiveNotes(event midi.Event) {
	var ch, key, vel uint8

	if event.MessageType == "Note On" {
		if event.Message.GetNoteOn(&ch, &key, &vel) {
			if vel > 0 {
				// Note on with velocity > 0 = note started
				if e.activeNotes[ch] == nil {
					e.activeNotes[ch] = make(map[uint8]bool)
				}
				e.activeNotes[ch][key] = true
			} else {
				// Note on with velocity 0 = note off
				if e.activeNotes[ch] != nil {
					delete(e.activeNotes[ch], key)
				}
			}
		}
	} else if event.MessageType == "Note Off" {
		if event.Message.GetNoteOff(&ch, &key, &vel) {
			if e.activeNotes[ch] != nil {
				delete(e.activeNotes[ch], key)
			}
		}
	}
}

func (e EventViewer) replaceNoteNumbers(event midi.Event, data string) string {
	var ch, key, vel, pressure uint8

	if event.MessageType == "Note On" || event.MessageType == "Note Off" {
		if event.Message.GetNoteOn(&ch, &key, &vel) {
			return fmt.Sprintf("Note: %s, Velocity: %d", midi.NoteToName(key), vel)
		}
		if event.Message.GetNoteOff(&ch, &key, &vel) {
			return fmt.Sprintf("Note: %s, Velocity: %d", midi.NoteToName(key), vel)
		}
	}
	if event.MessageType == "Poly Aftertouch" {
		if event.Message.GetPolyAfterTouch(&ch, &key, &pressure) {
			return fmt.Sprintf("Note: %s, Pressure: %d", midi.NoteToName(key), pressure)
		}
	}
	return data
}

// parseEventData extracts note, velocity, controller, and value from an event
func (e EventViewer) parseEventData(event midi.Event) (note, vel, ctrl, val string) {
	var ch, key, velocity, controller, value, pressure uint8
	var pitchValue int16

	switch event.MessageType {
	case "Note On", "Note Off":
		if event.Message.GetNoteOn(&ch, &key, &velocity) || event.Message.GetNoteOff(&ch, &key, &velocity) {
			if e.filter.ShowMusicalNotes {
				note = midi.NoteToName(key)
			} else {
				note = fmt.Sprintf("%d", key)
			}
			vel = fmt.Sprintf("%d", velocity)
		}
	case "Poly Aftertouch":
		if event.Message.GetPolyAfterTouch(&ch, &key, &pressure) {
			if e.filter.ShowMusicalNotes {
				note = midi.NoteToName(key)
			} else {
				note = fmt.Sprintf("%d", key)
			}
			val = fmt.Sprintf("%d", pressure)
		}
	case "CC":
		if event.Message.GetControlChange(&ch, &controller, &value) {
			ctrl = fmt.Sprintf("%d", controller)
			val = fmt.Sprintf("%d", value)
		}
	case "Program Change":
		if event.Message.GetProgramChange(&ch, &value) {
			val = fmt.Sprintf("%d", value)
		}
	case "Aftertouch":
		if event.Message.GetAfterTouch(&ch, &pressure) {
			val = fmt.Sprintf("%d", pressure)
		}
	case "Pitch Bend":
		var pitchLSB, pitchMSB uint16
		if event.Message.GetPitchBend(&ch, &pitchValue, &pitchLSB) {
			val = fmt.Sprintf("%d", pitchValue)
		} else {
			_ = pitchMSB // unused
		}
	}

	return
}

// MIDIEventMsg is sent when a new MIDI event is received
type MIDIEventMsg struct {
	Message gomidi.Message
}

// OpenOptionsModalMsg is sent to open the options modal
type OpenOptionsModalMsg struct{}

// FilterUpdatedMsg is sent when the filter has been updated
type FilterUpdatedMsg struct {
	Filter models.Filter
}

// BackToDeviceSelectionMsg is sent to return to device selection
type BackToDeviceSelectionMsg struct{}
