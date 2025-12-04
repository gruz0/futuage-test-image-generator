package generator

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"github.com/chai2010/webp"
)

// EncodeImage encodes an image to the specified format and quality
func EncodeImage(img image.Image, outputPath, format string, quality int) error {
	// Ensure output directory exists
	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Create output file
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer func() {
		if cerr := file.Close(); cerr != nil && err == nil {
			err = fmt.Errorf("failed to close output file: %w", cerr)
		}
	}()

	// Encode based on format
	format = strings.ToLower(format)
	switch format {
	case "jpeg", "jpg":
		return encodeJPEG(file, img, quality)
	case "png":
		return encodePNG(file, img)
	case "webp":
		return encodeWebP(file, img, quality)
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}
}

// encodeJPEG encodes image to JPEG format
func encodeJPEG(file *os.File, img image.Image, quality int) error {
	opts := &jpeg.Options{
		Quality: quality,
	}
	if err := jpeg.Encode(file, img, opts); err != nil {
		return fmt.Errorf("failed to encode JPEG: %w", err)
	}
	return nil
}

// encodePNG encodes image to PNG format
func encodePNG(file *os.File, img image.Image) error {
	encoder := &png.Encoder{
		CompressionLevel: png.BestCompression,
	}
	if err := encoder.Encode(file, img); err != nil {
		return fmt.Errorf("failed to encode PNG: %w", err)
	}
	return nil
}

// encodeWebP encodes image to WebP format
func encodeWebP(file *os.File, img image.Image, quality int) error {
	// Convert quality to WebP quality (0-100 scale, but use float32)
	opts := &webp.Options{
		Lossless: false,
		Quality:  float32(quality),
	}

	if err := webp.Encode(file, img, opts); err != nil {
		return fmt.Errorf("failed to encode WebP: %w", err)
	}

	return nil
}

// GetFileExtension returns the appropriate file extension for a format
func GetFileExtension(format string) string {
	format = strings.ToLower(format)
	switch format {
	case "jpeg", "jpg":
		return ".jpg"
	case "png":
		return ".png"
	case "webp":
		return ".webp"
	default:
		return ".jpg"
	}
}
