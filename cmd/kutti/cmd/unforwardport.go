/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"fmt"

	"github.com/rajch/kutti/internal/pkg/kuttilog"
	"github.com/spf13/cobra"
)

// unforwardportCmd represents the unforwardport command
var unforwardportCmd = &cobra.Command{
	Use:           "unforwardport NODENAME",
	Aliases:       []string{"unpublish", "unforward", "unmap"},
	Short:         "Unforward a node port",
	Long:          `Unforward a node port.`,
	Run:           unforwardportCommand,
	Args:          nodenameonlyargs,
	SilenceErrors: true,
}

func init() {
	nodeCmd.AddCommand(unforwardportCmd)

	unforwardportCmd.Flags().IntP("nodeport", "n", 0, "Node port to unmap")
}

func unforwardportCommand(cmd *cobra.Command, args []string) {
	cluster, err := getCluster(cmd)
	if err != nil {
		kuttilog.Printf(0, "Error: %v", err)
		return
	}

	nodename := args[0]
	node, ok := cluster.Nodes[nodename]
	if !ok {
		kuttilog.Printf(0, "Error: node '%v' not found.\n", nodename)
		return
	}

	nodeport, _ := cmd.Flags().GetInt("nodeport")
	if nodeport == 0 {
		fmt.Println("Error: Please provide a valid nodeport.")
		return
	}

	err = node.UnforwardPort(nodeport)
	if err != nil {
		fmt.Printf("Error: Cannot unforward node port %v: %v.\n", nodeport, err)
		return
	}

	fmt.Printf("Node port %v unforwarded.\n", nodeport)
}
