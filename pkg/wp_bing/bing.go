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
		fmt.Printf("Skipping BING processing; no folder found\n")
		return
	}
	fmt.Printf("Updating BING images:\n")
	os.MkdirAll(folder, 0777)

	copied := 0
	skipped := 0
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
		success, err := wallpaper.Copy(source, target)
		if success {
			copied++
		} else if err != nil {
			fmt.Printf("error copying file %s: %v", file.Name(), err)
		} else {
			skipped++
		}
	}

	fmt.Printf("  %d new wallpapers copied; %d existing files skipped\n", copied, skipped)
}
