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

const version = "1.3.3"

// TODO: We have all these maps keyed by assetName; perhaps it would be better
// to have an Asset struct?
var assetBySize map[int64]map[string]string
var toBeCopied map[string]bool
var photoData map[string]map[string]string
var fileExt map[string]string
var config spotlight.Config

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
	fileExt = make(map[string]string)
	photoData = make(map[string]map[string]string)

	for _, file := range files {
		assetName := file.Name()
		assetPath := sourcePath + "/" + assetName
		if isWallpaper(assetName, width, height) {
			fileSize := file.Size()
			if _, ok := assetBySize[fileSize]; !ok {
				assetBySize[fileSize] = make(map[string]string)
			}
			assetBySize[fileSize][assetName] = cryptoSum(assetPath)
			toBeCopied[assetName] = true
			for _, image := range metadata.Images {
				if image.FileSize() == fileSize {
					// the chances of anything coming from Mars are a million to one
					// the chances of there being two images of this filesize are miniscule
					// however, death rays etc! So:
					// TODO: we need to look at comparing with sha256 value too
					photoData[assetName] = make(map[string]string)
					photoData[assetName]["copyright"] = image.Copyright()
					photoData[assetName]["description"] = image.Description()
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

func isWallpaper(assetName string, width, height int) bool {
	assetPath := config.SourcePath + "/" + assetName
	asset, err := os.Open(assetPath)
	if err != nil {
		return false // Cannot read, so not interested in it
	}
	defer asset.Close()

	image, err := jpeg.DecodeConfig(asset)
	if err == nil {
		fileExt[assetName] = "jpg"
	} else {
		image, err = png.DecodeConfig(asset)
		if err == nil {
			fileExt[assetName] = "png"
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
			for assetName, assetHash := range assetBySize[fileSize] {
				if wpHash == assetHash {
					toBeCopied[assetName] = false
					matchesFound++
					log.Printf("** %s matched with %s", assetName, file.Name())
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

	for assetName, tbc := range toBeCopied {
		if tbc {
			newPath := targetPath + "/"
			newName := assetName
			if _, ok := photoData[assetName]; ok {
				desc := photoData[assetName]["description"]
				cr := photoData[assetName]["copyright"]
				newName = newFilename(desc, cr)
				if !smart {
					newPath += prefix
				}
			} else {
				newPath += prefix
			}
			newName += "." + fileExt[assetName]
			newPath += newName
			nbytes, err := copyFile(assetName, newPath)
			if err == nil {
				log.Printf("New image: %s (copied from %s)", newName, assetName)
				fmt.Printf("Copied %d bytes of %s to %s\n", nbytes, assetName, newName)
				copied++
			} else {
				if nbytes == 0 {
					fmt.Printf("Error copying file: %v\n", err)
				} else {
					fmt.Printf("Copied %d bytes of %s to %s; unable to set file time\n", nbytes, assetName, newName)
				}
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

func copyFile(sName, dst string) (int64, error) {
	src := config.SourcePath + "/" + sName
	file, err := os.Stat(src)
	if err != nil {
		return 0, err
	}
	srcMTime := file.ModTime()

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	dest, err := os.Create(dst)
	if err != nil {
		return 0, err
	}

	nbytes, err := io.Copy(dest, source)
	dest.Close()

	err = os.Chtimes(dst, srcMTime, srcMTime)

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
