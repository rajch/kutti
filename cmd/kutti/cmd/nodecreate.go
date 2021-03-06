package cmd

import (
	"github.com/rajch/kutti/cmd/kutti/defaults"
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
	Run:           nodecreateCommand,
	SilenceErrors: true,
}

func init() {
	nodeCmd.AddCommand(nodecreateCmd)

	nodecreateCmd.Flags().StringP("cluster", "c", defaults.Getdefault("cluster"), "cluster name")
	nodecreateCmd.Flags().IntP("sshport", "p", 0, "host port to forward node SSH port")
}

func nodecreateCommand(cmd *cobra.Command, args []string) {
	// Get cluster to create node in
	cluster, err := getCluster(cmd)
	if err != nil {
		kuttilog.Printf(0, "Error: %v.", err)
		return
	}

	// Check validity of node name
	nodename := args[0]
	err = cluster.ValidateNodeName(nodename)
	if err != nil {
		kuttilog.Printf(0, "Error: Could not create node '%s': %v.", nodename, err)
		return
	}

	// Check for sshport for drivers that require it
	driver, _ := clustermanager.GetDriver(cluster.DriverName)
	sshport, _ := cmd.Flags().GetInt("sshport")
	if driver.RequiresPortForwarding() && sshport == 0 {
		kuttilog.Printf(0, "Error: SSH port forwarding required for nodes in the '%s' cluster.", cluster.Name)
		return
	}

	// Check if sshport is occupied
	if sshport != 0 {
		err = cluster.CheckHostport(sshport)
		if err != nil {
			kuttilog.Printf(0, "Error: Cannot use host port %v: %v.", sshport, err)
			return
		}
	}

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
			kuttilog.Printf(0, "Warning: Could not forward SSH port: %v.", err)
			// Don't fail node creation
			kuttilog.Printf(0, "Warning: Try manually mapping the SSH port, or delete and re-create this node.")
		}
	}

	if kuttilog.V(1) {
		kuttilog.Printf(1, "Node '%s' created.", nodename)
	} else {
		kuttilog.Println(0, nodename)
	}
}
