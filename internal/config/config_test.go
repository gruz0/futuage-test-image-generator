package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseRatio(t *testing.T) {
	tests := []struct {
		name        string
		ratio       string
		wantWidth   int
		wantHeight  int
		wantDecimal float64
		wantDisplay string
		wantErr     bool
	}{
		{
			name:        "standard ratio 2:3",
			ratio:       "2:3",
			wantWidth:   2,
			wantHeight:  3,
			wantDecimal: 0.6666666666666666,
			wantDisplay: "2-3",
			wantErr:     false,
		},
		{
			name:        "square ratio 1:1",
			ratio:       "1:1",
			wantWidth:   1,
			wantHeight:  1,
			wantDecimal: 1.0,
			wantDisplay: "1-1",
			wantErr:     false,
		},
		{
			name:        "landscape ratio 16:9",
			ratio:       "16:9",
			wantWidth:   16,
			wantHeight:  9,
			wantDecimal: 1.7777777777777777,
			wantDisplay: "16-9",
			wantErr:     false,
		},
		{
			name:        "decimal ratio 1.91:1",
			ratio:       "1.91:1",
			wantWidth:   191,
			wantHeight:  100,
			wantDecimal: 1.91,
			wantDisplay: "1_91-1",
			wantErr:     false,
		},
		{
			name:        "ratio with spaces",
			ratio:       " 4 : 5 ",
			wantWidth:   4,
			wantHeight:  5,
			wantDecimal: 0.8,
			wantDisplay: "4-5", // all spaces removed
			wantErr:     false,
		},
		{
			name:    "invalid format - no colon",
			ratio:   "2x3",
			wantErr: true,
		},
		{
			name:    "invalid format - too many colons",
			ratio:   "2:3:4",
			wantErr: true,
		},
		{
			name:    "invalid width",
			ratio:   "abc:3",
			wantErr: true,
		},
		{
			name:    "invalid height",
			ratio:   "2:def",
			wantErr: true,
		},
		{
			name:    "zero width",
			ratio:   "0:3",
			wantErr: true,
		},
		{
			name:    "negative height",
			ratio:   "2:-3",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseRatio(tt.ratio)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ParseRatio(%q) expected error, got nil", tt.ratio)
				}
				return
			}

			if err != nil {
				t.Errorf("ParseRatio(%q) unexpected error: %v", tt.ratio, err)
				return
			}

			if got.Width != tt.wantWidth {
				t.Errorf("ParseRatio(%q).Width = %d, want %d", tt.ratio, got.Width, tt.wantWidth)
			}
			if got.Height != tt.wantHeight {
				t.Errorf("ParseRatio(%q).Height = %d, want %d", tt.ratio, got.Height, tt.wantHeight)
			}
			if got.Decimal != tt.wantDecimal {
				t.Errorf("ParseRatio(%q).Decimal = %f, want %f", tt.ratio, got.Decimal, tt.wantDecimal)
			}
			if got.DisplayName != tt.wantDisplay {
				t.Errorf("ParseRatio(%q).DisplayName = %q, want %q", tt.ratio, got.DisplayName, tt.wantDisplay)
			}
		})
	}
}

func TestCalculateDimensions(t *testing.T) {
	tests := []struct {
		name       string
		ratio      string
		baseSize   int
		wantWidth  int
		wantHeight int
	}{
		{
			name:       "portrait 2:3 with base 1500",
			ratio:      "2:3",
			baseSize:   1500,
			wantWidth:  1000,
			wantHeight: 1500,
		},
		{
			name:       "square 1:1 with base 1000",
			ratio:      "1:1",
			baseSize:   1000,
			wantWidth:  1000,
			wantHeight: 1000,
		},
		{
			name:       "landscape 16:9 with base 1920",
			ratio:      "16:9",
			baseSize:   1920,
			wantWidth:  1920,
			wantHeight: 1080,
		},
		{
			name:       "portrait 4:5 with base 1350",
			ratio:      "4:5",
			baseSize:   1350,
			wantWidth:  1080,
			wantHeight: 1350,
		},
		{
			name:       "decimal ratio 1.91:1 with base 1200",
			ratio:      "1.91:1",
			baseSize:   1200,
			wantWidth:  1200,
			wantHeight: 628,
		},
		{
			name:       "tiny size 100",
			ratio:      "2:3",
			baseSize:   100,
			wantWidth:  66,
			wantHeight: 100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ratioInfo, err := ParseRatio(tt.ratio)
			if err != nil {
				t.Fatalf("ParseRatio(%q) failed: %v", tt.ratio, err)
			}

			gotWidth, gotHeight := CalculateDimensions(ratioInfo, tt.baseSize)

			if gotWidth != tt.wantWidth {
				t.Errorf("CalculateDimensions(%q, %d) width = %d, want %d", tt.ratio, tt.baseSize, gotWidth, tt.wantWidth)
			}
			if gotHeight != tt.wantHeight {
				t.Errorf("CalculateDimensions(%q, %d) height = %d, want %d", tt.ratio, tt.baseSize, gotHeight, tt.wantHeight)
			}
		})
	}
}

func TestLoadConfig_Default(t *testing.T) {
	cfg, err := LoadConfig("")
	if err != nil {
		t.Fatalf("LoadConfig('') failed: %v", err)
	}

	if cfg.Version == "" {
		t.Error("LoadConfig('') returned config with empty version")
	}

	if len(cfg.Presets) == 0 {
		t.Error("LoadConfig('') returned config with no presets")
	}

	if len(cfg.Sizes) == 0 {
		t.Error("LoadConfig('') returned config with no sizes")
	}

	if len(cfg.Formats) == 0 {
		t.Error("LoadConfig('') returned config with no formats")
	}

	// Check for expected presets
	expectedPresets := []string{"platform", "common", "edge"}
	for _, preset := range expectedPresets {
		if _, ok := cfg.Presets[preset]; !ok {
			t.Errorf("LoadConfig('') missing expected preset %q", preset)
		}
	}
}

func TestLoadConfig_CustomFile(t *testing.T) {
	// Create a temporary config file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test-config.json")

	configJSON := `{
		"version": "test-1.0",
		"presets": {
			"test": {
				"description": "Test preset",
				"ratios": ["1:1", "2:3"]
			}
		},
		"sizes": {
			"test": {
				"description": "Test size",
				"base_sizes": [100, 200]
			}
		},
		"formats": {
			"jpeg": {
				"qualities": [85],
				"mime_type": "image/jpeg",
				"extension": ".jpg"
			}
		},
		"targets": {},
		"edge_cases": []
	}`

	if err := os.WriteFile(configPath, []byte(configJSON), 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	cfg, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig(%q) failed: %v", configPath, err)
	}

	if cfg.Version != "test-1.0" {
		t.Errorf("LoadConfig() version = %q, want %q", cfg.Version, "test-1.0")
	}

	if _, ok := cfg.Presets["test"]; !ok {
		t.Error("LoadConfig() missing 'test' preset")
	}
}

func TestLoadConfig_InvalidFile(t *testing.T) {
	_, err := LoadConfig("/nonexistent/path/config.json")
	if err == nil {
		t.Error("LoadConfig() expected error for nonexistent file, got nil")
	}
}

func TestConfigValidate(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr bool
	}{
		{
			name: "valid config",
			config: Config{
				Version: "1.0.0",
				Presets: map[string]Preset{
					"test": {Ratios: []string{"1:1"}},
				},
				Sizes: map[string]SizeConfig{
					"test": {BaseSizes: []int{100}},
				},
				Formats: map[string]Format{
					"jpeg": {Qualities: []int{85}},
				},
			},
			wantErr: false,
		},
		{
			name: "missing version",
			config: Config{
				Presets: map[string]Preset{"test": {Ratios: []string{"1:1"}}},
				Sizes:   map[string]SizeConfig{"test": {BaseSizes: []int{100}}},
				Formats: map[string]Format{"jpeg": {Qualities: []int{85}}},
			},
			wantErr: true,
		},
		{
			name: "missing presets",
			config: Config{
				Version: "1.0.0",
				Sizes:   map[string]SizeConfig{"test": {BaseSizes: []int{100}}},
				Formats: map[string]Format{"jpeg": {Qualities: []int{85}}},
			},
			wantErr: true,
		},
		{
			name: "invalid ratio in preset",
			config: Config{
				Version: "1.0.0",
				Presets: map[string]Preset{
					"test": {Ratios: []string{"invalid"}},
				},
				Sizes:   map[string]SizeConfig{"test": {BaseSizes: []int{100}}},
				Formats: map[string]Format{"jpeg": {Qualities: []int{85}}},
			},
			wantErr: true,
		},
		{
			name: "invalid target dimensions",
			config: Config{
				Version: "1.0.0",
				Presets: map[string]Preset{"test": {Ratios: []string{"1:1"}}},
				Sizes:   map[string]SizeConfig{"test": {BaseSizes: []int{100}}},
				Formats: map[string]Format{"jpeg": {Qualities: []int{85}}},
				Targets: map[string]Target{
					"bad": {Dimensions: []int{0, 100}},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Config.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetCategoryForRatio(t *testing.T) {
	cfg := &Config{
		Presets: map[string]Preset{
			"platform": {Ratios: []string{"2:3", "4:5", "1:1"}},
			"common":   {Ratios: []string{"3:2", "16:9"}},
		},
	}

	tests := []struct {
		ratio    string
		expected string
	}{
		{"2:3", "platform"},
		{"4:5", "platform"},
		{"16:9", "common"},
		{"99:99", "edge"}, // unknown ratio defaults to edge
	}

	for _, tt := range tests {
		t.Run(tt.ratio, func(t *testing.T) {
			got := cfg.GetCategoryForRatio(tt.ratio)
			if got != tt.expected {
				t.Errorf("GetCategoryForRatio(%q) = %q, want %q", tt.ratio, got, tt.expected)
			}
		})
	}
}

