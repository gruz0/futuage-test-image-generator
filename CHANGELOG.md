# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2025-12-04

### Added

- Initial release of FutuAge Test Image Generator
- CLI tool with `generate`, `list`, and `version` commands
- Support for 3 ratio categories (platform, common, edge)
- Support for 5 size categories (tiny, small, medium, large, xlarge)
- Support for 3 formats (JPEG, PNG, WebP)
- 7 platform targets (Pinterest, Instagram, LinkedIn, TikTok)
- 5 edge case scenarios (too-small, max-res, extreme ratios)
- Grid pattern backgrounds with category-based colors
- Text overlay with image metadata
- Corner markers (TL, TR, BL, BR)
- 2px category-colored borders
- Parallel image generation with goroutines (10 workers)
- manifest.json generation with complete metadata
- Filter support for ratios, sizes, and formats
- Custom configuration file support
- Embedded default configuration
- Real-time progress reporting
- Performance metrics display

### Performance

- Generates 246 images in 4.58 seconds (~54 images/sec)
- 2× faster than 10-second target
- Total output size: ~14 MB

### Documentation

- Comprehensive README with usage examples
- Example configuration files

## [Unreleased]

### Planned Features

- GitHub Actions workflow for releases
- Pre-built binaries for macOS, Linux, Windows
- Docker image for containerized usage
- Additional platform targets (Twitter, Facebook, etc.)
- Quality presets (low, medium, high)
- Batch configuration support
- Progress bar with ETA
- Dry-run mode
- Image comparison utilities

---

**Legend:**

- `Added` for new features
- `Changed` for changes in existing functionality
- `Deprecated` for soon-to-be removed features
- `Removed` for now removed features
- `Fixed` for any bug fixes
- `Security` in case of vulnerabilities
