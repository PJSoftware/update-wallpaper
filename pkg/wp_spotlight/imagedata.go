package wp_spotlight

// ImageData is the information extracted from the JSON files
type ImageData struct {
	fileSize      int64
	width, height int64
	sha256        string
	url           string
	entityID      string
	copyright     string
	description   string
	metadataSrc   string
}

// FileSize returns the file size of the image
func (id *ImageData) FileSize() int64 {
	return id.fileSize
}

// Description returns the image description
func (id *ImageData) Description() string {
	return id.description
}

// Copyright returns the image copyright information
func (id *ImageData) Copyright() string {
	return id.copyright
}
