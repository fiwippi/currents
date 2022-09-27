package session

import (
	"encoding/binary"

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

			buf := make([]byte, 8)
			_, err = s.port.Read(buf)
			if err != nil {
				continue
			}
			log.Debug().Bytes("buf", buf).Uint8("val", buf[0]).Msg("arduino buffer")
		}
	}()

	return nil
}

func (s *Server) Disconnect() error {
	defer func() {
		s.port = nil
	}()

	if s.port != nil {
		// TODO better way to close and stop colours? wait for an ack with timeout?
		s.SendColour(colorful.Color{R: 0, G: 0, B: 0})
		s.SendColour(colorful.Color{R: 0, G: 0, B: 0})
		s.SendColour(colorful.Color{R: 0, G: 0, B: 0})
		return s.port.Close()
	}
	return nil
}

// TODO mutexes with the ports

func (s *Server) SendColour(clr colorful.Color) {
	if s.port == nil {
		return
	}

	var packed uint64
	packed += 0xAA << 56
	packed += 0xBB << 48
	packed += uint64(packColour(clr)) << 16
	packed += 0xCC << 8
	packed += 0xDD

	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, packed)
	s.port.Write(buf)
}

func packColour(clr colorful.Color) uint32 {
	r, g, b := clr.RGB255()

	var packed uint32
	packed += uint32(r) << 16
	packed += uint32(g) << 8
	packed += uint32(b)

	return packed
}
