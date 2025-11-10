#! /bin/sh

if [ $# -ge 1 ]
then
    ver=$1
else
    echo "Enter a version number"
    cur=$(cat main.go | grep -i "appVersion" | grep "=" | awk '{print $3}' | tr -d '"')
    echo "    current: $cur"
    read ver
    if [ -z "$ver" ]
    then
        echo "Enter a version!"
        echo "No version change detected, continuing to allow compile to continue"
        exit
    else
        echo "Version: $ver"
        # exit
    fi
fi

echo "version: $ver"
echo "main.go"
# Match appVersion with any whitespace (tabs/spaces) before =, preserve original spacing
sed -i '' "s/\(appVersion[[:space:]]*\)= \"[^\"]*\"/\1= \"$ver\"/" main.go

echo "FyneApp.toml"
sed -i '' "s/Version = \"[^\"]*\"/Version = \"$ver\"/" FyneApp.toml

echo "Inno Setup Inno/KrankyBearNotify.iss"
sed -i '' "s/MyAppVersion \"[^\"]*\"/MyAppVersion \"$ver\"/" ./Inno/KrankyBearNotify.iss

echo "Inno Setup winres/winres.json"
sed -i '' "s/\"file_version\": \"[^\"]*\"/\"file_version\": \"$ver\"/" ./winres/winres.json
sed -i '' "s/\"product_version\": \"[^\"]*\"/\"product_version\": \"$ver\"/" ./winres/winres.json
sed -i '' "s/\"FileVersion\": \"[^\"]*\"/\"FileVersion\": \"$ver\"/" ./winres/winres.json
sed -i '' "s/\"ProductVersion\": \"[^\"]*\"/\"ProductVersion\": \"$ver\"/" ./winres/winres.json

echo "No direct Info.plist updates - updating Info-plist.txt which can be renamed if wanted"
# Only update version-related keys in Info-plist.txt (key on one line, value on next)
sed -i '' "/CFBundleShortVersionString\|CFBundleVersion/{N;s/<string>[^<]*<\/string>/<string>$ver<\/string>/;}" ./Info-plist.txt

echo "Update LICENSE and ReleaseNotes.txt"
cp LICENSE Resources
cp ReleaseNotes.txt ./Resources
cp LICENSE KrankyBearNotify.app/Contents/Resources
cp ReleaseNotes.txt ./KrankyBearNotify.app/Contents/Resources
cp ./Info-plist.txt ./KrankyBearNotify.app/Contents/Info-plist.txt

echo "Update package.sh"
sed -i '' "s/VERSION:-[^}]*}/VERSION:-$ver}/" ./package.sh

# "Now this is not the end. It is not even the beginning of the end. But it is, perhaps, the end of the beginning." Winston Churchill, November 10, 1942
