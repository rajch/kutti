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
}
