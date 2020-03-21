package cmd

import (
	"fmt"

	"github.com/rajch/kutti/pkg/clustermanager"
	"github.com/spf13/cobra"
)

// clusterlsCmd represents the clusterls command
var clusterlsCmd = &cobra.Command{
	Use:     "ls",
	Aliases: []string{"list"},
	Short:   "List clusters",
	Long:    `List clusters.`,
	Run:     clusterls,
}

func init() {
	clusterCmd.AddCommand(clusterlsCmd)
}

func clusterls(cmd *cobra.Command, args []string) {
	fmt.Printf("%-21.21s  %-10.10s  %-11.11s  %-10.10s\n", "NAME", "DRIVER", "K8S VERSION", "NODES")
	clustermanager.Load()
	clustermanager.ForEachCluster(func(cluster *clustermanager.Cluster) bool {
		fmt.Printf("%-21.21s  %-10.10s  %-11.11s  %-10d\n", cluster.Name, cluster.DriverName, cluster.K8sVersion, len(cluster.Nodes))
		return false
	})
}
