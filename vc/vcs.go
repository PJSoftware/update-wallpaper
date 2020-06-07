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
	detected vcs
}

// Detect which version control method (if any) is active on specified folder
func Detect(folder string) *Software {
	s := new(Software)
	s.folder = folder
	s.detected = none
	return s
}

// IsActive returns true if VC is active; otherwise returns false
func (s *Software) IsActive() bool {
	if s.detected == none {
		return false
	}
	return true
}

// Rename performs our file renaming via appropriate VC commands
func (s *Software) Rename(oldFN, newFN string, targetPath string) {
	old := targetPath + "/" + oldFN
	new := targetPath + "/" + newFN
	switch s.detected {
	case none:
		os.Rename(old, new)
	default:
		log.Printf("Unsupported VC software: %v", s.detected)
	}
}
