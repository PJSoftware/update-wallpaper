package spotlight

// ImageData is the information extracted from the JSON files
type ImageData struct {
	fileSize      string
	width, height string
	sha256        string
	url           string
	entityID      string
	copyright     string
	description   string
}
