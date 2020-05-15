package cmd

import (
	"errors"

	"github.com/spf13/cobra"
)

// driverCmd represents the driver command
var driverCmd = &cobra.Command{
	Use:   "driver",
	Short: "View drivers",
	Long:  `View drivers`,
}

func init() {
	rootCmd.AddCommand(driverCmd)
}

func drivernameonlyargs(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return errors.New("DRIVERNAME is required")
	}

	if len(args) > 1 {
		cmd.SilenceUsage = true
		return errors.New("only DRIVERNAME is required")
	}

	return nil
}
