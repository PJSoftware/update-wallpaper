package splashscreen

import "fmt"

// Show displays our splash "screen" with version & copyright info
func Show(cmd string) {
	fmt.Printf("Windows Spotlight Tools v%s\n", version)
	fmt.Printf("%s: Copyright © 2020 by PJSoftware\n\n", cmd)
}
