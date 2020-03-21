package cmd

import (
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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// driverCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// driverCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
