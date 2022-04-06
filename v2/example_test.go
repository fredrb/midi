package midi_test

import (
	"fmt"
	"os"
	"time"

	. "gitlab.com/gomidi/midi/v2"

	// testdrv has one in port and one out port which is connected to the in port
	_ "gitlab.com/gomidi/midi/v2/drivers/testdrv"
	//_ "gitlab.com/gomidi/midi/v2/drivers/rtmididrv"
	// when using rtmidi ("for real"), replace with the line above
)

func Example() {

	var eachMessage = func(msg Message, timestampms int32) {
		if msg.Is(RealTimeMsg) {
			// ignore realtime messages
			return
		}
		var channel, key, velocity, cc, val uint8
		switch {

		// is better, than to use GetNoteOn (handles note on messages with velocity of 0 as expected)
		case msg.GetNoteStart(&channel, &key, &velocity):
			fmt.Printf("note started at %vms channel: %v key: %v velocity: %v\n", timestampms, channel, key, velocity)

		// is better, than to use GetNoteOff (handles note on messages with velocity of 0 as expected)
		case msg.GetNoteEnd(&channel, &key):
			fmt.Printf("note ended at %vms channel: %v key: %v\n", timestampms, channel, key)

		case msg.GetControlChange(&channel, &cc, &val):
			fmt.Printf("control change %v %q channel: %v value: %v at %vms\n", cc, ControlChangeName[cc], channel, val, timestampms)

		default:
			fmt.Printf("received %s at %vms\n", msg, timestampms)
		}
	}

	// always good to close the driver at the end
	defer CloseDriver()

	// allows you to get the ports when using "real" drivers like rtmididrv or portmididrv
	if len(os.Args) == 2 && os.Args[1] == "list" {
		fmt.Printf("MIDI IN Ports\n")
		for i, port := range InPorts() {
			fmt.Printf("no: %v %q\n", i, port)
		}
		fmt.Printf("\n\nMIDI OUT Ports\n")
		for i, port := range OutPorts() {
			fmt.Printf("no: %v %q\n", i, port)
		}
		fmt.Printf("\n\n")
		return
	}

	var out int = 0
	// takes the first out port, for real, consider
	// var out = OutByName("my synth")

	// creates a sender function to the out port
	send, _ := SendTo(out)

	var in int = 0
	// here we take first in port, for real, consider
	// var in = InByName("my midi keyboard")

	// listens to the in port and calls eachMessage for every message.
	// any running status bytes are converted and only complete messages are passed to the eachMessage.
	stop, _ := ListenTo(in, eachMessage)

	{ // send some messages
		send(NoteOn(0, Db(4), 100))
		time.Sleep(time.Millisecond * 30)
		send(NoteOff(0, Db(4)))
		send(Pitchbend(0, -12))
		time.Sleep(time.Millisecond * 20)
		send(ProgramChange(1, 12))
		send(ControlChange(2, FootPedalMSB, On))
	}

	// stops listening
	stop()

	// Output:
	// note started at 0ms channel: 0 key: 61 velocity: 100
	// note ended at 30ms channel: 0 key: 61
	// received PitchBend channel: 0 pitch: -12 (8180) at 30ms
	// received ProgramChange channel: 1 program: 12 at 50ms
	// control change 4 "Foot Pedal (MSB)" channel: 2 value: 127 at 50ms

}