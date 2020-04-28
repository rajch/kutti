package cmd

import (
	"fmt"

	"github.com/rajch/kutti/pkg/clustermanager"
	"github.com/spf13/cobra"
)

// clustercleardefaultCmd represents the clustercleardefault command
var clustercleardefaultCmd = &cobra.Command{
	Use:   "cleardefault",
	Short: "Clears the default cluster.",
	Long:  `Clears the default cluster.`,
	Args:  cobra.NoArgs,
	Run:   clustercleardefault,
}

func init() {
	clusterCmd.AddCommand(clustercleardefaultCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// clustercleardefaultCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// clustercleardefaultCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func clustercleardefault(cmd *cobra.Command, args []string) {
	clustermanager.ClearDefaultCluster()
	fmt.Println("Default cluster cleared.")
}
