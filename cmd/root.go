package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	version = "1.0.0"
)

var rootCmd = &cobra.Command{
	Use:   "futuage-test-image-gen",
	Short: "Test image generator for FutuAge asset processing pipeline",
	Long: `A fast, standalone CLI tool for generating comprehensive test images
with various aspect ratios, sizes, and formats for testing the FutuAge
asset processing pipeline.`,
	Version: version,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(generateCmd)
	rootCmd.AddCommand(listCmd)
}
