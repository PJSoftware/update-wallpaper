package wp_spotlight

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/pjsoftware/update-wallpaper/pkg/sha"
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
					if isUnidentified(file.Name()) && as.byName[name].hasName() {
						log.Printf("** '%s' will replace existing '%s'", name, file.Name())
						as.byName[name].replace = filePath
					} else {
						as.byName[name].toBeCopied = false
						matchesFound++
						log.Printf("** '%s' matched with '%s'", name, file.Name())
					}
				}
			}
		}

	}
	as.matches = matchesFound
	return wpFound, matchesFound
}

func isUnidentified(fn string) bool {
	badPrefix := []string{
		NO_DESCRIPTION,
		"ZZZ_",
	}

	for _, prefix := range badPrefix {
		if startsWith(fn, prefix) {
			return true
		}
	}
	return false
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



func (a *asset) hasName() bool {
	return !startsWith(a.description, NO_DESCRIPTION)
}

func startsWith(testing string, target string) bool {
	lenTarget := len(target)
	if len(testing) < lenTarget {
		return false
	}

	return testing[0:lenTarget] == target
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
		log.Printf("* Skipped copying '%s'; different version already exists\n", a.newName)
		return 0, 0
	}

	if a.replace != "" {
		old := targetPath + "/" + a.replace
		new := targetPath + "/" + a.newName
		os.Rename(old, new)
		log.Printf("New image %s replaced existing %s", a.newName, a.replace)
		fmt.Printf("New name %s for existing unidentified image\n", a.newName)
		return 0, 1
	}

	err := a.copyFile(sourcePath)
	if err == nil {
		log.Printf("New image: %s (copied from %s)", a.newName, a.name)
		fmt.Printf("Copied %s to %s\n", a.name, a.newName)
		return 1, 0
	} else {
		fmt.Printf("Error copying file: %v\n", err)
		return 0, 0
	}

}

func (a *asset) copyFile(fromFolder string) error {
	return wallpaper.Copy(fromFolder + "/" + a.name, a.newPath)
}
