package cmd

import (
	"fmt"

	"github.com/rajch/kutti/pkg/clustermanager"
	"github.com/spf13/cobra"
)

// driverfetchversionCmd represents the driversfetchversion command
var driverfetchversionCmd = &cobra.Command{
	Use:     "fetchversion DRIVERNAME",
	Aliases: []string{"getversion"},
	Short:   "Fetch the image for a version from the internet, or a local file.",
	Long:    `Fetch the image for a version from the internet, or a local file.`,
	Args:    drivernameonlyargs,
	Run:     driverfetchversion,
}

func init() {
	driverCmd.AddCommand(driverfetchversionCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// driversfetchversionCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// driversfetchversionCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	driverfetchversionCmd.Flags().StringP("version", "v", "1.17", "K8s version to fetch image for")
	driverfetchversionCmd.Flags().StringP("fromfile", "f", "", "Image file path")
}

func driverfetchversion(cmd *cobra.Command, args []string) {

	drivername := args[0]

	clustermanager.Load()
	driver, ok := clustermanager.GetDriver(drivername)
	if !ok {
		fmt.Printf("Error: Driver '%v' not found.\n", drivername)
		return
	}

	version, _ := cmd.Flags().GetString("version")

	image, err := driver.GetImage(version)
	if err != nil {
		fmt.Printf("Error:%v.\n", err)
		return
	}

	filename, err := cmd.Flags().GetString("fromfile")
	if err != nil || filename == "" {
		fmt.Println("Remote fetching not yet implemented.")
		return
	}

	err = image.FromFile(filename)
	if err != nil {
		fmt.Printf("Error:%v.\n", err)
		return
	}

	fmt.Printf("Image for version %v imported.\n", image.K8sVersion())
}
