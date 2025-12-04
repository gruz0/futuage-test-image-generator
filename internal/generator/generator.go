package generator

import (
	"fmt"
	"image"
	"image/color"
)

// ImageSpec defines the specification for generating a test image
type ImageSpec struct {
	Width        int
	Height       int
	Ratio        string
	RatioDecimal float64
	Format       string
	Quality      int
	SizeCategory string
	Category     string // platform, common, edge
	OutputPath   string
	Filename     string
}

// CategoryColors defines the background colors for each category
var CategoryColors = map[string]color.RGBA{
	"platform": {R: 74, G: 144, B: 226, A: 255}, // Blue #4A90E2
	"common":   {R: 126, G: 211, B: 33, A: 255}, // Green #7ED321
	"edge":     {R: 245, G: 166, B: 35, A: 255}, // Orange #F5A623
}

// Generate creates a test image based on the provided specification
func Generate(spec ImageSpec) error {
	// 1. Create blank image canvas
	img := image.NewRGBA(image.Rect(0, 0, spec.Width, spec.Height))

	// 2. Draw grid pattern background
	if err := DrawGridBackground(img, spec.Category); err != nil {
		return fmt.Errorf("failed to draw grid background: %w", err)
	}

	// 3. Draw 2px border
	DrawBorder(img, spec.Category, 2)

	// 4. Render centered text overlay
	lines := []string{
		fmt.Sprintf("%d×%d", spec.Width, spec.Height),
		fmt.Sprintf("%s (%.3f)", spec.Ratio, spec.RatioDecimal),
		fmt.Sprintf("%s Q%d", spec.Format, spec.Quality),
		spec.SizeCategory,
	}
	if err := DrawTextOverlay(img, lines, spec.Width, spec.Height); err != nil {
		return fmt.Errorf("failed to draw text overlay: %w", err)
	}

	// 5. Draw corner markers
	DrawCornerMarkers(img, spec.Width, spec.Height)

	// 6. Encode to target format
	if err := EncodeImage(img, spec.OutputPath, spec.Format, spec.Quality); err != nil {
		return fmt.Errorf("failed to encode image: %w", err)
	}

	return nil
}

// GetFontSize returns adaptive font size based on image dimensions
func GetFontSize(width, height int) float64 {
	maxDim := width
	if height > maxDim {
		maxDim = height
	}

	switch {
	case maxDim <= 200:
		return 10.0
	case maxDim <= 800:
		return 14.0
	case maxDim <= 1500:
		return 18.0
	case maxDim <= 3000:
		return 24.0
	default:
		return 32.0
	}
}

// GetGridSize returns adaptive grid size based on image dimensions
func GetGridSize(width, height int) int {
	maxDim := width
	if height > maxDim {
		maxDim = height
	}

	// Base grid size is 100px, scale proportionally for smaller images
	switch {
	case maxDim <= 200:
		return 20
	case maxDim <= 500:
		return 50
	case maxDim <= 1000:
		return 75
	default:
		return 100
	}
}
