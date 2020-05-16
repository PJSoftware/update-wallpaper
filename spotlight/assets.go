package spotlight

import (
	"fmt"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/pjsoftware/win-spotlight/config"
	"github.com/pjsoftware/win-spotlight/errors"
	"github.com/pjsoftware/win-spotlight/util"
)

// Assets gives us a better way to handle our Asset collection
type Assets struct {
	meta      MetaData
	config    config.Config
	matches   int
	byName    map[string]*Asset
	sumBySize map[int64]map[string]string
}

// Asset provides an interface to the contents of the Windows
// Spotlight Assets folder, some of which we are interested in.
type Asset struct {
	name        string
	path        string
	extension   string
	copyright   string
	description string
	copyThis    bool
	newName     string
	newPath     string
}

// Init scans the asset folder to find all the valid assets
func (as *Assets) Init(config config.Config) {
	as.config = config
	as.meta.ImportAll()

	as.byName = make(map[string]*Asset)
	as.sumBySize = make(map[int64]map[string]string)
	as.browse()
}

// Count returns the number of valid Assets found
func (as *Assets) Count() int {
	return len(as.byName)
}

// Compare scans targetPath folder and compares to Assets
func (as *Assets) Compare() (int, int) {
	wpFound := 0
	matchesFound := 0
	files, err := ioutil.ReadDir(as.config.TargetPath)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		filePath := as.config.TargetPath + "/" + file.Name()
		fileSize := file.Size()
		wpFound++
		if _, ok := as.sumBySize[fileSize]; ok {
			wpHash := util.FileHash(filePath)
			for name, hash := range as.sumBySize[fileSize] {
				if wpHash == hash {
					as.byName[name].copyThis = false
					matchesFound++
					log.Printf("** '%s' matched with '%s'", name, file.Name())
				}
			}
		}

	}
	as.matches = matchesFound
	return wpFound, matchesFound
}

// Copy copies all new, non-matched assets to wallpaper
func (as *Assets) Copy() int {
	copied := 0

	if as.Count() <= as.matches {
		return copied
	}

	for _, asset := range as.byName {
		copied += asset.publish(as.config)
	}
	return copied
}

// browse called by Init()
func (as *Assets) browse() {
	files, err := ioutil.ReadDir(as.config.SourcePath)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		asset := new(Asset)
		asset.name = file.Name()
		asset.path = as.config.SourcePath + "/" + asset.name
		asset.extension = as.wpExtension(asset.path)

		if asset.extension != "" {
			fileSize := file.Size()
			if _, ok := as.sumBySize[fileSize]; !ok {
				as.sumBySize[fileSize] = make(map[string]string)
			}
			as.sumBySize[fileSize][asset.name] = util.FileHash(asset.path)
			asset.copyThis = true
			for _, image := range as.meta.Images {
				if image.FileSize() == fileSize {
					// TODO: we should look at comparing with sha256 value too
					// on the billion-to-one chance we get two assets with an
					// identical size
					asset.copyright = image.Copyright()
					asset.description = image.Description()
					if asset.description == "Unidentified Photo" {
						asset.description += " (" + asset.name + ")"
					}
				}
			}
			as.byName[asset.name] = asset
		}
	}
}

// wpExtension called by browse()
func (as *Assets) wpExtension(assetPath string) string {
	nullExt := ""
	ext := nullExt

	asset, err := os.Open(assetPath)
	if err != nil {
		return nullExt // Cannot read, so not interested in it
	}
	defer asset.Close()

	image, err := jpeg.DecodeConfig(asset)
	if err == nil {
		ext = "jpg"
	} else {
		image, err = png.DecodeConfig(asset)
		if err == nil {
			ext = "png"
		} else {
			return nullExt // Neither a JPEG nor a PNG, so not interested in it
		}
	}

	if image.Width != as.config.Width || image.Height != as.config.Height {
		return nullExt
	}

	return ext
}

func (a *Asset) setNewName(cfg config.Config) {
	a.newPath = cfg.TargetPath + "/"
	a.newName = a.name
	if a.copyright != "" {
		a.newFilename()
		if !cfg.SmartPrefix {
			a.newPath += cfg.Prefix
		}
	} else {
		a.newPath += cfg.Prefix
	}
	a.newName += "." + a.extension
	a.newPath += a.newName
}

func (a *Asset) newFilename() {
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

func (a *Asset) publish(cfg config.Config) int {
	if !a.copyThis {
		return 0
	}

	a.setNewName(cfg)
	if _, err := os.Stat(a.newPath); err == nil {
		log.Printf("* Skipped copying '%s'; different version already exists\n", a.newName)
		return 0
	}

	nbytes, err := a.copyFile(cfg.SourcePath)
	if err == nil {
		log.Printf("New image: %s (copied from %s)", a.newName, a.name)
		fmt.Printf("Copied %d bytes of %s to %s\n", nbytes, a.name, a.newName)
		return 1
	}

	if nbytes == 0 {
		fmt.Printf("Error copying file: %v\n", err)
		return 0
	}

	fmt.Printf("Copied %d bytes of '%s' to '%s'; unable to set file time\n", nbytes, a.name, a.newName)
	return 1
}

func (a *Asset) copyFile(fromFolder string) (int64, error) {
	src := fromFolder + "/" + a.name
	file, err := os.Stat(src)
	if err != nil {
		return 0, errors.E{Code: errors.EFILENOTFOUND, Message: "Source file not found"}
	}
	srcMTime := file.ModTime()

	source, err := os.Open(src)
	if err != nil {
		return 0, errors.E{Code: errors.EREADERROR, Message: "Could not read source file"}
	}
	defer source.Close()

	dest, err := os.Create(a.newPath)
	if err != nil {
		return 0, errors.E{Code: errors.EWRITEERROR, Message: "Could not create target file"}
	}

	nbytes, err := io.Copy(dest, source)
	dest.Close()

	err = os.Chtimes(a.newPath, srcMTime, srcMTime)

	return nbytes, err
}
