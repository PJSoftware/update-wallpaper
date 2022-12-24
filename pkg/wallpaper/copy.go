package wallpaper

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/pjsoftware/update-wallpaper/pkg/wperr"
)

func Copy(sourcePath, targetPath string) error {
	if fileExists(targetPath) {
		fmt.Printf("  Skipping existing %s\n", targetPath)
		return nil
	}

	fmt.Printf("  Copying to %s\n", targetPath)
	file, err := os.Stat(sourcePath)
	if err != nil {
		return &wperr.E{Code: wperr.EFileNotFound, Message: "Source file not found"}
	}
	srcMTime := file.ModTime()

	source, err := os.Open(sourcePath)
	if err != nil {
		return &wperr.E{Code: wperr.EReadError, Message: "Could not read source file"}
	}
	defer source.Close()

	dest, err := os.Create(targetPath)
	if err != nil {
		return &wperr.E{Code: wperr.EWriteError, Message: "Could not create target file"}
	}
	defer dest.Close()

	_, err = io.Copy(dest, source)
	if err != nil {
		return &wperr.E{Code: wperr.ECopyError, Message: "Could not copy file"}
	}

	err = os.Chtimes(targetPath, srcMTime, srcMTime)
	return err
}

func fileExists(path string) bool {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return false
	}

	return true
}