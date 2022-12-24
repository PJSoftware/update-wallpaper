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
		source := filepath.Join(bingWallpaperFolder, file.Name())
		res, err := wallpaper.Resolution(source)
		if err != nil {
			continue
		}

		targetFolder := filepath.Join(folder, res.Name)
		err = os.MkdirAll(targetFolder, 0777)
		if err != nil {
			fmt.Printf("  Could not create %s\n", targetFolder)
			continue
		}

		target := filepath.Join(targetFolder, file.Name())
		wallpaper.Copy(source, target)
	}
}
