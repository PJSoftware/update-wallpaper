package wallpaper

import (
	"io"
	"os"

	"github.com/pjsoftware/update-wallpaper/pkg/errors"
)

func Copy(sourcePath, targetPath string) (int64, error) {
	file, err := os.Stat(sourcePath)
	if err != nil {
		return 0, &errors.E{Code: errors.EFileNotFound, Message: "Source file not found"}
	}
	srcMTime := file.ModTime()

	source, err := os.Open(sourcePath)
	if err != nil {
		return 0, &errors.E{Code: errors.EReadError, Message: "Could not read source file"}
	}
	defer source.Close()

	dest, err := os.Create(targetPath)
	if err != nil {
		return 0, &errors.E{Code: errors.EWriteError, Message: "Could not create target file"}
	}

	numBytes, _ := io.Copy(dest, source)
	dest.Close()

	err = os.Chtimes(targetPath, srcMTime, srcMTime)

	return numBytes, err
}
