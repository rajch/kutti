package cmd

import (
	"github.com/rajch/kutti/internal/pkg/kuttilog"
	"github.com/rajch/kutti/pkg/clustermanager"
	"github.com/spf13/cobra"
)

// clusterrmCmd represents the clusterrm command
var clusterrmCmd = &cobra.Command{
	Use:           "rm CLUSTERNAME",
	Aliases:       []string{"delete"},
	Short:         "Delete a cluster.",
	Long:          `Delete a cluster.`,
	Run:           clusterrm,
	Args:          clusternameonlyargs,
	SilenceErrors: true,
}

func init() {
	clusterCmd.AddCommand(clusterrmCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// clusterrmCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// clusterrmCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func clusterrm(cmd *cobra.Command, args []string) {
	clustername := args[0]

	clustermanager.Load()
	kuttilog.Printf(2, "Deleting cluster %s...\n", clustername)
	err := clustermanager.DeleteCluster(clustername)
	if err != nil {
		kuttilog.Printf(0, "Error: Could not delete cluster %s: %v.\n", clustername, err)
		return
	}

	if kuttilog.V(1) {
		kuttilog.Printf(1, "Cluster '%s' deleted.\n", clustername)
	} else {
		kuttilog.Println(0, clustername)
	}

}
