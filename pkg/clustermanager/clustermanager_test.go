package clustermanager_test

// package core_test

import (
	"os"
	"path"
	"testing"

	"github.com/rajch/kutti/pkg/clustermanager"
	"github.com/rajch/kutti/pkg/core"
	_ "github.com/rajch/kutti/pkg/vboxdriver"
)

func TestInit(t *testing.T) {
	configdir, _ := core.ConfigDir()
	configfile := path.Join(configdir, "clusters.json")
	_, err := os.Stat(configfile)
	if err != nil {
		t.Logf("Config file not found after init. This is a problem:%v\n", err)
		t.FailNow()
	}

	//t.Logf("%+v\n", clustermanager)
}

func TestDrivers(t *testing.T) {
	var result bool
	resultp := &result

	clustermanager.ForEachDriver(func(d core.VMDriver) bool {
		if d.Name() == "vbox" {
			*resultp = true
			return true
		}
		return false
	})

	if !result {
		t.Logf("Driver test failed.")
		t.FailNow()
	}
}

func TestNewCluster(t *testing.T) {
	err := clustermanager.NewEmptyCluster("testclust1", "1.17", "vbox")
	if err != nil {
		t.Logf("NewEmptyCluster failed with error:%v", err)
		t.FailNow()
	}

	cluster, _ := clustermanager.GetCluster("testclust1")
	t.Log(*cluster)

	t.Log("As of now, you will have to remove the cluster artifacts yourself.")
	t.Log("Remember to remove the VMs, the NAT network, and the DHCP server (VBoxManage dhcpserver remove --netname testcluster1net\n")
}

func TestAddNewNode(t *testing.T) {
	cluster, ok := clustermanager.GetCluster("testclust1")
	if !ok {
		t.Log("Cluster 'testclust1' not foumd. This test is supposed to run after TestNewCluster.")
		t.FailNow()
	}
	node, err := cluster.AddUninitializedNode("testnode1")
	if err != nil {
		t.Logf("AddUninitializedNode failed with error:%v", err)
		t.FailNow()
	}

	t.Logf("Node:%+v", node)
	t.Logf("Cluster:%+v", cluster)
}

func TestDeleteNode(t *testing.T) {
	cluster, ok := clustermanager.GetCluster("testclust1")
	if !ok {
		t.Log("Cluster 'testclust1' not foumd. This test is supposed to run after TestNewCluster.")
		t.FailNow()
	}
	err := cluster.DeleteNode("testnode1", true)
	if err != nil {
		t.Logf("DeleteNode failed with error:%v", err)
		t.FailNow()
	}

	t.Logf("Cluster:%+v", cluster)
}

func TestDeleteCluster(t *testing.T) {
	err := clustermanager.DeleteCluster("testclust1")
	if err != nil {
		t.Logf("DeleteCluster failed with error:%v", err)
		t.FailNow()
	}

}
