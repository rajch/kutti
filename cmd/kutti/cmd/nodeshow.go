package cmd

import (
	"fmt"

	"github.com/rajch/kutti/internal/pkg/kuttilog"
	"github.com/spf13/cobra"
)

// nodeshowCmd represents the nodeshow command
var nodeshowCmd = &cobra.Command{
	Use:           "show NODENAME",
	Aliases:       []string{"describe", "inspect", "get"},
	Short:         "Show node details",
	Long:          `Show node details.`,
	Run:           nodeshowCommand,
	Args:          nodenameonlyargs,
	SilenceErrors: true,
}

func init() {
	nodeCmd.AddCommand(nodeshowCmd)

	nodeshowCmd.Flags().StringP("cluster", "c", "", "cluster name")
}

func nodeshowCommand(cmd *cobra.Command, args []string) {
	cluster, err := getCluster(cmd)
	if err != nil {
		kuttilog.Printf(0, "Error: %v", err)
		return
	}

	nodename := args[0]
	node, ok := cluster.Nodes[nodename]
	if !ok {
		fmt.Printf("Error: Node '%s' does not exist.\n", nodename)
		return
	}

	fmt.Printf(
		"Name: %v\nType: %v\nPorts:\n",
		node.Name,
		node.Type,
	)

	for nodeport, hostport := range node.Ports {
		fmt.Printf(
			"  - NodePort: %v\n    HostPort: %v\n",
			nodeport,
			hostport,
		)
	}

	fmt.Print("Status: ")
	fmt.Printf("%v\n", node.Status())
}
