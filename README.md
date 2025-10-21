# MIDI Viewer

A simple terminal-based MIDI debug viewer featuring real-time event monitoring, flexible filtering, and an intuitive interface.

## Features

- **Device Selection**: Choose from available MIDI input devices
- **Real-time Event Display**: View MIDI events as they happen with newest events at the top
- **Flexible Filtering**:
  - Filter by MIDI channels (1-16)
  - Filter by message types (Note On/Off, CC, Program Change, Pitch Bend, etc.)
  - Toggle column visibility (Time, Channel, Event, Note, Velocity, Controller, Value)
  - Show musical note names (C4, D#5) or MIDI note numbers (60, 63)
- **Active Notes Display**: See which notes are currently playing
- **Pause/Resume**: Pause event capture to examine current events
- **Theme Support**: Choose between dark and light themes
- **Clean TUI**: Built with [Bubble Tea](https://github.com/charmbracelet/bubbletea) and [Lipgloss](https://github.com/charmbracelet/lipgloss)

## Installation

### Prerequisites

- Go 1.21 or later
- A C++ compiler (for RtMidi)
  - macOS: Xcode Command Line Tools (`xcode-select --install`)
  - Linux: `build-essential` package
  - Windows: MinGW or Visual Studio

### Build from Source

```bash
git clone <repository-url>
cd midi-viewer
make build
```

This will create a `midi-viewer` executable in the current directory.

Alternatively, build manually:

```bash
CGO_CXXFLAGS="-Wno-vla-cxx-extension" go build -o midi-viewer ./cmd/midi-viewer
```

## Usage

### Basic Usage

```bash
./midi-viewer
```

### Theme Selection

```bash
# Use dark theme (default)
./midi-viewer --theme dark

# Use light theme
./midi-viewer --theme light
```

## Keyboard Controls

### Device Selection Screen
- `↑/↓` or `k/j`: Navigate device list
- `Enter`: Select device
- `q` or `Esc`: Quit

### Event Viewer Screen
- `Space`: Pause/unpause event capture
- `o`: Open options modal (filtering and settings)
- `c`: Clear all captured events
- `Esc`: Return to device selection
- `q` or `Ctrl+C`: Quit application

### Options Modal
- `←/→` or `h/l`: Switch between sections (Channels, Event Types, Columns, Settings)
- `↑/↓` or `k/j`: Navigate items within a section
- `Space` or `Enter`: Toggle selected item
- `c`: Clear all filters (show everything)
- `o` or `Esc`: Close modal and apply changes

## Event Display

The event viewer shows MIDI events in a columnar format:

| Column | Description |
|--------|-------------|
| **Time** | Timestamp of the event (HH:MM:SS.mmm) |
| **Chan** | MIDI channel (1-16) |
| **Event** | Message type (Note On, Note Off, CC, etc.) |
| **Note** | Note name/number (for note events) |
| **Vel** | Velocity (for note events) |
| **Ctrl** | Controller number (for CC events) |
| **Val** | Controller/pitch/pressure value |

The display shows the most recent events at the top, keeping up to 1000 events in memory.

### Active Notes

At the bottom of the event viewer, you'll see a line showing currently playing notes (notes that have received Note On but not yet Note Off). This is helpful for debugging stuck notes or understanding chord progression.

## Filtering

Access the options modal by pressing `o` in the event viewer.

### Channels
Toggle visibility for individual MIDI channels (1-16). By default, all channels are visible.

### Event Types
Filter the types of MIDI messages to display:
- **Note On**: Note pressed
- **Note Off**: Note released
- **CC**: Control Change messages
- **Program Change**: Instrument/patch changes
- **Pitch Bend**: Pitch wheel movements
- **Poly Aftertouch**: Per-note pressure
- **Aftertouch**: Channel pressure
- **SysEx**: System Exclusive messages
- **Clock/Start/Stop/Continue**: MIDI timing messages

### Columns
Show or hide specific columns in the event display to focus on relevant information.

### Settings
- **Musical Notes**: Toggle between musical note names (C4, D#5) and MIDI note numbers (60, 63)

## Development

### Build

```bash
make build
```

### Clean

```bash
make clean
```

### Run Tests

```bash
make test
```

Or manually:

```bash
go test ./... -v
```

## Project Structure

```
midi-viewer/
├── cmd/
│   └── midi-viewer/     # Main application entry point
├── internal/
│   ├── midi/            # MIDI device handling and parsing
│   ├── models/          # Data models and filtering logic
│   └── ui/
│       ├── components/  # UI components (device selector, event viewer, options modal)
│       └── theme/       # Color themes
├── Makefile
└── README.md
```

## Architecture

The application follows clean, composable architecture principles:

- **Components**: Self-contained UI components using the Elm Architecture (Model-Update-View)
- **Models**: Shared data structures and business logic
- **MIDI Layer**: Abstraction over the gomidi library for device management
- **Themes**: Centralized color schemes for consistent styling

## Dependencies

- [gomidi/midi](https://gitlab.com/gomidi/midi) - MIDI library for Go
- [gomidi/rtmididrv](https://gitlab.com/gomidi/rtmididrv) - RtMidi driver for cross-platform MIDI support
- [charmbracelet/bubbletea](https://github.com/charmbracelet/bubbletea) - Terminal UI framework
- [charmbracelet/lipgloss](https://github.com/charmbracelet/lipgloss) - Terminal styling library

## System Requirements

This application only works on MacOS at the moment because of the dependency on RtMidi. In theory it should be compatible with any platform with a compatible MIDI library, so pull requests are welcome!

### Windows
Requires a C++ compiler like MinGW or Visual Studio.

## Message Types Supported

- **Note On/Off**: Key press and release events
- **Control Change (CC)**: Knobs, sliders, pedals
- **Program Change**: Instrument/patch selection
- **Pitch Bend**: Pitch wheel movements
- **Poly Aftertouch**: Per-note pressure
- **Aftertouch**: Channel pressure
- **SysEx**: System Exclusive messages
- **Clock**: MIDI timing clock
- **Start/Stop/Continue**: Transport controls
- **Active Sense**: Keep-alive messages
- **Reset**: System reset

## Known Limitations

- The application keeps the last 1000 events in memory. Older events are automatically discarded.

## License

MIT

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
