package cmd

import (
	"fmt"

	"github.com/rajch/kutti/internal/pkg/kuttilog"

	"github.com/spf13/cobra"
)

// nodelsCmd represents the nodels command
var nodelsCmd = &cobra.Command{
	Use:                   "ls",
	Aliases:               []string{"list"},
	Short:                 "List nodes",
	Long:                  `List nodes.`,
	Run:                   nodels,
	SilenceUsage:          true,
	SilenceErrors:         true,
	DisableFlagsInUseLine: true,
}

func init() {
	nodeCmd.AddCommand(nodelsCmd)

	nodelsCmd.Flags().StringP("cluster", "c", "", "Cluster name")
}

func nodels(cmd *cobra.Command, args []string) {
	cluster, err := getCluster(cmd)
	if err != nil {
		kuttilog.Printf(0, "Error: %v", err)
		return
	}
	kuttilog.Printf(2, "Nodes for cluster %s:\n", cluster.Name)
	fmt.Printf("%-12.12s  %-10.10s  %s\n", "NAME", "TYPE", "STATUS")
	for _, node := range cluster.Nodes {
		status := node.Status()
		fmt.Printf("%-12.12s  %-10.10s  %s\n", node.Name, node.Type, status)
	}

}
