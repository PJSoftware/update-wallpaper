package wp_spotlight

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"path/filepath"
	"regexp"
)

// imageMetadata is the information extracted from the JSON files
type imageMetadata struct {
	fileSize      int64
	width, height int64
	sha256        string
	url           string
	entityID      string
	copyright     string
	description   string
	metadataSrc   string
}

// assetMetadata is the interface/container for imageMetadata entries
type assetMetadata struct {
	imageMD []imageMetadata
}

const NO_DESCRIPTION string = "Unidentified Photo"
const NO_COPYRIGHT string = "© Unknown Photographer"

// read is the entrypoint to all asset metadata; it reads all relevant files
func (m *assetMetadata) read() {
	mdPath := spotlightMetadataFolder
	err := filepath.WalkDir(mdPath,
		func(path string, info fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			
			fi, _ := info.Info()
			if fi.Mode().IsRegular() {
				match, _ := regexp.MatchString(`[\\/]\d+[\\/]\d+$`, path) // vs "[\\\\/]\\d+[\\\\/]\\d+$"
				if match {
					for _, jsonItem := range locateSpotlightMetadata(path) {
						m.parseMetadata(jsonItem, path)
					}
				}
			}

			return nil
		})

	if err != nil {
		log.Fatalf("error in ImportAll: %v\n", err)
	}
}

// locateSpotlightMetadata returns a slice of relevant Spotlight JSON items if found in file
func locateSpotlightMetadata(fileName string) []map[string]interface{} {
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
func (m *assetMetadata) parseMetadata(data map[string]interface{}, src string) {
	if _, ok := data["properties"]; !ok {
		return
	}
	if _, ok := data["items"]; !ok {
		return
	}

	var metadata imageMetadata
	pOK := parseProperties(data["properties"].(map[string]interface{}), &metadata)

	if pOK {
		parseItems(data["items"].([]interface{}), &metadata)
		metadata.metadataSrc = src
		m.imageMD = append(m.imageMD, metadata)
	}
}

func parseProperties(prop map[string]interface{}, metadata *imageMetadata) bool {
	if _, ok := prop["landscapeImage"]; !ok {
		return false
	}
	landscape := prop["landscapeImage"].(map[string]interface{})

	metadata.fileSize = int64(float64From(landscape, "fileSize"))
	metadata.height = int64(float64From(landscape, "height"))
	metadata.width = int64(float64From(landscape, "width"))
	metadata.sha256 = stringFrom(landscape, "sha256")
	metadata.url = stringFrom(landscape, "image")

	return metadata.fileSize > 0
}

func parseItems(item []interface{}, metadata *imageMetadata) bool {
	var ok bool

	if len(item) == 0 {
		metadata.entityID = "UNKNOWN"
		metadata.description = NO_DESCRIPTION
		metadata.copyright = NO_COPYRIGHT
		return false
	}

	item1 := item[0].(map[string]interface{})

	var prop, copyright, desc map[string]interface{}
	if prop, ok = item1["properties"].(map[string]interface{}); !ok {
		return false
	}
	if copyright, ok = prop["copyright"].(map[string]interface{}); !ok {
		return false
	}
	if desc, ok = prop["description"].(map[string]interface{}); !ok {
		return false
	}

	metadata.entityID = stringFrom(item1, "entityId")
	metadata.copyright = stringFrom(copyright, "text")
	metadata.description = stringFrom(desc, "text")

	return metadata.description != ""
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
