package wp_bing

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pjsoftware/update-wallpaper/pkg/wallpaper"
)

func Update(folder string) {	
	files, err := os.ReadDir(bingWallpaperFolder)
	if err != nil {
		fmt.Printf("Skipping BING processing; no folder found\n");
		return
	}
	fmt.Printf("Updating BING images:\n")

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
