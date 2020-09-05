package cmd

import (
	"fmt"

	"github.com/rajch/kutti/internal/pkg/kuttilog"
	"github.com/spf13/cobra"
)

// nodeunforwardportCmd represents the unforwardport command
var nodeunforwardportCmd = &cobra.Command{
	Use:           "unforwardport NODENAME",
	Aliases:       []string{"unpublish", "unforward", "unmap"},
	Short:         "Unforward a node port",
	Long:          `Unforward a node port.`,
	Run:           nodeunforwardportCommand,
	Args:          nodenameonlyargs,
	SilenceErrors: true,
}

func init() {
	nodeCmd.AddCommand(nodeunforwardportCmd)

	nodeunforwardportCmd.Flags().IntP("nodeport", "n", 0, "node port to unmap")
}

func nodeunforwardportCommand(cmd *cobra.Command, args []string) {
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

	nodeport, _ := cmd.Flags().GetInt("nodeport")
	if nodeport == 0 {
		fmt.Println("Error: Please provide a valid nodeport.")
		return
	}

	err = node.UnforwardPort(nodeport)
	if err != nil {
		fmt.Printf("Error: Cannot unforward node port %v: %v.\n", nodeport, err)
		return
	}

	fmt.Printf("Node port %v unforwarded.\n", nodeport)
}
