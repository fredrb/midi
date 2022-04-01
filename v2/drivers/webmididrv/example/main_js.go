// +build js,wasm

package main

import (
	"bytes"
	"fmt"
	"syscall/js"
	"time"

	"gitlab.com/gomidi/midi/v2"
	_ "gitlab.com/gomidi/midi/v2/drivers/webmididrv"
)

/*
to build, run

GOOS=js GOARCH=wasm go build -o main.wasm main_js.go
*/

func log(message string) {
	document := js.Global().Get("document")
	p := document.Call("createElement", "p")
	p.Set("innerHTML", message)
	document.Get("body").Call("appendChild", p)
}

func main() {
	defer midi.CloseDriver()
	var bf bytes.Buffer

	for i, in := range midi.InPorts() {
		fmt.Fprintf(&bf, "found MIDI in port: %v: %s<br />", i, in)
	}

	fmt.Fprintf(&bf, "<br><br>")

	for i, out := range midi.OutPorts() {
		fmt.Fprintf(&bf, "found MIDI out port: %v: %s<br />", i, out)
	}

	log(bf.String())

	stop, err := midi.ListenTo(0, midi.ReceiverFunc(func(msg midi.Message, timestamp int32) {
		log(fmt.Sprintf("got: %s<br />", msg))
	}))
	e(err)

	s, err := midi.SendTo(0)
	e(err)

	log(fmt.Sprintf("send: NoteOn key: %v veloctiy: %v on channel %v<br />", 60, 120, 3))

	// do some writing: if you are using a loopback midi device on your os, you will see
	// this messages in the browser window
	s.Send(midi.NoteOn(3, 60, 120))
	time.Sleep(time.Second)
	log(fmt.Sprintf("send: NoteOff key: %v on channel %v<br />", 60, 3))
	s.Send(midi.NoteOff(3, 60))

	qsynth := midi.FindOutPort("qsynth")

	if qsynth >= 0 {
		qs, err := midi.SendTo(qsynth)
		e(err)

		qs.Send(midi.NoteOn(0, 60, 120))
		time.Sleep(time.Millisecond * 500)
		qs.Send(midi.NoteOff(0, 60))
	}

	// stay alive
	ch := make(chan bool)
	<-ch
	stop()
}

func e(err error) {
	if err != nil {
		log(fmt.Sprintf("ERROR: %s", err.Error()))
	}
}
