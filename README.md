# Win-Spotlight
Automate copying of Windows Spotlight images into Wallpaper Folder

The task is simple:
* Look in the deeply nested Assets folder.
* Examine all files therein.
* Identify which ones are actually JPGs.
* Filter to only include the wallpaper-sized photos.
* Compare these with the ones already in the Wallpaper folder.
* Copy across the ones we don't already have.

## Lessons Learned
Based on my testing of this on a whole two machines, one at 1920x1080 and one at 2560x1080, I am concluding that the service only delivers 1920x1080 images. However, YMMV.

From reading I have done elsewhere, it will apparently match the images to your desktop resolution--but it seems it only looks at height.

Either way, it *does* seem to provide both Portrait and Landscape versions of the images.

Also, I have only ever found JPG images, but apparently it will also deliver PNG files. At the moment my code specifically only looks for JPGs. That may change.
