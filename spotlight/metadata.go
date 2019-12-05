package spotlight

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

// MetaData is the interface/container for ImageData entries
type MetaData struct {
	size       int
	currentIdx int
	Images     []ImageData
}

// ImportAll is the entrypoint to all MetaData; it reads all relevant files
func (m *MetaData) ImportAll() {
	path := GetPaths()

	err := filepath.Walk(path.Metadata(),
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.Mode().IsRegular() {
				match, _ := regexp.MatchString(`[\\/]\d+[\\/]\d+$`, path) // vs "[\\\\/]\\d+[\\\\/]\\d+$"
				if match {
					for _, jsonItem := range jsonItemsFromFile(path) {
						m.parseJSON(jsonItem, path)
					}
				}
			}
			// data, err := readJSON(info.filepath)

			return nil
		})

	if err != nil {
		log.Println(err)
	}

}

// jsonItemsFromFile returns a slice of relevant Spotlight JSON items if found in file
func jsonItemsFromFile(fileName string) []map[string]interface{} {
	var rv []map[string]interface{}

	rawData, err := ReadUTF16(fileName)
	if err != nil {
		return nil // Just skip file if we have trouble reading it
	}

	if !json.Valid(rawData) {
		return nil
	}

	var jsonData map[string]interface{}
	json.Unmarshal(rawData, &jsonData)

	if _, ok := jsonData["batchrsp"]; !ok {
		return nil
	}
	batchRSP := jsonData["batchrsp"].(map[string]interface{})

	if _, ok := batchRSP["items"]; !ok {
		return nil
	}
	items := batchRSP["items"].([]interface{})

	for _, obj := range items {
		itemObj := obj.(map[string]interface{})
		if _, ok := itemObj["item"]; !ok {
			continue
		}
		itemStr := itemObj["item"]
		itemBytes := []byte(itemStr.(string))

		var itemMap map[string]interface{}

		if !json.Valid(itemBytes) {
			continue
		}

		err := json.Unmarshal(itemBytes, &itemMap)
		if err != nil {
			fmt.Println(err)
			continue
		}

		if _, ok := itemMap["ad"]; !ok {
			continue
		}
		adObj := itemMap["ad"].(map[string]interface{})

		if _, ok := adObj["name"]; !ok {
			continue
		}
		if adObj["name"] != "LockScreen" {
			continue
		}

		rv = append(rv, adObj)
	}

	return rv
}

// read
func (m *MetaData) parseJSON(data map[string]interface{}, src string) {
	if _, ok := data["properties"]; !ok {
		return
	}
	if _, ok := data["items"]; !ok {
		return
	}

	var image ImageData

	pOK := parseProperties(data["properties"].(map[string]interface{}), &image)

	item := data["items"].([]interface{})
	iOK := parseItems(item[0].(map[string]interface{}), &image)

	if pOK && iOK {
		image.metadataSrc = src

		m.size++
		m.currentIdx = m.size - 1
		m.Images = append(m.Images, image)
	}
}

func parseProperties(prop map[string]interface{}, image *ImageData) bool {
	if _, ok := prop["landscapeImage"]; !ok {
		return false
	}
	landscape := prop["landscapeImage"].(map[string]interface{})

	image.fileSize = int64(float64From(landscape, "fileSize"))
	image.height = int64(float64From(landscape, "height"))
	image.width = int64(float64From(landscape, "width"))
	image.sha256 = stringFrom(landscape, "sha256")
	image.url = stringFrom(landscape, "image")

	return image.fileSize > 0
}

func parseItems(item map[string]interface{}, image *ImageData) bool {
	var ok bool

	var prop, copyright, desc map[string]interface{}
	if prop, ok = item["properties"].(map[string]interface{}); !ok {
		return false
	}
	if copyright, ok = prop["copyright"].(map[string]interface{}); !ok {
		return false
	}
	if desc, ok = prop["description"].(map[string]interface{}); !ok {
		return false
	}

	image.entityID = stringFrom(item, "entityId")
	image.copyright = stringFrom(copyright, "text")
	image.description = stringFrom(desc, "text")

	return image.description != ""
}

func float64From(dat map[string]interface{}, key string) float64 {
	rv, ok := dat[key]
	if ok {
		return rv.(float64)
	}
	return 0.0
}

func stringFrom(dat map[string]interface{}, key string) string {
	rv, ok := dat[key]
	if ok {
		return rv.(string)
	}
	return ""
}
