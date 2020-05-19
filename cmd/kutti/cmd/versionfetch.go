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

	"github.com/spf13/cobra"
)

// versionfetchCmd represents the versionfetch command
var versionfetchCmd = &cobra.Command{
	Use:   "fetch VERSION",
	Short: "Fetch an image for a specified version of Kubernetes",
	Long:  `Fetch an image for a specified version of Kubernetes.`,
	Run:   versionfetch,
	Args:  versiononlyargs,
}

func init() {
	versionCmd.AddCommand(versionfetchCmd)

	// versionfetchCmd.Flags().StringP("version", "v", defaults.Getdefault("version"), "K8s version to fetch image for")
	versionfetchCmd.Flags().StringP("fromfile", "f", "", "Image file path")
}

func versionfetch(cmd *cobra.Command, args []string) {
	driver, err := getDriver(cmd)
	if err != nil {
		kuttilog.Printf(0, "Error: %v", err)
		return
	}

	versionstring := args[0]

	version, err := driver.GetVersion(versionstring)
	if err != nil {
		kuttilog.Printf(0, "Error:%v.", err)
		return
	}

	filename, err := cmd.Flags().GetString("fromfile")
	if err != nil || filename == "" {
		kuttilog.Printf(1, "Fetching version %s...", versionstring)
		err = version.Fetch()
		if err != nil {
			kuttilog.Printf(0, "Error: Could not download version %s: %v.", versionstring, err)
			return
		}

		if kuttilog.V(1) {
			kuttilog.Printf(1, "Fetched version %s.", versionstring)
		} else {
			kuttilog.Println(0, versionstring)
		}
		return
	}

	kuttilog.Printf(2, "Importing local image for version %v...", versionstring)
	err = version.FromFile(filename)
	if err != nil {
		kuttilog.Printf(0, "Error: Could not import local image: %v.", err)
		return
	}

	if kuttilog.V(1) {
		kuttilog.Printf(1, "Image for version %v imported.", version.K8sversion())
	} else {
		kuttilog.Println(0, version.K8sversion())
	}

}
