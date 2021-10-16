package main

import (
	"currents/pkg/gui"
)

// TODO better instructions on how to setup
// TODO interpolation with bezier curve
// TODO stop flashes with no colour (i.e. use a different code for shutdown) and ignore black colours on the arduino's end

func main() {
	gui.Run()
}
