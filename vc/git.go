package vc

import (
	"os/exec"
)

func isGit() bool {
	cmd := exec.Command("git", "status")
	err := cmd.Run()
	if err != nil {
		return false
	}
	return true
}

func gitRename(old, new string) {
}
