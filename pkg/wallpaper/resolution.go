package wallpaper

import (
	"image/jpeg"
	"os"
)

type ImageResolution struct {
	Width int
	Height int
	Name string
}

func Resolution(path string) (*ImageResolution, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	imageData, err := jpeg.DecodeConfig(file)
	if err != nil {
		return nil, err
	}

	res := &ImageResolution{}
	res.Width = imageData.Width
	res.Height = imageData.Height

	if res.matches(1920, 1080) {
		res.Name = "HD"
	} else if res.matches(2560, 1080) {
		res.Name = "UW"
	} else if res.matches(3840, 2160) {
		res.Name = "4K"
	} else if res.matches(5120, 2160) {
		res.Name = "5K"
	} else {
		res.Name = "Other"
	}

	return res, nil
}

func (ir *ImageResolution) matches(width, height int) bool {
	return ir.Width == width && ir.Height == height
}
