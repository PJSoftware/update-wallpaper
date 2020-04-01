package ini

import (
	"bufio"
	"log"
	"os"
	"regexp"

	"../errors"
)

const unnamedSection = "**PARENT**"
const noSectionNames = "**FLAT**"

// File object provides interface to an ini file
type File struct {
	sections   map[string]*Section
	fileName   string
	canFlatten bool
}

// Parse reads an ini file and creates Section/Value objects as required
func (f *File) Parse(fileName string) error {

	// Compile our regexp strings first; failure here is a typing error
	// and therefore panic on failure is justified
	reSect := regexp.MustCompile(`^[[](.+)[]]`)
	reValue := regexp.MustCompile(`^(\S+)=(.*)$`)
	reExtract := regexp.MustCompile(`^"(.*)"$|^'(.*)'$`)

	file, err := os.Open(fileName)
	if err != nil {
		return errors.E{Code: errors.EFILENOTFOUND}
	}
	defer file.Close()

	f.sections = make(map[string]*Section)
	f.fileName = fileName

	var currSect *Section
	currSectName := unnamedSection
	f.canFlatten = true
	valuesRead := make(map[string]bool)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line[0:0] == "#" {
			continue
		}

		sectFound := reSect.FindStringSubmatch(line)
		if sectFound != nil {
			currSectName = sectFound[1]
			currSect = f.addSection(currSectName)
			continue
		}

		valFound := reValue.FindStringSubmatch(line)
		if valFound != nil {
			if currSect == nil {
				currSect = f.addSection(currSectName)
			}
			key, value := valFound[1], valFound[2]

			strFound := reExtract.FindStringSubmatch(value)
			if strFound != nil {
				value = strFound[1]
			}

			if _, ok := valuesRead[key]; ok {
				f.canFlatten = false
			} else {
				valuesRead[key] = true
			}

			currSect.addValue(key, value)
		}
	}

	if err := scanner.Err(); err != nil {
		return errors.E{Code: errors.EREADERROR}
	}

	return nil
}

func (f *File) addSection(sectName string) *Section {
	sect := newSection(sectName)
	f.sections[sectName] = sect
	return sect
}

// Section returns named Section object
func (f *File) Section(sectName string) *Section {
	if sect, ok := f.sections[sectName]; ok {
		return sect
	}
	log.Printf("INI.Section: section '%s' not found", sectName)
	return nil
}

// Sections returns slice of section names
func (f *File) Sections() []string {
	keys := make([]string, len(f.sections))

	i := 0
	for key := range f.sections {
		keys[i] = key
		i++
	}
	return keys
}
