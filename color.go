package lights

import (
	"fmt"
	"image/color"
	"strconv"
)

// ParseColorCode parses a color code in either #RGB or #RRGGBB formats
// following CSS color specifications. It returns the color found or
// an error if the color is not a valid color code.
func ParseColorCode(colorCode string) (color.Color, error) {
	switch len(colorCode) {
	case 4:
		r, err := strconv.ParseInt(string(colorCode[1])+string(colorCode[1]), 16, 0)
		if err != nil {
			return nil, err
		}
		g, err := strconv.ParseInt(string(colorCode[2])+string(colorCode[2]), 16, 0)
		if err != nil {
			return nil, err
		}
		b, err := strconv.ParseInt(string(colorCode[3])+string(colorCode[3]), 16, 0)
		if err != nil {
			return nil, err
		}
		// Should be uint8s, just cast from int64
		return color.RGBA{uint8(r), uint8(g), uint8(b), uint8(0)}, nil
	case 7:
		r, err := strconv.ParseInt(string(colorCode[1])+string(colorCode[2]), 16, 0)
		if err != nil {
			return nil, err
		}
		g, err := strconv.ParseInt(string(colorCode[3])+string(colorCode[4]), 16, 0)
		if err != nil {
			return nil, err
		}
		b, err := strconv.ParseInt(string(colorCode[5])+string(colorCode[6]), 16, 0)
		if err != nil {
			return nil, err
		}
		// Should be uint8s, just cast from int64
		return color.RGBA{uint8(r), uint8(g), uint8(b), uint8(0)}, nil
	default:
		return nil, fmt.Errorf("Color code must be 4 or 6 characters - found %d", len(colorCode))
	}
}
