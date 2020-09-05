package cmd

import (
	"errors"
	"fmt"

	"github.com/rajch/kutti/cmd/kutti/defaults"
	"github.com/rajch/kutti/pkg/clustermanager"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Manage versions",
	Long:  `Manage versions.`,
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

// Common functions for version subcommands

func getDriver(cmd *cobra.Command) (*clustermanager.Driver, error) {
	var driver *clustermanager.Driver
	drivername, _ := cmd.Flags().GetString("driver")

	if drivername == "" {
		drivername = defaults.Getdefault("driver")
	}

	if drivername == "" {
		return nil,
			fmt.Errorf("no driver specified and default driver not set. Use --driver, or select a default driver using 'kutti driver setdefault'")
	}

	driver, _ = clustermanager.GetDriver(drivername)
	if driver == nil {
		return nil,
			fmt.Errorf("driver '%v' not found", drivername)
	}

	return driver, nil
}

func versiononlyargs(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return errors.New("VERSION is required")
	}

	if len(args) > 1 {
		cmd.SilenceUsage = true
		return errors.New("only VERSION is required")
	}

	return nil
}
