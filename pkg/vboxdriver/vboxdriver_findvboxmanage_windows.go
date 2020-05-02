// +build windows

package vboxdriver

import (
	"os"
	"os/exec"
	"path"
)

func findvboxmanage() (string, error) {

	// First, try looking it up on the path
	toolpath, err := exec.LookPath("VBoxManage.exe")
	if err == nil {
		return toolpath, nil
	}

	// Then try looking in well-known places
	// Try %ProgramFiles% first, then hardcode
	progfileslocation := os.Getenv("ProgramFiles")
	if progfileslocation != "" {
		toolpath = path.Join(progfileslocation, "Oracle", "VirtualBox", "VBoxManage.exe")
	} else {
		toolpath = "C:\\Program Files\\Oracle\\VirtualBox\\VBoxManage.exe"
	}

	// If it exists. we're good
	if _, err = os.Stat(toolpath); err == nil {
		return toolpath, nil
	}

	// Give up
	return "", err
}
