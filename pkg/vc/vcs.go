package vc

import (
	"fmt"
	"log"
	"os"
)

type vcs string

const (
	none vcs = "none"
	git  vcs = "git"
	svn  vcs = "svn"
)

// Software is the data for our Version Control handler
type Software struct {
	folder    string
	CanCommit bool
	Detected  vcs
}

// Detect which version control method (if any) is active on specified folder
func Detect(folder string) *Software {
	s := new(Software)
	s.folder = folder
	s.CanCommit = false
	s.Detected = none

	pushFolder(folder)
	if isGit() {
		s.Detected = git
	}
	popFolder()

	return s
}

// IsActive returns true if VC is active; otherwise returns false
func (s *Software) IsActive() bool {
	if s.Detected == none {
		return false
	}
	return true
}

// Rename performs our file renaming via appropriate VC commands
func (s *Software) Rename(oldFN, newFN string, targetPath string) {
	if s.Detected == none {
		old := targetPath + "/" + oldFN
		new := targetPath + "/" + newFN
		os.Rename(old, new)
		return
	}

	pushFolder(targetPath)
	switch s.Detected {
	case git:
		s.gitRename(oldFN, newFN)
	default:
		log.Printf("Unsupported VC software: %v", s.Detected)
	}
	popFolder()
}

// Update to latest state
func (s *Software) Update() {
	if s.Detected == none {
		return
	}

	fmt.Printf("Updating to latest version\n")
	switch s.Detected {
	case git:
		s.gitPull()
	default:
		log.Printf("Unsupported VC software: %v", s.Detected)
	}
}

// Add a file for committing
func (s *Software) Add(filePath string) {
	if s.Detected == none {
		return
	}

	switch s.Detected {
	case git:
		s.gitAdd(filePath)
	default:
		log.Printf("Unsupported VC software: %v", s.Detected)
	}
}

// Commit new/changed files to Version Control
func (s *Software) Commit(commitMsg string) {
	if s.Detected == none || !s.CanCommit {
		return
	}

	cn := os.Getenv("COMPUTERNAME")
	msg := commitMsg + " (" + cn + ")"

	fmt.Printf("Committing changes to repo\n")
	switch s.Detected {
	case git:
		s.gitCommit(msg)
	default:
		log.Printf("Unsupported VC software: %v", s.Detected)
	}
}

var folderStack []string

func pushFolder(folder string) {
	cwd, err := os.Getwd()
	if err != nil {
		log.Printf("Error: %s", err)
		return
	}

	err = os.Chdir(folder)
	if err != nil {
		log.Printf("Error: %s", err)
		return
	}

	folderStack = append(folderStack, cwd)
}

func popFolder() {
	if len(folderStack) == 0 {
		return
	}

	index := len(folderStack) - 1
	folder := folderStack[index]
	folderStack = folderStack[:index]
	os.Chdir(folder)
}
