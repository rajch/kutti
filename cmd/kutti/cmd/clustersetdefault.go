package cmd

import (
	"github.com/rajch/kutti/cmd/kutti/defaults"
	"github.com/rajch/kutti/internal/pkg/kuttilog"
	"github.com/rajch/kutti/pkg/clustermanager"
	"github.com/spf13/cobra"
)

// clustersetdefaultCmd represents the clustersetdefault command
var clustersetdefaultCmd = &cobra.Command{
	Use:           "setdefault CLUSTERNAME",
	Aliases:       []string{"select"},
	Short:         "Sets the default cluster.",
	Long:          `Sets the default cluster.`,
	Args:          clusternameonlyargs,
	Run:           clustersetdefault,
	SilenceErrors: true,
}

func init() {
	clusterCmd.AddCommand(clustersetdefaultCmd)
}

func clustersetdefault(cmd *cobra.Command, args []string) {
	clustername := args[0]
	_, ok := clustermanager.GetCluster(clustername)
	if !ok {
		kuttilog.Printf(0, "Error: Cluster '%s' not found.", clustername)
		return
	}

	defaults.Setdefault("cluster", clustername)
	if kuttilog.V(1) {
		kuttilog.Printf(1, "Default cluster set to '%s'.", clustername)
	} else {
		kuttilog.Println(0, clustername)
	}
}
