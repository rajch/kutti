package cmd

import (
	"errors"

	"github.com/rajch/kutti/cmd/kutti/defaults"

	"github.com/rajch/kutti/internal/pkg/kuttilog"

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
	Args:          clusternameonlyargs,
	Run:           clustercreateCommand,
	SilenceErrors: true,
}

func init() {
	clusterCmd.AddCommand(clustercreateCmd)

	clustercreateCmd.Flags().StringP(
		"driver",
		"d",
		defaults.Getdefault("driver"),
		"driver name",
	)

	clustercreateCmd.Flags().StringP(
		"version",
		"v",
		defaults.Getdefault("version"),
		"Kubernetes version for the cluster",
	)

	clustercreateCmd.Flags().BoolP(
		"unmanaged",
		"u",
		false,
		"create an unmanaged cluster with no nodes",
	)

	clustercreateCmd.Flags().BoolP(
		"select",
		"s",
		false,
		"set the new cluster as default",
	)
}

func clusternameonlyargs(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return errors.New("CLUSTERNAME is required")
	}

	if len(args) > 1 {
		cmd.SilenceUsage = true
		return errors.New("only CLUSTERNAME is required")
	}

	return nil
}

func clustercreateCommand(cmd *cobra.Command, args []string) {
	clustername := args[0]

	driver, _ := cmd.Flags().GetString("driver")
	_, ok := clustermanager.GetDriver(driver)
	if !ok {
		kuttilog.Printf(
			0,
			"Error: Driver '%s' not found. Cannot create cluster.",
			driver,
		)
	}

	k8sversion, _ := cmd.Flags().GetString("version")
	if k8sversion == "" {
		kuttilog.Println(
			0,
			"Error: Kubernetes version not specified, and no default set. Please use --version, or select a default version.",
		)
		return
	}

	unmanaged, _ := cmd.Flags().GetBool("unmanaged")
	setdefault, _ := cmd.Flags().GetBool("select")

	var err error

	err = clustermanager.ValidateClusterName(clustername)
	if err == nil {
		if unmanaged {
			kuttilog.Printf(2, "Creating cluster '%s'...\n", clustername)
			err = clustermanager.NewEmptyCluster(
				clustername,
				k8sversion,
				driver,
			)
		} else {
			err = errors.New("managed cluster creation not yet implemented")
		}
	}

	if err != nil {
		kuttilog.Printf(0, "Error: Could not create cluster %s: %v.\n", clustername, err)
		return
	}

	if kuttilog.V(1) {
		kuttilog.Printf(1, "Cluster '%s' created.\n", clustername)
	} else {
		kuttilog.Print(0, clustername)
	}

	if setdefault {
		defaults.Setdefault("cluster", clustername)
		kuttilog.Printf(1, "Default cluster set to '%s'.", clustername)
	}
}
