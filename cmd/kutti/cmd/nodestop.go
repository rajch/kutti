package cmd

import (
	"github.com/rajch/kutti/cmd/kutti/defaults"
	"github.com/rajch/kutti/internal/pkg/kuttilog"
	"github.com/spf13/cobra"
)

// nodestopCmd represents the nodestop command
var nodestopCmd = &cobra.Command{
	Use:   "stop NODENAME",
	Short: "Stop a node",
	Long:  `Stop a node.`,
	Args:  nodenameonlyargs,
	Run:   nodestopCommand,
}

func init() {
	nodeCmd.AddCommand(nodestopCmd)

	nodestopCmd.Flags().StringP("cluster", "c", defaults.Getdefault("cluster"), "cluster name")
}

func nodestopCommand(cmd *cobra.Command, args []string) {
	cluster, err := getCluster(cmd)
	if err != nil {
		kuttilog.Printf(0, "Error: %v.", err)
		return
	}

	nodename := args[0]
	node, ok := cluster.Nodes[nodename]
	if !ok {
		kuttilog.Printf(0, "Error: node '%v' not found.\n", nodename)
		return
	}

	if node.Status() == "Stopped" {
		kuttilog.Printf(0, "Error: Node %s already stopped.", nodename)
		return
	}

	if node.Status() == "Unknown" {
		kuttilog.Printf(0, "Cannot stop node %s: node status unknown.", nodename)
		return
	}

	err = node.Stop()
	if err != nil {
		kuttilog.Printf(0, "Error: %v.", err)
		return
	}

	if kuttilog.V(1) {
		kuttilog.Printf(1, "Node '%s' stopped.", nodename)
	} else {
		kuttilog.Println(0, node.Name)
	}
}
