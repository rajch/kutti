package clustermanager_test

// package core_test

import (
	"os"
	"path"
	"testing"
	"time"

	"github.com/rajch/kutti/internal/pkg/kuttilog"

	"github.com/rajch/kutti/pkg/clustermanager"
	"github.com/rajch/kutti/pkg/core"
	_ "github.com/rajch/kutti/pkg/vboxdriver"
)

func TestInit(t *testing.T) {
	kuttilog.Setloglevel(4)

	configdir, _ := core.ConfigDir()
	configfile := path.Join(configdir, "clusters.json")
	_, err := os.Stat(configfile)
	if err != nil {
		t.Logf("Config file not found after init. This is a problem:%v\n", err)
		t.FailNow()
	}

}

func TestDrivers(t *testing.T) {
	var result bool
	resultp := &result

	clustermanager.ForEachDriver(func(d *clustermanager.Driver) bool {
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
	err := clustermanager.NewEmptyCluster("testclust1", "1.18", "vbox")
	if err != nil {
		t.Logf("NewEmptyCluster failed with error:%v", err)
		t.FailNow()
	}

	cluster, _ := clustermanager.GetCluster("testclust1")
	t.Log(*cluster)
}

func TestAddNewNode(t *testing.T) {
	cluster, ok := clustermanager.GetCluster("testclust1")
	if !ok {
		t.Log("Cluster 'testclust1' not foumd. This test is supposed to run after TestNewCluster.")
		t.FailNow()
	}
	node, err := cluster.NewUninitializedNode("testnode1")
	if err != nil {
		t.Logf("NewUninitializedNode failed with error:%v", err)
		t.FailNow()
	}

	t.Logf("Node:%+v", node)
	t.Logf("Cluster:%+v", cluster)
	t.Log("Waiting 10 seconds...")
	time.Sleep(time.Duration(10) * time.Second)
}

func TestForwardSSHPort(t *testing.T) {
	cluster, ok := clustermanager.GetCluster("testclust1")
	if !ok {
		t.Log("Cluster 'testclust1' not foumd. This test is supposed to run after TestNewCluster.")
		t.FailNow()
	}

	node, ok := cluster.Nodes["testnode1"]
	if !ok {
		t.Log("Cluster 'testnode1' not foumd. This test is supposed to run after TestAddNewNode.")
		t.FailNow()
	}

	/*
		err := node.Start()
		if err != nil {
			t.Logf("Error starting node testnode1: %v", err)
			t.FailNow()
		}

		t.Log("Waiting 50 seconds, then forwarding SSH Port...")
		time.Sleep(time.Duration(50) * time.Second)
	*/
	err := node.ForwardSSHPort(9091)
	if err != nil {
		t.Logf("Could not forward SSH port: %v", err)
		t.FailNow()
	}

	t.Log("SSH port forwarding successful. Waiting 10 seconds...")
	time.Sleep(time.Duration(10) * time.Second)
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
	err := clustermanager.DeleteCluster("testclust1", false)
	if err != nil {
		t.Logf("DeleteCluster failed with error:%v", err)
		t.FailNow()
	}
}
