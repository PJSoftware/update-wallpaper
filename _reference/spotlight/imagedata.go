package spotlight

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
