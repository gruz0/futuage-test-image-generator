package config

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/gruz0/futuage-test-image-generator/internal/generator"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// SpecBuilder builds ImageSpec instances from configuration
type SpecBuilder struct {
	Config  *Config
	Filters *Filters
	BaseDir string
}

// NewSpecBuilder creates a new SpecBuilder
func NewSpecBuilder(cfg *Config, filters *Filters, baseDir string) *SpecBuilder {
	return &SpecBuilder{
		Config:  cfg,
		Filters: filters,
		BaseDir: baseDir,
	}
}

// BuildSpecs generates all ImageSpec instances based on configuration and filters
func (b *SpecBuilder) BuildSpecs() ([]generator.ImageSpec, error) {
	var specs []generator.ImageSpec

	// 1. Generate ratio-based images (from presets)
	ratioSpecs, err := b.buildRatioSpecs()
	if err != nil {
		return nil, fmt.Errorf("failed to build ratio specs: %w", err)
	}
	specs = append(specs, ratioSpecs...)

	// 2. Generate platform target images
	targetSpecs, err := b.buildTargetSpecs()
	if err != nil {
		return nil, fmt.Errorf("failed to build target specs: %w", err)
	}
	specs = append(specs, targetSpecs...)

	// 3. Generate edge case images
	edgeSpecs, err := b.buildEdgeCaseSpecs()
	if err != nil {
		return nil, fmt.Errorf("failed to build edge case specs: %w", err)
	}
	specs = append(specs, edgeSpecs...)

	return specs, nil
}

// buildRatioSpecs builds specs for all ratio presets
func (b *SpecBuilder) buildRatioSpecs() ([]generator.ImageSpec, error) {
	var specs []generator.ImageSpec

	for presetName, preset := range b.Config.Presets {
		// Check if this ratio category should be included
		if !b.Filters.ShouldIncludeRatioCategory(presetName) {
			continue
		}

		for _, ratioStr := range preset.Ratios {
			ratioInfo, err := ParseRatio(ratioStr)
			if err != nil {
				return nil, fmt.Errorf("invalid ratio %s: %w", ratioStr, err)
			}

			// Generate for each size category
			for sizeName, sizeConfig := range b.Config.Sizes {
				if !b.Filters.ShouldIncludeSizeCategory(sizeName) {
					continue
				}

				// Pick one representative size from each category (first one)
				if len(sizeConfig.BaseSizes) == 0 {
					continue
				}
				baseSize := sizeConfig.BaseSizes[0]

				// Calculate dimensions
				width, height := CalculateDimensions(ratioInfo, baseSize)

				// Generate for each format
				for formatName, format := range b.Config.Formats {
					if !b.Filters.ShouldIncludeFormat(formatName) {
						continue
					}

					// Pick one representative quality (first one)
					if len(format.Qualities) == 0 {
						continue
					}
					quality := format.Qualities[0]

					// Build filename
					filename := fmt.Sprintf("%s_%dx%d_%s_q%d%s",
						strings.ToLower(sizeName),
						width, height,
						strings.ToLower(formatName),
						quality,
						format.Extension,
					)

					// Build output path
					outputPath := filepath.Join(
						b.BaseDir,
						"ratios",
						ratioInfo.DisplayName,
						filename,
					)

					spec := generator.ImageSpec{
						Width:        width,
						Height:       height,
						Ratio:        ratioStr,
						RatioDecimal: ratioInfo.Decimal,
						Format:       strings.ToUpper(formatName),
						Quality:      quality,
						SizeCategory: cases.Title(language.English).String(sizeName),
						Category:     presetName,
						OutputPath:   outputPath,
						Filename:     filename,
					}

					specs = append(specs, spec)
				}
			}
		}
	}

	return specs, nil
}

// buildTargetSpecs builds specs for platform targets
func (b *SpecBuilder) buildTargetSpecs() ([]generator.ImageSpec, error) {
	var specs []generator.ImageSpec

	for targetName, target := range b.Config.Targets {
		// Parse ratio to get category and decimal
		ratioInfo, err := ParseRatio(target.Ratio)
		if err != nil {
			return nil, fmt.Errorf("invalid ratio %s for target %s: %w", target.Ratio, targetName, err)
		}

		// Check if ratio category should be included
		category := b.Config.GetCategoryForRatio(target.Ratio)
		if !b.Filters.ShouldIncludeRatioCategory(category) {
			continue
		}

		width := target.Dimensions[0]
		height := target.Dimensions[1]

		// Determine size category based on dimensions
		maxDim := width
		if height > maxDim {
			maxDim = height
		}
		sizeCategory := b.getSizeCategoryForDimension(maxDim)

		// Check if size category should be included
		if !b.Filters.ShouldIncludeSizeCategory(sizeCategory) {
			continue
		}

		// Generate for each format (typically only JPEG for targets)
		for formatName, format := range b.Config.Formats {
			if !b.Filters.ShouldIncludeFormat(formatName) {
				continue
			}

			// Use Q85 for targets (typical platform quality)
			quality := 85
			if len(format.Qualities) > 0 {
				// Find closest to 85
				quality = format.Qualities[0]
				for _, q := range format.Qualities {
					if abs(q-85) < abs(quality-85) {
						quality = q
					}
				}
			}

			// Build filename
			filename := fmt.Sprintf("%s_%dx%d_%s_q%d%s",
				targetName,
				width, height,
				strings.ToLower(formatName),
				quality,
				format.Extension,
			)

			// Build output path
			outputPath := filepath.Join(
				b.BaseDir,
				"targets",
				filename,
			)

			spec := generator.ImageSpec{
				Width:        width,
				Height:       height,
				Ratio:        target.Ratio,
				RatioDecimal: ratioInfo.Decimal,
				Format:       strings.ToUpper(formatName),
				Quality:      quality,
				SizeCategory: cases.Title(language.English).String(sizeCategory),
				Category:     category,
				OutputPath:   outputPath,
				Filename:     filename,
			}

			specs = append(specs, spec)
		}
	}

	return specs, nil
}

// buildEdgeCaseSpecs builds specs for edge cases
func (b *SpecBuilder) buildEdgeCaseSpecs() ([]generator.ImageSpec, error) {
	var specs []generator.ImageSpec

	for _, edgeCase := range b.Config.EdgeCases {
		width := edgeCase.Dimensions[0]
		height := edgeCase.Dimensions[1]

		// Calculate ratio
		ratioDecimal := float64(width) / float64(height)
		ratioStr := fmt.Sprintf("%d:%d", width, height)

		// Simplify ratio if possible
		gcd := gcd(width, height)
		if gcd > 1 {
			ratioStr = fmt.Sprintf("%d:%d", width/gcd, height/gcd)
		}

		category := "edge"

		// Check if ratio category should be included
		if !b.Filters.ShouldIncludeRatioCategory(category) {
			continue
		}

		// Determine size category
		maxDim := width
		if height > maxDim {
			maxDim = height
		}
		sizeCategory := b.getSizeCategoryForDimension(maxDim)

		// Check if size category should be included
		if !b.Filters.ShouldIncludeSizeCategory(sizeCategory) {
			continue
		}

		// Generate for each format
		for formatName, format := range b.Config.Formats {
			if !b.Filters.ShouldIncludeFormat(formatName) {
				continue
			}

			// Use first quality
			if len(format.Qualities) == 0 {
				continue
			}
			quality := format.Qualities[0]

			// Build filename
			filename := fmt.Sprintf("%s_%dx%d_%s_q%d%s",
				edgeCase.Name,
				width, height,
				strings.ToLower(formatName),
				quality,
				format.Extension,
			)

			// Build output path
			outputPath := filepath.Join(
				b.BaseDir,
				"edge-cases",
				filename,
			)

			spec := generator.ImageSpec{
				Width:        width,
				Height:       height,
				Ratio:        ratioStr,
				RatioDecimal: ratioDecimal,
				Format:       strings.ToUpper(formatName),
				Quality:      quality,
				SizeCategory: cases.Title(language.English).String(sizeCategory),
				Category:     category,
				OutputPath:   outputPath,
				Filename:     filename,
			}

			specs = append(specs, spec)
		}
	}

	return specs, nil
}

// getSizeCategoryForDimension returns the size category for a given dimension
func (b *SpecBuilder) getSizeCategoryForDimension(dim int) string {
	switch {
	case dim <= 200:
		return "tiny"
	case dim <= 800:
		return "small"
	case dim <= 1500:
		return "medium"
	case dim <= 3000:
		return "large"
	default:
		return "xlarge"
	}
}

// abs returns the absolute value of an integer
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// gcd returns the greatest common divisor
func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}
