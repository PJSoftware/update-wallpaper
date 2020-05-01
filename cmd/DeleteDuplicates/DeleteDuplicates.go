package main

import (
	"github.com/pjsoftware/win-spotlight/config"
	"github.com/pjsoftware/win-spotlight/paths"
	"github.com/pjsoftware/win-spotlight/splashscreen"
	"github.com/pjsoftware/win-spotlight/wallpaper"
)

var cfg config.Config

func main() {
	splashscreen.Show("DeleteDuplicates", "1.0.0")

	exePath := paths.GetEXEFolder()
	cfg.Init(exePath)

	master := wallpaper.ImportFolder(cfg.TargetPath)
	archive := wallpaper.ImportFolder(cfg.TargetPath + "\\" + cfg.Archive)

	dh := cfg.DupHandler == "SVN-Rename"
	archive.DeleteDuplicates(master, dh)
}
