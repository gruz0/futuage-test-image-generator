package generator

import (
	"image"
	"image/color"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
)

// DrawTextOverlay draws centered text overlay with metadata
func DrawTextOverlay(img *image.RGBA, lines []string, width, height int) error {
	white := color.RGBA{R: 255, G: 255, B: 255, A: 255}
	black := color.RGBA{R: 0, G: 0, B: 0, A: 255}

	// Get adaptive font size
	fontSize := GetFontSize(width, height)

	// Use basic font (scaled appropriately)
	face := getFontFace(fontSize)

	// Calculate total text block height
	lineHeight := int(fontSize * 1.5)
	totalHeight := len(lines) * lineHeight

	// Start Y position (centered vertically)
	startY := (height - totalHeight) / 2

	// Draw each line centered
	for i, line := range lines {
		y := startY + (i+1)*lineHeight

		// Measure text width for centering
		textWidth := font.MeasureString(face, line).Ceil()
		x := (width - textWidth) / 2

		// Draw text with outline
		drawTextWithOutline(img, line, x, y, face, white, black)
	}

	return nil
}

// getFontFace returns appropriate font face based on size
func getFontFace(size float64) font.Face {
	// For now, use basic font
	// In production, you might want to load TrueType fonts based on size
	return basicfont.Face7x13
}

// DrawTextAt draws text at a specific position with outline
func DrawTextAt(img *image.RGBA, text string, x, y int, fontSize float64, textColor, outlineColor color.Color) {
	face := getFontFace(fontSize)
	drawTextWithOutline(img, text, x, y, face, textColor, outlineColor)
}

// MeasureText returns the width and height of rendered text
func MeasureText(text string, fontSize float64) (width, height int) {
	face := getFontFace(fontSize)
	width = font.MeasureString(face, text).Ceil()
	height = int(fontSize * 1.5)
	return width, height
}
