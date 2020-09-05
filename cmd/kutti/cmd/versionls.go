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
	"github.com/rajch/kutti/pkg/clustermanager"
	"github.com/spf13/cobra"
)

// versionlsCmd represents the versionls command
var versionlsCmd = &cobra.Command{
	Use:                   "ls",
	Aliases:               []string{"list"},
	Short:                 "List Kubernetes Versions",
	Long:                  `List Kubernetes Versions.`,
	Run:                   versionls,
	DisableFlagsInUseLine: true,
}

func init() {
	versionCmd.AddCommand(versionlsCmd)

	versionlsCmd.Flags().StringP("driver", "d", "", "Driver name")
}

func versionls(cmd *cobra.Command, args []string) {
	driver, err := getDriver(cmd)
	if err != nil {
		kuttilog.Printf(0, "Error: %v", err)
		return
	}

	kuttilog.Printf(2, "Versions for driver %s:\n", driver.Name())
	fmt.Println("VERSION  STATUS")
	driver.ForEachVersion(func(v *clustermanager.Version) bool {
		fmt.Printf("%7.7s  %s\n", v.K8sversion(), v.Status())
		return false
	})
}
