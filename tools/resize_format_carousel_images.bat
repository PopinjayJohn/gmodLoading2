mogrify -format jpg ..\static\carousel\*.*
mogrify -resize 1366x768! ..\static\carousel\*.jpg
ECHO dont forget to rename the images to 1, 2 and 3
PAUSE