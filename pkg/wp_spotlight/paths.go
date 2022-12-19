package wp_spotlight

import (
	"os"
)

// SpotlightPaths contains our path data and should be instantiated via GetPaths()
type SpotlightPaths struct {
	root     string
	assets   string
	metadata string
}

const SPOTLIGHTFOLDER = "Packages/Microsoft.Windows.ContentDeliveryManager_cw5n1h2txyewy"

// GetSpotlightPaths returns our singleton instance of the Paths struct
func GetSpotlightPaths() *SpotlightPaths {
		paths := &SpotlightPaths{}

		local := os.Getenv("LOCALAPPDATA")
		paths.root = local + "/" + SPOTLIGHTFOLDER
		paths.assets = "LocalState/Assets"
		paths.metadata = "LocalState/ContentManagementSDK/Creatives"
	return paths
}

// ContentRoot returns the spotlight ContentDelivery root folder
func (p *SpotlightPaths) ContentRoot() string {
	return p.root
}

// SetContentRoot allows us to use a different source folder tree
// primarily for debugging purposes because we do not usually need to do so
func (p *SpotlightPaths) SetContentRoot(newRoot string) {
	p.root = newRoot
}

// Assets returns the path to the spotlight assets folder
func (p *SpotlightPaths) Assets() string {
	return p.root + "/" + p.assets
}

// Metadata returns the path to the spotlight metadata parent folder
func (p *SpotlightPaths) Metadata() string {
	return p.root + "/" + p.metadata
}
