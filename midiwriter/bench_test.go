package midiwriter

import (
	"testing"

	"gitlab.com/gomidi/midi/midimessage/channel"
)

type writeNothing struct{}

func (w writeNothing) Write([]byte) (i int, err error) {
	return
}

// BenchmarkNoteOnOffSameChannel1000 writes 1000 channel messages per iteration
// which are noteon or noteoff messages on the same channel.
// running status is used, since they all have the same status byte
func BenchmarkNoteOnOffSameChannel1000(b *testing.B) {
	b.StopTimer()

	wr := New(writeNothing{})

	var (
		m1 = channel.Channel1.NoteOn(20, 100)
		m2 = channel.Channel1.NoteOff(20)
		m3 = channel.Channel1.NoteOn(23, 70)
		m4 = channel.Channel1.NoteOff(23)
	)

	var err error

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		for j := 0; j < 250; j++ {
			err = wr.Write(m1)
			if err != nil {
				b.Fatalf("Error: %v", err)
			}
			err = wr.Write(m2)
			if err != nil {
				b.Fatalf("Error: %v", err)
			}
			err = wr.Write(m3)
			if err != nil {
				b.Fatalf("Error: %v", err)
			}
			err = wr.Write(m4)
			if err != nil {
				b.Fatalf("Error: %v", err)
			}
		}
	}

}

// BenchmarkNoteOnOffAlternatingChannel1000 writes 1000 channel messages per iteration
// which are noteon or noteoff messages alternating on different channels.
// therefor running status can't be used, although it tries to
func BenchmarkNoteOnOffAlternatingChannel1000(b *testing.B) {
	b.StopTimer()

	wr := New(writeNothing{})

	var (
		m1 = channel.Channel1.NoteOn(20, 100)
		m2 = channel.Channel4.NoteOn(23, 70)
		m3 = channel.Channel1.NoteOff(20)
		m4 = channel.Channel4.NoteOff(23)
	)

	var err error

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		for j := 0; j < 250; j++ {
			err = wr.Write(m1)
			if err != nil {
				b.Fatalf("Error: %v", err)
			}
			err = wr.Write(m2)
			if err != nil {
				b.Fatalf("Error: %v", err)
			}
			err = wr.Write(m3)
			if err != nil {
				b.Fatalf("Error: %v", err)
			}
			err = wr.Write(m4)
			if err != nil {
				b.Fatalf("Error: %v", err)
			}
		}
	}

}

// BenchmarkNoteOnOffRunningStatusDisabled1000 writes 1000 channel messages per iteration
// which are noteon or noteoff messages alternating on different channels.
// running status is disabled as option
func BenchmarkNoteOnOffRunningStatusDisabled1000(b *testing.B) {
	b.StopTimer()

	wr := New(writeNothing{}, NoRunningStatus())

	var (
		m1 = channel.Channel1.NoteOn(20, 100)
		m2 = channel.Channel4.NoteOn(23, 70)
		m3 = channel.Channel1.NoteOff(20)
		m4 = channel.Channel4.NoteOff(23)
	)

	var err error

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		for j := 0; j < 250; j++ {
			err = wr.Write(m1)
			if err != nil {
				b.Fatalf("Error: %v", err)
			}
			err = wr.Write(m2)
			if err != nil {
				b.Fatalf("Error: %v", err)
			}
			err = wr.Write(m3)
			if err != nil {
				b.Fatalf("Error: %v", err)
			}
			err = wr.Write(m4)
			if err != nil {
				b.Fatalf("Error: %v", err)
			}
		}
	}

}
