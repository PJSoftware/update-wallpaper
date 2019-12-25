package main

import (
	"crypto/sha256"
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
	"strings"

	"./spotlight"
)

const version = "1.2"

var assetBySize map[int64]map[string]string
var toBeCopied map[string]bool
var fileName map[string]string // Let's cache the filenames so we don't need to re-extract
var photoData map[string]map[string]string
var fileExt map[string]string

func main() {
	fmt.Printf("UpdateSpotlight v%s -- by PJSoftware\n", version)
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

	found := browseAssets(config.SourcePath, config.Width, config.Height, metadata)
	total, dups := scanExisting(config.TargetPath)

	copied := 0
	if found > dups {
		copied = copyNewAssets(config.TargetPath, config.Prefix, config.SmartPrefix)
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
			assetBySize[fileSize][assetPath] = cryptoSum(assetPath)
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

func cryptoSum(filePath string) string {
	file, err := os.Open(filePath)
	defer file.Close()

	if err == nil {
		hash := sha256.New()
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
			wpHash := cryptoSum(filePath)
			for assetPath, assetHash := range assetBySize[fileSize] {
				if wpHash == assetHash {
					if toBeCopied[assetPath] {
						toBeCopied[assetPath] = false
						matchesFound++
					}
				}
			}
		}

	}

	fmt.Printf("%d Existing wallpapers found\n", wpFound)
	fmt.Printf("%d Spotlight assets match existing; skipping\n", matchesFound)
	return wpFound, matchesFound
}

func copyNewAssets(targetPath, prefix string, smart bool) int {
	copied := 0

	for assetPath, tbc := range toBeCopied {
		if tbc {
			newPath := targetPath + "/"
			newName := fileName[assetPath]
			if _, ok := photoData[assetPath]; ok {
				desc := photoData[assetPath]["description"]
				cr := photoData[assetPath]["copyright"]
				newName = newFilename(desc, cr)
				if !smart {
					newPath += prefix
				}
			} else {
				newPath += prefix
			}
			newName += "." + fileExt[assetPath]
			newPath += newName
			nbytes, err := copyFile(assetPath, newPath)
			if err == nil {
				log.Printf("New image: %s (copied from %s)", newName, fileName[assetPath])
				fmt.Printf("Copied %d bytes of %s to %s\n", nbytes, fileName[assetPath], newName)
				copied++
			} else {
				fmt.Printf("Error copying file: %v\n", err)
			}
		}
	}

	return copied
}

func newFilename(desc, cr string) string {
	re := regexp.MustCompile(` *[<>:"/\|?*]+ *`)
	desc = re.ReplaceAllString(desc, " + ")
	cr = re.ReplaceAllString(cr, " + ")

	re = regexp.MustCompile(` +`)
	desc = re.ReplaceAllString(desc, " ")
	cr = re.ReplaceAllString(cr, " ")

	desc = strings.TrimSpace(desc)
	cr = strings.TrimSpace(cr)

	hasSym, _ := regexp.MatchString(`^© `, cr)
	if !hasSym {
		cr = "© " + cr
	}

	return desc + " " + cr
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
