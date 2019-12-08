# Win-Spotlight

* INI does not correctly support blank/null values (ie, "Prefix =" or "Prefix = ''")
* Code currently assumes that output file does not exist. Perhaps there was a good chance of this when the filename was the same as the asset file, but now that we are renaming them per the metadata, there is actually a greater chance of a clash. This needs to be looked into. If this check is performed at the point of copying, we already know that if a filename matches, it contains something different than our asset, because only unique files will be copied.
* Need to use SHA256 value from metadata to confirm metadata matches asset file.
