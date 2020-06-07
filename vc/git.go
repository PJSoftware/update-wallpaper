package vc

import (
	"fmt"
	"os/exec"
)

func isGit(folder string) bool {
	cmd := exec.Command("git", "status")
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return false
	}
	return true
}

func gitRename(old, new string) {
}
