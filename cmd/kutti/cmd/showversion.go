package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// KuttiVersion contains the kutti CLI version string.
var KuttiVersion string

// showversionCmd represents the showversion command
var showversionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show kutti version",
	Long:  `Show kutti version.`,
	Run:   showversionCommand,
}

func init() {
	rootCmd.AddCommand(showversionCmd)
}

func showversionCommand(cmd *cobra.Command, args []string) {
	fmt.Printf("kutti version %s\n", KuttiVersion)
}
