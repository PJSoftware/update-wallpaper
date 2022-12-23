package wp_spotlight

import (
	"os"
	"path/filepath"
)

const spotlightFOLDER = "Packages\\Microsoft.Windows.ContentDeliveryManager_cw5n1h2txyewy"

var spotlightRootFolder = ""
var spotlightAssetFolder = ""
var spotlightMetadataFolder = ""

func init() {
	local := os.Getenv("LOCALAPPDATA")
	spotlightRootFolder = filepath.Join(local, spotlightFOLDER)
  spotlightAssetFolder = filepath.Join(spotlightRootFolder, "LocalState\\Assets")
  spotlightMetadataFolder = filepath.Join(spotlightRootFolder, "LocalState\\ContentManagementSDK\\Creatives")
}
