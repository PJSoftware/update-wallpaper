package wp_bing

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/pjsoftware/update-wallpaper/pkg/wallpaper"
)

func Update(folder string) {
	fmt.Printf("Updating BING images:\n")

	files, err := os.ReadDir(bingWallpaperFolder)
	if err != nil {
		log.Fatalf("bing: Update() error reading %s: %v", bingWallpaperFolder, err)
	}

	for _, file := range files {
		path := filepath.Join(bingWallpaperFolder, file.Name())
		res, err := wallpaper.Resolution(path)
		if err != nil {
			continue
		}

		fmt.Printf("- %s (%s - %dx%d)\n",file.Name(), res.Name, res.Width, res.Height)
	}
}