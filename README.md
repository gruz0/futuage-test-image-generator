# FutuAge Test Image Generator

A fast, standalone CLI tool for generating comprehensive test images with various aspect ratios, sizes, and formats for testing the FutuAge asset processing pipeline.

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

## Features

- ⚡ **Blazing Fast**: Generate 246 images in 4.6 seconds (~54 images/sec)
- 📊 **Comprehensive Coverage**: All platform ratios, sizes, formats, and edge cases
- 📝 **Self-Documenting**: Metadata baked into images and manifest.json
- 📁 **Organized Output**: Clean directory structure (ratios/, targets/, edge-cases/)
- 🚀 **Single Binary**: No runtime dependencies, easy distribution
- ⚙️ **Parallel Processing**: Concurrent generation using goroutines (10 workers)
- 🔧 **Flexible Configuration**: JSON-based, extensible for new platforms
- 🎨 **Visual Design**: Grid patterns, text overlays, corner markers, borders

## Performance

Real-world performance metrics:

```
Full Test Suite: 246 images in 4.58 seconds
Generation Rate: 53.8 images/second
Total Size:      14 MB
Concurrency:     10 parallel workers
Target Met:      ✅ (2× faster than 10s target)
```

## Installation

### Download Pre-built Binary

```bash
# macOS (ARM64)
curl -L https://github.com/gruz0/futuage-test-image-generator/releases/latest/download/futuage-test-image-gen-darwin-arm64 -o /usr/local/bin/futuage-test-image-gen
chmod +x /usr/local/bin/futuage-test-image-gen

# macOS (Intel)
curl -L https://github.com/gruz0/futuage-test-image-generator/releases/latest/download/futuage-test-image-gen-darwin-amd64 -o /usr/local/bin/futuage-test-image-gen
chmod +x /usr/local/bin/futuage-test-image-gen

# Linux
curl -L https://github.com/gruz0/futuage-test-image-generator/releases/latest/download/futuage-test-image-gen-linux-amd64 -o /usr/local/bin/futuage-test-image-gen
chmod +x /usr/local/bin/futuage-test-image-gen
```

### Build from Source

```bash
git clone https://github.com/gruz0/futuage-test-image-generator.git
cd futuage-test-image-generator
go build -o futuage-test-image-gen .
mv futuage-test-image-gen /usr/local/bin/
```

## Quick Start

```bash
# Generate all standard test images
futuage-test-image-gen generate --output ./test-images/

# List all available presets
futuage-test-image-gen list

# Show version
futuage-test-image-gen --version
```

## Usage

### Generate Command

```bash
# Generate all standard test images (default config)
futuage-test-image-gen generate --output ./test-images/

# Generate specific ratio categories
futuage-test-image-gen generate --ratios platform --output ./test-images/
futuage-test-image-gen generate --ratios common,edge --output ./test-images/

# Generate specific size categories
futuage-test-image-gen generate --sizes medium,large --output ./test-images/

# Generate specific formats only
futuage-test-image-gen generate --formats jpeg,png --output ./test-images/

# Generate with custom configuration file
futuage-test-image-gen generate --config ./custom-specs.json --output ./test-images/

# Combine filters
futuage-test-image-gen generate \
  --ratios platform \
  --sizes medium \
  --formats jpeg \
  --output ./test-images/
```

### List Command

```bash
# Show all available presets
futuage-test-image-gen list
```

Output:

```
Available Presets:

Ratio Presets:
  platform: 2:3, 4:5, 1:1, 9:16, 1.91:1 (Platform-recommended ratios)
  common:   3:2, 4:3, 16:9, 5:4, 21:9 (Common photo ratios)
  edge:     1:2, 2:1, 1:3, 3:1 (Edge case ratios)

Size Categories:
  tiny:    100, 150, 200 px
  small:   500, 640, 800 px
  medium:  1000, 1080, 1200, 1500 px
  large:   2000, 2160, 3000 px
  xlarge:  4000, 4096, 5000 px

Formats:
  jpeg: Q60, Q82, Q95
  png:  Q95
  webp: Q82, Q90

Platform Targets:
  PINTEREST_2_3  Pinterest  1000×1500 (2:3)
  IG_FEED_4_5    Instagram  1080×1350 (4:5)
  IG_FEED_1_1    Instagram  1080×1080 (1:1)
  IG_STORY       Instagram  1080×1920 (9:16)
  TIKTOK_9_16    TikTok     1080×1920 (9:16)
  LI_1_1         LinkedIn   1200×1200 (1:1)
  LI_1_91_1      LinkedIn   1200×628  (1.91:1)
```

## Output Structure

```
test-images/
├── ratios/                           # Organized by aspect ratio
│   ├── 2-3/
│   │   ├── tiny_200x300_jpeg_q82.jpg
│   │   ├── medium_1000x1500_jpeg_q85.jpg
│   │   └── ...
│   ├── 4-5/
│   ├── 1-1/
│   └── ...
│
├── targets/                          # Exact platform specifications
│   ├── PINTEREST_2_3_1000x1500_jpeg_q85.jpg
│   ├── IG_FEED_4_5_1080x1350_jpeg_q85.jpg
│   └── ...
│
├── edge-cases/                       # Edge case scenarios
│   ├── too-small_50x50_jpeg_q82.jpg
│   ├── max-res-square_4096x4096_jpeg_q95.jpg
│   └── ...
│
└── manifest.json                     # Complete metadata for all images
```

## Image Visual Design

Each generated image includes:

- **Grid Pattern Background**: 100px grid with category-based colors
  - Platform ratios: Blue (#4A90E2)
  - Common ratios: Green (#7ED321)
  - Edge cases: Orange (#F5A623)
- **Text Overlay**: Centered metadata (dimensions, ratio, format, quality, size category)
- **Corner Markers**: TL, TR, BL, BR labels for orientation
- **2px Border**: Clearly delineates image edges

## Manifest.json

The generator creates a comprehensive `manifest.json` file with metadata for all images:

```json
{
  "generated_at": "2025-12-04T18:54:37Z",
  "tool_version": "1.0.0",
  "config_version": "1.0.0",
  "total_images": 246,
  "images": [
    {
      "filename": "ratios/2-3/medium_666x1000_jpeg_q60.jpg",
      "category": "ratios",
      "subcategory": "2-3",
      "width": 666,
      "height": 1000,
      "ratio": "2:3",
      "ratio_decimal": 0.667,
      "format": "jpeg",
      "quality": 60,
      "file_size_bytes": 43010,
      "size_category": "medium"
    }
  ]
}
```

Use this manifest for programmatic test validation in your integration tests.

## Usage Examples

### Example 1: Quick Test Set for Development

Generate a minimal set for quick testing:

```bash
futuage-test-image-gen generate \
  --ratios platform \
  --sizes medium \
  --formats jpeg \
  --output ./quick-test/
```

Result: ~10 images in <1 second

### Example 2: Instagram-Specific Testing

Generate only Instagram-relevant images:

```bash
# Generate platform ratios (includes IG 4:5, 1:1, 9:16)
futuage-test-image-gen generate \
  --ratios platform \
  --output ./instagram-tests/
```

### Example 3: Edge Case Testing Only

Test extreme scenarios:

```bash
futuage-test-image-gen generate \
  --ratios edge \
  --output ./edge-case-tests/
```

Result: 1:2, 2:1, 1:3, 3:1 ratios + extreme dimensions

### Example 4: Performance Testing (All Images)

Full test suite for comprehensive validation:

```bash
futuage-test-image-gen generate --output ./test-images/
```

Result: 246 images covering all scenarios

### Example 5: PNG Only for Transparency Testing

```bash
futuage-test-image-gen generate \
  --formats png \
  --output ./png-tests/
```

Result: All ratios and sizes, PNG format only

## Configuration

See [configs/default.json](configs/default.json) for the complete default configuration.

### Custom Configuration Example

Create a custom configuration for specific needs:

```json
{
  "version": "1.0.0",
  "presets": {
    "custom": {
      "description": "Custom ratios for testing",
      "ratios": ["2:3", "1:1"]
    }
  },
  "sizes": {
    "test": {
      "description": "Test sizes",
      "base_sizes": [1000, 2000]
    }
  },
  "formats": {
    "jpeg": {
      "qualities": [85],
      "mime_type": "image/jpeg",
      "extension": ".jpg"
    }
  }
}
```

Use it with:

```bash
futuage-test-image-gen generate \
  --config custom-config.json \
  --output ./custom-tests/
```

## Development

### Prerequisites

- Go 1.21 or higher

### Building

```bash
# Build for current platform
go build -o futuage-test-image-gen .

# Cross-compile for all platforms
GOOS=darwin GOARCH=arm64 go build -o dist/futuage-test-image-gen-darwin-arm64 .
GOOS=darwin GOARCH=amd64 go build -o dist/futuage-test-image-gen-darwin-amd64 .
GOOS=linux GOARCH=amd64 go build -o dist/futuage-test-image-gen-linux-amd64 .
GOOS=windows GOARCH=amd64 go build -o dist/futuage-test-image-gen-windows-amd64.exe .
```

### Running

```bash
# Run without building
go run . generate --output ./test-images/

# Build and run
go build -o futuage-test-image-gen .
./futuage-test-image-gen generate --output ./test-images/
```

## Integration with FutuAge

### Generating Test Images

```bash
# In FutuAge repository
futuage-test-image-gen generate --output tests/fixtures/images/
```

### CI/CD Integration

```yaml
# .github/workflows/test.yml
- name: Install test image generator
  run: |
    curl -L https://github.com/gruz0/futuage-test-image-generator/releases/latest/download/futuage-test-image-gen-linux-amd64 -o /usr/local/bin/futuage-test-image-gen
    chmod +x /usr/local/bin/futuage-test-image-gen

- name: Generate test images
  run: |
    futuage-test-image-gen generate --output tests/fixtures/images/

- name: Run integration tests
  run: npm run test:integration
```

## License

MIT License - see [LICENSE](LICENSE) file for details

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
