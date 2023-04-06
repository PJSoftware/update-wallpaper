package wallpaper

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/pjsoftware/update-wallpaper/pkg/wperr"
)

func Copy(sourcePath, targetPath string) (bool, error) {
	if PathExists(targetPath) {
		return false, nil
	}

	fmt.Printf("  Copying to %s\n", targetPath)
	file, err := os.Stat(sourcePath)
	if err != nil {
		return false, &wperr.E{Code: wperr.EFileNotFound, Message: "Source file not found"}
	}
	srcMTime := file.ModTime()

	source, err := os.Open(sourcePath)
	if err != nil {
		return false, &wperr.E{Code: wperr.EReadError, Message: "Could not read source file"}
	}
	defer source.Close()

	dest, err := os.Create(targetPath)
	if err != nil {
		return false, &wperr.E{Code: wperr.EWriteError, Message: "Could not create target file"}
	}
	defer dest.Close()

	_, err = io.Copy(dest, source)
	if err != nil {
		return false, &wperr.E{Code: wperr.ECopyError, Message: "Could not copy file"}
	}

	err = os.Chtimes(targetPath, srcMTime, srcMTime)
	return true, err
}

func PathExists(path string) bool {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return false
	}

	return true
}