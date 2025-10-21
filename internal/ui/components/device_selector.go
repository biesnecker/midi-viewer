package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"midi-viewer/internal/midi"
	"midi-viewer/internal/ui/theme"
)

type deviceSelectorKeyMap struct {
	Up     key.Binding
	Down   key.Binding
	Select key.Binding
	Quit   key.Binding
}

var deviceSelectorKeys = deviceSelectorKeyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "down"),
	),
	Select: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c", "esc"),
		key.WithHelp("q/esc", "quit"),
	),
}

// DeviceSelector represents the device selection UI component
type DeviceSelector struct {
	devices      []midi.Device
	cursor       int
	theme        theme.Theme
	width        int
	height       int
	err          error
}

// NewDeviceSelector creates a new device selector
func NewDeviceSelector(t theme.Theme) DeviceSelector {
	return DeviceSelector{
		theme:   t,
		cursor:  0,
	}
}

// Init initializes the device selector
func (d DeviceSelector) Init() tea.Cmd {
	return func() tea.Msg {
		devices, err := midi.GetInputDevices()
		if err != nil {
			return DeviceErrorMsg{err}
		}
		return DevicesLoadedMsg{devices}
	}
}

// Update handles messages
func (d DeviceSelector) Update(msg tea.Msg) (DeviceSelector, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		d.width = msg.Width
		d.height = msg.Height

	case DevicesLoadedMsg:
		d.devices = msg.Devices
		if len(d.devices) == 0 {
			d.err = fmt.Errorf("no MIDI input devices found")
		}

	case DeviceErrorMsg:
		d.err = msg.Err

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, deviceSelectorKeys.Quit):
			return d, tea.Quit
		case key.Matches(msg, deviceSelectorKeys.Up):
			if d.cursor > 0 {
				d.cursor--
			}
		case key.Matches(msg, deviceSelectorKeys.Down):
			if d.cursor < len(d.devices)-1 {
				d.cursor++
			}
		case key.Matches(msg, deviceSelectorKeys.Select):
			if len(d.devices) > 0 {
				return d, func() tea.Msg {
					return DeviceSelectedMsg{d.devices[d.cursor]}
				}
			}
		}
	}

	return d, nil
}

// View renders the device selector
func (d DeviceSelector) View() string {
	if d.err != nil {
		return d.renderError()
	}

	if len(d.devices) == 0 {
		return d.renderLoading()
	}

	return d.renderDeviceList()
}

func (d DeviceSelector) renderError() string {
	errorStyle := lipgloss.NewStyle().
		Foreground(d.theme.Error).
		Bold(true).
		Padding(1, 2)

	return lipgloss.Place(
		d.width,
		d.height,
		lipgloss.Center,
		lipgloss.Center,
		errorStyle.Render(fmt.Sprintf("Error: %v", d.err)),
	)
}

func (d DeviceSelector) renderLoading() string {
	loadingStyle := lipgloss.NewStyle().
		Foreground(d.theme.Muted).
		Padding(1, 2)

	return lipgloss.Place(
		d.width,
		d.height,
		lipgloss.Center,
		lipgloss.Center,
		loadingStyle.Render("Loading MIDI devices..."),
	)
}

func (d DeviceSelector) renderDeviceList() string {
	titleStyle := lipgloss.NewStyle().
		Foreground(d.theme.Primary).
		Bold(true).
		Padding(1, 0)

	itemStyle := lipgloss.NewStyle().
		Foreground(d.theme.Foreground).
		Padding(0, 2)

	selectedStyle := lipgloss.NewStyle().
		Foreground(d.theme.Background).
		Background(d.theme.Primary).
		Padding(0, 2).
		Bold(true)

	helpStyle := lipgloss.NewStyle().
		Foreground(d.theme.Muted).
		Padding(1, 0)

	var b strings.Builder

	b.WriteString(titleStyle.Render("Select MIDI Input Device"))
	b.WriteString("\n\n")

	for i, device := range d.devices {
		cursor := "  "
		if i == d.cursor {
			cursor = "> "
		}

		line := fmt.Sprintf("%s%d. %s", cursor, i+1, device.Name)

		if i == d.cursor {
			b.WriteString(selectedStyle.Render(line))
		} else {
			b.WriteString(itemStyle.Render(line))
		}
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(helpStyle.Render("↑/↓: navigate • enter: select • q/esc: quit"))

	return b.String()
}

// DevicesLoadedMsg is sent when MIDI devices have been loaded
type DevicesLoadedMsg struct {
	Devices []midi.Device
}

// DeviceSelectedMsg is sent when a user selects a MIDI device
type DeviceSelectedMsg struct {
	Device midi.Device
}

// DeviceErrorMsg is sent when there's an error loading or selecting devices
type DeviceErrorMsg struct {
	Err error
}
