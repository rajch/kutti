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

// versionsetdefaultCmd represents the versionsetdefault command
var versionsetdefaultCmd = &cobra.Command{
	Use:           "setdefault VERSION",
	Short:         "Set default version",
	Long:          `Set default version.`,
	Run:           versionsetdefault,
	Args:          versiononlyargs,
	SilenceErrors: true,
}

func init() {
	versionCmd.AddCommand(versionsetdefaultCmd)

	versionsetdefaultCmd.Flags().StringP("driver", "d", defaults.Getdefault("driver"), "Driver name")
}

func versionsetdefault(cmd *cobra.Command, args []string) {
	driver, err := getDriver(cmd)
	if err != nil {
		kuttilog.Printf(0, "Error: %v", err)
		return
	}

	versionstring := args[0]
	_, err = driver.GetVersion(versionstring)
	if err != nil {
		kuttilog.Printf(0, "Error: Could not set default version: %v.\n", err)
		return
	}

	defaults.Setdefault("version", versionstring)

	if kuttilog.V(1) {
		kuttilog.Printf(1, "Default version set to %s.", versionstring)
	} else {
		kuttilog.Println(0, versionstring)
	}
}
