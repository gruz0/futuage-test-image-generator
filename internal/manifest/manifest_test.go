package manifest

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gruz0/futuage-test-image-generator/internal/generator"
)

func TestNewManifest(t *testing.T) {
	m := NewManifest("1.0.0", "2.0.0")

	if m.ToolVersion != "1.0.0" {
		t.Errorf("NewManifest().ToolVersion = %q, want %q", m.ToolVersion, "1.0.0")
	}

	if m.ConfigVersion != "2.0.0" {
		t.Errorf("NewManifest().ConfigVersion = %q, want %q", m.ConfigVersion, "2.0.0")
	}

	if m.GeneratedAt == "" {
		t.Error("NewManifest().GeneratedAt is empty")
	}

	if m.TotalImages != 0 {
		t.Errorf("NewManifest().TotalImages = %d, want 0", m.TotalImages)
	}

	if m.Images == nil {
		t.Error("NewManifest().Images is nil")
	}
}

func TestManifest_AddImage(t *testing.T) {
	m := NewManifest("1.0.0", "1.0.0")

	spec := generator.ImageSpec{
		Width:        1000,
		Height:       1500,
		Ratio:        "2:3",
		RatioDecimal: 0.667,
		Format:       "JPEG",
		Quality:      85,
		SizeCategory: "Medium",
		Category:     "platform",
		OutputPath:   "/tmp/test/ratios/2-3/test.jpg",
		Filename:     "test.jpg",
	}

	m.AddImage(spec, 12345)

	if m.TotalImages != 1 {
		t.Errorf("After AddImage(), TotalImages = %d, want 1", m.TotalImages)
	}

	if len(m.Images) != 1 {
		t.Fatalf("After AddImage(), len(Images) = %d, want 1", len(m.Images))
	}

	img := m.Images[0]

	if img.Width != 1000 {
		t.Errorf("Image.Width = %d, want 1000", img.Width)
	}

	if img.Height != 1500 {
		t.Errorf("Image.Height = %d, want 1500", img.Height)
	}

	if img.Format != "jpeg" {
		t.Errorf("Image.Format = %q, want %q", img.Format, "jpeg")
	}

	if img.FileSizeBytes != 12345 {
		t.Errorf("Image.FileSizeBytes = %d, want 12345", img.FileSizeBytes)
	}

	if img.Category != "ratios" {
		t.Errorf("Image.Category = %q, want %q", img.Category, "ratios")
	}

	if img.Subcategory != "2-3" {
		t.Errorf("Image.Subcategory = %q, want %q", img.Subcategory, "2-3")
	}
}

func TestManifest_Write(t *testing.T) {
	tmpDir := t.TempDir()
	outputPath := filepath.Join(tmpDir, "manifest.json")

	m := NewManifest("1.0.0", "1.0.0")
	m.AddImage(generator.ImageSpec{
		Width:        100,
		Height:       100,
		Ratio:        "1:1",
		RatioDecimal: 1.0,
		Format:       "PNG",
		Quality:      95,
		SizeCategory: "Tiny",
		Category:     "test",
		OutputPath:   "/tmp/targets/test.png",
		Filename:     "test.png",
	}, 5000)

	err := m.Write(outputPath)
	if err != nil {
		t.Fatalf("Manifest.Write() error = %v", err)
	}

	// Verify file was created
	data, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("Failed to read manifest file: %v", err)
	}

	// Verify it's valid JSON
	var parsed Manifest
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("Manifest file is not valid JSON: %v", err)
	}

	if parsed.TotalImages != 1 {
		t.Errorf("Parsed manifest TotalImages = %d, want 1", parsed.TotalImages)
	}

	if len(parsed.Images) != 1 {
		t.Errorf("Parsed manifest len(Images) = %d, want 1", len(parsed.Images))
	}
}

func TestExtractCategoryFromPath(t *testing.T) {
	tests := []struct {
		path            string
		wantCategory    string
		wantSubcategory string
	}{
		{
			path:            "/tmp/output/ratios/2-3/file.jpg",
			wantCategory:    "ratios",
			wantSubcategory: "2-3",
		},
		{
			path:            "/tmp/output/targets/file.jpg",
			wantCategory:    "targets",
			wantSubcategory: "",
		},
		{
			path:            "/tmp/output/edge-cases/file.jpg",
			wantCategory:    "edge-cases",
			wantSubcategory: "",
		},
		{
			path:            "/tmp/output/ratios/16-9/subdir/file.jpg",
			wantCategory:    "ratios",
			wantSubcategory: "16-9",
		},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			gotCat, gotSub := extractCategoryFromPath(tt.path)
			if gotCat != tt.wantCategory {
				t.Errorf("extractCategoryFromPath(%q) category = %q, want %q", tt.path, gotCat, tt.wantCategory)
			}
			if gotSub != tt.wantSubcategory {
				t.Errorf("extractCategoryFromPath(%q) subcategory = %q, want %q", tt.path, gotSub, tt.wantSubcategory)
			}
		})
	}
}

func TestGetRelativePath(t *testing.T) {
	tests := []struct {
		path     string
		expected string
	}{
		{
			path:     "/tmp/output/ratios/2-3/file.jpg",
			expected: "ratios/2-3/file.jpg",
		},
		{
			path:     "/tmp/output/targets/file.jpg",
			expected: "targets/file.jpg",
		},
		{
			path:     "/tmp/output/edge-cases/file.jpg",
			expected: "edge-cases/file.jpg",
		},
		{
			path:     "/unknown/path/file.jpg",
			expected: "file.jpg",
		},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			got := getRelativePath(tt.path)
			if got != tt.expected {
				t.Errorf("getRelativePath(%q) = %q, want %q", tt.path, got, tt.expected)
			}
		})
	}
}

func TestManifest_Summary(t *testing.T) {
	m := NewManifest("1.0.0", "1.0.0")

	// Empty manifest
	summary := m.Summary()
	if !strings.Contains(summary, "No images") {
		t.Errorf("Empty manifest summary should contain 'No images', got %q", summary)
	}

	// Add some images
	m.AddImage(generator.ImageSpec{
		Width:      100,
		Height:     100,
		OutputPath: "/tmp/ratios/1-1/a.jpg",
	}, 1000)
	m.AddImage(generator.ImageSpec{
		Width:      200,
		Height:     200,
		OutputPath: "/tmp/ratios/1-1/b.jpg",
	}, 2000)
	m.AddImage(generator.ImageSpec{
		Width:      300,
		Height:     300,
		OutputPath: "/tmp/targets/c.jpg",
	}, 3000)

	summary = m.Summary()

	if !strings.Contains(summary, "3 images") {
		t.Errorf("Summary should mention 3 images, got %q", summary)
	}

	if !strings.Contains(summary, "ratios") {
		t.Errorf("Summary should mention ratios category, got %q", summary)
	}

	if !strings.Contains(summary, "targets") {
		t.Errorf("Summary should mention targets category, got %q", summary)
	}
}

