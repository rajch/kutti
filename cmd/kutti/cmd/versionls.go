package cmd

import (
	"fmt"

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
	Run:                   versionls,
	DisableFlagsInUseLine: true,
}

func init() {
	versionCmd.AddCommand(versionlsCmd)

	versionlsCmd.Flags().StringP("driver", "d", "", "driver name")
}

func versionls(cmd *cobra.Command, args []string) {
	driver, err := getDriver(cmd)
	if err != nil {
		kuttilog.Printf(0, "Error: %v", err)
		return
	}

	kuttilog.Printf(2, "Versions for driver %s:\n", driver.Name())
	fmt.Println("VERSION  LOCAL COPY")
	driver.ForEachVersion(func(v *clustermanager.Version) bool {
		fmt.Printf("%7.7s  %s\n", v.K8sversion(), v.Status())
		return false
	})
}
