package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// nodermCmd represents the noderm command
var nodermCmd = &cobra.Command{
	Use:     "rm NODENAME",
	Aliases: []string{"delete"},
	Short:   "Delete a node",
	Long:    `Delete a node.`,
	Args:    nodenameonlyargs,
	Run:     noderm,
}

func init() {
	nodeCmd.AddCommand(nodermCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// nodermCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// nodermCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	nodermCmd.Flags().StringP("cluster", "c", "", "Cluster name")
	nodermCmd.Flags().BoolP("force", "f", false, "Forcibly delete running nodes.")
}

func noderm(cmd *cobra.Command, args []string) {
	cluster, err := getCluster(cmd)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	nodename := args[0]
	forceflag, _ := cmd.Flags().GetBool("force")

	err = cluster.DeleteNode(nodename, forceflag)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println(nodename)
	}
}
