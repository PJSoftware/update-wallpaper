package main

import (
	"fmt"
	"log"
	"os"

	"github.com/pjsoftware/win-spotlight/config"
	"github.com/pjsoftware/win-spotlight/paths"
	"github.com/pjsoftware/win-spotlight/splashscreen"
	"github.com/pjsoftware/win-spotlight/spotlight"
)

const version = "1.4.3"

var assets spotlight.Assets
var cfg config.Config

func main() {
	splashscreen.Show("UpdateSpotlight")

	logFile, exePath := initFiles()
	defer logFile.Close()

	// Must initialise config before assets
	cfg.Init(exePath)
	assets.Init(cfg)

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
	exePath := paths.GetEXEFolder()
	logfn := exePath + "UpdateSpotlight.log"
	_ = os.Remove(logfn)

	logFile, err := os.OpenFile(logfn, os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(logFile)
	return logFile, exePath
}
