package app

import "fmt"

// SplashScreen displays our "splash screen" with version & copyright info
func (a *App) SplashScreen() {
	fmt.Printf("Windows Wallpaper Toolset\n")
	fmt.Printf("%s (v%s)\n", a.name, a.version)
	fmt.Printf("Copyright © 2022 by PJSoftware\n\n")
}
