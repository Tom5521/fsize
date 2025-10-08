Output screenshots/demo.gif

Require fsize

Set Shell fish
Set FontSize 32
Set Width 1300
Set Height 800

Env LANG en

Type "# Let's try fsize on its own project directory..." Enter
Type "fsize" Sleep 500ms Enter
Wait
Screenshot screenshots/demo.png
Sleep 4s
Type "clear" Enter
Type "# What if we run fsize on a more huge directory?" Sleep 500ms Enter
Type "fsize /usr" Sleep 1s Enter
Wait@40s
Sleep 5s
Type "# I want to skip some directories..." Sleep 500ms Enter
Type "fsize -i builds" Sleep 1s Enter
Sleep 3s
Type "# I think it would be better with a wildcard" Sleep 500ms Enter
Type "fsize -wi 'builds/*'" Sleep 1s Enter
Wait@20s
Sleep 5s
