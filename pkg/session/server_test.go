package session

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/lucasb-eyer/go-colorful"
	"github.com/stretchr/testify/assert"
)

var (
	white = colorful.LinearRgb(1, 1, 1)
	red   = colorful.LinearRgb(1, 0, 0)
	green = colorful.LinearRgb(0, 1, 0)
	blue  = colorful.LinearRgb(0, 0, 1)
)

func showBinary(n uint32) string {
	return fmt.Sprintf("%032s", strconv.FormatInt(int64(n), 2))
}

func TestPackColour(t *testing.T) {
	assert.Equal(t, "00000000111111111111111111111111", showBinary(packColour(white)))
	assert.Equal(t, "00000000111111110000000000000000", showBinary(packColour(red)))
	assert.Equal(t, "00000000000000001111111100000000", showBinary(packColour(green)))
	assert.Equal(t, "00000000000000000000000011111111", showBinary(packColour(blue)))
}
