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

// forwardportCmd represents the forwardport command
var forwardportCmd = &cobra.Command{
	Use:           "forwardport NODENAME",
	Aliases:       []string{"publish", "forward", "map"},
	Short:         "Forwards a node port to a host port",
	Long:          `Forwards a node port to a host port.`,
	Run:           forwardportCommand,
	Args:          nodenameonlyargs,
	SilenceErrors: true,
}

func init() {
	nodeCmd.AddCommand(forwardportCmd)

	forwardportCmd.Flags().IntP("hostport", "p", 0, "Port on the host")
	forwardportCmd.Flags().IntP("nodeport", "n", 0, "Port on the node")
}

func forwardportCommand(cmd *cobra.Command, args []string) {
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

	hostport, _ := cmd.Flags().GetInt("hostport")
	if hostport == 0 {
		fmt.Println("Error: Please provide a valid hostport.")
		return
	}

	err = cluster.CheckHostport(hostport)
	if err != nil {
		fmt.Printf("Error: Cannot forward to host port %v: %v.\n", hostport, err)
		return
	}

	err = node.ForwardPort(hostport, nodeport)
	if err != nil {
		fmt.Printf(
			"Error: Could not forward node port %v to host port %v: %v\n",
			nodeport,
			hostport,
			err,
		)
		return
	}

	kuttilog.Printf(
		2,
		"Forwarded node port %v to host port %v.\n",
		nodeport,
		hostport,
	)
}
