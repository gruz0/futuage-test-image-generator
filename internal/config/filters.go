package config

import (
	"fmt"
	"strings"
)

// Filters represents the active filters for generation
type Filters struct {
	Ratios  []string // e.g., ["platform", "common"]
	Sizes   []string // e.g., ["medium", "large"]
	Formats []string // e.g., ["jpeg", "png"]
}

// NewFilters creates a new Filters instance
func NewFilters(ratios, sizes, formats []string) *Filters {
	return &Filters{
		Ratios:  normalizeStrings(ratios),
		Sizes:   normalizeStrings(sizes),
		Formats: normalizeStrings(formats),
	}
}

// IsEmpty returns true if no filters are active
func (f *Filters) IsEmpty() bool {
	return len(f.Ratios) == 0 && len(f.Sizes) == 0 && len(f.Formats) == 0
}

// ShouldIncludeRatioCategory checks if a ratio category should be included
func (f *Filters) ShouldIncludeRatioCategory(category string) bool {
	if len(f.Ratios) == 0 {
		return true // No filter, include all
	}

	category = strings.ToLower(category)
	for _, r := range f.Ratios {
		if strings.ToLower(r) == category {
			return true
		}
	}
	return false
}

// ShouldIncludeSizeCategory checks if a size category should be included
func (f *Filters) ShouldIncludeSizeCategory(category string) bool {
	if len(f.Sizes) == 0 {
		return true // No filter, include all
	}

	category = strings.ToLower(category)
	for _, s := range f.Sizes {
		if strings.ToLower(s) == category {
			return true
		}
	}
	return false
}

// ShouldIncludeFormat checks if a format should be included
func (f *Filters) ShouldIncludeFormat(format string) bool {
	if len(f.Formats) == 0 {
		return true // No filter, include all
	}

	format = strings.ToLower(format)
	for _, fmt := range f.Formats {
		if strings.ToLower(fmt) == format {
			return true
		}
	}
	return false
}

// Validate checks if the filters are valid against the configuration
func (f *Filters) Validate(cfg *Config) error {
	// Validate ratio categories
	for _, ratio := range f.Ratios {
		ratio = strings.ToLower(ratio)
		found := false
		for presetName := range cfg.Presets {
			if strings.ToLower(presetName) == ratio {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("unknown ratio category: %s", ratio)
		}
	}

	// Validate size categories
	for _, size := range f.Sizes {
		size = strings.ToLower(size)
		found := false
		for sizeName := range cfg.Sizes {
			if strings.ToLower(sizeName) == size {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("unknown size category: %s", size)
		}
	}

	// Validate formats
	for _, format := range f.Formats {
		format = strings.ToLower(format)
		found := false
		for formatName := range cfg.Formats {
			if strings.ToLower(formatName) == format {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("unknown format: %s", format)
		}
	}

	return nil
}

// normalizeStrings converts strings to lowercase and removes empty strings
func normalizeStrings(strs []string) []string {
	var result []string
	for _, s := range strs {
		s = strings.TrimSpace(s)
		if s != "" {
			result = append(result, strings.ToLower(s))
		}
	}
	return result
}

// Summary returns a human-readable summary of active filters
func (f *Filters) Summary() string {
	if f.IsEmpty() {
		return "No filters (generating all)"
	}

	var parts []string
	if len(f.Ratios) > 0 {
		parts = append(parts, fmt.Sprintf("ratios: %s", strings.Join(f.Ratios, ", ")))
	}
	if len(f.Sizes) > 0 {
		parts = append(parts, fmt.Sprintf("sizes: %s", strings.Join(f.Sizes, ", ")))
	}
	if len(f.Formats) > 0 {
		parts = append(parts, fmt.Sprintf("formats: %s", strings.Join(f.Formats, ", ")))
	}

	return strings.Join(parts, " | ")
}
