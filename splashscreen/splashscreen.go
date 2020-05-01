package splashscreen

import "fmt"

// Show displays our splash "screen" with version & copyright info
func Show(cmd string, oldver string) {
	fmt.Printf("Windows Spotlight Toolset v%s\n", version)
	fmt.Printf("%s: Copyright © 2020 by PJSoftware\n", cmd)
	fmt.Printf("(Was '%s v%s')\n", cmd, oldver)
	fmt.Println()
}
