package spotlight

import (
	"os"
	"sync"
)

type paths struct {
	assets   string
	metadata string
}

const spotlight = "Packages/Microsoft.Windows.ContentDeliveryManager_cw5n1h2txyewy"

// Implementation of Singleton via http://marcio.io/2015/07/singleton-pattern-in-go/
// As presented, results in lint warning #210; exporting Paths prevents this but loses
// the guarantee that the struct cannot be used before initialising
var instance *paths
var once sync.Once

// GetInstance returns our singleton instance of the paths struct
func GetInstance() *paths {
	once.Do(func() {
		instance = &paths{}

		local := os.Getenv("LOCALAPPDATA")
		instance.assets = local + "/" + spotlight + "/LocalState/Assets"
		instance.metadata = local + "/" + spotlight + "/LocalState/ContentManagementSDK/Creatives"
	})
	return instance
}

func (p *paths) Assets() string {
	return p.assets
}

func (p *paths) Metadata() string {
	return p.metadata
}
