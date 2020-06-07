package vc

import (
	"fmt"
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

func (s *Software) gitRename(old, new string) {
	cmd := exec.Command("git", "mv", old, new)
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error renaming with git: %s\n", err)
		return
	}
	s.CanCommit = true
}

func (s *Software) gitPull() {
	cmd := exec.Command("git", "pull")
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error pulling from origin: %s\n", err)
		return
	}
}

func (s *Software) gitAdd(file string) {
	cmd := exec.Command("git", "add", file)
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error adding file with git: %s\n", err)
		return
	}
	s.CanCommit = true
}

func (s *Software) gitCommit(msg string) {
	cmd := exec.Command("git", "commit", "-m", msg)
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error committing to local repo: %s\n", err)
		return
	}
	s.CanCommit = false

	cmd = exec.Command("git", "push")
	err = cmd.Run()
	if err != nil {
		fmt.Printf("Error pushing to remote repo: %s\n", err)
	}
}
