package generator

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
)

// DrawGridBackground draws a grid pattern background with category-based color
func DrawGridBackground(img *image.RGBA, category string) error {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// Get background color based on category
	bgColor, ok := CategoryColors[category]
	if !ok {
		return fmt.Errorf("unknown category: %s", category)
	}

	// Fill entire background with category color
	draw.Draw(img, bounds, &image.Uniform{bgColor}, image.Point{}, draw.Src)

	// Get adaptive grid size
	gridSize := GetGridSize(width, height)

	// Draw grid lines (semi-transparent white)
	gridColor := color.RGBA{R: 255, G: 255, B: 255, A: 51} // 20% opacity (51/255)

	// Draw vertical grid lines
	for x := gridSize; x < width; x += gridSize {
		for y := 0; y < height; y++ {
			if x < width {
				img.Set(x, y, gridColor)
			}
		}
	}

	// Draw horizontal grid lines
	for y := gridSize; y < height; y += gridSize {
		for x := 0; x < width; x++ {
			if y < height {
				img.Set(x, y, gridColor)
			}
		}
	}

	return nil
}
