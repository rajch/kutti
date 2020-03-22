package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

// nodecreateCmd represents the nodecreate command
var nodecreateCmd = &cobra.Command{
	Use:     "create NODENAME",
	Aliases: []string{"add"},
	Short:   "Create a node",
	Long:    `Create a node.`,
	Args:    nodecreateargs,
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

func nodecreateargs(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return errors.New("NODENAME is required")
	}

	if len(args) > 1 {
		cmd.SilenceUsage = true
		return errors.New("only NODENAME is required")
	}

	return nil
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
