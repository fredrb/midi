package smf

import (
	"bytes"
	"fmt"

	// "os"
	//"log"
	"testing"
)

func testRead(t *testing.T, input []byte) string {
	var out bytes.Buffer
	out.WriteString("\n")
	smf, err := ReadAll(bytes.NewReader(input))
	if err != nil {
		t.Fatalf("can't read: %v", err)
	}

	/*
		err := rd.ReadHeader()
		if err != nil {
			t.Fatalf("can't read header: %v", err)
		}
	*/

	//hd := rd.Header()

	out.WriteString(fmt.Sprintf("SMF%v\n", smf.Format()))
	out.WriteString(fmt.Sprintf("%v Track(s)\n", smf.NumTracks()))
	out.WriteString(fmt.Sprintf("TimeFormat: %s\n", smf.TimeFormat))

	tracks := smf.Tracks()

	for i, track := range tracks {
		for _, ev := range track.Events {
			//out.WriteString(fmt.Sprintf("Track %v@%v %s\n", i, ev.Delta, ev.MessageType()))
			//m := midi2.NewMessage(ev.Data)
			//m.Type = midi2.GetMessageType(ev.Data)
			out.WriteString(fmt.Sprintf("Track %v@%v %v\n", i, ev.Delta, ev.Message()))
		}
	}
	/*
		var _ = hd
		var msg midi.Message

		for {
			msg, err = rd.Read()

			if err != nil {
				break
			}

			out.WriteString(fmt.Sprintf("Track %v@%v %s\n", rd.Track(), rd.Delta(), msg))
		}
	*/

	return out.String()

}

func TestReadSMF0(t *testing.T) {
	var expected = `
SMF0
1 Track(s)
TimeFormat: 96 MetricTicks
Track 0@0 MetaTimeSigMsg meter: 4/4
Track 0@0 MetaTempoMsg bpm: 120
Track 0@0 Channel0Msg & ProgramChangeMsg program: 5
Track 0@0 Channel1Msg & ProgramChangeMsg program: 46
Track 0@0 Channel2Msg & ProgramChangeMsg program: 70
Track 0@0 Channel2Msg & NoteOnMsg key: 48 velocity: 96
Track 0@0 Channel2Msg & NoteOnMsg key: 60 velocity: 96
Track 0@96 Channel1Msg & NoteOnMsg key: 67 velocity: 64
Track 0@96 Channel0Msg & NoteOnMsg key: 76 velocity: 32
Track 0@192 Channel2Msg & NoteOffMsg key: 48 velocity: 64
Track 0@0 Channel2Msg & NoteOffMsg key: 60 velocity: 64
Track 0@0 Channel1Msg & NoteOffMsg key: 67 velocity: 64
Track 0@0 Channel0Msg & NoteOffMsg key: 76 velocity: 64
Track 0@0 MetaEndOfTrackMsg
`
	//l := log.Default()

	//if got, want := testRead(t, SpecSMF0, Debug(l)), expected; got != want {
	if got, want := testRead(t, SpecSMF0), expected; got != want {
		t.Errorf("got:\n%v\n\nwanted\n%v\n\n", got, want)
	}

}

/*
func TestReadSMF1Missing(t *testing.T) {

	rd := New(bytes.NewReader(examples.SpecSMF1Missing))
	err := rd.ReadHeader()

	for err == nil {
		_, err = rd.Read()
	}

	if err != ErrMissing {
		t.Errorf("expected ErrMissing, got: %#v", err)
	}

}
*/

func TestReadSMF1(t *testing.T) {
	var expected = `
SMF1
4 Track(s)
TimeFormat: 96 MetricTicks
Track 0@0 MetaTimeSigMsg meter: 4/4
Track 0@0 MetaTempoMsg bpm: 120
Track 0@384 MetaEndOfTrackMsg
Track 1@0 Channel0Msg & ProgramChangeMsg program: 5
Track 1@192 Channel0Msg & NoteOnMsg key: 76 velocity: 32
Track 1@192 Channel0Msg & NoteOnMsg key: 76 velocity: 0
Track 1@0 MetaEndOfTrackMsg
Track 2@0 Channel1Msg & ProgramChangeMsg program: 46
Track 2@96 Channel1Msg & NoteOnMsg key: 67 velocity: 64
Track 2@288 Channel1Msg & NoteOnMsg key: 67 velocity: 0
Track 2@0 MetaEndOfTrackMsg
Track 3@0 Channel2Msg & ProgramChangeMsg program: 70
Track 3@0 Channel2Msg & NoteOnMsg key: 48 velocity: 96
Track 3@0 Channel2Msg & NoteOnMsg key: 60 velocity: 96
Track 3@384 Channel2Msg & NoteOnMsg key: 48 velocity: 0
Track 3@0 Channel2Msg & NoteOnMsg key: 60 velocity: 0
Track 3@0 MetaEndOfTrackMsg
`
	if got, want := testRead(t, SpecSMF1), expected; got != want {
		t.Errorf("got:\n%v\n\nwanted\n%v\n\n", got, want)
	}

}

/*
func TestReadSMF1NoteOffPedantic(t *testing.T) {
	var expected = `
SMF1
4 Track(s)
TimeFormat: 96 MetricTicks
Track 0@0 meta.TimeSig 4/4 clocksperclick 24 dsqpq 8
Track 0@0 meta.Tempo BPM: 120.00
Track 0@384 meta.EndOfTrack
Track 1@0 channel.ProgramChange channel 0 program 5
Track 1@192 channel.NoteOn channel 0 key 76 velocity 32
Track 1@192 channel.NoteOff channel 0 key 76
Track 1@0 meta.EndOfTrack
Track 2@0 channel.ProgramChange channel 1 program 46
Track 2@96 channel.NoteOn channel 1 key 67 velocity 64
Track 2@288 channel.NoteOff channel 1 key 67
Track 2@0 meta.EndOfTrack
Track 3@0 channel.ProgramChange channel 2 program 70
Track 3@0 channel.NoteOn channel 2 key 48 velocity 96
Track 3@0 channel.NoteOn channel 2 key 60 velocity 96
Track 3@384 channel.NoteOff channel 2 key 48
Track 3@0 channel.NoteOff channel 2 key 60
Track 3@0 meta.EndOfTrack
`

	if got, want := testRead(t, examples.SpecSMF1, NoteOffVelocity()), expected; got != want {
		t.Errorf("got:\n%v\n\nwanted\n%v\n\n", got, want)
	}

}

func TestReadSMF0NoteOffPedantic(t *testing.T) {
	var expected = `
SMF0
1 Track(s)
TimeFormat: 96 MetricTicks
Track 0@0 meta.TimeSig 4/4 clocksperclick 24 dsqpq 8
Track 0@0 meta.Tempo BPM: 120.00
Track 0@0 channel.ProgramChange channel 0 program 5
Track 0@0 channel.ProgramChange channel 1 program 46
Track 0@0 channel.ProgramChange channel 2 program 70
Track 0@0 channel.NoteOn channel 2 key 48 velocity 96
Track 0@0 channel.NoteOn channel 2 key 60 velocity 96
Track 0@96 channel.NoteOn channel 1 key 67 velocity 64
Track 0@96 channel.NoteOn channel 0 key 76 velocity 32
Track 0@192 channel.NoteOffVelocity channel 2 key 48 velocity 64
Track 0@0 channel.NoteOffVelocity channel 2 key 60 velocity 64
Track 0@0 channel.NoteOffVelocity channel 1 key 67 velocity 64
Track 0@0 channel.NoteOffVelocity channel 0 key 76 velocity 64
Track 0@0 meta.EndOfTrack
`

	if got, want := testRead(t, examples.SpecSMF0, NoteOffVelocity()), expected; got != want {
		t.Errorf("got:\n%v\n\nwanted\n%v\n\n", got, want)
	}

}

func TestReadSysEx(t *testing.T) {
	var bf bytes.Buffer

	wr := smfwriter.New(&bf)
	wr.Write(sysex.Escape(realtime.Start.Raw()))
	wr.SetDelta(0)
	wr.Write(channel.Channel2.NoteOn(65, 90))
	wr.SetDelta(10)
	wr.Write(sysex.SysEx([]byte{0x90, 0x51}))
	wr.SetDelta(1)
	wr.Write(channel.Channel2.NoteOff(65))
	wr.Write(sysex.Start([]byte{0x90, 0x51}))
	wr.SetDelta(5)
	wr.Write(sysex.Continue([]byte{0x90, 0x51}))
	wr.SetDelta(5)
	wr.Write(sysex.End([]byte{0x90, 0x51}))
	wr.Write(meta.EndOfTrack)

	rd := New(bytes.NewReader(bf.Bytes()))

	var m midi.Message
	var err error

	var res bytes.Buffer
	res.WriteString("\n")
	for {
		m, err = rd.Read()

		// breaking at least with io.EOF
		if err != nil {
			break
		}

		switch v := m.(type) {
		case sysex.Escape:
			fmt.Fprintf(&res, "[%v] Sysex Escape: % X\n", rd.Delta(), v.Data())
		case sysex.Start:
			fmt.Fprintf(&res, "[%v] Sysex Start: % X\n", rd.Delta(), v.Data())
		case sysex.Continue:
			fmt.Fprintf(&res, "[%v] Sysex Continue: % X\n", rd.Delta(), v.Data())
		case sysex.End:
			fmt.Fprintf(&res, "[%v] Sysex End: % X\n", rd.Delta(), v.Data())
		case sysex.SysEx:
			fmt.Fprintf(&res, "[%v] Sysex: % X\n", rd.Delta(), v.Data())
		case channel.NoteOn:
			fmt.Fprintf(&res, "[%v] NoteOn at channel %v: key %v velocity %v\n", rd.Delta(), v.Channel(), v.Key(), v.Velocity())
		case channel.NoteOff:
			fmt.Fprintf(&res, "[%v] NoteOff at channel %v: key %v\n", rd.Delta(), v.Channel(), v.Key())
		}

	}

	expected := `
[0] Sysex Escape: FA
[0] NoteOn at channel 2: key 65 velocity 90
[10] Sysex: 90 51
[1] NoteOff at channel 2: key 65
[0] Sysex Start: 90 51
[5] Sysex Continue: 90 51
[5] Sysex End: 90 51
`

	if got, want := res.String(), expected; got != want {
		t.Errorf("got\n%v\n\nwant\n%v\n\n", got, want)
	}

}

func TestX(t *testing.T) {
	src := []byte{0x4D, 0x54, 0x68, 0x64, 0x00, 0x00, 0x00, 0x06, 0x00, 0x00, 0x00, 0x01, 0x03, 0xC0, 0x4D, 0x54, 0x72, 0x6B, 0x00, 0x00, 0x00, 0x0B, 0x00, 0x90, 0x32, 0x21, 0x02, 0x32, 0x00, 0x00, 0xFF, 0x2F, 0x00}
	_ = src

	rd := New(bytes.NewReader(src))

	err := rd.ReadHeader()

	if err != nil {
		t.Fatalf("Error: %s", err.Error())
	}
	_ = rd

	// fmt.Printf("%v\n", rd.Header())

	var msg midi.Message
	msg, err = rd.Read()

	if err != nil {
		t.Fatalf("Error: %s", err.Error())
	}

	_ = msg
	// fmt.Printf("%s\n", msg)
}
*/
