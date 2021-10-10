package audio

import (
	"errors"
	"math"

	"github.com/lucasb-eyer/go-colorful"
)

var ErrInvalidInterpMode = errors.New("interpolation mode specified is invalid")

type InterpolateMode int

const (
	// Blended returns a HCL-blend between the two colors around `t`,
	// Note: It relies heavily on the fact that the gradient keypoints are sorted
	Blended InterpolateMode = iota
	// Blocky return the colour nearest to t instead of blending
	// colours together
	Blocky
)

func (im InterpolateMode) String() string {
	return [...]string{"Blended", "Blocky"}[im]
}

func (im InterpolateMode) Interpolate(t float64, g Gradient) colorful.Color {
	if len(g) == 0 {
		return colorful.Color{}
	}

	switch im {
	case Blended:
		for i := 0; i < len(g)-1; i++ {
			c1 := g[i]
			c2 := g[i+1]
			if c1.Pos <= t && t <= c2.Pos {
				// We are in between c1 and c2. Go blend them!
				t := (t - c1.Pos) / (c2.Pos - c1.Pos)
				return c1.Col.BlendHcl(c2.Col, t).Clamped()
			}
		}

		// Nothing found? Means we're at (or past) the last gradient keypoint.
		return g[len(g)-1].Col
	case Blocky:
		var index int
		var min float64 = 1.0

		for i, c := range g {
			dist := math.Abs(c.Pos - t)
			if dist < min {
				min = dist
				index = i
			}
		}

		return g[index].Col
	}

	panic(ErrInvalidInterpMode)
}
