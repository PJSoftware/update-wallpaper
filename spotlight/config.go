package spotlight

import (
	"log"

	"../ini"
)

// TODO: When reading INI file, we should log if no INI found, but not if a
//  particular value is missing. It is common for values to be commented out of
//  INI files when the default is to be used!

// Config provides interface to values from ini file
type Config struct {
	Width, Height int
	TargetPath    string
	Prefix        string
	iniFile       ini.File
}

// Init sets values to those from ini file, or to defaults if an error occurs
func (s *Config) Init(exePath string) {
	s.setDefaults()

	err := s.iniFile.Parse(exePath + "UpdateSpotlight.ini")
	if err != nil {
		log.Print("world.Init: Error reading INI file: " + err.Error())
		log.Print("world.Init: using Default parameters instead")
		return
	}

	spotlight := s.iniFile.Section("Spotlight")

	s.Width = spotlight.Value("ImageWidth").AsInt64(1920)

	_, err = s.iniFile.Section(section)
	if err != nil {
		log.Printf(err.Error())
		return
	}

	s.readWidth(section)
	s.readHeight(section)
	s.readPath(section)
	s.readPrefix(section)
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

func (s *Config) readPrefix(sectName string) {
	key := "Prefix"
	val, err := s.iniFile.Value(sectName, key)
	if err == nil {
		s.Prefix = val
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
	s.Prefix = ""
}
