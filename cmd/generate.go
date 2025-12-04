package cmd

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/gruz0/futuage-test-image-generator/internal/config"
	"github.com/gruz0/futuage-test-image-generator/internal/filesystem"
	"github.com/gruz0/futuage-test-image-generator/internal/generator"
	"github.com/gruz0/futuage-test-image-generator/internal/manifest"
	"github.com/spf13/cobra"
)

var (
	outputDir  string
	configFile string
	ratios     []string
	sizes      []string
	formats    []string
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate test images",
	Long: `Generate test images based on configuration.

Examples:
  # Generate all standard test images
  futuage-test-image-gen generate --output ./test-images/

  # Generate specific ratio categories
  futuage-test-image-gen generate --ratios platform --output ./test-images/

  # Generate with custom configuration
  futuage-test-image-gen generate --config ./custom-config.json --output ./test-images/`,
	RunE: runGenerate,
}

func init() {
	generateCmd.Flags().StringVarP(&outputDir, "output", "o", "./test-images", "Output directory for generated images")
	generateCmd.Flags().StringVarP(&configFile, "config", "c", "", "Custom configuration file (optional)")
	generateCmd.Flags().StringSliceVar(&ratios, "ratios", []string{}, "Ratio categories to generate (platform, common, edge)")
	generateCmd.Flags().StringSliceVar(&sizes, "sizes", []string{}, "Size categories to generate (tiny, small, medium, large, xlarge)")
	generateCmd.Flags().StringSliceVar(&formats, "formats", []string{}, "Format types to generate (jpeg, png, webp)")
}

func runGenerate(cmd *cobra.Command, args []string) error {
	startTime := time.Now()

	fmt.Println("🖼  FutuAge Test Image Generator")
	fmt.Println()

	// 1. Load configuration
	fmt.Printf("Loading configuration...\n")
	cfg, err := config.LoadConfig(configFile)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	if configFile != "" {
		fmt.Printf("  ✓ Loaded custom config: %s\n", configFile)
	} else {
		fmt.Printf("  ✓ Loaded default config (version %s)\n", cfg.Version)
	}

	// 2. Create and validate filters
	filters := config.NewFilters(ratios, sizes, formats)
	if err := filters.Validate(cfg); err != nil {
		return fmt.Errorf("invalid filters: %w", err)
	}

	if !filters.IsEmpty() {
		fmt.Printf("  ✓ Filters: %s\n", filters.Summary())
	} else {
		fmt.Printf("  ✓ No filters (generating all)\n")
	}
	fmt.Println()

	// 3. Ensure output directory structure
	fmt.Printf("Setting up output directory: %s\n", outputDir)
	if err := filesystem.EnsureDirectoryStructure(outputDir); err != nil {
		return fmt.Errorf("failed to create directory structure: %w", err)
	}
	fmt.Printf("  ✓ Directory structure created\n")
	fmt.Println()

	// 4. Build image specifications
	fmt.Printf("Building image specifications...\n")
	builder := config.NewSpecBuilder(cfg, filters, outputDir)
	specs, err := builder.BuildSpecs()
	if err != nil {
		return fmt.Errorf("failed to build specs: %w", err)
	}
	fmt.Printf("  ✓ Generated %d image specifications\n", len(specs))
	fmt.Println()

	// 5. Generate images in parallel
	fmt.Printf("Generating images...\n")
	orchestrator := generator.NewOrchestrator(10) // 10 concurrent workers

	// Set up progress callback
	orchestrator.SetProgressCallback(func(completed, total int, elapsed time.Duration) {
		percentage := float64(completed) / float64(total) * 100
		fmt.Printf("\r  Progress: %d/%d (%.1f%%) - %.1fs elapsed", completed, total, percentage, elapsed.Seconds())
	})

	results, err := orchestrator.GenerateAll(specs)
	fmt.Println() // New line after progress
	if err != nil {
		fmt.Printf("  ⚠ Warning: %v\n", err)
	}

	fmt.Printf("  ✓ Completed: %d images\n", orchestrator.Stats.Completed)
	if orchestrator.Stats.Failed > 0 {
		fmt.Printf("  ✗ Failed: %d images\n", orchestrator.Stats.Failed)
	}
	fmt.Printf("  ⏱ Duration: %.2fs (%.1f images/sec)\n",
		orchestrator.Stats.Duration().Seconds(),
		orchestrator.Stats.ImagesPerSecond())
	fmt.Println()

	// 6. Generate manifest
	fmt.Printf("Generating manifest...\n")
	mf := manifest.NewManifest(version, cfg.Version)
	for _, result := range results {
		if result.Error == nil {
			mf.AddImage(result.Spec, result.FileSize)
		}
	}

	manifestPath := filepath.Join(outputDir, "manifest.json")
	if err := mf.Write(manifestPath); err != nil {
		return fmt.Errorf("failed to write manifest: %w", err)
	}
	fmt.Printf("  ✓ Manifest written to: %s\n", manifestPath)
	fmt.Println()

	// 7. Print summary
	fmt.Println("Summary:")
	fmt.Println("--------")
	fmt.Print(mf.Summary())
	fmt.Printf("\nTotal time: %.2fs\n", time.Since(startTime).Seconds())
	fmt.Printf("Output directory: %s\n", outputDir)
	fmt.Println()
	fmt.Println("✓ Done!")

	return nil
}
