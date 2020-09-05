package cmd

import (
	"github.com/rajch/kutti/internal/pkg/kuttilog"

	"github.com/spf13/cobra"
)

// versionfetchCmd represents the versionfetch command
var versionfetchCmd = &cobra.Command{
	Use:   "fetch VERSION",
	Short: "Fetch local copy of a specified Kubernetes version image",
	Long:  `Fetch local copy of a specified Kubernetes version image.`,
	Run:   versionfetch,
	Args:  versiononlyargs,
}

func init() {
	versionCmd.AddCommand(versionfetchCmd)

	versionfetchCmd.Flags().StringP("fromfile", "f", "", "fetch image from specified file path")
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
		kuttilog.Printf(1, "Downloading local copy image for Kubernetes version %s...", versionstring)
		err = version.Fetch()
		if err != nil {
			kuttilog.Printf(0, "Error: Could not download local copy image for Kubernetes version %s: %v.", versionstring, err)
			return
		}

		if kuttilog.V(1) {
			kuttilog.Printf(1, "Downloaded local copy image for Kubernetes version %s.", versionstring)
		} else {
			kuttilog.Println(0, versionstring)
		}
		return
	}

	kuttilog.Printf(2, "Importing local copy image for version %v...", versionstring)
	err = version.FromFile(filename)
	if err != nil {
		kuttilog.Printf(0, "Error: Could not import local copy image: %v.", err)
		return
	}

	if kuttilog.V(1) {
		kuttilog.Printf(1, "Local copy image for version %v imported.", version.K8sversion())
	} else {
		kuttilog.Println(0, version.K8sversion())
	}

}
