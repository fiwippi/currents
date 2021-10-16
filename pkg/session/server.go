package session

import (
	"encoding/binary"
	"fmt"

	"github.com/lucasb-eyer/go-colorful"
	"go.bug.st/serial"

	"currents/internal/log"
)

type Server struct {
	port serial.Port
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Running() bool {
	return s.port != nil
}

func (s *Server) Connect(port string) error {
	temp, err := serial.Open(port, &serial.Mode{BaudRate: 9600})
	if err != nil {
		return err
	}
	s.port = temp

	go func() {
		for {
			if s.port == nil {
				return
			}

			buf := make([]byte, 4)
			_, err = s.port.Read(buf)
			if err != nil {
				continue
			}
			log.Trace().Bytes("buf", buf).Uint32("val", binary.BigEndian.Uint32(buf)).Msg("arduino buffer")
		}
	}()

	return nil
}

func (s *Server) Disconnect() error {
	defer func() {
		s.port = nil
	}()

	if s.port != nil {
		s.SendColour(colorful.Color{R: 0, G: 0, B: 0})
		return s.port.Close()
	}
	return nil
}

func (s *Server) SendColour(clr colorful.Color) {
	// Get components in range 0-255
	r, g, b := clr.RGB255()

	// Send over the serial port
	s.SendString("-1,") // Signal start of new colour
	s.SendString(fmt.Sprintf("%d,", r))
	s.SendString(fmt.Sprintf("%d,", g))
	s.SendString(fmt.Sprintf("%d,", b))
}

func (s *Server) SendString(data string) error {
	if s.port == nil {
		return nil
	}

	_, err := s.port.Write([]byte(data))
	return err
}
