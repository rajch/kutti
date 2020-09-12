package cmd

import (
	"github.com/rajch/kutti/cmd/kutti/defaults"
	"github.com/rajch/kutti/internal/pkg/kuttilog"
	"github.com/rajch/kutti/pkg/clustermanager"
	"github.com/spf13/cobra"
)

// versionlsCmd represents the versionls command
var versionlsCmd = &cobra.Command{
	Use:                   "ls",
	Aliases:               []string{"list"},
	Short:                 "List supported Kubernetes versions",
	Long:                  `List supported Kubernetes versions.`,
	Run:                   versionlsCommand,
	DisableFlagsInUseLine: true,
}

func init() {
	versionCmd.AddCommand(versionlsCmd)

	versionlsCmd.Flags().StringP("driver", "d", defaults.Getdefault("driver"), "driver name")
}

func versionlsCommand(cmd *cobra.Command, args []string) {
	driver, err := getDriver(cmd)
	if err != nil {
		kuttilog.Printf(0, "Error: %v.", err)
		return
	}

	kuttilog.Printf(0, "Versions for driver %s:\n", driver.Name())
	kuttilog.Println(0, "VERSION   LOCAL COPY")
	driver.ForEachVersion(func(v *clustermanager.Version) bool {
		kuttilog.Printf(
			0,
			"%8.8s  %s\n",
			defaultdecorate(v.K8sversion(), "version"),
			v.Status(),
		)
		return false
	})
}
