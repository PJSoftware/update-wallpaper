package wp_bing

import (
	"os"
	"path/filepath"
)

const bingFOLDER = "Microsoft\\BingWallpaperApp\\WPImages"

var bingWallpaperFolder = ""

func init() {
	local := os.Getenv("LOCALAPPDATA")
  bingWallpaperFolder = filepath.Join(local, bingFOLDER)
}
