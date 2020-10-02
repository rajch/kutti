package cmd

import (
	"github.com/rajch/kutti/cmd/kutti/defaults"
	"github.com/rajch/kutti/internal/pkg/kuttilog"
	"github.com/spf13/cobra"
)

// nodestartCmd represents the nodestart command
var nodestartCmd = &cobra.Command{
	Use:           "start NODENAME",
	Short:         "Start a node",
	Long:          `Start a node.`,
	Args:          nodenameonlyargs,
	Run:           nodestartCommand,
	SilenceErrors: true,
}

func init() {
	nodeCmd.AddCommand(nodestartCmd)

	nodestartCmd.Flags().StringP("cluster", "c", defaults.Getdefault("cluster"), "cluster name")
}

func nodestartCommand(cmd *cobra.Command, args []string) {
	cluster, err := getCluster(cmd)
	if err != nil {
		kuttilog.Printf(0, "Error: %v.", err)
		return
	}

	nodename := args[0]
	node, ok := cluster.Nodes[nodename]
	if !ok {
		kuttilog.Printf(0, "Error: node '%v' not found.", nodename)
		return
	}

	if node.Status() == "Running" {
		kuttilog.Printf(0, "Error: Node %s already running.", nodename)
		return
	}

	if node.Status() == "Unknown" {
		kuttilog.Printf(0, "Cannot start node %s: node status unknown.", nodename)
		return
	}

	kuttilog.Printf(2, "Starting node %s...", nodename)
	err = node.Start()
	if err != nil {
		kuttilog.Printf(0, "Error: %v", err)
		return
	}

	if kuttilog.V(1) {
		kuttilog.Printf(1, "Node '%s' started.", nodename)
	} else {
		kuttilog.Println(0, nodename)
	}
}
