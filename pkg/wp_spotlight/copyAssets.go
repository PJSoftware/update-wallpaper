package wp_spotlight

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/pjsoftware/update-wallpaper/pkg/sha"
	"github.com/pjsoftware/update-wallpaper/pkg/util"
	"github.com/pjsoftware/update-wallpaper/pkg/wallpaper"
)

/////////////////////////////////////////////////////////////////////////////
// Code above works; code below not so much
/////////////////////////////////////////////////////////////////////////////

// count returns the number of valid Assets found
func (as *assets) count() int {
	return len(as.byName)
}

// compareWithExisting scans targetPath folder and compares to Assets
func (as *assets) compareWithExisting() (int, int) {
	wpFound := 0
	matchesFound := 0
	files, err := os.ReadDir(as.targetFolder)
	if err != nil {
		log.Fatalf("compare error reading %s: %v", as.targetFolder, err)
	}

	for _, file := range files {
		filePath := as.targetFolder + "/" + file.Name()
		fileInfo, _ := file.Info()
		fileSize := fileInfo.Size()
		wpFound++
		if _, ok := as.sumBySize[fileSize]; ok {
			existingHash, err := sha.FileHash(filePath)
			if err != nil {
				log.Fatalf("Error calculating hash: %s, %v", filePath, err)
			}
			for name, assetHash := range as.sumBySize[fileSize] {
				if existingHash == assetHash {
					as.byName[name].replace = file.Name()
				}
			}
		}

	}
	as.matches = matchesFound
	return wpFound, matchesFound
}

// Copy copies all new, non-matched assets to wallpaper
func (as *assets) Copy() (int, int) {
	copied := 0
	renamed := 0

	if as.count() <= as.matches {
		return copied, renamed
	}

	for _, asset := range as.byName {
		cc, rc := asset.publish(as.sourceFolder, as.targetFolder)
		copied += cc
		renamed += rc
	}
	return copied, renamed
}

func (a *asset) setNewName(targetPath string) {
	a.newPath = targetPath + "/"
	a.newName = a.name
	a.newFilename()
	a.newName += ".jpg"
	a.newPath += a.newName
}

func (a *asset) newFilename() {
	type repl struct {
		in  string
		out string
	}
	rs := []repl{
		{`[<>:"/\|?*]+`, ` + `},
		{`\s+`, ` `},
		{`(\s*©)+\s*`, ` © `},
		{`(\s+[+])+\s+`, ` + `},
	}

	nfn := strings.TrimSpace(a.description + " © " + a.copyright)

	for _, r := range rs {
		re := regexp.MustCompile(r.in)
		nfn = re.ReplaceAllString(nfn, r.out)
	}

	a.newName = nfn
}

func (a *asset) publish(sourcePath, targetPath string) (int, int) {
	if !a.toBeCopied {
		return 0, 0
	}

	a.setNewName(targetPath)
	if _, err := os.Stat(a.newPath); err == nil {
		return 0, 0
	}

	if a.replace != "" {
		if util.FirstN(a.newName, len(NO_DESCRIPTION)) == NO_DESCRIPTION {
			// fmt.Printf("- Skip attempt to rename to %s\n", a.newName)
			return 0, 0
		}

		old := targetPath + "/" + a.replace
		new := targetPath + "/" + a.newName
		fmt.Printf("- Renaming: %s\n        to: %s\n", a.replace, a.newName)
		os.Rename(old, new)
		setTime(new)
		return 0, 1
	}

	_, err := a.copyFile(sourcePath)
	if err == nil {
		fmt.Printf("- Copying: %s\n", a.newName)
		setTime(targetPath + "/" + a.newName)
		return 1, 0
	}

	fmt.Printf("Error copying file: %v\n", err)
	return 0, 0
}

func (a *asset) copyFile(fromFolder string) (bool, error) {
	return wallpaper.Copy(fromFolder+"/"+a.name, a.newPath)
}

func setTime(fileName string) {
	currentTime := time.Now().Local()
	err := os.Chtimes(fileName, currentTime, currentTime)
	if err != nil {
		fmt.Printf("  > error changing file time: %v\n", err)
	}
}
