package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"./spotlight"
)

const version = "1.4.2"

var assets spotlight.Assets
var config spotlight.Config

func main() {
	logFile, exePath := initFiles()
	defer logFile.Close()

	welcomeMsg := fmt.Sprintf("UpdateSpotlight v%s -- by PJSoftware\n", version)
	fmt.Printf(welcomeMsg)
	log.Printf(welcomeMsg)

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

func initFiles() (*os.File, string) {
	exePath := getEXEFolder()
	logfn := exePath + "UpdateSpotlight.log"
	_ = os.Remove(logfn)

	logFile, err := os.OpenFile(logfn, os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(logFile)
	return logFile, exePath
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
