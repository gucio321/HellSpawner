package hsutil

import "image/color"

// Color converts an rgba uint32 to a colorEnabled.RGBA
func Color(rgba uint32) color.RGBA {
	const (
		a, b, g, r = 0, 1, 2, 3
		byteWidth  = 8
		byteMask   = 0xff
	)

	//nolint:gosec // this uses bitmask so these itagers will never actually overflow - this is intended.
	result := color.RGBA{
		R: uint8((rgba >> (r * byteWidth)) & byteMask),
		G: uint8((rgba >> (g * byteWidth)) & byteMask),
		B: uint8((rgba >> (b * byteWidth)) & byteMask),
		A: uint8((rgba >> (a * byteWidth)) & byteMask),
	}

	return result
}
