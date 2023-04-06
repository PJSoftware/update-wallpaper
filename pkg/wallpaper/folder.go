package wallpaper

import (
	"log"
	"os"
)

// Folder contains a collection of files
type Folder struct {
	path   string
	files  []*File
	bySize map[int64][]*File
}

// File holds the data for each file
type File struct {
	name string
	size int64
	// hash string
}

// ImportFolder is our constructor method
func ImportFolder(fPath string) *Folder {
	f := new(Folder)
	f.path = fPath
	f.bySize = make(map[int64][]*File)

	files, err := os.ReadDir(f.path)
	if err != nil {
		log.Fatalf("ImportFolder error reading %s: %v", f.path, err)
	}

	for _, file := range files {
		fs := new(File)
		fs.name = file.Name()
		info, _ := file.Info()
		fs.size = info.Size()
		f.files = append(f.files, fs)
		f.bySize[fs.size] = append(f.bySize[fs.size], fs)
	}

	return f
}
