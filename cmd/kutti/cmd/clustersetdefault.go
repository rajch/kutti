package cmd

import (
	"github.com/rajch/kutti/cmd/kutti/defaults"
	"github.com/rajch/kutti/internal/pkg/kuttilog"
	"github.com/spf13/cobra"
)

// clustersetdefaultCmd represents the clustersetdefault command
var clustersetdefaultCmd = &cobra.Command{
	Use:           "setdefault CLUSTERNAME",
	Short:         "Sets the default cluster.",
	Long:          `Sets the default cluster.`,
	Args:          clusternameonlyargs,
	Run:           clustersetdefault,
	SilenceErrors: true,
}

func init() {
	clusterCmd.AddCommand(clustersetdefaultCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// clustersetdefaultCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// clustersetdefaultCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func clustersetdefault(cmd *cobra.Command, args []string) {
	clustername := args[0]
	defaults.Setdefault("cluster", clustername)
	kuttilog.Print(2, "Default cluster set to: ")
	kuttilog.Println(0, clustername)

}
