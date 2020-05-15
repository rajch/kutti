package cmd

import (
	"github.com/spf13/cobra"
)

// clusterCmd represents the cluster command
var clusterCmd = &cobra.Command{
	Use:           "cluster",
	Short:         "Manage clusters",
	SilenceErrors: true,
}

func init() {
	rootCmd.AddCommand(clusterCmd)
}
