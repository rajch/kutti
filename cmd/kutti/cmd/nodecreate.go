package cmd

import (
	"github.com/rajch/kutti/internal/pkg/kuttilog"

	"github.com/spf13/cobra"
)

// nodecreateCmd represents the nodecreate command
var nodecreateCmd = &cobra.Command{
	Use:           "create NODENAME",
	Aliases:       []string{"add"},
	Short:         "Create a node",
	Long:          `Create a node.`,
	Args:          nodenameonlyargs,
	Run:           nodecreate,
	SilenceErrors: true,
}

func init() {
	nodeCmd.AddCommand(nodecreateCmd)

	nodecreateCmd.Flags().StringP("cluster", "c", "", "Cluster name")

}

func nodecreate(cmd *cobra.Command, args []string) {
	cluster, err := getCluster(cmd)
	if err != nil {
		kuttilog.Printf(0, "Error: %v", err)
		return
	}

	nodename := args[0]
	kuttilog.Printf(1, "Creating node '%s' on cluster %s...", nodename, cluster.Name)
	_, err = cluster.NewUninitializedNode(nodename)
	if err != nil {
		kuttilog.Printf(0, "Error: Could not create node %v: %v.", nodename, err)
		return
	}

	if kuttilog.V(1) {
		kuttilog.Printf(1, "Node '%s' created.", nodename)
	} else {
		kuttilog.Println(0, nodename)
	}
}
