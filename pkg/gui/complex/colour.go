package complex

import (
	"image/color"

	"github.com/lucasb-eyer/go-colorful"
)

func convertNRGBA(c colorful.Color) color.NRGBA {
	r, g, b := c.RGB255()
	return color.NRGBA{R: r, G: g, B: b, A: 0xFF}
}
