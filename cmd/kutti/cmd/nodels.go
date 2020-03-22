package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// nodelsCmd represents the nodels command
var nodelsCmd = &cobra.Command{
	Use:           "ls",
	Aliases:       []string{"list"},
	Short:         "List nodes",
	Long:          `List nodes.`,
	Run:           nodels,
	SilenceUsage:  true,
	SilenceErrors: true,
}

func init() {
	nodeCmd.AddCommand(nodelsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// nodelsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// nodelsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	nodelsCmd.Flags().StringP("cluster", "c", "", "Cluster name")
}

func nodels(cmd *cobra.Command, args []string) {
	cluster, err := getCluster(cmd)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Printf("%-12.12s  %-10.10s  %s\n", "NAME", "TYPE", "STATUS")
	for _, node := range cluster.Nodes {
		status := node.Status()
		fmt.Printf("%-12.12s  %-10.10s  %s\n", node.Name, node.Type, status)
	}

}
