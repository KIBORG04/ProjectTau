package chartjs

import (
	"fmt"
	"image/color"
	"ssstatistics/internal/utils"
)

type RGBA color.RGBA

// MarshalJSON satisfies the json.Marshaller interface.
func (c *RGBA) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"rgba(%d, %d, %d, %.3f)\"", c.R, c.G, c.B, float64(c.A)/255)), nil
}

func RandomColor(count int) []*RGBA {
	colors := make([]*RGBA, 0, count)
	for i := 0; i < count; i++ {
		colors = append(colors, utils.Pick(NiceColors))
	}
	return colors
}

func ObcurceHex(color RGBA, coeff float32) *RGBA {
	color.R = color.R - uint8(float32(color.R)*coeff)
	color.G = color.G - uint8(float32(color.G)*coeff)
	color.B = color.B - uint8(float32(color.B)*coeff)
	return &color
}

func RandomBorder(backColors []*RGBA) []*RGBA {
	colors := make([]*RGBA, 0, len(backColors))
	for _, v := range backColors {
		rgb := ObcurceHex(*v, 0.2)
		colors = append(colors, rgb)
	}
	return colors
}
