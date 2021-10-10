package audio

import (
	"bytes"

	"github.com/gen2brain/malgo"
)

type Device struct {
	Name string
	id   malgo.DeviceID
}

func newDevice(di malgo.DeviceInfo) Device {
	nameBytes := []byte(di.Name())
	nameTrimmed := bytes.Trim(nameBytes, "\x00")

	return Device{
		Name: string(nameTrimmed),
		id:   di.ID,
	}
}

func (d *Device) config(conf *Config) malgo.DeviceConfig {
	deviceConfig := malgo.DefaultDeviceConfig(malgo.Capture)
	deviceConfig.Capture.DeviceID = d.id.Pointer()
	deviceConfig.Capture.Format = malgo.FormatF32
	deviceConfig.Capture.Channels = conf.Channels
	deviceConfig.SampleRate = conf.SampleRate
	deviceConfig.PerformanceProfile = malgo.LowLatency
	return deviceConfig
}
