package cmd

import (
	"github.com/rajch/kutti/cmd/kutti/defaults"
	"github.com/rajch/kutti/internal/pkg/kuttilog"
	"github.com/spf13/cobra"
)

// nodeunforwardportCmd represents the unforwardport command
var nodeunforwardportCmd = &cobra.Command{
	Use:           "unforwardport NODENAME",
	Aliases:       []string{"unpublish", "unforward", "portunforward", "unmap"},
	Short:         "Unforward a node port",
	Long:          `Unforward a node port.`,
	Run:           nodeunforwardportCommand,
	Args:          nodenameonlyargs,
	SilenceErrors: true,
}

func init() {
	nodeCmd.AddCommand(nodeunforwardportCmd)

	nodeunforwardportCmd.Flags().StringP("cluster", "c", defaults.Getdefault("cluster"), "cluster name")
	nodeunforwardportCmd.Flags().IntP("nodeport", "n", 0, "node port to unmap")
}

func nodeunforwardportCommand(cmd *cobra.Command, args []string) {
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

	nodeport, _ := cmd.Flags().GetInt("nodeport")
	if nodeport == 0 {
		kuttilog.Println(0, "Error: Please provide a valid nodeport.")
		return
	}

	err = node.UnforwardPort(nodeport)
	if err != nil {
		kuttilog.Printf(0, "Error: Cannot unforward node port %v: %v.\n", nodeport, err)
		return
	}

	if kuttilog.V(1) {
		kuttilog.Printf(1, "Node port %v unforwarded.\n", nodeport)
	} else {
		kuttilog.Println(0, nodeport)
	}
}
