package spotlight

import (
	"log"
)

// Config provides interface to values from ini file
type Config struct {
	Width, Height int
	TargetPath    string
	iniFile       INI
}

// Init sets values to those from ini file, or to defaults if an error occurs
func (s *Config) Init() {
	s.setDefaults()
	err := s.iniFile.Parse("UpdateSpotlight.ini")
	if err != nil {
		log.Print("world.Init: Error reading INI file: " + err.Error())
		log.Print("world.Init: using Default parameters instead")
		return
	}

	section := "Spotlight"
	_, err = s.iniFile.Section(section)
	if err != nil {
		log.Printf(err.Error())
		return
	}

	s.readWidth(section)
	s.readHeight(section)
	s.readPath(section)
}

func (s *Config) readWidth(sectName string) {
	key := "ImageWidth"
	val, err := s.iniFile.IntValue(sectName, key)
	if err == nil {
		s.Width = val
	} else {
		log.Print(err.Error())
		log.Printf("Config: %s not found, using default", key)
	}
}

func (s *Config) readHeight(sectName string) {
	key := "ImageHeight"
	val, err := s.iniFile.IntValue(sectName, key)
	if err == nil {
		s.Height = val
	} else {
		log.Print(err.Error())
		log.Printf("Config: %s not found, using default", key)
	}
}

func (s *Config) readPath(sectName string) {
	key := "DestinationFolder"
	val, err := s.iniFile.Value(sectName, key)
	if err == nil {
		s.TargetPath = val
	} else {
		log.Print(err.Error())
		log.Printf("Config: %s not found, using default", key)
	}
}

func (s *Config) setDefaults() {
	// I believe Spotlight delivers 1920x1080 by default anyway
	s.Width = 1920
	s.Height = 1080
	s.TargetPath = "C:\\Wallpaper"
}
