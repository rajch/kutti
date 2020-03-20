package core_test

import (
	"os"
	"testing"

	"github.com/rajch/kutti/pkg/core"
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

func TestConfigDir(t *testing.T) {
	t.Log("Testing ConfigDir()...")
	dir, err := core.ConfigDir()
	if err != nil {
		t.Logf("ConfigDir failed with error:%v\n", err)
		t.FailNow()
	}

	t.Logf("ConfigDir returned %v. Checking for its presence...\n", dir)

	cdinfo, err := os.Stat(dir)
	if err != nil {
		t.Logf("ConfigDir was not actually found. Error:%v\n", err)
		t.FailNow()
	}

	if !cdinfo.IsDir() {
		t.Log("ConfigDir was not a directory.")
		t.FailNow()
	}
}
