@echo off

echo Spotlight update
UpdateWallpaper

echo Renumber
"C:\Program Files\Git\usr\bin\perl.exe" renumber.pl

echo Bing
"C:\Program Files\Git\usr\bin\perl.exe" update_bing.pl

pause
