package generator

import "fmt"

// GenerateTestImage creates a single test image for verification
func GenerateTestImage(outputPath string) error {
	spec := ImageSpec{
		Width:        1000,
		Height:       1500,
		Ratio:        "2:3",
		RatioDecimal: 0.667,
		Format:       "jpeg",
		Quality:      85,
		SizeCategory: "Medium",
		Category:     "platform",
		OutputPath:   outputPath,
		Filename:     "test_1000x1500_jpeg_q85.jpg",
	}

	if err := Generate(spec); err != nil {
		return fmt.Errorf("failed to generate test image: %w", err)
	}

	return nil
}
