package core

import (
	"os"
	"testing"
)

func TestCacheDir(t *testing.T) {
	t.Log("Testing CacheDir()...")
	dir, err := CacheDir()
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
