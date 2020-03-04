package spotlight

import (
	"log"

	"../ini"
)

// Config provides interface to values from ini file
type Config struct {
	Width, Height int
	TargetPath    string
	SourcePath    string
	Prefix        string
	SmartPrefix   bool
	iniFile       ini.File
}

// Init sets values to those from ini file, or to defaults if an error occurs
func (s *Config) Init(exePath string) {
	err := s.iniFile.Parse(exePath + "UpdateSpotlight.ini")
	if err != nil {
		log.Print("world.Init: Error reading INI file: " + err.Error())
		log.Print("world.Init: using Default parameters instead")
		return
	}

	sectWallpaper := s.iniFile.Section("Wallpaper")
	s.Width = sectWallpaper.Value("ImageWidth").AsInt(1920)
	s.Height = sectWallpaper.Value("ImageHeight").AsInt(1080)
	s.TargetPath = sectWallpaper.Value("DestinationFolder").AsString(`C:\Wallpaper`, false)

	sectPrefix := s.iniFile.Section("Prefix")
	s.Prefix = sectPrefix.Value("Prefix").AsString("ZZZ_Unsorted-", false)
	s.SmartPrefix = sectPrefix.Value("SmartPrefix").AsBool(true)

	// SpotlightContentFolder should only be specified in testing
	contentRoot := sectWallpaper.ValueOptional("SpotlightContentFolder").AsString(GetPaths().ContentRoot(), false)
	GetPaths().SetContentRoot(contentRoot)
	s.SourcePath = GetPaths().Assets()
}
