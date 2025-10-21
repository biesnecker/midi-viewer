package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"midi-viewer/internal/models"
	"midi-viewer/internal/ui/theme"
)

type optionsModalKeyMap struct {
	Up     key.Binding
	Down   key.Binding
	Left   key.Binding
	Right  key.Binding
	Toggle key.Binding
	Close  key.Binding
	Clear  key.Binding
}

var optionsModalKeys = optionsModalKeyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "down"),
	),
	Left: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←/h", "left"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→/l", "right"),
	),
	Toggle: key.NewBinding(
		key.WithKeys(" ", "enter"),
		key.WithHelp("space/enter", "toggle"),
	),
	Close: key.NewBinding(
		key.WithKeys("o", "esc"),
		key.WithHelp("o/esc", "close"),
	),
	Clear: key.NewBinding(
		key.WithKeys("c"),
		key.WithHelp("c", "clear all"),
	),
}

type optionsSection int

const (
	sectionChannels optionsSection = iota
	sectionMessageTypes
	sectionColumns
	sectionSettings
)

// OptionsModal allows filtering events by channel, message type, and column visibility
type OptionsModal struct {
	filter         models.Filter
	theme          theme.Theme
	width          int
	height         int
	cursor         int
	currentSection optionsSection
	messageTypes   []string
	columns        []string
}

// NewOptionsModal creates a new options modal
func NewOptionsModal(filter models.Filter, t theme.Theme) OptionsModal {
	return OptionsModal{
		filter:         filter,
		theme:          t,
		currentSection: sectionChannels,
		cursor:         0,
		messageTypes: []string{
			"Note On",
			"Note Off",
			"CC",
			"Program Change",
			"Pitch Bend",
			"Poly Aftertouch",
			"Aftertouch",
			"SysEx",
			"Clock",
			"Start",
			"Stop",
			"Continue",
		},
		columns: []string{
			"Time",
			"Chan",
			"Event",
			"Note",
			"Vel",
			"Ctrl",
			"Val",
		},
	}
}

// Update handles messages
func (o OptionsModal) Update(msg tea.Msg) (OptionsModal, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		o.width = msg.Width
		o.height = msg.Height

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, optionsModalKeys.Up):
			if o.cursor > 0 {
				o.cursor--
			}
		case key.Matches(msg, optionsModalKeys.Down):
			maxCursor := 0
			if o.currentSection == sectionChannels {
				maxCursor = 15 // 16 channels (0-15)
			} else if o.currentSection == sectionMessageTypes {
				maxCursor = len(o.messageTypes) - 1
			} else if o.currentSection == sectionColumns {
				maxCursor = len(o.columns) - 1
			} else if o.currentSection == sectionSettings {
				maxCursor = 0 // Only one setting for now
			}
			if o.cursor < maxCursor {
				o.cursor++
			}
		case key.Matches(msg, optionsModalKeys.Left):
			if o.currentSection == sectionMessageTypes {
				o.currentSection = sectionChannels
				o.cursor = 0
			} else if o.currentSection == sectionColumns {
				o.currentSection = sectionMessageTypes
				o.cursor = 0
			} else if o.currentSection == sectionSettings {
				o.currentSection = sectionColumns
				o.cursor = 0
			}
		case key.Matches(msg, optionsModalKeys.Right):
			if o.currentSection == sectionChannels {
				o.currentSection = sectionMessageTypes
				o.cursor = 0
			} else if o.currentSection == sectionMessageTypes {
				o.currentSection = sectionColumns
				o.cursor = 0
			} else if o.currentSection == sectionColumns {
				o.currentSection = sectionSettings
				o.cursor = 0
			}
		case key.Matches(msg, optionsModalKeys.Toggle):
			if o.currentSection == sectionChannels {
				o.filter.ToggleChannel(uint8(o.cursor))
			} else if o.currentSection == sectionMessageTypes {
				o.filter.ToggleMessageType(o.messageTypes[o.cursor])
			} else if o.currentSection == sectionColumns {
				o.filter.ToggleColumn(o.columns[o.cursor])
			} else if o.currentSection == sectionSettings {
				// Toggle musical notes setting
				o.filter.ShowMusicalNotes = !o.filter.ShowMusicalNotes
			}
		case key.Matches(msg, optionsModalKeys.Clear):
			o.filter = models.NewFilter()
		case key.Matches(msg, optionsModalKeys.Close):
			return o, func() tea.Msg {
				return CloseOptionsModalMsg{o.filter}
			}
		}
	}

	return o, nil
}

// View renders the options modal
func (o OptionsModal) View() string {
	modalWidth := 110
	modalHeight := 25

	modalStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(o.theme.ModalBorder).
		Background(o.theme.ModalBackground).
		Padding(1, 2).
		Width(modalWidth).
		Height(modalHeight)

	titleStyle := lipgloss.NewStyle().
		Foreground(o.theme.Primary).
		Bold(true).
		Align(lipgloss.Center).
		Width(modalWidth - 4)

	sectionTitleStyle := lipgloss.NewStyle().
		Foreground(o.theme.Secondary).
		Bold(true).
		Underline(true).
		Padding(0, 0, 1, 0)

	itemStyle := lipgloss.NewStyle().
		Foreground(o.theme.Foreground)

	selectedStyle := lipgloss.NewStyle().
		Foreground(o.theme.Background).
		Background(o.theme.Primary).
		Bold(true)

	activeStyle := lipgloss.NewStyle().
		Foreground(o.theme.Success).
		Bold(true)

	helpStyle := lipgloss.NewStyle().
		Foreground(o.theme.Muted).
		Align(lipgloss.Center).
		Width(modalWidth - 4).
		Padding(1, 0, 0, 0)

	var b strings.Builder

	b.WriteString(titleStyle.Render("Options"))
	b.WriteString("\n\n")

	// Four column layout
	channelsCol := o.renderChannelSection(sectionTitleStyle, itemStyle, selectedStyle, activeStyle)
	typesCol := o.renderMessageTypeSection(sectionTitleStyle, itemStyle, selectedStyle, activeStyle)
	columnsCol := o.renderColumnsSection(sectionTitleStyle, itemStyle, selectedStyle, activeStyle)
	settingsCol := o.renderSettingsSection(sectionTitleStyle, itemStyle, selectedStyle, activeStyle)

	columns := lipgloss.JoinHorizontal(
		lipgloss.Top,
		lipgloss.NewStyle().Width(22).Render(channelsCol),
		lipgloss.NewStyle().Width(22).Render(typesCol),
		lipgloss.NewStyle().Width(22).Render(columnsCol),
		lipgloss.NewStyle().Width(22).Render(settingsCol),
	)

	b.WriteString(columns)
	b.WriteString("\n")

	helpText := "↑/↓: navigate • ←/→: switch section • space: toggle • c: clear all • o/esc: close"
	b.WriteString(helpStyle.Render(helpText))

	modal := modalStyle.Render(b.String())

	// Center the modal
	return lipgloss.Place(
		o.width,
		o.height,
		lipgloss.Center,
		lipgloss.Center,
		modal,
	)
}

func (o OptionsModal) renderChannelSection(titleStyle, itemStyle, selectedStyle, activeStyle lipgloss.Style) string {
	var b strings.Builder

	sectionActive := o.currentSection == sectionChannels
	title := "Channels"
	if sectionActive {
		title = "> " + title
	} else {
		title = "  " + title
	}
	b.WriteString(titleStyle.Render(title))
	b.WriteString("\n")

	for ch := 0; ch < 16; ch++ {
		cursor := "  "
		if sectionActive && o.cursor == ch {
			cursor = "> "
		}

		checkbox := "☐"
		if o.filter.IsChannelVisible(uint8(ch)) {
			checkbox = "☑"
		}

		label := fmt.Sprintf("%s%s Ch %d", cursor, checkbox, ch+1)

		if sectionActive && o.cursor == ch {
			if o.filter.IsChannelVisible(uint8(ch)) {
				b.WriteString(activeStyle.Render(label))
			} else {
				b.WriteString(selectedStyle.Render(label))
			}
		} else {
			if o.filter.IsChannelVisible(uint8(ch)) {
				b.WriteString(activeStyle.Render(label))
			} else {
				b.WriteString(itemStyle.Render(label))
			}
		}
		b.WriteString("\n")
	}

	return b.String()
}

func (o OptionsModal) renderMessageTypeSection(titleStyle, itemStyle, selectedStyle, activeStyle lipgloss.Style) string {
	var b strings.Builder

	sectionActive := o.currentSection == sectionMessageTypes
	title := "Event Types"
	if sectionActive {
		title = "> " + title
	} else {
		title = "  " + title
	}
	b.WriteString(titleStyle.Render(title))
	b.WriteString("\n")

	for i, msgType := range o.messageTypes {
		cursor := "  "
		if sectionActive && o.cursor == i {
			cursor = "> "
		}

		checkbox := "☐"
		if o.filter.IsMessageTypeVisible(msgType) {
			checkbox = "☑"
		}

		label := fmt.Sprintf("%s%s %s", cursor, checkbox, msgType)

		if sectionActive && o.cursor == i {
			if o.filter.IsMessageTypeVisible(msgType) {
				b.WriteString(activeStyle.Render(label))
			} else {
				b.WriteString(selectedStyle.Render(label))
			}
		} else {
			if o.filter.IsMessageTypeVisible(msgType) {
				b.WriteString(activeStyle.Render(label))
			} else {
				b.WriteString(itemStyle.Render(label))
			}
		}
		b.WriteString("\n")
	}

	return b.String()
}

func (o OptionsModal) renderColumnsSection(titleStyle, itemStyle, selectedStyle, activeStyle lipgloss.Style) string {
	var b strings.Builder

	sectionActive := o.currentSection == sectionColumns
	title := "Columns"
	if sectionActive {
		title = "> " + title
	} else {
		title = "  " + title
	}
	b.WriteString(titleStyle.Render(title))
	b.WriteString("\n")

	for i, col := range o.columns {
		cursor := "  "
		if sectionActive && o.cursor == i {
			cursor = "> "
		}

		checkbox := "☐"
		if o.filter.IsColumnVisible(col) {
			checkbox = "☑"
		}

		label := fmt.Sprintf("%s%s %s", cursor, checkbox, col)

		if sectionActive && o.cursor == i {
			if o.filter.IsColumnVisible(col) {
				b.WriteString(activeStyle.Render(label))
			} else {
				b.WriteString(selectedStyle.Render(label))
			}
		} else {
			if o.filter.IsColumnVisible(col) {
				b.WriteString(activeStyle.Render(label))
			} else {
				b.WriteString(itemStyle.Render(label))
			}
		}
		b.WriteString("\n")
	}

	return b.String()
}

func (o OptionsModal) renderSettingsSection(titleStyle, itemStyle, selectedStyle, activeStyle lipgloss.Style) string {
	var b strings.Builder

	sectionActive := o.currentSection == sectionSettings
	title := "Settings"
	if sectionActive {
		title = "> " + title
	} else {
		title = "  " + title
	}
	b.WriteString(titleStyle.Render(title))
	b.WriteString("\n")

	cursor := "  "
	if sectionActive && o.cursor == 0 {
		cursor = "> "
	}

	checkbox := "☐"
	if o.filter.ShowMusicalNotes {
		checkbox = "☑"
	}

	label := fmt.Sprintf("%s%s Musical Notes", cursor, checkbox)

	if sectionActive && o.cursor == 0 {
		if o.filter.ShowMusicalNotes {
			b.WriteString(activeStyle.Render(label))
		} else {
			b.WriteString(selectedStyle.Render(label))
		}
	} else {
		if o.filter.ShowMusicalNotes {
			b.WriteString(activeStyle.Render(label))
		} else {
			b.WriteString(itemStyle.Render(label))
		}
	}
	b.WriteString("\n")

	return b.String()
}

// CloseOptionsModalMsg is sent when the options modal is closed
type CloseOptionsModalMsg struct {
	Filter models.Filter
}
