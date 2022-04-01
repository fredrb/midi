package smf

import (
	"bytes"
	"fmt"

	"gitlab.com/gomidi/midi/v2"
)

type Message []byte

func (m Message) Bytes() []byte {
	return []byte(m)
}

func (m Message) IsPlayable() bool {
	return m.Type().IsPlayable()
}

/*
func (m Message) Type() midi.Type {
	return GetMetaType(m[1])
}
*/

func (m Message) Type() midi.Type {
	return GetType(m)
}

func GetType(msg []byte) midi.Type {
	if IsMeta(msg) {
		return GetMetaType(msg[1])
	} else {
		return midi.GetType(msg)
	}
}

func (m Message) Is(t midi.Type) bool {
	return m.Type().Is(t)
}

// NoteOn returns true if (and only if) the message is a NoteOnMsg.
// Then it also extracts the data to the given arguments
func (m Message) ScanNoteOn(channel, key, velocity *uint8) (is bool) {
	return midi.Message(m).ScanNoteOn(channel, key, velocity)
}

// NoteStart returns true if (and only if) the message is a NoteOnMsg with a velocity > 0.
// Then it also extracts the data to the given arguments
func (m Message) ScanNoteStart(channel, key, velocity *uint8) (is bool) {
	return midi.Message(m).ScanNoteStart(channel, key, velocity)
}

// NoteOff returns true if (and only if) the message is a NoteOffMsg.
// Then it also extracts the data to the given arguments
func (m Message) ScanNoteOff(channel, key, velocity *uint8) (is bool) {
	return midi.Message(m).ScanNoteOff(channel, key, velocity)
}

// Channel returns true if (and only if) the message is a ChannelMsg.
// Then it also extracts the data to the given arguments
func (m Message) ScanChannel(channel *uint8) (is bool) {
	return midi.Message(m).ScanChannel(channel)
}

// NoteEnd returns true if (and only if) the message is a NoteOnMsg with a velocity == 0 or a NoteOffMsg.
// Then it also extracts the data to the given arguments
func (m Message) ScanNoteEnd(channel, key, velocity *uint8) (is bool) {
	return midi.Message(m).ScanNoteEnd(channel, key, velocity)
}

// PolyAfterTouch returns true if (and only if) the message is a PolyAfterTouchMsg.
// Then it also extracts the data to the given arguments
func (m Message) ScanPolyAfterTouch(channel, key, pressure *uint8) (is bool) {
	return midi.Message(m).ScanPolyAfterTouch(channel, key, pressure)
}

// AfterTouch returns true if (and only if) the message is a AfterTouchMsg.
// Then it also extracts the data to the given arguments
func (m Message) ScanAfterTouch(channel, pressure *uint8) (is bool) {
	return midi.Message(m).ScanAfterTouch(channel, pressure)
}

// ProgramChange returns true if (and only if) the message is a ProgramChangeMsg.
// Then it also extracts the data to the given arguments
func (m Message) ScanProgramChange(channel, program *uint8) (is bool) {
	return midi.Message(m).ScanProgramChange(channel, program)
}

// PitchBend returns true if (and only if) the message is a PitchBendMsg.
// Then it also extracts the data to the given arguments
// Either relative or absolute may be nil, if not needed.
func (m Message) ScanPitchBend(channel *uint8, relative *int16, absolute *uint16) (is bool) {
	return midi.Message(m).ScanPitchBend(channel, relative, absolute)
}

// ControlChange returns true if (and only if) the message is a ControlChangeMsg.
// Then it also extracts the data to the given arguments
func (m Message) ScanControlChange(channel, controller, value *uint8) (is bool) {
	return midi.Message(m).ScanControlChange(channel, controller, value)
}

// String represents the Message as a string that contains the MsgType and its properties.
func (m Message) String() string {

	if IsMeta(m) {
		var bf bytes.Buffer
		fmt.Fprintf(&bf, m.Type().String())

		var val1 uint8
		var val2 uint8
		var text string
		var bpm float64
		// TODO: complete
		switch {
		case m.ScanTempo(&bpm):
			fmt.Fprintf(&bf, " bpm: %0.2f", bpm)
		case m.ScanMeter(&val1, &val2):
			fmt.Fprintf(&bf, " meter: %v/%v", val1, val2)
		default:
			switch m.Type() {
			case MetaLyric, MetaMarker, MetaCopyright, MetaText, MetaCuepoint, MetaDevice, MetaInstrument, MetaProgramName, MetaTrackName:
				m.text(&text)
				fmt.Fprintf(&bf, " text: %q", text)
			}
		}

		return bf.String()
	} else {
		return midi.Message(m).String()
	}

}
