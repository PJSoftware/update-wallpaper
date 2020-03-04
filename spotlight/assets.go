package spotlight

import (
	"crypto/sha256"
	"encoding/hex"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"log"
	"os"
)

// Assets gives us a better way to handle our Asset collection
type Assets struct {
	Folder    string
	Meta      MetaData
	ByName    map[string]*Asset
	SumBySize map[int64]map[string]string
}

// Asset provides an interface to the contents of the Windows
// Spotlight Assets folder, some of which we are interested in.
type Asset struct {
	Name        string
	Path        string
	Extension   string
	Copyright   string
	Description string
	CopyThis    bool
}

// Init scans the asset folder to find all the valid assets
func (as *Assets) Init(assetFolder string, width, height int) {
	as.Folder = assetFolder
	as.Meta.ImportAll()

	as.ByName = make(map[string]*Asset)
	as.SumBySize = make(map[int64]map[string]string)
	as.browse(width, height)
}

// Count returns the number of valid Assets found
func (as *Assets) Count() int {
	return len(as.ByName)
}

// Compare scans targetPath folder and compares to Assets
func (as *Assets) Compare(targetPath string) (int, int) {
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
		if _, ok := as.SumBySize[fileSize]; ok {
			wpHash := cryptoSum(filePath)
			for name, hash := range as.SumBySize[fileSize] {
				if wpHash == hash {
					as.ByName[name].CopyThis = false
					matchesFound++
					log.Printf("** '%s' matched with '%s'", name, file.Name())
				}
			}
		}

	}
	return wpFound, matchesFound
}

// browse called by Init()
func (as *Assets) browse(width, height int) {
	files, err := ioutil.ReadDir(as.Folder)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		asset := new(Asset)
		asset.Name = file.Name()
		asset.Path = as.Folder + "/" + asset.Name
		asset.Extension = wpExtension(asset.Path, width, height)

		if asset.Extension != "" {
			fileSize := file.Size()
			if _, ok := as.SumBySize[fileSize]; !ok {
				as.SumBySize[fileSize] = make(map[string]string)
			}
			as.SumBySize[fileSize][asset.Name] = cryptoSum(asset.Path)
			asset.CopyThis = true
			for _, image := range as.Meta.Images {
				if image.FileSize() == fileSize {
					// TODO: we should look at comparing with sha256 value too
					// on the billion-to-one chance we get two assets with an
					// identical size
					asset.Copyright = image.Copyright()
					asset.Description = image.Description()
				}
			}
			as.ByName[asset.Name] = asset
		}
	}
}

// wpExtension called by browse()
func wpExtension(assetPath string, width, height int) string {
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

	if image.Width != width || image.Height != height {
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
