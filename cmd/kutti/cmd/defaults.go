package cmd

import (
	"fmt"

	"github.com/rajch/kutti/cmd/kutti/defaults"
	"github.com/spf13/cobra"
)

// defaultsCmd represents the defaults command
var defaultsCmd = &cobra.Command{
	Use:     "defaults",
	Aliases: []string{"default", "showdefaults"},
	Short:   "View defaults",
	Long:    `View defaults.`,
	Run:     defaultsCommand,
}

func init() {
	rootCmd.AddCommand(defaultsCmd)
}

func defaultsCommand(cmd *cobra.Command, args []string) {
	fmt.Printf(
		"Driver: %v\nVersion: %v\nCluster: %v\n",
		defaults.Getdefault("driver"),
		defaults.Getdefault("version"),
		defaults.Getdefault("cluster"),
	)
}
