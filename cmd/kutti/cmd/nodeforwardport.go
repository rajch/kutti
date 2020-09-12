package cmd

import (
	"github.com/rajch/kutti/cmd/kutti/defaults"
	"github.com/rajch/kutti/internal/pkg/kuttilog"
	"github.com/rajch/kutti/pkg/clustermanager"
	"github.com/spf13/cobra"
)

// nodeforwardportCmd represents the forwardport command
var nodeforwardportCmd = &cobra.Command{
	Use:           "forwardport NODENAME",
	Aliases:       []string{"publish", "forward", "portforward", "map"},
	Short:         "Forward a node port to a host port",
	Long:          `Forward a node port to a host port.`,
	Run:           nodeforwardportCommand,
	Args:          nodenameonlyargs,
	SilenceErrors: true,
}

func init() {
	nodeCmd.AddCommand(nodeforwardportCmd)

	nodeforwardportCmd.Flags().StringP("cluster", "c", defaults.Getdefault("cluster"), "cluster name")
	nodeforwardportCmd.Flags().IntP("hostport", "p", 0, "port on the host")
	nodeforwardportCmd.Flags().IntP("nodeport", "n", 0, "port on the node")
}

func nodeforwardportCommand(cmd *cobra.Command, args []string) {
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

	nodeport, _ := cmd.Flags().GetInt("nodeport")
	if !clustermanager.IsValidPort(nodeport) {
		kuttilog.Println(0, "Error: Please provide a valid nodeport. Valid ports are between 1 and 65535.")
		return
	}

	hostport, _ := cmd.Flags().GetInt("hostport")
	if !clustermanager.IsValidPort(hostport) {
		kuttilog.Println(0, "Error: Please provide a valid hostport. Valid ports are between 1 and 65535.")
		return
	}

	err = cluster.CheckHostport(hostport)
	if err != nil {
		kuttilog.Printf(0, "Error: Cannot forward to host port %v: %v.\n", hostport, err)
		return
	}

	err = node.ForwardPort(hostport, nodeport)
	if err != nil {
		kuttilog.Printf(
			0,
			"Error: Could not forward node port %v to host port %v: %v.\n",
			nodeport,
			hostport,
			err,
		)
		return
	}

	if kuttilog.V(1) {
		kuttilog.Printf(
			1,
			"Forwarded node port %v to host port %v.\n",
			nodeport,
			hostport,
		)
	} else {
		kuttilog.Println(0, hostport)
	}
}
