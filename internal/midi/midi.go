package midi

import (
	"fmt"
	"time"

	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/drivers"
	"gitlab.com/gomidi/midi/v2/drivers/rtmididrv"
)

// Device represents a MIDI input device
type Device struct {
	Name   string
	Number int
	Port   drivers.In
}

// Event represents a MIDI event with metadata
type Event struct {
	Timestamp   time.Time
	Message     midi.Message
	Channel     uint8
	MessageType string
	Data        string
	RawBytes    []byte
}

// InitDriver initializes the MIDI driver
func InitDriver() (func(), error) {
	drv, err := rtmididrv.New()
	if err != nil {
		return nil, fmt.Errorf("could not create MIDI driver: %w", err)
	}

	drivers.Register(drv)

	cleanup := func() {
		drv.Close()
	}

	return cleanup, nil
}

// GetInputDevices returns all available MIDI input devices
func GetInputDevices() ([]Device, error) {
	ins := midi.GetInPorts()
	devices := make([]Device, len(ins))

	for i, port := range ins {
		devices[i] = Device{
			Name:   port.String(),
			Number: i,
			Port:   port,
		}
	}

	return devices, nil
}

// ParseMessage parses a MIDI message into an Event
func ParseMessage(msg midi.Message) Event {
	event := Event{
		Timestamp:   time.Now(),
		Message:     msg,
		RawBytes:    msg.Bytes(),
		MessageType: getMessageType(msg),
		Data:        formatMessageData(msg),
	}

	// Extract channel if applicable
	var ch uint8
	if msg.GetChannel(&ch) {
		event.Channel = ch
	}

	return event
}

func getMessageType(msg midi.Message) string {
	var ch, key, vel, controller, value, program, pressure uint8
	var rel int16
	var abs uint16
	var bt []byte

	switch {
	case msg.GetNoteOn(&ch, &key, &vel):
		return "Note On"
	case msg.GetNoteOff(&ch, &key, &vel):
		return "Note Off"
	case msg.GetControlChange(&ch, &controller, &value):
		return "CC"
	case msg.GetProgramChange(&ch, &program):
		return "Program Change"
	case msg.GetPitchBend(&ch, &rel, &abs):
		return "Pitch Bend"
	case msg.GetPolyAfterTouch(&ch, &key, &pressure):
		return "Poly Aftertouch"
	case msg.GetAfterTouch(&ch, &pressure):
		return "Aftertouch"
	case msg.GetSysEx(&bt):
		return "SysEx"
	case msg.Is(midi.TimingClockMsg):
		return "Clock"
	case msg.Is(midi.StartMsg):
		return "Start"
	case msg.Is(midi.StopMsg):
		return "Stop"
	case msg.Is(midi.ContinueMsg):
		return "Continue"
	case msg.Is(midi.ActiveSenseMsg):
		return "Active Sense"
	case msg.Is(midi.ResetMsg):
		return "Reset"
	default:
		return "Unknown"
	}
}

func formatMessageData(msg midi.Message) string {
	var ch, key, vel, controller, value, program, pressure uint8
	var rel int16
	var abs uint16

	switch {
	case msg.GetNoteOn(&ch, &key, &vel):
		return fmt.Sprintf("Note: %d, Velocity: %d", key, vel)
	case msg.GetNoteOff(&ch, &key, &vel):
		return fmt.Sprintf("Note: %d, Velocity: %d", key, vel)
	case msg.GetControlChange(&ch, &controller, &value):
		return fmt.Sprintf("Controller: %d, Value: %d", controller, value)
	case msg.GetProgramChange(&ch, &program):
		return fmt.Sprintf("Program: %d", program)
	case msg.GetPitchBend(&ch, &rel, &abs):
		return fmt.Sprintf("Value: %d (abs: %d)", rel, abs)
	case msg.GetPolyAfterTouch(&ch, &key, &pressure):
		return fmt.Sprintf("Note: %d, Pressure: %d", key, pressure)
	case msg.GetAfterTouch(&ch, &pressure):
		return fmt.Sprintf("Pressure: %d", pressure)
	default:
		return fmt.Sprintf("%v", msg.Bytes())
	}
}

// NoteToName converts a MIDI note number to its musical name (e.g., 60 -> C4)
func NoteToName(note uint8) string {
	noteNames := []string{"C", "C#", "D", "D#", "E", "F", "F#", "G", "G#", "A", "A#", "B"}
	octave := int(note/12) - 1
	noteName := noteNames[note%12]
	return fmt.Sprintf("%s%d", noteName, octave)
}
