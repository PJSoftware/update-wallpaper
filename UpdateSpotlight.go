package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"./spotlight"
)

const version = "1.4"

var assets spotlight.Assets
var config spotlight.Config

func main() {
	fmt.Printf("UpdateSpotlight v%s -- by PJSoftware\n", version)
	// First determine exepath and set LOG file location
	exePath := getEXEFolder()
	logfn := exePath + "UpdateSpotlight.log"
	_ = os.Remove(logfn)

	logFile, err := os.OpenFile(logfn, os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)
	log.Printf("UpdateSpotlight v%s -- by PJSoftware", version)

	// Must initialise config before assets
	config.Init(exePath)
	assets.Init(config)

	found := assets.Count()
	fmt.Printf("%d Spotlight images found\n", found)

	total, dups := assets.Compare()
	fmt.Printf("%d Existing wallpapers found\n", found)
	fmt.Printf("%d Spotlight assets match existing; skipping\n", dups)

	copied := assets.Copy()
	fmt.Printf("%d new images copied\n", copied)
	log.Printf("Existing: %d; Incoming: %d; New: %d", total, found, copied)
}

func getEXEFolder() string {
	exeFilename := os.Args[0]
	exeFolder := filepath.Dir(exeFilename)
	exeAbsFolder, err := filepath.Abs(exeFolder)
	if err != nil {
		log.Printf("Unable to determine EXE folder: %v", err)
		return ""
	}

	return exeAbsFolder + "\\"
}
