package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"./spotlight"
)

const sourceFolder = "Packages/Microsoft.Windows.ContentDeliveryManager_cw5n1h2txyewy/LocalState/Assets"

var localAppData = os.Getenv("LOCALAPPDATA")
var sourcePath = localAppData + "/" + sourceFolder

var assetBySize map[int64]map[string]string
var toBeCopied map[string]bool
var fileName map[string]string // Let's cache the filenames so we don't need to re-extract
var fileExt map[string]string

func main() {
	exePath := getEXEFolder()
	logFile, err := os.OpenFile(exePath+"UpdateSpotlight.log", os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	var config spotlight.Config
	config.Init(exePath)

	found := browseAssets(sourcePath, config.Width, config.Height)
	total, dups := scanExisting(config.TargetPath)

	copied := 0
	if found > dups {
		copied = copyNewAssets(config.TargetPath, config.Prefix)
	}
	fmt.Printf("%d new images copied\n", copied)
	log.Printf("Existing: %d; Incoming: %d; New: %d", total, found, copied)
}

func browseAssets(sourcePath string, width, height int) int {
	assetsFound := 0
	files, err := ioutil.ReadDir(sourcePath)
	if err != nil {
		log.Fatal(err)
	}

	assetBySize = make(map[int64]map[string]string)
	toBeCopied = make(map[string]bool)
	fileName = make(map[string]string)
	fileExt = make(map[string]string)

	for _, file := range files {
		assetPath := sourcePath + "/" + file.Name()
		if isWallpaper(assetPath, width, height) {
			fileSize := file.Size()
			if _, ok := assetBySize[fileSize]; !ok {
				assetBySize[fileSize] = make(map[string]string)
			}
			assetBySize[fileSize][assetPath] = md5String(assetPath)
			toBeCopied[assetPath] = true
			fileName[assetPath] = file.Name()
			assetsFound++
		}
	}

	fmt.Printf("%d Spotlight images found\n", assetsFound)
	return assetsFound
}

func md5String(filePath string) string {
	file, err := os.Open(filePath)
	defer file.Close()

	if err == nil {
		hash := md5.New()
		if _, err := io.Copy(hash, file); err == nil {
			return hex.EncodeToString(hash.Sum(nil))
		}
	}

	return ""
}

func isWallpaper(filePath string, width, height int) bool {
	asset, err := os.Open(filePath)
	if err != nil {
		return false // Cannot read, so not interested in it
	}
	defer asset.Close()

	image, err := jpeg.DecodeConfig(asset)
	if err == nil {
		fileExt[filePath] = "jpg"
	} else {
		image, err = png.DecodeConfig(asset)
		if err == nil {
			fileExt[filePath] = "png"
		} else {
			return false // Neither a JPEG nor a PNG, so not interested in it
		}
	}

	if image.Width != width || image.Height != height {
		return false
	}

	return true
}

func scanExisting(targetPath string) (int, int) {
	wpFound := 0
	matchesFound := 0
	files, err := ioutil.ReadDir(targetPath)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		filePath := targetPath + "/" + file.Name()
		fileSize := file.Size()
		wpFound++
		if _, ok := assetBySize[fileSize]; ok {
			wpHash := md5String(filePath)
			for assetPath, assetHash := range assetBySize[fileSize] {
				if wpHash == assetHash {
					toBeCopied[assetPath] = false
					matchesFound++
				}
			}
		}

	}

	fmt.Printf("%d Existing wallpapers found\n", wpFound)
	fmt.Printf("%d Spotlight assets match existing; skipping\n", matchesFound)
	return wpFound, matchesFound
}

func copyNewAssets(targetPath, prefix string) int {
	copied := 0

	for assetPath, tbc := range toBeCopied {
		if tbc {
			newPath := targetPath + "/" + prefix + fileName[assetPath] + "." + fileExt[assetPath]
			nbytes, err := copyFile(assetPath, newPath)
			if err == nil {
				fmt.Printf("Copied %d bytes of %s to %s\n", nbytes, fileName[assetPath], targetPath)
				copied++
			}
		}
	}

	return copied
}

func copyFile(src, dst string) (int64, error) {
	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	dest, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer dest.Close()

	nbytes, err := io.Copy(dest, source)
	return nbytes, err
}

func getEXEFolder() string {
	exeFilename := os.Args[0]
	exeFolder := filepath.Dir(exeFilename)
	exeAbsFolder, err := filepath.Abs(exeFolder)
	if err != nil {
		log.Printf("Unable to determine EXE folder: %v", err)
		return ""
	}

	return exeAbsFolder + "\\"
}