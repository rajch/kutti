// +build !windows

package vboxdriver

import "os/exec"

func findvboxmanage() (string, error) {

	// Try looking it up on the path
	toolpath, err := exec.LookPath("VBoxManage")
	if err == nil {
		return toolpath, nil
	}

	// Give up
	return "", err
}
