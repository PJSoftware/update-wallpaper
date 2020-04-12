#!bash

cd cmd/UpdateSpotlight
go build
ls -l *.exe
cd ../..

cp cmd/UpdateSpotlight/UpdateSpotlight.exe /c/Wallpaper
