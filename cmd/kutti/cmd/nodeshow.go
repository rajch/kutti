package cmd

import (
	"github.com/rajch/kutti/cmd/kutti/defaults"
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

	nodeshowCmd.Flags().StringP("cluster", "c", defaults.Getdefault("cluster"), "cluster name")
}

func nodeshowCommand(cmd *cobra.Command, args []string) {
	cluster, err := getCluster(cmd)
	if err != nil {
		kuttilog.Printf(0, "Error: %v.", err)
		return
	}

	nodename := args[0]
	node, ok := cluster.Nodes[nodename]
	if !ok {
		kuttilog.Printf(0, "Error: Node '%s' does not exist.\n", nodename)
		return
	}

	kuttilog.Printf(
		0,
		"Name: %v\nType: %v\nPorts:\n",
		node.Name,
		node.Type,
	)

	for nodeport, hostport := range node.Ports {
		kuttilog.Printf(
			0,
			"  - NodePort: %v\n    HostPort: %v\n",
			nodeport,
			hostport,
		)
	}

	kuttilog.Printf(0, "Status: %v\n", node.Status())
}
