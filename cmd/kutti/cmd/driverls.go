package cmd

import (
	"fmt"

	"github.com/rajch/kutti/pkg/clustermanager"

	"github.com/spf13/cobra"
)

// driverlsCmd represents the driverls command
var driverlsCmd = &cobra.Command{
	Use:                   "ls",
	Short:                 "List drivers",
	Long:                  `List drivers.`,
	Run:                   driverls,
	DisableFlagsInUseLine: true,
}

func init() {
	driverCmd.AddCommand(driverlsCmd)
}

func driverls(cmd *cobra.Command, args []string) {
	fmt.Printf("%-12.12s  %-40.40s  %-12.12s\n", "NAME", "DESCRIPTION", "STATUS")
	clustermanager.ForEachDriver(func(driver *clustermanager.Driver) bool {
		fmt.Printf("%-12.12s  %-40.40s  %-12.12s\n", driver.Name(), driver.Description(), driver.Status())
		return false
	})
}
