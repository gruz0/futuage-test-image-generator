package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Config represents the complete configuration
type Config struct {
	Version   string                `json:"version"`
	Presets   map[string]Preset     `json:"presets"`
	Sizes     map[string]SizeConfig `json:"sizes"`
	Formats   map[string]Format     `json:"formats"`
	Targets   map[string]Target     `json:"targets"`
	EdgeCases []EdgeCase            `json:"edge_cases"`
}

// Preset represents a ratio preset category
type Preset struct {
	Description string   `json:"description"`
	Ratios      []string `json:"ratios"`
}

// SizeConfig represents a size category configuration
type SizeConfig struct {
	Description string `json:"description"`
	BaseSizes   []int  `json:"base_sizes"`
}

// Format represents an image format specification
type Format struct {
	Qualities []int  `json:"qualities"`
	MimeType  string `json:"mime_type"`
	Extension string `json:"extension"`
}

// Target represents a platform target specification
type Target struct {
	Platform    string `json:"platform"`
	Dimensions  []int  `json:"dimensions"`
	Ratio       string `json:"ratio"`
	Description string `json:"description"`
}

// EdgeCase represents an edge case test scenario
type EdgeCase struct {
	Name        string `json:"name"`
	Dimensions  []int  `json:"dimensions"`
	Description string `json:"description"`
}

// LoadConfig loads configuration from file or returns default
func LoadConfig(configPath string) (*Config, error) {
	var cfg Config

	if configPath == "" {
		// Load embedded default configuration
		if err := json.Unmarshal(defaultConfigJSON, &cfg); err != nil {
			return nil, fmt.Errorf("failed to parse default config: %w", err)
		}
	} else {
		// Load custom configuration file
		data, err := os.ReadFile(configPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}

		if err := json.Unmarshal(data, &cfg); err != nil {
			return nil, fmt.Errorf("failed to parse config file: %w", err)
		}
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &cfg, nil
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if c.Version == "" {
		return fmt.Errorf("version is required")
	}

	if len(c.Presets) == 0 {
		return fmt.Errorf("at least one preset is required")
	}

	if len(c.Sizes) == 0 {
		return fmt.Errorf("at least one size category is required")
	}

	if len(c.Formats) == 0 {
		return fmt.Errorf("at least one format is required")
	}

	// Validate ratios
	for presetName, preset := range c.Presets {
		for _, ratio := range preset.Ratios {
			if _, err := ParseRatio(ratio); err != nil {
				return fmt.Errorf("invalid ratio %s in preset %s: %w", ratio, presetName, err)
			}
		}
	}

	// Validate targets
	for targetName, target := range c.Targets {
		if len(target.Dimensions) != 2 {
			return fmt.Errorf("target %s must have exactly 2 dimensions", targetName)
		}
		if target.Dimensions[0] <= 0 || target.Dimensions[1] <= 0 {
			return fmt.Errorf("target %s dimensions must be positive", targetName)
		}
	}

	// Validate edge cases
	for _, edgeCase := range c.EdgeCases {
		if len(edgeCase.Dimensions) != 2 {
			return fmt.Errorf("edge case %s must have exactly 2 dimensions", edgeCase.Name)
		}
		if edgeCase.Dimensions[0] <= 0 || edgeCase.Dimensions[1] <= 0 {
			return fmt.Errorf("edge case %s dimensions must be positive", edgeCase.Name)
		}
	}

	return nil
}

// RatioInfo represents parsed ratio information
type RatioInfo struct {
	Ratio       string
	Width       int
	Height      int
	Decimal     float64
	IsDecimal   bool
	DisplayName string // e.g., "2-3", "1-1", "16-9"
}

// ParseRatio parses a ratio string (e.g., "2:3", "1.91:1", "16:9") into components
func ParseRatio(ratio string) (*RatioInfo, error) {
	ratio = strings.TrimSpace(ratio)

	// Split by ":"
	parts := strings.Split(ratio, ":")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid ratio format: %s (expected W:H)", ratio)
	}

	// Parse width and height
	widthStr := strings.TrimSpace(parts[0])
	heightStr := strings.TrimSpace(parts[1])

	// Check if it's a decimal ratio (e.g., "1.91:1")
	isDecimal := strings.Contains(widthStr, ".") || strings.Contains(heightStr, ".")

	var width, height int
	var decimal float64

	if isDecimal {
		// Parse as float then convert to ratio
		w, err := strconv.ParseFloat(widthStr, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid ratio width: %s", widthStr)
		}
		h, err := strconv.ParseFloat(heightStr, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid ratio height: %s", heightStr)
		}

		decimal = w / h

		// Convert to integer ratio for calculations (multiply by 100)
		width = int(w * 100)
		height = int(h * 100)
	} else {
		// Parse as integers
		w, err := strconv.Atoi(widthStr)
		if err != nil {
			return nil, fmt.Errorf("invalid ratio width: %s", widthStr)
		}
		h, err := strconv.Atoi(heightStr)
		if err != nil {
			return nil, fmt.Errorf("invalid ratio height: %s", heightStr)
		}

		width = w
		height = h
		decimal = float64(w) / float64(h)
	}

	if width <= 0 || height <= 0 {
		return nil, fmt.Errorf("ratio dimensions must be positive: %s", ratio)
	}

	// Create display name (replace : with -, remove spaces)
	displayName := strings.ReplaceAll(ratio, ":", "-")
	displayName = strings.ReplaceAll(displayName, ".", "_")
	displayName = strings.ReplaceAll(displayName, " ", "")

	return &RatioInfo{
		Ratio:       ratio,
		Width:       width,
		Height:      height,
		Decimal:     decimal,
		IsDecimal:   isDecimal,
		DisplayName: displayName,
	}, nil
}

// CalculateDimensions calculates actual image dimensions from ratio and base size
func CalculateDimensions(ratioInfo *RatioInfo, baseSize int) (width, height int) {
	// Determine if ratio is portrait (height > width) or landscape
	if ratioInfo.Height > ratioInfo.Width {
		// Portrait: base size is the height
		height = baseSize
		width = int(float64(height) * float64(ratioInfo.Width) / float64(ratioInfo.Height))
	} else {
		// Landscape or square: base size is the width
		width = baseSize
		height = int(float64(width) * float64(ratioInfo.Height) / float64(ratioInfo.Width))
	}

	// Ensure dimensions are at least 1
	if width < 1 {
		width = 1
	}
	if height < 1 {
		height = 1
	}

	return width, height
}

// GetCategoryForRatio determines which category a ratio belongs to
func (c *Config) GetCategoryForRatio(ratio string) string {
	for categoryName, preset := range c.Presets {
		for _, r := range preset.Ratios {
			if r == ratio {
				return categoryName
			}
		}
	}
	return "edge" // default to edge category
}

// GetSizeCategoryName returns the size category name for a given base size
func (c *Config) GetSizeCategoryName(baseSize int) string {
	for categoryName, sizeConfig := range c.Sizes {
		for _, size := range sizeConfig.BaseSizes {
			if size == baseSize {
				return categoryName
			}
		}
	}
	return "unknown"
}
