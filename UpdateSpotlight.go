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
	"regexp"

	"./spotlight"
)

const sourceFolder = "Packages/Microsoft.Windows.ContentDeliveryManager_cw5n1h2txyewy/LocalState/Assets"

var localAppData = os.Getenv("LOCALAPPDATA")
var sourcePath = localAppData + "/" + sourceFolder

var assetBySize map[int64]map[string]string
var toBeCopied map[string]bool
var fileName map[string]string // Let's cache the filenames so we don't need to re-extract
var photoData map[string]map[string]string
var fileExt map[string]string

func main() {
	// First determine exepath and set LOG file location
	exePath := getEXEFolder()
	logFile, err := os.OpenFile(exePath+"UpdateSpotlight.log", os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	var config spotlight.Config
	config.Init(exePath)

	metadata := new(spotlight.MetaData)
	metadata.ImportAll()

	found := browseAssets(sourcePath, config.Width, config.Height, metadata)
	total, dups := scanExisting(config.TargetPath)

	copied := 0
	if found > dups {
		copied = copyNewAssets(config.TargetPath, config.Prefix)
	}
	fmt.Printf("%d new images copied\n", copied)
	log.Printf("Existing: %d; Incoming: %d; New: %d", total, found, copied)
}

func browseAssets(sourcePath string, width, height int, metadata *spotlight.MetaData) int {
	assetsFound := 0
	files, err := ioutil.ReadDir(sourcePath)
	if err != nil {
		log.Fatal(err)
	}

	assetBySize = make(map[int64]map[string]string)
	toBeCopied = make(map[string]bool)
	fileName = make(map[string]string)
	fileExt = make(map[string]string)
	photoData = make(map[string]map[string]string)

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
			for _, image := range metadata.Images {
				if image.FileSize() == fileSize {
					// the chances of anything coming from Mars are a million to one
					// the chances of there being two images of this filesize are miniscule
					// however, death rays etc! So:
					// TODO: we need to look at comparing with sha256 value too
					photoData[assetPath] = make(map[string]string)
					photoData[assetPath]["copyright"] = image.Copyright()
					photoData[assetPath]["description"] = image.Description()
				}
			}
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
			newPath := targetPath + "/" + prefix
			newName := fileName[assetPath]
			if _, ok := photoData[assetPath]; ok {
				// TODO: Take a closer look at our code to strip out invalid characters
				// Currently it is far too simplistic, and stripping out more than it should
				// Examples:
				//    'Grundarfjörður' was converted to 'Grundarfj r ur'
				//    'Kjölur' was converted to 'Kj lur'
				log.Printf("New image: %s (%s)", photoData[assetPath]["description"], photoData[assetPath]["copyright"])
				newName = photoData[assetPath]["description"] + " -- "
				re := regexp.MustCompile(` *[/|].*$`)
				newName += re.ReplaceAllString(photoData[assetPath]["copyright"], "")

				re = regexp.MustCompile(`[^- a-zA-Z0-9,]+`)
				newName = re.ReplaceAllString(newName, " ")

				re = regexp.MustCompile(` +`)
				newName = re.ReplaceAllString(newName, " ")

				re = regexp.MustCompile(`^ +`)
				newName = re.ReplaceAllString(newName, "")

				re = regexp.MustCompile(` +$`)
				newName = re.ReplaceAllString(newName, "")
				log.Printf("Geneated filename: %s", newName)
			}
			newName += "." + fileExt[assetPath]
			newPath += newName
			nbytes, err := copyFile(assetPath, newPath)
			if err == nil {
				fmt.Printf("Copied %d bytes of %s to %s\n", nbytes, fileName[assetPath], newName)
				copied++
			} else {
				fmt.Printf("Error copying file: %v\n", err)
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
