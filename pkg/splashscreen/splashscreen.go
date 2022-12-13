package splashscreen

import "fmt"

// Show displays our splash "screen" with version & copyright info
func Show(cmd string, version string) {
	fmt.Printf("Windows Wallpaper Toolset v%s\n", version)
	fmt.Printf("%s: Copyright © 2022 by PJSoftware\n", cmd)
	fmt.Println()
}
