package sdui

import (
	"fmt"
	"strconv"
)

type RGB struct {
	R, G, B uint8
}

// PaletteFromSeed parses a hex seed color and returns CSS custom properties
// for Screen.props.style.
func PaletteFromSeed(hex string) (map[string]string, error) {
	c, err := parseRGB(hex)
	if err != nil {
		return nil, err
	}

	return map[string]string{
		"--color-primary": hex,
		"--color-light":   blendRGB(c, true, 0.3),
		"--color-dark":    blendRGB(c, false, 0.3),
	}, nil
}

func parseRGB(hex string) (RGB, error) {
	if len(hex) != 7 || hex[0] != '#' {
		return RGB{}, fmt.Errorf("sdui: invalid hex color %q", hex)
	}

	r, errR := strconv.ParseUint(hex[1:3], 16, 8)
	g, errG := strconv.ParseUint(hex[3:5], 16, 8)
	b, errB := strconv.ParseUint(hex[5:7], 16, 8)
	if errR != nil || errG != nil || errB != nil {
		return RGB{}, fmt.Errorf("sdui: invalid hex characters in %q", hex)
	}

	return RGB{R: uint8(r), G: uint8(g), B: uint8(b)}, nil
}

// blendRGB mixes rgb toward white or black by factor (0–1) and returns #RRGGBB.
func blendRGB(c RGB, towardWhite bool, factor float64) string {
	return fmt.Sprintf("#%02X%02X%02X",
		blendChannel(c.R, towardWhite, factor),
		blendChannel(c.G, towardWhite, factor),
		blendChannel(c.B, towardWhite, factor),
	)
}

func blendChannel(channel uint8, towardWhite bool, factor float64) uint8 {
	var blended float64
	if towardWhite {
		blended = float64(channel) + (255-float64(channel))*factor
	} else {
		blended = float64(channel) * (1 - factor)
	}
	return uint8(blended)
}
