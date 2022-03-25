package midi

import (
	"bytes"
	"fmt"
	"testing"
)

func TestChannelString(t *testing.T) {

	tests := []struct {
		input    Message
		expected string
	}{
		{
			Channel(1).Aftertouch(120),
			"AfterTouchMsg channel: 1 pressure: 120",
		},
		{
			Channel(8).ControlChange(7, 110),
			//"ControlChangeMsg channel: 8 controller: 7 (\"Volume (MSB)\") value 110",
			"ControlChangeMsg channel: 8 controller: 7 value: 110",
		},
		{
			Channel(2).NoteOn(100, 80),
			"NoteOnMsg channel: 2 key: 100 velocity: 80",
		},
		{
			Channel(3).NoteOff(80),
			"NoteOffMsg channel: 3 key: 80",
		},
		{
			Channel(4).NoteOffVelocity(80, 20),
			"NoteOffMsg channel: 4 key: 80 velocity: 20",
		},
		{
			Channel(4).Pitchbend(300),
			"PitchBendMsg channel: 4 pitch: 300 (8492)",
		},
		{
			Channel(4).PolyAftertouch(86, 109),
			"PolyAfterTouchMsg channel: 4 key: 86 pressure: 109",
		},
		{
			Channel(4).ProgramChange(83),
			"ProgramChangeMsg channel: 4 program: 83",
		},

		// too high values
		{
			Channel(1).Aftertouch(130),
			"AfterTouchMsg channel: 1 pressure: 127",
		},
		{
			Channel(8).ControlChange(137, 130),
			//"ControlChangeMsg channel: 8 controller: 127 (\"Poly Operation\") value 127",
			"ControlChangeMsg channel: 8 controller: 127 value: 127",
		},
		{
			Channel(2).NoteOn(130, 130),
			"NoteOnMsg channel: 2 key: 127 velocity: 127",
		},
		{
			Channel(3).NoteOff(180),
			"NoteOffMsg channel: 3 key: 127",
		},
		{
			Channel(4).NoteOffVelocity(180, 220),
			"NoteOffMsg channel: 4 key: 127 velocity: 127",
		},
		{
			Channel(4).Pitchbend(12300),
			"PitchBendMsg channel: 4 pitch: 8191 (16383)",
		},
		{
			Channel(4).PolyAftertouch(186, 190),
			"PolyAfterTouchMsg channel: 4 key: 127 pressure: 127",
		},
		{
			Channel(4).ProgramChange(183),
			"ProgramChangeMsg channel: 4 program: 127",
		},
	}

	for _, test := range tests {

		var bf bytes.Buffer

		bf.WriteString(test.input.String())

		if got, want := bf.String(), test.expected; got != want {
			t.Errorf("got: %#v; wanted %#v", got, want)
		}
	}

}

func TestChannelRaw(t *testing.T) {

	tests := []struct {
		input    Message
		expected string
	}{
		{ // 0
			Channel(1).Aftertouch(120),
			"D1 78",
		},
		{ // 1
			Channel(8).ControlChange(7, 110),
			"B8 07 6E",
		},
		{ // 2
			Channel(2).NoteOn(100, 80),
			"92 64 50",
		},
		{ // 3
			Channel(3).NoteOff(80),
			"83 50 00",
		},
		{
			Channel(4).NoteOffVelocity(80, 20),
			"84 50 14",
		},
		{
			Channel(4).Pitchbend(300),
			"E4 2C 42",
		},
		{
			Channel(4).PolyAftertouch(86, 109),
			"A4 56 6D",
		},
		{
			Channel(4).ProgramChange(83),
			"C4 53",
		},
	}

	for i, test := range tests {

		var bf bytes.Buffer

		bf.Write(test.input.Bytes())

		if got, want := fmt.Sprintf("% X", bf.Bytes()), test.expected; got != want {
			t.Errorf("[%v] got: %#v; wanted %#v", i, got, want)
		}
	}

}

/*
func TestSetChannel(t *testing.T) {

	tests := []struct {
		input     Message
		toChannel uint8
		expected  string
	}{
		{
			Channel1.Aftertouch(120),
			5,
			"channel.Aftertouch channel 5 pressure 120",
		},
		{
			Channel8.ControlChange(7, 110),
			9,
			"channel.ControlChange channel 9 controller 7 (\"Volume (MSB)\") value 110",
		},
		{
			Channel2.NoteOn(100, 80),
			0,
			"channel.NoteOn channel 0 key 100 velocity 80",
		},
		{
			Channel3.NoteOff(80),
			2,
			"channel.NoteOff channel 2 key 80",
		},
		{
			Channel4.NoteOffVelocity(80, 20),
			11,
			"channel.NoteOffVelocity channel 11 key 80 velocity 20",
		},
		{
			Channel4.Pitchbend(300),
			14,
			"channel.Pitchbend channel 14 value 300 absValue 0",
		},
		{
			Channel4.PolyAftertouch(86, 109),
			2,
			"channel.PolyAftertouch channel 2 key 86 pressure 109",
		},
		{
			Channel4.ProgramChange(83),
			0,
			"channel.ProgramChange channel 0 program 83",
		},
	}

	for _, test := range tests {

		var bf bytes.Buffer

		msg := SetChannel(test.input, test.toChannel)

		bf.WriteString(msg.String())

		if got, want := bf.String(), test.expected; got != want {
			t.Errorf("got: %#v; wanted %#v", got, want)
		}
	}

}
*/