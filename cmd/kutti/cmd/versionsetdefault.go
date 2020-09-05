package cmd

import (
	"github.com/rajch/kutti/cmd/kutti/defaults"
	"github.com/rajch/kutti/internal/pkg/kuttilog"

	"github.com/spf13/cobra"
)

// versionsetdefaultCmd represents the versionsetdefault command
var versionsetdefaultCmd = &cobra.Command{
	Use:           "setdefault VERSION",
	Aliases:       []string{"select"},
	Short:         "Set default version",
	Long:          `Set default version.`,
	Run:           versionsetdefault,
	Args:          versiononlyargs,
	SilenceErrors: true,
}

func init() {
	versionCmd.AddCommand(versionsetdefaultCmd)

	versionsetdefaultCmd.Flags().StringP("driver", "d", defaults.Getdefault("driver"), "driver name")
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
