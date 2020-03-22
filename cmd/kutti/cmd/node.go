package cmd

import (
	"fmt"

	"github.com/rajch/kutti/pkg/clustermanager"
	"github.com/spf13/cobra"
)

// nodeCmd represents the node command
var nodeCmd = &cobra.Command{
	Use:   "node",
	Short: "Manage nodes",
	Long:  `Manage nodes.`,
}

func init() {
	rootCmd.AddCommand(nodeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// nodeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// nodeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func getCluster(cmd *cobra.Command) (*clustermanager.Cluster, error) {
	var cluster *clustermanager.Cluster
	clustername, _ := cmd.Flags().GetString("cluster")

	if clustername == "" {
		cluster = clustermanager.DefaultCluster()
		if cluster == nil {
			return nil,
				fmt.Errorf("no cluster specified and default cluster not set. Use --cluster, or select a default cluster using 'kutti cluster setdefault'")

		}
	} else {
		cluster, _ = clustermanager.GetCluster(clustername)
		if cluster == nil {
			return nil,
				fmt.Errorf("cluster '%v' not found", clustername)

		}
	}

	return cluster, nil
}
