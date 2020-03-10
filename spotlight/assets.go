package spotlight

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
	"regexp"
	"strings"
)

// Assets gives us a better way to handle our Asset collection
type Assets struct {
	meta      MetaData
	config    Config
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
func (as *Assets) Init(config Config) {
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
			wpHash := cryptoSum(filePath)
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
			as.sumBySize[fileSize][asset.name] = cryptoSum(asset.path)
			asset.copyThis = true
			for _, image := range as.meta.Images {
				if image.FileSize() == fileSize {
					// TODO: we should look at comparing with sha256 value too
					// on the billion-to-one chance we get two assets with an
					// identical size
					asset.copyright = image.Copyright()
					asset.description = image.Description()
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

// cryptoSum called by browse() and Compare()
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

func (a *Asset) setNewName(cfg Config) {
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
	desc := a.description
	cr := a.copyright

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

	a.newName = desc + " " + cr
}

func (a *Asset) publish(cfg Config) int {
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
		return 0, err
	}
	srcMTime := file.ModTime()

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	dest, err := os.Create(a.newPath)
	if err != nil {
		return 0, err
	}

	nbytes, err := io.Copy(dest, source)
	dest.Close()

	err = os.Chtimes(a.newPath, srcMTime, srcMTime)

	return nbytes, err
}
