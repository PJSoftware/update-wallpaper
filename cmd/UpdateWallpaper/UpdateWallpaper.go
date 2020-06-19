package main

import (
	"fmt"
	"log"
	"os"

	"github.com/pjsoftware/win-spotlight/config"
	"github.com/pjsoftware/win-spotlight/paths"
	"github.com/pjsoftware/win-spotlight/splashscreen"
	"github.com/pjsoftware/win-spotlight/spotlight"
	"github.com/pjsoftware/win-spotlight/vc"
)

var assets spotlight.Assets
var cfg config.Config

func main() {
	splashscreen.Show("UpdateWallpaper")

	logFile, exePath := initFiles()
	defer logFile.Close()
	cfg.Init(exePath)

	updateSpotlight()
	updateMomentum()
}

func initFiles() (*os.File, string) {
	exePath := paths.GetEXEFolder()
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
	useVC := vc.Detect(cfg.TargetPath)
	useVC.Update()

	found := assets.Count()
	fmt.Printf("%d Spotlight images found\n", found)

	total, duplicates := assets.Compare()
	fmt.Printf("%d Existing wallpapers found\n", found)
	fmt.Printf("%d Spotlight assets match existing; skipping\n", duplicates)

	copied, replaced := assets.Copy(useVC)
	fmt.Printf("%d new images copied\n", copied)
	if replaced > 0 {
		fmt.Printf("%d existing images replaced\n", replaced)
	}
	log.Printf("Existing: %d; Incoming: %d; New: %d; Replaced: %d", total, found, copied, replaced)
	useVC.Commit("Add new Spotlight Files")
}

func updateMomentum() {
}
