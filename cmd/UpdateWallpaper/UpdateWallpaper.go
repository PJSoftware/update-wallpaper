package main

import (
	"github.com/pjsoftware/update-wallpaper/pkg/app"
	"github.com/pjsoftware/update-wallpaper/pkg/ui"
	"github.com/pjsoftware/update-wallpaper/pkg/wp_bing"
	"github.com/pjsoftware/update-wallpaper/pkg/wp_spotlight"
)

func main() {
	app := app.NewApp(TITLE, VERSION)
	app.InitFolder()
	app.SplashScreen()

	wp_spotlight.Update(app.WallpaperFolder("Spotlight"))
	wp_bing.Update(app.WallpaperFolder("Bing"))
	ui.Pause()
}
