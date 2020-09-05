package cmd

import (
	"fmt"

	"github.com/rajch/kutti/pkg/clustermanager"
	"github.com/spf13/cobra"
)

// clustershowCmd represents the clustershow command
var clustershowCmd = &cobra.Command{
	Use:           "show CLUSTERNAME",
	Aliases:       []string{"describe", "inspect", "get"},
	Short:         "Show cluster details",
	Long:          `Show cluster details.`,
	Run:           clustershowCommand,
	Args:          clusternameonlyargs,
	SilenceErrors: true,
}

func init() {
	clusterCmd.AddCommand(clustershowCmd)
}

func clustershowCommand(cmd *cobra.Command, args []string) {
	clustername := args[0]
	cluster, ok := clustermanager.GetCluster(clustername)
	if !ok {
		fmt.Printf("Error: Cluster '%s' does not exist.\n", clustername)
		return
	}

	fmt.Printf(
		"Name: %v\nType: %v\nK8sVersion: %v\nDriver: %v\nNodes:\n",
		cluster.Name,
		cluster.Type,
		cluster.K8sVersion,
		cluster.DriverName,
	)
	for nodename, node := range cluster.Nodes {
		fmt.Printf("  - %v:\n      SSHPort: %v\n", nodename, node.Ports[22])
	}
}
