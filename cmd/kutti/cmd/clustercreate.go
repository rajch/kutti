package cmd

import (
	"errors"
	"fmt"

	"github.com/rajch/kutti/pkg/clustermanager"
	_ "github.com/rajch/kutti/pkg/vboxdriver" // Virtual Box driver
	"github.com/spf13/cobra"
)

// clustercreateCmd represents the clustercreate command
var clustercreateCmd = &cobra.Command{
	Use:           "create CLUSTERNAME",
	Aliases:       []string{"add"},
	Short:         "Create a new cluster",
	Long:          `Create a new cluster.`,
	Args:          clustercreateargs,
	Run:           clustercreate,
	SilenceErrors: true,
}

func init() {
	clusterCmd.AddCommand(clustercreateCmd)

	clustercreateCmd.Flags().StringP(
		"version",
		"v",
		"1.17",
		"Kubernetes version for the cluster",
	)

	clustercreateCmd.Flags().StringP(
		"driver",
		"d",
		"vbox",
		"Cluster management driver",
	)

	clustercreateCmd.Flags().BoolP(
		"unmanaged",
		"u",
		false,
		"Create an unmanaged cluster with no nodes",
	)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// clustercreateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// clustercreateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func clustercreateargs(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return errors.New("CLUSTERNAME is required")
	}

	if len(args) > 1 {
		cmd.SilenceUsage = true
		return errors.New("only CLUSTERNAME is required")
	}

	return nil
}

func clustercreate(cmd *cobra.Command, args []string) {
	clustername := args[0]
	driver, _ := cmd.Flags().GetString("driver")
	// TODO: Validate driver
	k8sversion, _ := cmd.Flags().GetString("version")
	// TODO: Validate version

	unmanaged, _ := cmd.Flags().GetBool("unmanaged")
	var err error

	if unmanaged {
		err = clustermanager.NewEmptyCluster(
			clustername,
			k8sversion,
			driver,
		)
	} else {
		err = errors.New("managed cluster creation not yet implemented")
	}

	if err != nil {
		fmt.Printf("Could not create cluster %s: %v.\n", clustername, err)
	}
}
