package main

import (
	"fmt"

	"github.com/pjsoftware/win-spotlight/config"
	"github.com/pjsoftware/win-spotlight/paths"
	"github.com/pjsoftware/win-spotlight/wallpaper"
)

const version = "1.0.0"

var cfg config.Config

func main() {
	welcomeMsg := fmt.Sprintf("DeleteDuplicates v%s -- by PJSoftware\n", version)
	fmt.Printf(welcomeMsg)

	exePath := paths.GetEXEFolder()
	cfg.Init(exePath)

	master := wallpaper.ImportFolder(cfg.TargetPath)
	archive := wallpaper.ImportFolder(cfg.TargetPath + "\\" + cfg.Archive)

	dh := cfg.DupHandler == "SVN-Rename"
	archive.DeleteDuplicates(master, dh)
}
