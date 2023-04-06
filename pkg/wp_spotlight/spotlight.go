package wp_spotlight

import (
	"fmt"
)

type spotlight struct {
	assets *assets
}

func Update(folder string) {
	fmt.Printf("Updating SPOTLIGHT images:\n")
	wp := newSpotlight()
	wp.assets = readAssets(folder)
	found := wp.assets.count()

	/////////////////////////////////////////////////////////////////////////////
	// Code above works; code below not so much
	/////////////////////////////////////////////////////////////////////////////

	total, _ := wp.assets.compareWithExisting()
	copied, replaced := wp.assets.Copy()

	fmt.Printf("* %d new images copied", copied)
	if replaced > 0 {
		fmt.Printf("; %d existing images replaced", replaced)
	}
	fmt.Println()

	fmt.Printf("* Existing: %d; Incoming: %d; New: %d; Replaced: %d\n", total, found, copied, replaced)
}

func newSpotlight() *spotlight {
	return new(spotlight)
}
