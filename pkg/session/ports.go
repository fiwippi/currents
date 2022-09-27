package session

import (
	"io"
	"time"

	"go.bug.st/serial"
	"go.bug.st/serial/enumerator"
)

func GetAvailablePorts() ([]string, error) {
	portData, err := enumerator.GetDetailedPortsList()
	if err != nil {
		return nil, err
	}

	ports := make([]string, 0)
	for _, port := range portData {
		if port.IsUSB {
			ports = append(ports, port.Name)
		}
	}
	return ports, nil
}

type fakePort struct{}

func (fp fakePort) SetMode(mode *serial.Mode) error                      { return nil }
func (fp fakePort) Read(p []byte) (n int, err error)                     { return 0, io.EOF }
func (fp fakePort) Write(p []byte) (n int, err error)                    { return 0, nil }
func (fp fakePort) ResetInputBuffer() error                              { return nil }
func (fp fakePort) ResetOutputBuffer() error                             { return nil }
func (fp fakePort) SetDTR(dtr bool) error                                { return nil }
func (fp fakePort) SetRTS(rts bool) error                                { return nil }
func (fp fakePort) GetModemStatusBits() (*serial.ModemStatusBits, error) { return nil, nil }
func (fp fakePort) SetReadTimeout(t time.Duration) error                 { return nil }
func (fp fakePort) Close() error                                         { return nil }
