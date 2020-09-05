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

// nodeshowCmd represents the nodeshow command
var nodeshowCmd = &cobra.Command{
	Use:           "show NODENAME",
	Aliases:       []string{"describe", "inspect", "get"},
	Short:         "Show node details",
	Long:          `Show node details.`,
	Run:           nodeshowCommand,
	Args:          nodenameonlyargs,
	SilenceErrors: true,
}

func init() {
	nodeCmd.AddCommand(nodeshowCmd)

	nodeshowCmd.Flags().StringP("cluster", "c", "", "Cluster name")
}

func nodeshowCommand(cmd *cobra.Command, args []string) {
	cluster, err := getCluster(cmd)
	if err != nil {
		kuttilog.Printf(0, "Error: %v", err)
		return
	}

	nodename := args[0]
	node, ok := cluster.Nodes[nodename]
	if !ok {
		fmt.Printf("Error: Node '%s' does not exist.\n", nodename)
		return
	}

	fmt.Printf(
		"Name: %v\nType: %v\nPorts:\n",
		node.Name,
		node.Type,
	)

	for containerport, hostport := range node.Ports {
		fmt.Printf(
			"  - HostPort: %v\n    NodePort: %v\n",
			hostport,
			containerport,
		)
	}

	fmt.Print("Status: ")
	fmt.Printf("%v\n", node.Status())
}
