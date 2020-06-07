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
- If an incoming file matches an existing image but the name is different, rename with the _new_ MetaData name.

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

### Version Control Support

`UpdateSpotlight` detects whether an incoming image is identical to an existing image even if the filename is different. This may indicate that the existing image was imported by an earlier version before Metadata-naming was added as an option. In such a case, the existing file is renamed to the new, "correct" filename.

If the wallpaper folder is under version control, this renaming will be performed via the VC software (where supported) because a `move` or `rename` of a file is typically much more efficient when the VC is aware the two are the same file!

Additionally, `UpdateSpotlight` will commit the new (or renamed) wallpapers if a VC is detected.

Currently supported Version Control software includes:

- --in development--
