package cmd

import (
	"fmt"

	"github.com/rajch/kutti/pkg/clustermanager"
	"github.com/rajch/kutti/pkg/core"

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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// driverlsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// driverlsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func driverls(cmd *cobra.Command, args []string) {
	fmt.Println("Name\tDescription\tStatus")
	clustermanager.Load()
	clustermanager.ForEachDriver(func(driver core.VMDriver) bool {
		fmt.Printf("%v\t%v\t%v\n", driver.Name(), driver.Description(), driver.Status())
		return false
	})
}
