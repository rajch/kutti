package cmd

import (
	"github.com/rajch/kutti/internal/pkg/kuttilog"
	"github.com/spf13/cobra"
)

// nodermCmd represents the noderm command
var nodermCmd = &cobra.Command{
	Use:           "rm NODENAME",
	Aliases:       []string{"delete", "remove"},
	Short:         "Delete a node",
	Long:          `Delete a node.`,
	Args:          nodenameonlyargs,
	Run:           noderm,
	SilenceErrors: true,
}

func init() {
	nodeCmd.AddCommand(nodermCmd)

	nodermCmd.Flags().StringP("cluster", "c", "", "Cluster name")
	nodermCmd.Flags().BoolP("force", "f", false, "Forcibly delete running nodes.")
}

func noderm(cmd *cobra.Command, args []string) {
	cluster, err := getCluster(cmd)
	if err != nil {
		kuttilog.Printf(0, "Error: %v", err)
		return
	}

	nodename := args[0]
	forceflag, _ := cmd.Flags().GetBool("force")

	err = cluster.DeleteNode(nodename, forceflag)
	if err != nil {
		kuttilog.Printf(0, "Error: Could not delete node '%s': %v.", nodename, err)
		return
	}

	if kuttilog.V(1) {
		kuttilog.Printf(1, "Node '%s' deleted.", nodename)
	} else {
		kuttilog.Println(0, nodename)
	}

}
