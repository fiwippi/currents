package audio

import (
	"bytes"
	"fmt"
	"testing"
	"time"
)

func TestExampleBasic(t *testing.T) {
	//  1. Browse available devices
	//  2. Start capture for a device
	//  3. Stop capture for a device
	//  4. ability to change capture device
	//  5. Destroy the context of the audio device

	// 1.
	a, err := NewAudio()
	if err != nil {
		panic(err)
	}

	devices, err := a.Devices()
	if err != nil {
		panic(err)
	}

	device := devices[0]
	fmt.Printf("SIMPLE Chosen device: %s\n", device.Name)

	// 2.
	buf := bytes.NewBuffer(nil)
	err = a.StartCapture(device, buf, DefaultConfig())
	fmt.Println("SIMPLE Started capture 1")
	if err != nil {
		panic(err)
	}

	// 3.
	go func() {
		time.Sleep(3 * time.Second)
		a.StopCapture()
	}()

	err = <-a.Done
	if err != nil {
		panic(err)
	}
	fmt.Println("SIMPLE Ended capture 1")

	// 4.
	err = a.StartCapture(device, buf, DefaultConfig())
	if err != nil {
		panic(err)
	}
	fmt.Println("SIMPLE Started capture 2")

	go func() {
		time.Sleep(3 * time.Second)
		err = a.StartCapture(devices[1], buf, DefaultConfig())
		if err != nil {
			panic(err)
		}
		fmt.Println("SIMPLE Device changed 2")
		time.Sleep(3 * time.Second)
		a.StopCapture()
	}()

	err = <-a.Done
	if err != nil {
		panic(err)
	}
	fmt.Println("SIMPLE Ended capture 2")

	// 5.
	err = a.Destroy()
	if err != nil {
		panic(err)
	}
	fmt.Println("SIMPLE Context destroyed")
}

func TestExampleComplex(t *testing.T) {
	a, err := NewAudio()
	if err != nil {
		t.Error(err)
	}

	devices, err := a.Devices()
	if err != nil {
		t.Error(err)
	}
	device := devices[1]
	fmt.Printf("COMPLEX Chosen device: %s\n", device.Name)

	w, err := NewFFT(DefaultConfig())
	if err != nil {
		t.Error(err)
	}
	err = a.StartCapture(device, w, DefaultConfig())
	fmt.Println("COMPLEX Started capture 1")
	if err != nil {
		t.Error(err)
	}

	go func() {
		time.Sleep(10 * time.Second)
		a.StopCapture()
	}()

	// 3.
loop:
	for {
		select {
		case err := <-a.Done:
			w.Stop()
			if err != nil {
				panic(err)
			}
			fmt.Println("COMPLEX Ended capture 1")
			break loop
		case hue := <-w.Hues:
			fmt.Println(hue.Hex())
		}
	}

	// 5.
	err = a.Destroy()
	if err != nil {
		panic(err)
	}
	fmt.Println("COMPLEX Context destroyed")
}
