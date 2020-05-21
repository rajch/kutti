package cmd

import (
	"github.com/rajch/kutti/internal/pkg/kuttilog"
	"github.com/rajch/kutti/pkg/clustermanager"

	"github.com/spf13/cobra"
)

// nodecreateCmd represents the nodecreate command
var nodecreateCmd = &cobra.Command{
	Use:           "create NODENAME",
	Aliases:       []string{"add"},
	Short:         "Create a node",
	Long:          `Create a node.`,
	Args:          nodenameonlyargs,
	Run:           nodecreate,
	SilenceErrors: true,
}

func init() {
	nodeCmd.AddCommand(nodecreateCmd)

	nodecreateCmd.Flags().StringP("cluster", "c", "", "Cluster name")
	nodecreateCmd.Flags().Int("sshport", 0, "Host port to forward SSH")
}

func nodecreate(cmd *cobra.Command, args []string) {
	cluster, err := getCluster(cmd)
	if err != nil {
		kuttilog.Printf(0, "Error: %v", err)
		return
	}

	// Check for sshport for drivers that require it
	driver, _ := clustermanager.GetDriver(cluster.DriverName)
	sshport, _ := cmd.Flags().GetInt("sshport")
	if driver.RequiresPortForwarding() && sshport == 0 {
		kuttilog.Printf(0, "Error: SSH forward port required for nodes in the '%s' cluster.", cluster.Name)
		return
	}

	nodename := args[0]
	kuttilog.Printf(1, "Creating node '%s' on cluster %s...", nodename, cluster.Name)
	newnode, err := cluster.NewUninitializedNode(nodename)
	if err != nil {
		kuttilog.Printf(0, "Error: Could not create node %v: %v.", nodename, err)
		return
	}

	// Forward SSH port
	// Belt and suspenders if condition
	if driver.RequiresPortForwarding() && sshport != 0 {
		err = newnode.ForwardSSHPort(sshport)
		if err != nil {
			kuttilog.Printf(0, "Error: Could not forward SSH port: %v", err)
			// Don't fail node creation
		}
	}

	if kuttilog.V(1) {
		kuttilog.Printf(1, "Node '%s' created.", nodename)
	} else {
		kuttilog.Println(0, nodename)
	}
}
