package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// nodestartCmd represents the nodestart command
var nodestartCmd = &cobra.Command{
	Use:   "start NODENAME",
	Short: "Starts a node",
	Long:  `Starts a node.`,
	Args:  nodenameonlyargs,
	Run:   nodestart,
}

func init() {
	nodeCmd.AddCommand(nodestartCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// nodestartCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// nodestartCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	nodestartCmd.Flags().StringP("cluster", "c", "", "Cluster name")
}

func nodestart(cmd *cobra.Command, args []string) {
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

	err = node.Start()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println(node.Name)
	}

}
