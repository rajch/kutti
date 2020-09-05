package cmd

import (
	"fmt"

	"github.com/rajch/kutti/internal/pkg/kuttilog"
	"github.com/spf13/cobra"
)

// nodeforwardportCmd represents the forwardport command
var nodeforwardportCmd = &cobra.Command{
	Use:           "forwardport NODENAME",
	Aliases:       []string{"publish", "forward", "map"},
	Short:         "Forward a node port to a host port",
	Long:          `Forward a node port to a host port.`,
	Run:           nodeforwardportCommand,
	Args:          nodenameonlyargs,
	SilenceErrors: true,
}

func init() {
	nodeCmd.AddCommand(nodeforwardportCmd)

	nodeforwardportCmd.Flags().IntP("hostport", "p", 0, "port on the host")
	nodeforwardportCmd.Flags().IntP("nodeport", "n", 0, "port on the node")
}

func nodeforwardportCommand(cmd *cobra.Command, args []string) {
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

	hostport, _ := cmd.Flags().GetInt("hostport")
	if hostport == 0 {
		fmt.Println("Error: Please provide a valid hostport.")
		return
	}

	err = cluster.CheckHostport(hostport)
	if err != nil {
		fmt.Printf("Error: Cannot forward to host port %v: %v.\n", hostport, err)
		return
	}

	err = node.ForwardPort(hostport, nodeport)
	if err != nil {
		fmt.Printf(
			"Error: Could not forward node port %v to host port %v: %v\n",
			nodeport,
			hostport,
			err,
		)
		return
	}

	kuttilog.Printf(
		2,
		"Forwarded node port %v to host port %v.\n",
		nodeport,
		hostport,
	)
}
