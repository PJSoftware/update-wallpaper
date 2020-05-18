# Win-Spotlight

Automate copying of Windows Spotlight images into Wallpaper Folder

## UpdateSpotlight

The task is simple:

- Look in the predefined Assets folder.
- Examine all files therein.
- Identify which ones are actually JPGs.
- Filter to only include the wallpaper-sized photos.
- Compare these with the ones already in the Wallpaper folder.
- Copy across the ones we don't already have.

**UpdateSpotlight.exe** reads its configuration from the **Win-Spotlight.ini** file which looks as follows:

```ini
    [Wallpaper]

    # The resolution of image that we are looking for
    ImageWidth=1920
    ImageHeight=1080

    # This is where the Spotlight images should be delivered
    DestinationFolder="C:\Wallpaper"
```

If the INI file is not found by the program, it will default to using the above values.

The program was developed and tested on two computers, one with a screen resolution of 1920x1080, one with 2560x1080. In both cases, Spotlight delivered 1920x1080 images (and 1080x1920 Portrait variants.) It is possible that your system may be receiving different resolution images, in which case you will need to modify the **ImageWidth** and **ImageHeight** values to match your requirements.

**DestinationFolder** determines where the Spotlight assets, renamed to JPG (or PNG) files, should be placed. **UpdateSpotlight** does not merely look at filenames when determining whether an incoming Spotlight image already exists, so it is safe to rename them if required.
