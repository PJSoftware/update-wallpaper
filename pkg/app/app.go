package app

import (
	"log"
	"os"
	"path/filepath"
)

type App struct {
	name     string
	version  string
	wpFolder string
}

func NewApp(name, version string) *App {
	app := new(App)
	app.name = name
	app.version = version
	return app
}

func (a *App) InitFolder() {
	exeFilename := os.Args[0]
	exeFolder := filepath.Dir(exeFilename)
	exeAbsFolder, err := filepath.Abs(exeFolder)
	if err != nil {
		log.Fatalf("Unable to determine EXE folder: %v", err)
	}

	a.wpFolder = exeAbsFolder
}

func (a *App) WallpaperFolder(name string) string {
	return filepath.Join(a.wpFolder, name)
}
