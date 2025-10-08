#! /bin/sh

./krankybearnotify -title "Image Test" -message "Testing the -image alias flag with a long message that should wrap nicely" -image Resources/Images/KrankyBearBeret.png -timeout 5

./krankybearnotify -title "Test" -message "This is a very long message that should wrap properly instead of being truncated to match the title width. The text should flow nicely across multiple lines with proper word wrapping." -timeout 0

./krankybearnotify -title "Test" -message "This is a very long message that should wrap properly instead of being truncated to match the title width. The text should flow nicely across multiple lines with proper word wrapping." -timeout 0 -image KrankyBearBeret.png

./krankybearnotify -title "Test" -message "This is a very long message that should wrap properly instead of being truncated to match the title width. The text should flow nicely across multiple lines with proper word wrapping. This one does include an image and has height 400 and width 800." -timeout 0 -image KrankyBearBeret.png -height 400 -width 800