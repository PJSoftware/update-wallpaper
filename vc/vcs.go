package vc

import (
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
	folder   string
	Detected vcs
}

// Detect which version control method (if any) is active on specified folder
func Detect(folder string) *Software {
	s := new(Software)
	s.folder = folder
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
		gitRename(oldFN, newFN)
	default:
		log.Printf("Unsupported VC software: %v", s.Detected)
	}
	popFolder()
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
	} else {
		index := len(folderStack) - 1
		folder := folderStack[index]
		folderStack = folderStack[:index]
		os.Chdir(folder)
	}
}
