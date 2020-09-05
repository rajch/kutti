package cmd

import (
	"github.com/rajch/kutti/cmd/kutti/defaults"
	"github.com/rajch/kutti/internal/pkg/kuttilog"
	"github.com/rajch/kutti/pkg/clustermanager"
	"github.com/spf13/cobra"
)

// clusterrmCmd represents the clusterrm command
var clusterrmCmd = &cobra.Command{
	Use:           "rm CLUSTERNAME",
	Aliases:       []string{"delete", "remove"},
	Short:         "Delete a cluster.",
	Long:          `Delete a cluster.`,
	Run:           clusterrm,
	Args:          clusternameonlyargs,
	SilenceErrors: true,
}

func init() {
	clusterCmd.AddCommand(clusterrmCmd)
}

func clusterrm(cmd *cobra.Command, args []string) {
	clustername := args[0]

	kuttilog.Printf(2, "Deleting cluster %s...\n", clustername)
	err := clustermanager.DeleteCluster(clustername)
	if err != nil {
		kuttilog.Printf(0, "Error: Could not delete cluster %s: %v.\n", clustername, err)
		return
	}

	if defaults.Getdefault("cluster") == clustername {
		kuttilog.Printf(2, "Resetting default cluster.")
		defaults.Setdefault("cluster", "")
	}

	if kuttilog.V(1) {
		kuttilog.Printf(1, "Cluster '%s' deleted.\n", clustername)
	} else {
		kuttilog.Println(0, clustername)
	}

}
