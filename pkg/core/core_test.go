package core_test

import (
	"os"
	"testing"

	"github.com/rajch/kutti/pkg/core"
	_ "github.com/rajch/kutti/pkg/vboxdriver"
)

func TestCacheDir(t *testing.T) {
	t.Log("Testing CacheDir()...")
	dir, err := core.CacheDir()
	if err != nil {
		t.Logf("CacheDir failed with error:%v\n", err)
		t.FailNow()
	}

	t.Logf("CacheDir returned %v. Checking for its presence...\n", dir)

	cdinfo, err := os.Stat(dir)
	if err != nil {
		t.Logf("CacheDir was not actually found. Error:%v\n", err)
		t.FailNow()
	}

	if !cdinfo.IsDir() {
		t.Log("CacheDir was not a directory.")
		t.FailNow()
	}
}

func TestNewCluster(t *testing.T) {
	driver, err := core.NewDriver("vbox")
	if err != nil {
		t.Logf("Could not fetch virtualbox driver:%v", err)
		t.FailNow()
	}

	cluster, err := core.NewCluster(driver, "feck", "manager", 10001)
	if err != nil {
		t.Logf("Error while creating cluster:%v", err)
		t.FailNow()
	}

	t.Log(cluster)
	t.Log("As of now, you will have to remove the cluster artifacts yourself.")
	t.Log("Remember to remove the VMs, the NAT network, and the DHCP server (VBoxManage dhcpserver remove --netname fecknet\n")
}

func TestSSHClient(t *testing.T) {
	t.Logf("Testing SSH client, assuming a server at localhost:10001")
	results, err := core.DefaultClient.RunWithResults("localhost:10001", "echo HOSTNAME: $(hostname) && echo PWD: $(pwd) && ls -l")
	if err != nil {
		t.Logf("SSH Client failed with error:%v", err)
		t.FailNow()
	}

	t.Logf("Results were:\n%s", results)
	t.Logf("Now Testing for failure, assuming a server at localhost:10001")
	results, err = core.DefaultClient.RunWithResults("localhost:10001", "nosuchcommand available;")
	if err == nil {
		t.Log("SSH Client should have failed")
		t.FailNow()
	}

	t.Logf("Error was:%v\nResults were:\n%s", err, results)

}
