package audio

import (
	"github.com/lucasb-eyer/go-colorful"
)

// Gradient contains the "keypoints" of the colour gradient you want to generate.
// The position of each keypoint has to live in the range [0,1]
type Gradient []struct {
	Col colorful.Color `json:"colour"`
	Pos float64        `json:"position"`
}

func DefaultGradient() Gradient {
	return Gradient{{Col: colorful.Color{}, Pos: 1.0}}
}

// MustParseHex ensures hex strings can be parsed into colorful.Color and
// panics otherwise, useful for creating custom Gradient colours and ensuring
// they are valid at runtime
func MustParseHex(s string) colorful.Color {
	c, err := colorful.Hex(s)
	if err != nil {
		panic("MustParseHex: " + err.Error())
	}
	return c
}
