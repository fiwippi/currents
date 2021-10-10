package session

import (
	"encoding/binary"

	"github.com/lucasb-eyer/go-colorful"
	"go.bug.st/serial"

	"currents/internal/log"
)

type Server struct {
	port      serial.Port
	blackOnly bool
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
	s.blackOnly = true
	defer func() {
		s.blackOnly = false
		s.port = nil
	}()

	if s.port != nil {
		s.SendColour(colorful.Color{R: 0, G: 0, B: 0})
		return s.port.Close()
	}
	return nil
}

func (s *Server) SendColour(clr colorful.Color) {
	r, g, b := clr.RGB255()
	s.SendByte(0) // Signal start of new colour
	s.SendByte(r)
	s.SendByte(g)
	s.SendByte(b)
}

func (s *Server) SendByte(data byte) error {
	if s.port == nil || s.blackOnly && data != 0 {
		return nil
	}

	_, err := s.port.Write([]byte{data})
	return err
}
