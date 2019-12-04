package spotlight

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

// MetaData is the interface/container for ImageData entries
type MetaData struct {
	size       int
	currentIdx int
	image      []ImageData
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
					for _, jsonItem := range jsonItems(path) {
						m.parseJSON(jsonItem)
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

// readJSON attempts to open specified file, returns contents if valid json
func jsonItems(fileName string) []map[string]interface{} {
	fmt.Printf("Extracting JSON items from %s\n", fileName)
	var rv []map[string]interface{}

	rawData, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil // Just skip file if we have trouble reading it
	}
	fmt.Printf("* successfully read file: %d bytes\n", len(rawData))

	if !json.Valid(rawData) {
		return nil
	}
	fmt.Println("* file is valid JSON")

	var jsonData map[string]interface{}
	json.Unmarshal(rawData, &jsonData)

	if _, ok := jsonData["batchrsp"]; !ok {
		return nil
	}
	batchRSP := jsonData["batchrsp"].(map[string]interface{})
	fmt.Println("* batchRSP")

	if _, ok := batchRSP["items"]; !ok {
		return nil
	}
	items := batchRSP["items"].([]interface{})
	fmt.Println("* items")

	for _, obj := range items {
		itemObj := obj.(map[string]interface{})
		if _, ok := itemObj["item"]; !ok {
			return nil
		}
		itemStr := itemObj["item"]
		itemBytes := []byte(itemStr.(string))

		var itemMap map[string]interface{}
		fmt.Println(" > itemMap")
		json.Unmarshal(itemBytes, &itemMap)

		if _, ok := itemMap["ad"]; !ok {
			return nil
		}
		adObj := itemMap["ad"].(map[string]interface{})

		if _, ok := adObj["name"]; !ok {
			return nil
		}
		if adObj["name"] != "LockScreen" {
			return nil
		}

		fmt.Printf("  item found!\n")
		rv = append(rv, adObj)
	}

	return rv
}

// read
func (m *MetaData) parseJSON(data map[string]interface{}) {
	if _, ok := data["properties"]; !ok {
		return
	}
	if _, ok := data["items"]; !ok {
		return
	}

	var image ImageData

	pOK := parseProperties(data["properties"].(map[string]interface{}), &image)
	iOK := parseItems(data["items"].(map[string]interface{}), &image)

	if pOK && iOK {
		m.size++
		m.currentIdx = m.size - 1
		m.image = append(m.image, image)

		fmt.Printf("%s: %s x %s; %s\n", image.url, image.width, image.height, image.fileSize)
	}
}

func parseProperties(prop map[string]interface{}, image *ImageData) bool {
	if _, ok := prop["landscapeImage"]; !ok {
		return false
	}
	landscape := prop["landscapeImage"].(map[string]interface{})

	var ok bool
	if image.fileSize, ok = landscape["fileSize"].(string); !ok {
		return false
	}
	if image.height, ok = landscape["height"].(string); !ok {
		return false
	}
	if image.width, ok = landscape["width"].(string); !ok {
		return false
	}
	if image.sha256, ok = landscape["sha256"].(string); !ok {
		return false
	}
	if image.url, ok = landscape["image"].(string); !ok {
		return false
	}

	return true
}

func parseItems(item map[string]interface{}, image *ImageData) bool {
	var ok bool
	if image.entityID, ok = item["entityID"].(string); !ok {
		return false
	}

	var prop, copyright, desc map[string]interface{}

	if prop, ok = item["properties"].(map[string]interface{}); !ok {
		return false
	}

	if copyright, ok = prop["copyright"].(map[string]interface{}); !ok {
		return false
	}
	if image.copyright, ok = copyright["text"].(string); !ok {
		return false
	}

	if desc, ok = prop["description"].(map[string]interface{}); !ok {
		return false
	}
	if image.description, ok = desc["text"].(string); !ok {
		return false
	}

	return true
}
