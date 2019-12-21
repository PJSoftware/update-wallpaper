package spotlight

import (
	"os"
	"sync"
)

// Paths contains our path data and should be instantiated via GetPaths()
type Paths struct {
	root     string
	assets   string
	metadata string
}

const spotlight = "Packages/Microsoft.Windows.ContentDeliveryManager_cw5n1h2txyewy"

// Implementation of Singleton via http://marcio.io/2015/07/singleton-pattern-in-go/
// As presented, results in lint warning #210; exporting Paths prevents this but loses
// the guarantee that the struct cannot be used before initialising
var instance *Paths
var once sync.Once

// GetPaths returns our singleton instance of the Paths struct
func GetPaths() *Paths {
	once.Do(func() {
		instance = &Paths{}

		local := os.Getenv("LOCALAPPDATA")
		instance.root = local + "/" + spotlight
		instance.assets = "LocalState/Assets"
		instance.metadata = "LocalState/ContentManagementSDK/Creatives"
	})
	return instance
}

// ContentRoot returns the spotlight ContentDelivery root folder
func (p *Paths) ContentRoot() string {
	return p.root
}

// SetContentRoot allows us to use a different source folder tree
// primarily for debugging purposes because we do not usually need to do so
func (p *Paths) SetContentRoot(newRoot string) {
	p.root = newRoot
}

// Assets returns the path to the spotlight assets folder
func (p *Paths) Assets() string {
	return p.root + "/" + p.assets
}

// Metadata returns the path to the spotlight metadata parent folder
func (p *Paths) Metadata() string {
	return p.root + "/" + p.metadata
}
