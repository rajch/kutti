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
	"github.com/rajch/kutti/cmd/kutti/defaults"
	"github.com/rajch/kutti/internal/pkg/kuttilog"
	"github.com/spf13/cobra"
)

// versionrmCmd represents the versionrm command
var versionrmCmd = &cobra.Command{
	Use:                   "rm VERSION",
	Short:                 "Remove local copy of a version",
	Long:                  `Remove local copy of a version.`,
	Run:                   versionrm,
	Args:                  versiononlyargs,
	DisableFlagsInUseLine: true,
	SilenceErrors:         true,
}

func init() {
	versionCmd.AddCommand(versionrmCmd)

}

func versionrm(cmd *cobra.Command, args []string) {
	driver, err := getDriver(cmd)
	if err != nil {
		kuttilog.Printf(0, "Error: %v", err)
		return
	}

	versionstring := args[0]
	version, err := driver.GetVersion(versionstring)
	if err != nil {
		kuttilog.Printf(0, "Error:%v.\n", err)
		return
	}

	if version.Status() == "Available" {
		kuttilog.Printf(2, "Removing cached image of version %s...", versionstring)
		err = version.PurgeLocal()
		if err != nil {
			kuttilog.Printf(0, "Error: %v", err)
		}

		if defaults.Getdefault("version") == versionstring {
			kuttilog.Println(2, "Clearing default version.")
			defaults.Setdefault("version", "")
		}

		if kuttilog.V(1) {
			kuttilog.Printf(1, "Cached image of version %s removed.", versionstring)
		} else {
			kuttilog.Println(0, versionstring)
		}
	} else {
		kuttilog.Printf(0, "Error: cached image not found for version %s.", versionstring)
	}
}
