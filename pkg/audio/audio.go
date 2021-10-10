package audio

import (
	"errors"
	"io"
	"sync"

	"github.com/gen2brain/malgo"
)

var errDeviceChange = errors.New("the device has changed") // Signifies the capture func should not exit
var sampleSizeInBytes = malgo.SampleSizeInBytes(malgo.FormatF32)

type Audio struct {
	*sync.Mutex

	// Malgo context
	context *malgo.AllocatedContext

	// Capture context
	Done        chan error
	initialised bool
	device      *malgo.Device
	abortChan   chan error
}

// Creation / Deletion ---------------------------------------------

// NewAudio creates a new Audio struct and allocates the context for it.
func NewAudio() (*Audio, error) {
	c, err := malgo.InitContext(nil, malgo.ContextConfig{}, nil)
	if err != nil {
		return nil, err
	}

	a := &Audio{
		Mutex:     &sync.Mutex{},
		context:   c,
		abortChan: make(chan error),
		Done:      make(chan error),
	}
	go a.drainDoneChan()

	return a, nil
}

// MustCreateNewAudio panics if the audio was created unsuccessfully
func MustCreateNewAudio() *Audio {
	a, err := NewAudio()
	if err != nil {
		panic(err)
	}
	return a
}

// Destroy un-initialises the context and then frees it
func (a *Audio) Destroy() error {
	a.Lock()
	defer a.Unlock()

	err := a.context.Uninit()
	if err != nil {
		return err
	}
	a.context.Free()
	return nil
}

// Info ---------------------------------------------

// Devices returns a list of malgo.Capture devices which can be used
func (a *Audio) Devices() ([]Device, error) {
	a.Lock()
	defer a.Unlock()

	infos, err := a.context.Devices(malgo.Capture)
	if err != nil {
		return nil, err
	}

	devices := make([]Device, len(infos))
	for i := range infos {
		full, err := a.context.DeviceInfo(malgo.Capture, infos[i].ID, malgo.Shared)
		if err != nil {
			return nil, err
		}
		devices[i] = newDevice(full)
	}

	return devices, nil
}

func (a *Audio) MustParseDevices() []Device {
	d, err := a.Devices()
	if err != nil {
		panic(err)
	}
	return d
}

// Capture State ---------------------------------------------

func (a *Audio) StartCapture(d Device, w io.Writer, conf *Config) error {
	a.stopCapture(errDeviceChange)

	a.Lock()
	defer a.Unlock()

	// Create channels to abort the process
	aborted := false

	// Writes the data to the writer
	deviceConfig := d.config(conf)
	deviceCallbacks := malgo.DeviceCallbacks{
		Data: func(outputSamples, inputSamples []byte, frameCount uint32) {
			if aborted {
				return
			}

			_, err := w.Write(inputSamples)
			if err != nil {
				aborted = true
				a.abortChan <- err
			}
		},
	}

	// Stream the data to the writer
	device, err := malgo.InitDevice(a.context.Context, deviceConfig, deviceCallbacks)
	if err != nil {
		return err
	}
	a.initialised = true
	a.device = device

	err = device.Start()
	if err != nil {
		return err
	}

	return nil
}

func (a *Audio) StopCapture() {
	a.stopCapture(nil)
}

func (a *Audio) stopCapture(err error) {
	a.Lock()
	if a.initialised {
		a.device.Uninit()
	}
	a.initialised = false
	a.Unlock()
	a.abortChan <- err
}

func (a *Audio) drainDoneChan() {
	for e := range a.abortChan {
		if e == errDeviceChange {
			continue
		}
		go func(e error) { a.Done <- e }(e)
	}
}
