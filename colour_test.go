package chroma

import (
	"math"
	"testing"

	assert "github.com/alecthomas/assert/v2"
)

func TestColourRGB(t *testing.T) {
	colour := ParseColour("#8913af")
	assert.Equal(t, uint8(0x89), colour.Red())
	assert.Equal(t, uint8(0x13), colour.Green())
	assert.Equal(t, uint8(0xaf), colour.Blue())
}

func TestColourString(t *testing.T) {
	assert.Equal(t, "#8913af", ParseColour("#8913af").String())
}

func distance(a, b uint8) uint8 {
	if a < b {
		return b - a
	}
	return a - b
}

func TestColourBrighten(t *testing.T) {
	actual := NewColour(128, 128, 128).Brighten(0.5)
	// Closeish to what we expect is fine.
	assert.True(t, distance(192, actual.Red()) <= 2)
	assert.True(t, distance(192, actual.Blue()) <= 2)
	assert.True(t, distance(192, actual.Green()) <= 2)
	actual = NewColour(128, 128, 128).Brighten(-0.5)
	assert.True(t, distance(65, actual.Red()) <= 2)
	assert.True(t, distance(65, actual.Blue()) <= 2)
	assert.True(t, distance(65, actual.Green()) <= 2)
}

func TestColourBrightess(t *testing.T) {
	actual := NewColour(128, 128, 128).Brightness()
	assert.True(t, distance(128, uint8(actual*255.0)) <= 2)
}

// hue returns c's hue. See https://stackoverflow.com/a/23094494.
func hue(c Colour) float64 {
	r := float64(c.Red()) / 255
	g := float64(c.Green()) / 255
	b := float64(c.Blue()) / 255

	min := math.Min(math.Min(r, g), b)
	max := math.Max(math.Max(r, g), b)

	switch {
	case r == min:
		return (g - b) / (max - min)
	case g == min:
		return 2 + (b-r)/(max-min)
	default:
		return 4 + (r-g)/(max-min)
	}
}

func TestColourClampBrightness(t *testing.T) {
	// Start with a colour with a brightness close to 0.5.
	initial := NewColour(0, 128, 255)
	br := initial.Brightness()
	assertInDelta(t, 0.5, br)

	// Passing a range that includes the colour's brightness should be a no-op.
	assert.Equal(t, initial.String(), initial.ClampBrightness(br-0.01, br+0.01).String())

	// Clamping to [0, 0] or [1, 1] should produce black or white, respectively.
	assert.Equal(t, "#000000", initial.ClampBrightness(0, 0).String())
	assert.Equal(t, "#ffffff", initial.ClampBrightness(1, 1).String())

	// Clamping to a brighter or darker range should produce the requested
	// brightness while preserving the colour's hue.
	brighter := initial.ClampBrightness(0.75, 1)
	assertInDelta(t, 0.75, brighter.Brightness())
	assertInDelta(t, hue(initial), hue(brighter))

	darker := initial.ClampBrightness(0, 0.25)
	assertInDelta(t, 0.25, darker.Brightness())
	assertInDelta(t, hue(initial), hue(darker))
}

func assertInDelta(t *testing.T, expected, actual float64) {
	const delta = 0.01 // used for brightness and hue comparisons
	assert.True(t, actual > (expected-delta) && actual < (expected+delta))
}

func TestNewTermColour(t *testing.T) {
	c := NewTermColour(0)
	assert.True(t, c.IsSet())
	assert.True(t, c.IsTermColour())
	assert.Equal(t, uint8(0), c.TermIndex())

	c = NewTermColour(255)
	assert.True(t, c.IsSet())
	assert.True(t, c.IsTermColour())
	assert.Equal(t, uint8(255), c.TermIndex())

	c = NewTermColour(3)
	assert.Equal(t, uint8(3), c.TermIndex())
}

func TestTermColourIsNotTermColour(t *testing.T) {
	c := ParseColour("#ff0000")
	assert.False(t, c.IsTermColour())

	c = Colour(0)
	assert.False(t, c.IsTermColour())
}

func TestTermColourPaletteRGB(t *testing.T) {
	// Index 0 = black
	c := NewTermColour(0)
	assert.Equal(t, uint8(0), c.Red())
	assert.Equal(t, uint8(0), c.Green())
	assert.Equal(t, uint8(0), c.Blue())

	// Index 1 = red (maroon)
	c = NewTermColour(1)
	assert.Equal(t, uint8(0x80), c.Red())
	assert.Equal(t, uint8(0), c.Green())
	assert.Equal(t, uint8(0), c.Blue())

	// Index 15 = bright white
	c = NewTermColour(15)
	assert.Equal(t, uint8(0xff), c.Red())
	assert.Equal(t, uint8(0xff), c.Green())
	assert.Equal(t, uint8(0xff), c.Blue())

	// Index 232 = greyscale first entry (#080808)
	c = NewTermColour(232)
	assert.Equal(t, uint8(0x08), c.Red())
	assert.Equal(t, uint8(0x08), c.Green())
	assert.Equal(t, uint8(0x08), c.Blue())
}

func TestTermColourString(t *testing.T) {
	assert.Equal(t, "term-0", NewTermColour(0).String())
	assert.Equal(t, "term-3", NewTermColour(3).String())
	assert.Equal(t, "term-255", NewTermColour(255).String())
}

func TestTermColourGoString(t *testing.T) {
	assert.Equal(t, "TermColour(0)", NewTermColour(0).GoString())
	assert.Equal(t, "TermColour(255)", NewTermColour(255).GoString())
}

func TestParseTermColour(t *testing.T) {
	c := ParseColour("term-0")
	assert.True(t, c.IsTermColour())
	assert.Equal(t, uint8(0), c.TermIndex())

	c = ParseColour("term-3")
	assert.True(t, c.IsTermColour())
	assert.Equal(t, uint8(3), c.TermIndex())

	c = ParseColour("term-255")
	assert.True(t, c.IsTermColour())
	assert.Equal(t, uint8(255), c.TermIndex())

	// Invalid indices
	c = ParseColour("term-256")
	assert.False(t, c.IsSet())

	c = ParseColour("term-")
	assert.False(t, c.IsSet())

	c = ParseColour("term-abc")
	assert.False(t, c.IsSet())
}

func TestParseTermColourRoundTrip(t *testing.T) {
	for i := 0; i < 256; i++ {
		original := NewTermColour(uint8(i))
		parsed := ParseColour(original.String())
		assert.Equal(t, original, parsed)
	}
}

func TestTermColourBrightness(t *testing.T) {
	c := NewTermColour(15) // white
	assert.True(t, c.Brightness() > 0.99)

	c = NewTermColour(0) // black
	assert.True(t, c.Brightness() < 0.01)
}

func TestTermColourBrightenMaterialisesToRGB(t *testing.T) {
	c := NewTermColour(1) // red/maroon
	brightened := c.Brighten(0.5)
	assert.False(t, brightened.IsTermColour())
	assert.True(t, brightened.IsSet())
}
