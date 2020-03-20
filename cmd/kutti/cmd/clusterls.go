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
	fmt.Println("Name\tDriver\tK8s Version\tNodes")
	clustermanager.Load()
	clusters := clustermanager.Clusters()
	for _, cluster := range clusters {
		fmt.Printf("%v\t%v\t%v\t%v\n", cluster.Name, cluster.DriverName, cluster.K8sVersion, len(cluster.Nodes))
	}
}
