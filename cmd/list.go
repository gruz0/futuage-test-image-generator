package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available presets",
	Long: `Display all available presets including:
- Ratio presets (platform, common, edge)
- Size categories (tiny, small, medium, large, xlarge)
- Format specifications (jpeg, png, webp)
- Platform targets (Pinterest, Instagram, LinkedIn, TikTok)`,
	Run: runList,
}

func runList(cmd *cobra.Command, args []string) {
	fmt.Println("Available Presets:")
	fmt.Println()
	fmt.Println("Ratio Presets:")
	fmt.Println("  platform: 2:3, 4:5, 1:1, 9:16, 1.91:1 (Platform-recommended ratios)")
	fmt.Println("  common:   3:2, 4:3, 16:9, 5:4, 21:9 (Common photo ratios)")
	fmt.Println("  edge:     1:2, 2:1, 1:3, 3:1 (Edge case ratios)")
	fmt.Println()
	fmt.Println("Size Categories:")
	fmt.Println("  tiny:    100, 150, 200 px")
	fmt.Println("  small:   500, 640, 800 px")
	fmt.Println("  medium:  1000, 1080, 1200, 1500 px")
	fmt.Println("  large:   2000, 2160, 3000 px")
	fmt.Println("  xlarge:  4000, 4096, 5000 px")
	fmt.Println()
	fmt.Println("Formats:")
	fmt.Println("  jpeg: Q60, Q82, Q95")
	fmt.Println("  png:  Q95")
	fmt.Println("  webp: Q82, Q90")
	fmt.Println()
	fmt.Println("Platform Targets:")
	fmt.Println("  PINTEREST_2_3  Pinterest  1000×1500 (2:3)")
	fmt.Println("  IG_FEED_4_5    Instagram  1080×1350 (4:5)")
	fmt.Println("  IG_FEED_1_1    Instagram  1080×1080 (1:1)")
	fmt.Println("  IG_STORY       Instagram  1080×1920 (9:16)")
	fmt.Println("  TIKTOK_9_16    TikTok     1080×1920 (9:16)")
	fmt.Println("  LI_1_1         LinkedIn   1200×1200 (1:1)")
	fmt.Println("  LI_1_91_1      LinkedIn   1200×628  (1.91:1)")
}
