package cmd

import (
	"github.com/rajch/kutti/internal/pkg/kuttilog"
	"github.com/spf13/cobra"
)

// nodestartCmd represents the nodestart command
var nodestartCmd = &cobra.Command{
	Use:           "start NODENAME",
	Short:         "Starts a node",
	Long:          `Starts a node.`,
	Args:          nodenameonlyargs,
	Run:           nodestart,
	SilenceErrors: true,
}

func init() {
	nodeCmd.AddCommand(nodestartCmd)

	nodestartCmd.Flags().StringP("cluster", "c", "", "Cluster name")
}

func nodestart(cmd *cobra.Command, args []string) {
	cluster, err := getCluster(cmd)
	if err != nil {
		kuttilog.Printf(0, "Error: %v", err)
		return
	}

	nodename := args[0]
	node, ok := cluster.Nodes[nodename]
	if !ok {
		kuttilog.Printf(0, "Error: node '%v' not found.", nodename)
		return
	}

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
