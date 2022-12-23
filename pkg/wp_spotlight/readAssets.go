package wp_spotlight

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/pjsoftware/update-wallpaper/pkg/sha"
	"github.com/pjsoftware/update-wallpaper/pkg/wallpaper"
)

// Assets gives us a better way to handle our Asset collection
type assets struct {
	metadata      assetMetadata
	matches   		int
	byName    		map[string]*asset
	sumBySize 		map[int64]map[string]string
	sourceFolder 	string
	targetFolder 	string
}

// Asset provides an interface to the contents of the Windows
// Spotlight Assets folder, some of which we are interested in.
type asset struct {
	name        string
	path        string
	isWallpaper bool

	copyright   string
	description string
	matched     bool

	toBeCopied  bool
	newName     string
	newPath     string
	replace     string
}

// readAssets prepares a list of currently available Spotlight assets
func readAssets(folder string) *assets {
	as := new(assets)
	as.sourceFolder = spotlightAssetFolder
	as.targetFolder = folder

	as.metadata.read()
	as.locateWallpapers()
	return as
}

// locateWallpapers called by Init()
func (as *assets) locateWallpapers() {
	as.byName = make(map[string]*asset)
	as.sumBySize = make(map[int64]map[string]string)

	files, err := os.ReadDir(as.sourceFolder)
	if err != nil {
		log.Fatalf("browse() error reading %s: %v", as.sourceFolder, err)
	}

	for _, file := range files {
		asset := new(asset)
		
		asset.name = file.Name()
		asset.path = filepath.Join(as.sourceFolder, asset.name)
		asset.identify()

		if !asset.isWallpaper {
			continue
		}

		as.addAsset(asset, file)
	}
}

func (as *assets) addAsset(asset *asset, file fs.DirEntry) {
	fileInfo, err := file.Info()
	if err != nil {
		return // do not add
	}
	fileSize := fileInfo.Size()
	fileHash, err := sha.FileHash(asset.path)
	if err != nil {
		log.Fatalf("Error calculating hash: %s, %v", asset.path, err)
	}

	if _, ok := as.sumBySize[fileSize]; !ok {
		as.sumBySize[fileSize] = make(map[string]string)
	}
	as.sumBySize[fileSize][asset.name] = fileHash

	for _, metadata := range as.metadata.imageMD {
		if metadata.fileSize == fileSize {
			if asset.matched {
				fmt.Printf("ASSET ALREADY MATCHED!!!\n")
			}
			// TODO: we should look at comparing with sha256 value too
			// on the billion-to-one chance we get two assets with an
			// identical size
			asset.copyright = metadata.copyright
			asset.description = metadata.description
			asset.matched = true
			if asset.description == NO_DESCRIPTION {
				asset.description += " (" + metadata.entityID + ")"
			}
		}
	}

	if asset.description == "" {
		asset.description = NO_DESCRIPTION + " (" + asset.name + ")"
		asset.copyright = NO_COPYRIGHT
	}

	asset.toBeCopied = true
	as.byName[asset.name] = asset
}

// identify examines asset files to determine whether they are wallpapers
func (a *asset) identify() {
	res, err := wallpaper.Resolution(a.path)
	if err != nil {
		a.isWallpaper = false
		return
	}

	a.isWallpaper = (res.Name == "HD")
}
