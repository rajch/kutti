package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// nodecreateCmd represents the nodecreate command
var nodecreateCmd = &cobra.Command{
	Use:     "create NODENAME",
	Aliases: []string{"add"},
	Short:   "Create a node",
	Long:    `Create a node.`,
	Args:    nodenameonlyargs,
	Run:     nodecreate,
}

func init() {
	nodeCmd.AddCommand(nodecreateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// nodecreateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// nodecreateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	nodecreateCmd.Flags().StringP("cluster", "c", "", "Cluster name")

}

func nodecreate(cmd *cobra.Command, args []string) {
	cluster, err := getCluster(cmd)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	nodename := args[0]
	_, err = cluster.AddUninitializedNode(nodename)
	if err != nil {
		fmt.Printf("Could not create node %v: %v.\n", nodename, err)
		return
	}

	fmt.Println(nodename)
}
