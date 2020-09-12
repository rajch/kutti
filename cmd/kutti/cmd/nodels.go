package cmd

import (
	"github.com/rajch/kutti/cmd/kutti/defaults"
	"github.com/rajch/kutti/internal/pkg/kuttilog"

	"github.com/spf13/cobra"
)

// nodelsCmd represents the nodels command
var nodelsCmd = &cobra.Command{
	Use:                   "ls",
	Aliases:               []string{"list"},
	Short:                 "List nodes",
	Long:                  `List nodes.`,
	Run:                   nodelsCommand,
	SilenceUsage:          true,
	SilenceErrors:         true,
	DisableFlagsInUseLine: true,
}

func init() {
	nodeCmd.AddCommand(nodelsCmd)

	nodelsCmd.Flags().StringP("cluster", "c", defaults.Getdefault("cluster"), "cluster name")
}

func nodelsCommand(cmd *cobra.Command, args []string) {
	cluster, err := getCluster(cmd)
	if err != nil {
		kuttilog.Printf(0, "Error: %v", err)
		return
	}
	kuttilog.Printf(0, "Nodes for cluster %s:\n", cluster.Name)
	kuttilog.Printf(0, "%-12.12s  %-10.10s  %s\n", "NAME", "TYPE", "STATUS")
	for _, node := range cluster.Nodes {
		status := node.Status()
		kuttilog.Printf(0, "%-12.12s  %-10.10s  %s\n", node.Name, node.Type, status)
	}

}
