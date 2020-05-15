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
	"errors"
	"fmt"

	"github.com/rajch/kutti/cmd/kutti/defaults"
	"github.com/rajch/kutti/pkg/clustermanager"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Manage versions",
	Long:  `Manage versions.`,
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

// Common functions for version subcommands

func getDriver(cmd *cobra.Command) (*clustermanager.Driver, error) {
	var driver *clustermanager.Driver
	drivername, _ := cmd.Flags().GetString("driver")

	if drivername == "" {
		drivername = defaults.Getdefault("driver")
	}

	if drivername == "" {
		return nil,
			fmt.Errorf("no driver specified and default driver not set. Use --driver, or select a default driver using 'kutti driver setdefault'")
	}

	driver, _ = clustermanager.GetDriver(drivername)
	if driver == nil {
		return nil,
			fmt.Errorf("driver '%v' not found", drivername)
	}

	return driver, nil
}

func versiononlyargs(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return errors.New("VERSION is required")
	}

	if len(args) > 1 {
		cmd.SilenceUsage = true
		return errors.New("only VERSION is required")
	}

	return nil
}
