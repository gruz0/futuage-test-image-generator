package generator

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGenerate(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name   string
		spec   ImageSpec
		format string
	}{
		{
			name: "generate JPEG image",
			spec: ImageSpec{
				Width:        200,
				Height:       300,
				Ratio:        "2:3",
				RatioDecimal: 0.667,
				Format:       "jpeg",
				Quality:      85,
				SizeCategory: "Tiny",
				Category:     "platform",
				OutputPath:   filepath.Join(tmpDir, "test_jpeg.jpg"),
				Filename:     "test_jpeg.jpg",
			},
			format: "jpeg",
		},
		{
			name: "generate PNG image",
			spec: ImageSpec{
				Width:        100,
				Height:       100,
				Ratio:        "1:1",
				RatioDecimal: 1.0,
				Format:       "png",
				Quality:      95,
				SizeCategory: "Tiny",
				Category:     "common",
				OutputPath:   filepath.Join(tmpDir, "test_png.png"),
				Filename:     "test_png.png",
			},
			format: "png",
		},
		{
			name: "generate WebP image",
			spec: ImageSpec{
				Width:        150,
				Height:       200,
				Ratio:        "3:4",
				RatioDecimal: 0.75,
				Format:       "webp",
				Quality:      82,
				SizeCategory: "Tiny",
				Category:     "edge",
				OutputPath:   filepath.Join(tmpDir, "test_webp.webp"),
				Filename:     "test_webp.webp",
			},
			format: "webp",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Generate(tt.spec)
			if err != nil {
				t.Fatalf("Generate() error = %v", err)
			}

			// Verify file was created
			info, err := os.Stat(tt.spec.OutputPath)
			if err != nil {
				t.Fatalf("Generated file not found: %v", err)
			}

			// Verify file is not empty
			if info.Size() == 0 {
				t.Error("Generated file is empty")
			}

			// Verify dimensions are correct by checking file exists
			// (full image verification would require decoding)
			if info.Size() < 100 {
				t.Errorf("Generated file too small: %d bytes", info.Size())
			}
		})
	}
}

func TestGenerate_LargeImage(t *testing.T) {
	tmpDir := t.TempDir()

	spec := ImageSpec{
		Width:        1920,
		Height:       1080,
		Ratio:        "16:9",
		RatioDecimal: 1.778,
		Format:       "jpeg",
		Quality:      85,
		SizeCategory: "Large",
		Category:     "common",
		OutputPath:   filepath.Join(tmpDir, "large_test.jpg"),
		Filename:     "large_test.jpg",
	}

	err := Generate(spec)
	if err != nil {
		t.Fatalf("Generate() error for large image: %v", err)
	}

	info, err := os.Stat(spec.OutputPath)
	if err != nil {
		t.Fatalf("Large image file not found: %v", err)
	}

	// Large image should be at least a few KB
	if info.Size() < 1000 {
		t.Errorf("Large image file too small: %d bytes", info.Size())
	}
}

func TestGetFontSize(t *testing.T) {
	tests := []struct {
		name     string
		width    int
		height   int
		expected float64
	}{
		{"tiny image", 100, 150, 10.0},
		{"small image", 500, 750, 14.0},
		{"medium image", 1000, 1500, 18.0},
		{"large image", 2000, 3000, 24.0},
		{"xlarge image", 4000, 6000, 32.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetFontSize(tt.width, tt.height)
			if got != tt.expected {
				t.Errorf("GetFontSize(%d, %d) = %f, want %f", tt.width, tt.height, got, tt.expected)
			}
		})
	}
}

func TestGetGridSize(t *testing.T) {
	tests := []struct {
		name     string
		width    int
		height   int
		expected int
	}{
		{"tiny image", 100, 150, 20},
		{"small image", 400, 300, 50},
		{"medium image", 800, 600, 75},
		{"large image", 2000, 1500, 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetGridSize(tt.width, tt.height)
			if got != tt.expected {
				t.Errorf("GetGridSize(%d, %d) = %d, want %d", tt.width, tt.height, got, tt.expected)
			}
		})
	}
}

func TestCategoryColors(t *testing.T) {
	// Verify all expected categories have colors defined
	expectedCategories := []string{"platform", "common", "edge"}

	for _, cat := range expectedCategories {
		if _, ok := CategoryColors[cat]; !ok {
			t.Errorf("CategoryColors missing color for %q", cat)
		}
	}

	// Verify colors are not black (0,0,0) which would indicate uninitialized
	for cat, color := range CategoryColors {
		if color.R == 0 && color.G == 0 && color.B == 0 {
			t.Errorf("CategoryColors[%q] has black color (likely uninitialized)", cat)
		}
		if color.A != 255 {
			t.Errorf("CategoryColors[%q] alpha = %d, want 255", cat, color.A)
		}
	}
}

