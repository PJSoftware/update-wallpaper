package spotlight

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
)

// INI object provides interface to an ini file
type INI struct {
	section  map[string]bool
	value    map[string]map[string]string
	fileName string
}

// Parse reads an ini file and populates the INI object
func (i *INI) Parse(fileName string) error {
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	i.section = make(map[string]bool)
	i.value = make(map[string]map[string]string)

	i.fileName = fileName
	patternSect := `^[[](.+)[]]`
	patternValue := `^(\S+)\s*=\s*(\S.*)$`
	currSect := ""
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		isSect, err := regexp.MatchString(patternSect, line)
		if err != nil {
			return err
		}
		if isSect {
			reSect, err := regexp.Compile(patternSect)
			if err == nil {
				currSect = reSect.FindStringSubmatch(line)[1]
				i.section[currSect] = true
				i.value[currSect] = make(map[string]string)
			} else {
				log.Panic("INI: Could not compile regexp to read Section")
			}
		} else {
			isValue, err := regexp.MatchString(patternValue, line)
			if err != nil {
				return err
			}
			if isValue {
				reVal, err := regexp.Compile(patternValue)
				if err == nil {
					key := reVal.FindStringSubmatch(line)[1]
					value := reVal.FindStringSubmatch(line)[2]
					i.value[currSect][key] = value
				} else {
					log.Panic("INI: Could not compile regexp to read Value")
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

// Section determines whether the named section exists in the ini file
func (i *INI) Section(section string) (bool, error) {
	if i.section[section] {
		return true, nil
	}

	return false, fmt.Errorf("INI: Section '%s' not found in '%s'", section, i.fileName)
}

// Value returns value of named key in named section, as a string
func (i *INI) Value(section string, key string) (string, error) {
	_, err := i.Section(section)
	if err != nil {
		return "", err
	}

	if val, ok := i.value[section][key]; ok {
		return val, nil
	}

	return "", fmt.Errorf("INI: Value '%s' not found in section '%s' of '%s'", key, section, i.fileName)
}

// IntValue returns value of named key in named section, as an int
func (i *INI) IntValue(section string, key string) (int, error) {
	_, err := i.Section(section)
	if err != nil {
		return 0, err
	}

	if sval, ok := i.value[section][key]; ok {
		ival, err := strconv.Atoi(sval)
		if err != nil {
			return 0, fmt.Errorf("INI: IntValue '%s' has value '%s'; error converting to int", key, sval)
		}
		return ival, nil

	}

	return 0, fmt.Errorf("INI: IntValue '%s' not found in section '%s' of '%s'", key, section, i.fileName)
}
