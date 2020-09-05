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
	Short:                 "Update supported Kubernetes versions",
	Long:                  `Update supported Kubernetes versions.`,
	Run:                   driverupdateversions,
	Args:                  drivernameonlyargs,
	DisableFlagsInUseLine: true,
	SilenceErrors:         true,
	SilenceUsage:          true,
}

func init() {
	driverCmd.AddCommand(driverupdateversionsCmd)

	driverupdateversionsCmd.Flags().String("updateurl", "", "location of supported version data")
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
