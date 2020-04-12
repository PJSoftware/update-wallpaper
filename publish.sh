#!bash

cd cmd/UpdateSpotlight
go build
ls -l *.exe
cd ../..

cd cmd/DeleteDuplicates
go build
ls -l *.exe
cd ../..

cp cmd/UpdateSpotlight/UpdateSpotlight.exe /c/Wallpaper
cp cmd/DeleteDuplicates/DeleteDuplicates.exe /c/Wallpaper
