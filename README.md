# Win-Spotlight

Automate copying of Windows Spotlight images into Wallpaper Folder

## UpdateSpotlight

The task is simple:

* Look in the predefined Assets folder.
* Examine all files therein.
* Identify which ones are actually JPGs.
* Filter to only include the wallpaper-sized photos.
* Compare these with the ones already in the Wallpaper folder.
* Copy across the ones we don't already have.

**UpdateSpotlight.exe** reads its configuration from the **UpdateSpotlight.ini** file which looks as follows:

```ini
    [Wallpaper]

    # The resolution of image that we are looking for
    ImageWidth=1920
    ImageHeight=1080

    # This is where the Spotlight images should be delivered
    DestinationFolder="C:\Wallpaper"

    [Prefix]
    # Spotlight images copied into above folder wil have Prefix
    # added to their name based on the following values:

    # If prefix is not specified, or blank, no prefix will be used
    Prefix="ZZZ-"

    # If SmartPrefix is True, the prefix will only be applied
    # if no name is found in metadata. If SmartPrefix is False
    # or undefined, Prefix will be applied to all copied images.
    SmartPrefix=True

    [Archive]
    # Archive subfolder of DestinationFolder, contains old images
    Archive="_Archive"

    # Method can be either 'Delete' or 'SVN-Rename'. For most people
    # it will be enough to delete the duplicates from the archive folder
    # (if they even have one) but I keep my wallpapers in SVN, so it
    # makes more sense to svn-rename ratther than delete!
    Method="Delete"
```

If the INI file is not found by the program, it will default to using the above values.

The program was developed and tested on two computers, one with a screen resolution of 1920x1080, one with 2560x1080. In both cases, Spotlight delivered 1920x1080 images (and 1080x1920 Portrait variants.) It is possible that your system may be receiving different resolution images, in which case you will need to modify the **ImageWidth** and **ImageHeight** values to match your requirements.

**DestinationFolder** determines where the Spotlight assets, renamed to JPG (or PNG) files, should be placed. **UpdateSpotlight** does not merely look at filenames when determining whether an incoming Spotlight image already exists, so it is safe to rename them if required.

Any new Spotlight images will have the **Prefix** added to their filename to simplify any renaming you might wish to do.

If **SmartPrefix** is True, **Prefix** will only be used if the program could not find Spotlight metadata for the image. (The Spotlight delivery folders may contain a dozen or more images, but will likely only contain metadata for the last three or four images that were delivered to your computer.)

The `Archive` section is used by `DeleteDuplicates`.

## DeleteDuplicates

This can be used to compare your active Wallpaper folder with a non-active Archive folder. Any files in the Archive folder which are identical to a file in the Wallpaper folder will be deleted.

The `Archive` parameters from the INI file work as follows:

**Archive** specifies the name of the archive folder. It is treated as a subfolder of the **DestinationFolder**.

**Method** defaults to "Delete". If it is set to "SVN-Rename" `DeleteDuplicates` will, rather than delete the duplicate, rename it using `svn` to the new name.

## Note on Versioning

As I learn better ways to handle releases (not to mention better ways to use git branches) I am revising my approach. For Win-Spotlight, this means resetting my version to a "Toolset" version, and discarding any earlier version numbers that might have been in use. Apologies for any confusion. I shall do better going forward.
