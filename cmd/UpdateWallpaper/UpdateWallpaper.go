package main

import (
	"fmt"
	"log"
	"os"

	"github.com/pjsoftware/update-wallpaper/pkg/config"
	"github.com/pjsoftware/update-wallpaper/pkg/splashscreen"
	"github.com/pjsoftware/update-wallpaper/pkg/spotlight"
	"github.com/pjsoftware/update-wallpaper/pkg/util"
)

var assets spotlight.Assets
var cfg config.Config

func main() {
	splashscreen.Show("UpdateWallpaper", uwVersion)

	logFile, exePath := initFiles()
	defer logFile.Close()
	cfg.Init(exePath, "UpdateWallpaper.ini")

	updateSpotlight()
	updateMomentum()
}

func initFiles() (*os.File, string) {
	exePath := util.GetEXEFolder()
	logFN := exePath + "UpdateWallpaper.log"
	_ = os.Remove(logFN)

	logFile, err := os.OpenFile(logFN, os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(logFile)
	return logFile, exePath
}

func updateSpotlight() {
	assets.Init(cfg)

	found := assets.Count()
	fmt.Printf("%d Spotlight images found\n", found)

	total, duplicates := assets.Compare()
	fmt.Printf("%d Existing wallpapers found\n", found)
	fmt.Printf("%d Spotlight assets match existing; skipping\n", duplicates)

	copied, replaced := assets.Copy()
	fmt.Printf("%d new images copied\n", copied)
	if replaced > 0 {
		fmt.Printf("%d existing images replaced\n", replaced)
	}
	log.Printf("Existing: %d; Incoming: %d; New: %d; Replaced: %d", total, found, copied, replaced)
}

func updateMomentum() {
}
