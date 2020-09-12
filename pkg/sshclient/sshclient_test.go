package sshclient

import "testing"

func TestSSHClient(t *testing.T) {
	t.Logf("Testing SSH client, assuming a server at localhost:10001")
	client := New("kuttiadmin", "Pass@word1")
	results, err := client.RunWithResults("localhost:10001", "echo HOSTNAME: $(hostname) && echo PWD: $(pwd) && ls -l")
	if err != nil {
		t.Logf("SSH Client failed with error:%v", err)
		t.FailNow()
	}

	t.Logf("Results were:\n%s", results)
	t.Logf("Now Testing for failure, assuming a server at localhost:10001")
	results, err = client.RunWithResults("localhost:10001", "nosuchcommand available;")
	if err == nil {
		t.Log("SSH Client should have failed")
		t.FailNow()
	}

	t.Logf("Error was:%v\nResults were:\n%s", err, results)
}
