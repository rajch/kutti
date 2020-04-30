package cmd

import (
	"fmt"

	"github.com/rajch/kutti/pkg/clustermanager"
	"github.com/rajch/kutti/pkg/core"
	"github.com/spf13/cobra"
)

// driverversionsCmd represents the driverversions command
var driverversionsCmd = &cobra.Command{
	Use:   "versions DRIVERNAME",
	Short: "Show available kubernetes versions.",
	Long:  `Show available kubernetes versions.`,
	Args:  drivernameonlyargs,
	Run:   driverversions,
}

func init() {
	driverCmd.AddCommand(driverversionsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// driverversionsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// driverversionsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func driverversions(cmd *cobra.Command, args []string) {
	drivername := args[0]

	fmt.Printf("Available versions for driver:%s\n", drivername)
	clustermanager.Load()
	err := clustermanager.ForEachImage(drivername, func(image core.VMImage) bool {
		fmt.Printf("%s : %s\n", image.K8sVersion(), image.Status())
		return false
	})

	if err != nil {
		fmt.Printf("Error:%v\n", err)
	}
}
