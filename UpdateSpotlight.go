package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"./spotlight"
)

const version = "1.4"

var assets spotlight.Assets
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

	// Must initialise config before assets
	config.Init(exePath)
	assets.Init(config.SourcePath, config.Width, config.Height)

	found := assets.Count()
	fmt.Printf("%d Spotlight images found\n", found)

	total, dups := assets.Compare(config.TargetPath)
	fmt.Printf("%d Existing wallpapers found\n", found)
	fmt.Printf("%d Spotlight assets match existing; skipping\n", dups)

	copied := 0
	if found > dups {
		copied = copyNewAssets(config.TargetPath, config.Prefix, config.SmartPrefix)
	}
	fmt.Printf("%d new images copied\n", copied)
	log.Printf("Existing: %d; Incoming: %d; New: %d", total, found, copied)
}

func copyNewAssets(targetPath, prefix string, smart bool) int {
	copied := 0

	for _, asset := range assets.ByName {
		if asset.CopyThis {
			newPath := targetPath + "/"
			newName := asset.Name
			if asset.Copyright != "" {
				desc := asset.Description
				cr := asset.Copyright
				newName = newFilename(desc, cr)
				if !smart {
					newPath += prefix
				}
			} else {
				newPath += prefix
			}
			newName += "." + asset.Extension
			newPath += newName
			nbytes, err := copyFile(asset.Name, newPath)
			if err == nil {
				log.Printf("New image: %s (copied from %s)", newName, asset.Name)
				fmt.Printf("Copied %d bytes of %s to %s\n", nbytes, asset.Name, newName)
				copied++
			} else {
				if nbytes == 0 {
					fmt.Printf("Error copying file: %v\n", err)
				} else {
					fmt.Printf("Copied %d bytes of %s to %s; unable to set file time\n", nbytes, asset.Name, newName)
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
