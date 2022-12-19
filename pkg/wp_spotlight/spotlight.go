package wp_spotlight

import (
	"fmt"
	"log"
)

type spotlight struct {
	assets *assets
}

func Update(folder string) {
	wp := newSpotlight()
	wp.assets = readAssets()

	found := wp.assets.Count()
	fmt.Printf("%d Spotlight images found\n", found)

	total, duplicates := wp.assets.Compare()
	fmt.Printf("%d Existing wallpapers found\n", found)
	fmt.Printf("%d Spotlight assets match existing; skipping\n", duplicates)

	copied, replaced := wp.assets.Copy()
	fmt.Printf("%d new images copied\n", copied)
	if replaced > 0 {
		fmt.Printf("%d existing images replaced\n", replaced)
	}
	log.Printf("Existing: %d; Incoming: %d; New: %d; Replaced: %d", total, found, copied, replaced)
}

func newSpotlight() *spotlight {
	return new(spotlight)
}
