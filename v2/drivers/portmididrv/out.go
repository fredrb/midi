package portmididrv

import (
	"fmt"

	"sync"

	"gitlab.com/gomidi/midi/v2/drivers"
	"gitlab.com/gomidi/midi/v2/drivers/portmididrv/imported/portmidi"
)

func newOut(driver *Driver, deviceid portmidi.DeviceID, id int, name string) drivers.Out {
	return &out{driver: driver, id: id, name: name, deviceid: deviceid}
}

type out struct {
	deviceid portmidi.DeviceID
	id       int
	stream   *portmidi.Stream
	name     string
	mx       sync.RWMutex
	driver   *Driver
}

// IsOpen returns, wether the port is open
func (o *out) IsOpen() bool {
	o.mx.RLock()
	defer o.mx.RUnlock()
	return o.stream != nil
}

// Send writes a MIDI sysex message to the outut port
func (o *out) SendSysEx(data []byte) error {
	o.mx.RLock()
	if o.stream == nil {
		o.mx.RUnlock()
		return drivers.ErrPortClosed
	}
	o.mx.RUnlock()

	// since we always open the outputstream with a latency of 0
	// the timestamp is ignored
	var ts portmidi.Timestamp // or portmidi.Time()

	o.mx.Lock()
	defer o.mx.Unlock()
	//fmt.Printf("sending sysex % X\n", data)
	err := o.stream.WriteSysExBytes(ts, data)
	if err != nil {
		return fmt.Errorf("could not send message to MIDI out %v (%s): %v", o.Number(), o, err)
	}
	return nil
}

// Send writes a MIDI message to the outut port
// If the output port is closed, it returns midi.ErrPortClosed
func (o *out) Send(b []byte) error {
	o.mx.RLock()
	if o.stream == nil {
		o.mx.RUnlock()
		return drivers.ErrPortClosed
	}
	o.mx.RUnlock()

	/*
		if len(b) < 2 {
			return fmt.Errorf("cannot send less than two message bytes")
		}
	*/

	first := int64(b[0])

	var second int64
	if len(b) > 1 {
		second = int64(b[1])
	}

	var last int64
	// ProgramChange messages only have 2 bytes
	if len(b) > 2 {
		last = int64(b[2])
	}

	//	fmt.Printf("sending % X\n", b)

	o.mx.Lock()
	defer o.mx.Unlock()
	//o.driver.Lock()
	//defer o.driver.Unlock()
	err := o.stream.WriteShort(first, second, last)
	if err != nil {
		return fmt.Errorf("could not send message to MIDI out %v (%s): %v", o.Number(), o, err)
	}
	return nil
}

// Underlying returns the underlying *portmidi.Stream. It will be nil, if the port is closed.
// Use it with type casting:
//   portOut := o.Underlying().(*portmidi.Stream)
func (o *out) Underlying() interface{} {
	return o.stream
}

// Number returns the number of the MIDI out port.
// Since portmidis ports counting is confusing (out and in ports are counted together),
// we do our own counting.
func (o *out) Number() int {
	return o.id
}

// String returns the name of the MIDI out port.
func (o *out) String() string {
	return o.name
}

// Close closes the MIDI out port
func (o *out) Close() error {
	o.mx.RLock()
	if o.stream == nil {
		o.mx.RUnlock()
		return nil
	}
	o.mx.RUnlock()

	o.mx.Lock()
	defer o.mx.Unlock()
	err := o.stream.Close()
	if err != nil {
		return fmt.Errorf("can't close MIDI out %v (%s): %v", o.Number(), o, err)
	}
	o.stream = nil
	return nil
}

// Open opens the MIDI output port
func (o *out) Open() (err error) {
	o.mx.RLock()
	if o.stream != nil {
		o.mx.RUnlock()
		return nil
	}
	o.mx.RUnlock()

	o.mx.Lock()
	defer o.mx.Unlock()
	// we always open the outputstream with a latency of 0
	var latency int64
	o.stream, err = portmidi.NewOutputStream(o.deviceid, o.driver.buffersizeOut, latency)
	if err != nil {
		o.stream = nil
		return fmt.Errorf("can't open MIDI out port %v (%s): %v", o.Number(), o, err)
	}
	o.driver.Lock()
	defer o.driver.Unlock()
	o.driver.opened = append(o.driver.opened, o)
	return nil
}
