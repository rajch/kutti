package vboxdriver

import (
	"os"
	"testing"
)

func TestNew(t *testing.T) {
	t.Log("Testing New with default PATH...")

	drv, err := New()
	if err != nil {
		t.Logf("Error occured: %v\n", err)
		t.Fail()
		return
	}

	t.Logf("Returned path was: %v", drv.vboxmanagepath)

	t.Log("Testing New with bad PATH...")

	oldpath := os.Getenv("PATH")
	os.Setenv("PATH", "/")

	drv, err = New()
	if err == nil {
		t.Log("Vbox should not have been found.")
		t.Fail()
	} else {
		t.Logf("Error was: %v\n", err)
	}

	os.Setenv("PATH", oldpath)

}

func TestRunWithResults(t *testing.T) {

	t.Log("Testing runwithresults with 'hostname'...")
	output, err := runwithresults("hostname")
	if err != nil {
		t.Logf("Exec failed with error:%v\n", err)
		t.Fail()
	}
	t.Logf("Output was: \n'%v'\n", output)

	t.Log("Testing runwithresults with 'hostname -i'...")
	output, err = runwithresults("hostname", "-i")
	if err != nil {
		t.Logf("Exec failed with error:%v\n", err)
		t.Fail()
	}
	t.Logf("Output was: \n'%v'\n", output)
}

func TestListNetworks(t *testing.T) {
	drv, err := New()
	if err != nil {
		t.Logf("Error occured: %v\n", err)
		t.Fail()
		return
	}

	t.Logf("Returned path was: %v", drv.vboxmanagepath)

	err = listnetworks(drv.vboxmanagepath)
	if err != nil {
		t.Fatal(err)
	}
}
