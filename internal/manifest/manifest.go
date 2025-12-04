package manifest

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gruz0/futuage-test-image-generator/internal/generator"
)

// Manifest represents the complete metadata for all generated images
type Manifest struct {
	GeneratedAt   string        `json:"generated_at"`
	ToolVersion   string        `json:"tool_version"`
	ConfigVersion string        `json:"config_version"`
	TotalImages   int           `json:"total_images"`
	Images        []ImageRecord `json:"images"`
}

// ImageRecord represents metadata for a single generated image
type ImageRecord struct {
	Filename      string  `json:"filename"`
	Category      string  `json:"category"`
	Subcategory   string  `json:"subcategory"`
	Width         int     `json:"width"`
	Height        int     `json:"height"`
	Ratio         string  `json:"ratio"`
	RatioDecimal  float64 `json:"ratio_decimal"`
	Format        string  `json:"format"`
	Quality       int     `json:"quality"`
	FileSizeBytes int64   `json:"file_size_bytes"`
	SizeCategory  string  `json:"size_category"`
}

// NewManifest creates a new Manifest
func NewManifest(toolVersion, configVersion string) *Manifest {
	return &Manifest{
		GeneratedAt:   time.Now().UTC().Format(time.RFC3339),
		ToolVersion:   toolVersion,
		ConfigVersion: configVersion,
		Images:        []ImageRecord{},
	}
}

// AddImage adds an image record to the manifest
func (m *Manifest) AddImage(spec generator.ImageSpec, fileSize int64) {
	// Determine category and subcategory from output path
	category, subcategory := extractCategoryFromPath(spec.OutputPath)

	record := ImageRecord{
		Filename:      getRelativePath(spec.OutputPath),
		Category:      category,
		Subcategory:   subcategory,
		Width:         spec.Width,
		Height:        spec.Height,
		Ratio:         spec.Ratio,
		RatioDecimal:  spec.RatioDecimal,
		Format:        strings.ToLower(spec.Format),
		Quality:       spec.Quality,
		FileSizeBytes: fileSize,
		SizeCategory:  strings.ToLower(spec.SizeCategory),
	}

	m.Images = append(m.Images, record)
	m.TotalImages = len(m.Images)
}

// Write writes the manifest to a JSON file
func (m *Manifest) Write(outputPath string) error {
	// Ensure the manifest is up to date
	m.TotalImages = len(m.Images)

	// Marshal to JSON with pretty printing
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal manifest: %w", err)
	}

	// Write to file
	if err := os.WriteFile(outputPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write manifest file: %w", err)
	}

	return nil
}

// extractCategoryFromPath extracts category and subcategory from file path
// e.g., "/path/to/ratios/2-3/file.jpg" -> ("ratios", "2-3")
func extractCategoryFromPath(path string) (category, subcategory string) {
	// Get the path components
	dir := filepath.Dir(path)
	parts := strings.Split(filepath.ToSlash(dir), "/")

	// Find the category (ratios, targets, edge-cases)
	for i := len(parts) - 1; i >= 0; i-- {
		part := parts[i]
		if part == "ratios" || part == "targets" || part == "edge-cases" {
			category = part
			// Subcategory is the next part if it exists
			if i+1 < len(parts) {
				subcategory = parts[i+1]
			}
			break
		}
	}

	return category, subcategory
}

// getRelativePath returns the relative path from the base directory
// e.g., "/tmp/output/ratios/2-3/file.jpg" -> "ratios/2-3/file.jpg"
func getRelativePath(path string) string {
	// Find the category in the path
	parts := strings.Split(filepath.ToSlash(path), "/")

	for i, part := range parts {
		if part == "ratios" || part == "targets" || part == "edge-cases" {
			// Return from this point onwards
			return strings.Join(parts[i:], "/")
		}
	}

	// If no category found, return just the filename
	return filepath.Base(path)
}

// Summary returns a summary of the manifest
func (m *Manifest) Summary() string {
	if m.TotalImages == 0 {
		return "No images generated"
	}

	// Count by category
	categoryCounts := make(map[string]int)
	for _, img := range m.Images {
		categoryCounts[img.Category]++
	}

	// Calculate total size
	var totalSize int64
	for _, img := range m.Images {
		totalSize += img.FileSizeBytes
	}

	summary := fmt.Sprintf("Generated %d images (%.2f MB total)\n", m.TotalImages, float64(totalSize)/(1024*1024))
	summary += "By category:\n"
	for cat, count := range categoryCounts {
		summary += fmt.Sprintf("  %s: %d\n", cat, count)
	}

	return summary
}
