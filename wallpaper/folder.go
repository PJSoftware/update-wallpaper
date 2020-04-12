package wallpaper

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/pjsoftware/win-spotlight/util"
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
	hash string
}

// ImportFolder is our constructor method
func ImportFolder(fPath string) *Folder {
	f := new(Folder)
	f.path = fPath
	f.bySize = make(map[int64][]*File)

	files, err := ioutil.ReadDir(f.path)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		fs := new(File)
		fs.name = file.Name()
		fs.size = file.Size()
		f.files = append(f.files, fs)
		f.bySize[fs.size] = append(f.bySize[fs.size], fs)
	}

	return f
}

// DeleteDuplicates deletes files from f which match fc
func (f *Folder) DeleteDuplicates(fc *Folder) {
	for size, files := range f.bySize {
		if cf, ok := fc.bySize[size]; ok {
			for _, tbd := range files {
				tbd.hash = util.FileHash(f.path + "/" + tbd.name)
			}
			for _, tbc := range cf {
				tbc.hash = util.FileHash(fc.path + "/" + tbc.name)
				for _, tbd := range files {
					if tbc.hash == tbd.hash {
						fmt.Printf("The following files are identical:\n   '%s'\n-> '%s'\nDeleting indicated copy!\n\n", tbc.name, tbd.name)
						os.Remove(f.path + "/" + tbd.name)
					}
				}
			}
		}
	}
}
