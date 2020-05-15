package cmd

import (
	"github.com/rajch/kutti/cmd/kutti/defaults"
	"github.com/rajch/kutti/internal/pkg/kuttilog"
	"github.com/spf13/cobra"
)

// clustercleardefaultCmd represents the clustercleardefault command
var clustercleardefaultCmd = &cobra.Command{
	Use:           "cleardefault",
	Short:         "Clears the default cluster.",
	Long:          `Clears the default cluster.`,
	Args:          cobra.NoArgs,
	Run:           clustercleardefault,
	SilenceErrors: true,
}

func init() {
	clusterCmd.AddCommand(clustercleardefaultCmd)
}

func clustercleardefault(cmd *cobra.Command, args []string) {
	defaults.Setdefault("cluster", "")
	kuttilog.Println(1, "Default cluster cleared.")
}
