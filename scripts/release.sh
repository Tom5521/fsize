#!/usr/bin/env -S bash -x

latest_tag=$(git describe --tags --abbrev=0)

case "$1" in
"update-assets")
  ./do generate binaries
  gh release upload "$latest_tag" ./builds/* --generate-notes
  ;;
"")
  ./do generate binaries
  gh release create "$latest_tag" ./builds/* --generate-notes
  ;;
*) echo "Unrecognized option ($1)" ;;
esac
