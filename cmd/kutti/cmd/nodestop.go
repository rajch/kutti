package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// nodestopCmd represents the nodestop command
var nodestopCmd = &cobra.Command{
	Use:   "stop NODENAME",
	Short: "Stops a node",
	Long:  `Stops a node.`,
	Args:  nodecreateargs,
	Run:   nodestop,
}

func init() {
	nodeCmd.AddCommand(nodestopCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// nodestopCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// nodestopCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	nodestopCmd.Flags().StringP("cluster", "c", "", "Cluster name")
}

func nodestop(cmd *cobra.Command, args []string) {
	cluster, err := getCluster(cmd)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	nodename := args[0]
	node, ok := cluster.Nodes[nodename]
	if !ok {
		fmt.Printf("Error: node '%v' not found.\n", nodename)
		return
	}

	err = node.Stop()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println(node.Name)
	}
}
