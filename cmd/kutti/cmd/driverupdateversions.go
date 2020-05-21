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
	"github.com/rajch/kutti/internal/pkg/kuttilog"
	"github.com/rajch/kutti/pkg/clustermanager"
	"github.com/rajch/kutti/pkg/vboxdriver"
	"github.com/spf13/cobra"
)

// driverupdateversionsCmd represents the driverupdateversions command
var driverupdateversionsCmd = &cobra.Command{
	Use:                   "updateversions DRIVERNAME",
	Aliases:               []string{"updateversion", "update"},
	Short:                 "Update available versions",
	Long:                  `Update available versions.`,
	Run:                   driverupdateversions,
	Args:                  drivernameonlyargs,
	DisableFlagsInUseLine: true,
	SilenceErrors:         true,
	SilenceUsage:          true,
}

func init() {
	driverCmd.AddCommand(driverupdateversionsCmd)

	driverupdateversionsCmd.Flags().String("updateurl", "", "Update URL")
}

func driverupdateversions(cmd *cobra.Command, args []string) {
	drivername := args[0]
	driver, ok := clustermanager.GetDriver(drivername)
	if !ok {
		kuttilog.Printf(0, "Error: Driver '%s' not found.", drivername)
		return
	}

	updateurl, _ := cmd.Flags().GetString("updateurl")
	if updateurl != "" {
		vboxdriver.ImagesSourceURL = updateurl
	}

	kuttilog.Println(1, "Updating driver versions...")
	err := driver.UpdateVersions()
	if err != nil {
		kuttilog.Printf(0, "Error: %v.", err)
		return
	}

	kuttilog.Println(1, "Driver versions updated.")
}
