package cmd

import (
	"github.com/rajch/kutti/internal/pkg/kuttilog"
	"github.com/spf13/cobra"
)

// nodestopCmd represents the nodestop command
var nodestopCmd = &cobra.Command{
	Use:   "stop NODENAME",
	Short: "Stops a node",
	Long:  `Stops a node.`,
	Args:  nodenameonlyargs,
	Run:   nodestop,
}

func init() {
	nodeCmd.AddCommand(nodestopCmd)

	nodestopCmd.Flags().StringP("cluster", "c", "", "Cluster name")
}

func nodestop(cmd *cobra.Command, args []string) {
	cluster, err := getCluster(cmd)
	if err != nil {
		kuttilog.Printf(0, "Error: %v", err)
		return
	}

	nodename := args[0]
	node, ok := cluster.Nodes[nodename]
	if !ok {
		kuttilog.Printf(0, "Error: node '%v' not found.\n", nodename)
		return
	}

	err = node.Stop()
	if err != nil {
		kuttilog.Printf(0, "Error: %v", err)
		return
	}

	if kuttilog.V(1) {
		kuttilog.Printf(1, "Node '%s' stopped.", nodename)
	} else {
		kuttilog.Println(0, node.Name)
	}
}
