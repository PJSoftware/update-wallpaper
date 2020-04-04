package spotlight

import (
	"log"

	"github.com/pjsoftware/win-spotlight/ini"
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
		log.Print("config.Init: Error reading INI file: " + err.Error())
		log.Print("config.Init: using Default parameters instead")
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
	crv := sectWallpaper.ValueOptional("SpotlightContentFolder")
	if crv != nil {
		contentRoot := crv.AsString(GetPaths().ContentRoot(), false)
		GetPaths().SetContentRoot(contentRoot)
	}
	s.SourcePath = GetPaths().Assets()
}
