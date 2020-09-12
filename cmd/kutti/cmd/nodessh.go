package cmd

import (
	"fmt"

	"github.com/rajch/kutti/internal/pkg/kuttilog"
	"github.com/rajch/kutti/pkg/sshclient"
	"github.com/spf13/cobra"
)

// nodesshCmd represents the nodessh command
var nodesshCmd = &cobra.Command{
	Use:                   "ssh NODENAME",
	Short:                 "Open an SSH connection to the node",
	Long:                  `Open an SSH connection to the node.`,
	Run:                   nodesshCommand,
	Args:                  nodenameonlyargs,
	SilenceErrors:         true,
	DisableFlagsInUseLine: true,
}

func init() {
	nodeCmd.AddCommand(nodesshCmd)

	nodesshCmd.Flags().StringP("cluster", "c", "", "cluster name")
}

func nodesshCommand(cmd *cobra.Command, args []string) {
	cluster, err := getCluster(cmd)
	if err != nil {
		kuttilog.Printf(0, "Error: %v", err)
		return
	}

	nodename := args[0]
	node, ok := cluster.Nodes[nodename]
	if !ok {
		fmt.Printf("Error: Node '%s' not found.\n", nodename)
		return
	}

	if node.Status() != "Running" {
		fmt.Printf("Error: Node '%s' is not running.\n", nodename)
		return
	}

	sshport, ok := node.Ports[22]
	if !ok {
		fmt.Printf("Error: The SSH port of node '%s' has not been forwarded.\n", nodename)
		return
	}

	kuttilog.Printf(2, "Connecting to node %s...", nodename)
	address := fmt.Sprintf("localhost:%v", sshport)
	client := sshclient.New("user1", "Pass@word1")

	client.RunInterativeShell(address)
}
