package paths

import (
	"os"
	"sync"
)

// SpotlightPaths contains our path data and should be instantiated via GetPaths()
type SpotlightPaths struct {
	root     string
	assets   string
	metadata string
}

const spotlightFolder = "Packages/Microsoft.Windows.ContentDeliveryManager_cw5n1h2txyewy"

// Implementation of Singleton via http://marcio.io/2015/07/singleton-pattern-in-go/
// As presented, results in lint warning #210; exporting Paths prevents this but loses
// the guarantee that the struct cannot be used before initialising
// TODO: look into ways to resolve this
var instance *SpotlightPaths
var once sync.Once

// GetSpotlightPaths returns our singleton instance of the Paths struct
func GetSpotlightPaths() *SpotlightPaths {
	once.Do(func() {
		instance = &SpotlightPaths{}

		local := os.Getenv("LOCALAPPDATA")
		instance.root = local + "/" + spotlightFolder
		instance.assets = "LocalState/Assets"
		instance.metadata = "LocalState/ContentManagementSDK/Creatives"
	})
	return instance
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
