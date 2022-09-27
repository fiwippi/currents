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
	return &Server{port: fakePort{}}
}

func (s *Server) Connect(port string) error {
	temp, err := serial.Open(port, &serial.Mode{BaudRate: 9600})
	if err != nil {
		return err
	}
	if temp == nil {
		return fmt.Errorf("returned port is nil")
	}
	s.port = temp

	go func() {
		for {
			if s.port == (fakePort{}) {
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
		s.port = fakePort{}
	}()

	// TODO better way to close and stop colours? wait for an ack with timeout?
	s.SendColour(colorful.Color{R: 0, G: 0, B: 0})
	s.SendColour(colorful.Color{R: 0, G: 0, B: 0})
	s.SendColour(colorful.Color{R: 0, G: 0, B: 0})
	return s.port.Close()
}

func (s *Server) SendColour(clr colorful.Color) {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, markColour(packColour(clr)))
	s.port.Write(buf)
}

func markColour(clr uint32) uint64 {
	var marked uint64
	marked += 0xAA << 56
	marked += 0xBB << 48
	marked += uint64(clr) << 16
	marked += 0xCC << 8
	marked += 0xDD

	return marked
}

func packColour(clr colorful.Color) uint32 {
	r, g, b := clr.RGB255()

	var packed uint32
	packed += uint32(r) << 16
	packed += uint32(g) << 8
	packed += uint32(b)

	return packed
}
