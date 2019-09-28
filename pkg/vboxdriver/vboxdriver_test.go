package vboxdriver

import (
	"os"
	"testing"
)

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
		t.Log("VboxManage should not have been found.")
		t.Fail()
	} else {
		t.Logf("Error was: %v\n", err)
	}

	os.Setenv("PATH", oldpath)

}

func TestListNetworks(t *testing.T) {
	drv, err := New()
	if err != nil {
		t.Logf("Error occured: %v\n", err)
		t.Fail()
		return
	}

	t.Log("Testing ListNetworks...")

	_, err = drv.ListNetworks()
	if err != nil {
		t.Logf("Error in ListNetworks: %v\n", err)
		t.Fail()
	}
}

func TestListHosts(t *testing.T) {
	drv, err := New()
	if err != nil {
		t.Logf("Error occured: %v\n", err)
		t.Fail()
		return
	}

	t.Log("Testing ListHosts...")

	hosts, err := drv.ListHosts()
	if err != nil {
		t.Logf("Error in ListHots: %v\n", err)
		t.Fail()
	}

	for _, host := range hosts {
		t.Log(host.Name() + "," + host.Status())
	}
}

func TestNetworkOperations(t *testing.T) {
	t.Log("Creating New VBoxDriver...")

	drv, err := New()
	if err != nil {
		t.Logf("Error in New: %v\n", err)
		t.FailNow()
	}

	t.Log("Testing CreateNetwork...")

	nw, err := drv.CreateNetwork("zintakova")
	if err != nil {
		t.Logf("Error in CreateNetwork: %v\n", err)
		t.FailNow()
	}

	if nw.Name() != "zintakova" {
		t.Logf("Wrong name returned. Wanted zintakova, got %v.\n", nw.Name())
		t.FailNow()
	}

	t.Log("CreateNetwork worked as expected. Calling again with same parameters...")
	nw, err = drv.CreateNetwork("zintakova")
	if err == nil {
		t.Log("The second call to CreateNetwork should have failed.")
		t.FailNow()
	}

	t.Logf("Second call errored as expected, with %v. Calling TestListNetworks...", err)
	t.Log("Testing ListNetworks...")

	_, err = drv.ListNetworks()
	if err != nil {
		t.Logf("Error in ListNetworks: %v\n", err)
		t.FailNow()
	}

	t.Log("ListNetwork seems to have worked. Now calling CreateHost...")
	newnode, err := drv.CreateHost("champu", "zintakova", 0)
	if err != nil {
		t.Logf("Error from CreateHost: %v\n", err)
		t.FailNow()
	} else {

		t.Log(newnode)

		t.Logf("CreateHost seems to have worked. Now starting the new host...")
		err = newnode.Start()
		if err != nil {
			t.Logf("Host starting failed with error: %v", err)
			t.Fail()
		} else {

			t.Logf("Host starting worked. Now waiting twenty seconds, and forwarding SSH port...")
			newnode.WaitForStateChange(20)
			err = newnode.ForwardSSHPort(10001)
			if err != nil {
				t.Logf("SSH port forwarding failed with error: %v", err)
				t.Fail()
			}

			sshaddr := newnode.SSHAddress()
			if sshaddr != "localhost:10001" {
				t.Logf("SSH port mapping does not appear to be successful. Value was %s", sshaddr)
				t.Fail()
			} else {
				t.Log("SSH Port mapping successful.")
			}

		}

		t.Log("Now stopping host...")
		err = newnode.Stop()
		if err != nil {
			t.Logf("Error stopping host: %v", err)
			t.Fail()
		}

	}

	t.Logf("CreateNode seems to have created node with name %s and status %s. Now waiting 20 seconds and calling DeleteHost...", newnode.Name(), newnode.Status())
	newnode.WaitForStateChange(20)
	err = drv.DeleteHost("champu", "zintakova")
	if err != nil {
		t.Logf("Error from DeleteHost: %v\n", err)
		t.FailNow()
	}

	t.Log("DeleteNode seems to have worked. Now calling DeleteNetwork...")
	err = drv.DeleteNetwork("zintakova")
	if err != nil {
		t.Logf("Error from DeleteNetwork: %v\n", err)
		t.FailNow()
	}

}
