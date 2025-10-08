Output screenshots/demo.gif

Require fsize

Set Shell bash
Set FontSize 32
Set Width 1300
Set Height 800

Env LANG en

Type "# Let's try fsize on its own project directory..."
Enter
Type "fsize"
Sleep 500ms
Enter
Wait
Screenshot screenshots/demo.png
Sleep 4s
Type "clear"
Enter
Type "# What if we run fsize on a more huge directory?"
Enter
Type "fsize /usr"
Sleep 500ms
Enter
Wait+Screen /GID.*/
Sleep 5s
