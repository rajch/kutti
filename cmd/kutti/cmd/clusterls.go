package cmd

import (
	"github.com/rajch/kutti/internal/pkg/kuttilog"
	"github.com/rajch/kutti/pkg/clustermanager"
	"github.com/spf13/cobra"
)

// clusterlsCmd represents the clusterls command
var clusterlsCmd = &cobra.Command{
	Use:                   "ls",
	Aliases:               []string{"list"},
	Short:                 "List clusters",
	Long:                  `List clusters.`,
	Run:                   clusterlsCommand,
	DisableFlagsInUseLine: true,
}

func init() {
	clusterCmd.AddCommand(clusterlsCmd)
}

func clusterlsCommand(cmd *cobra.Command, args []string) {
	kuttilog.Printf(0, "%-21.21s  %-10.10s  %-11.11s  %-10.10s\n", "NAME", "DRIVER", "K8S VERSION", "NODES")
	clustermanager.ForEachCluster(func(cluster *clustermanager.Cluster) bool {
		kuttilog.Printf(
			0,
			"%-21.21s  %-10.10s  %-11.11s  %-10d\n",
			defaultdecorate(cluster.Name, "cluster"),
			cluster.DriverName,
			cluster.K8sVersion,
			len(cluster.Nodes),
		)
		return false
	})
}
