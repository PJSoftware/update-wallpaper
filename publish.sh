#!bash

cd cmd/UpdateSpotlight
go build
ls -l *.exe
cd ../..

cp cmd/UpdateSpotlight/UpdateSpotlight.exe /c/Wallpaper

# We want a python script, I think:
#
# git checkout dev
#   use git describe to determine current tag
# git checkout -b release-1.3
#   update splashscreen/version.go with new version
# git commit -m "bump to version 1.3"
#
# git checkout master
# git merge --no-ff release-1.3
# git branch -d release-1.3
# git tag -a v1.3 -m "Version 1.3: description of change"
# git push
# git push --tags
# ./publish.sh
# cd /c/Wallpaper
# git add UpdateSpotlight.exe
# git commit -m "upgrade UpdateSpotlight to v1.3"
# git push
# cd -
# git checkout dev
# git merge --no-ff master
