package cmd

import (
	"fmt"

	"github.com/rajch/kutti/pkg/clustermanager"

	"github.com/spf13/cobra"
)

// driverlsCmd represents the driverls command
var driverlsCmd = &cobra.Command{
	Use:                   "ls",
	Aliases:               []string{"list"},
	Short:                 "List drivers",
	Long:                  `List drivers.`,
	Run:                   driverlsCommand,
	DisableFlagsInUseLine: true,
}

func init() {
	driverCmd.AddCommand(driverlsCmd)
}

func driverlsCommand(cmd *cobra.Command, args []string) {
	fmt.Printf("%-12.12s  %-40.40s  %s\n", "NAME", "DESCRIPTION", "STATUS")
	clustermanager.ForEachDriver(func(driver *clustermanager.Driver) bool {
		fmt.Printf(
			"%-12.12s  %-40.40s  %s\n",
			defaultdecorate(driver.Name(), "driver"),
			driver.Description(),
			driver.Status(),
		)
		return false
	})
}
