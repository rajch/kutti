package cmd

import (
	"github.com/rajch/kutti/cmd/kutti/defaults"
	"github.com/rajch/kutti/internal/pkg/kuttilog"
	"github.com/spf13/cobra"
)

// versioncleardefaultCmd represents the versioncleardefault command
var versioncleardefaultCmd = &cobra.Command{
	Use:     "cleardefault",
	Aliases: []string{"unselect"},
	Short:   "Clear default version",
	Long:    `Clear default version.`,
	Run:     versioncleardefault,
}

func init() {
	versionCmd.AddCommand(versioncleardefaultCmd)

}

func versioncleardefault(cmd *cobra.Command, args []string) {
	defaults.Setdefault("version", "")
	kuttilog.Println(1, "Default version cleared.")
}
