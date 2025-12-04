package generator

import (
	"image"
	"image/color"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

// DrawBorder draws a border around the image
func DrawBorder(img *image.RGBA, category string, thickness int) {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// Get border color based on category
	borderColor, ok := CategoryColors[category]
	if !ok {
		borderColor = CategoryColors["platform"] // default to blue
	}

	// Draw top and bottom borders
	for t := 0; t < thickness; t++ {
		for x := 0; x < width; x++ {
			img.Set(x, t, borderColor)          // Top
			img.Set(x, height-1-t, borderColor) // Bottom
		}
	}

	// Draw left and right borders
	for t := 0; t < thickness; t++ {
		for y := 0; y < height; y++ {
			img.Set(t, y, borderColor)         // Left
			img.Set(width-1-t, y, borderColor) // Right
		}
	}
}

// DrawCornerMarkers draws corner labels (TL, TR, BL, BR)
func DrawCornerMarkers(img *image.RGBA, width, height int) {
	white := color.RGBA{R: 255, G: 255, B: 255, A: 255}
	black := color.RGBA{R: 0, G: 0, B: 0, A: 255}

	// Use basic font for corner markers
	face := basicfont.Face7x13

	offset := 10

	// TL - Top Left
	drawTextWithOutline(img, "TL", offset, offset+10, face, white, black)

	// TR - Top Right
	trWidth := font.MeasureString(face, "TR").Ceil()
	drawTextWithOutline(img, "TR", width-trWidth-offset, offset+10, face, white, black)

	// BL - Bottom Left
	drawTextWithOutline(img, "BL", offset, height-offset, face, white, black)

	// BR - Bottom Right
	brWidth := font.MeasureString(face, "BR").Ceil()
	drawTextWithOutline(img, "BR", width-brWidth-offset, height-offset, face, white, black)
}

// drawTextWithOutline draws text with a simple outline effect
func drawTextWithOutline(img *image.RGBA, text string, x, y int, face font.Face, textColor, outlineColor color.Color) {
	drawer := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(outlineColor),
		Face: face,
		Dot:  fixed.Point26_6{},
	}

	// Draw outline (8 directions)
	offsets := []struct{ dx, dy int }{
		{-1, -1}, {0, -1}, {1, -1},
		{-1, 0}, {1, 0},
		{-1, 1}, {0, 1}, {1, 1},
	}

	for _, offset := range offsets {
		drawer.Dot = fixed.Point26_6{
			X: fixed.I(x + offset.dx),
			Y: fixed.I(y + offset.dy),
		}
		drawer.DrawString(text)
	}

	// Draw main text
	drawer.Src = image.NewUniform(textColor)
	drawer.Dot = fixed.Point26_6{
		X: fixed.I(x),
		Y: fixed.I(y),
	}
	drawer.DrawString(text)
}
