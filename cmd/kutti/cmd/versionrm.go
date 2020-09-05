package cmd

import (
	"github.com/rajch/kutti/cmd/kutti/defaults"
	"github.com/rajch/kutti/internal/pkg/kuttilog"
	"github.com/spf13/cobra"
)

// versionrmCmd represents the versionrm command
var versionrmCmd = &cobra.Command{
	Use:                   "rm VERSION",
	Aliases:               []string{"delete", "remove"},
	Short:                 "Remove local copy of a Kuberenetes version image",
	Long:                  `Remove local copy of a Kuberenetes version image.`,
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
		kuttilog.Printf(2, "Removing local copy image of version %s...", versionstring)
		err = version.PurgeLocal()
		if err != nil {
			kuttilog.Printf(0, "Error: %v", err)
		}

		if defaults.Getdefault("version") == versionstring {
			kuttilog.Println(2, "Clearing default version.")
			defaults.Setdefault("version", "")
		}

		if kuttilog.V(1) {
			kuttilog.Printf(1, "Local copy image version %s removed.", versionstring)
		} else {
			kuttilog.Println(0, versionstring)
		}
	} else {
		kuttilog.Printf(0, "Error: local copy image not found for version %s.", versionstring)
	}
}
