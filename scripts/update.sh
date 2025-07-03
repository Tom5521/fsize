#!/usr/bin/env -S bash -x

long_latest_version=$(git describe --tags)

version() {
  echo "$long_latest_version" >./meta/version.txt
}

locales() {
  if ! command -v xgotext; then
    ./do install xgotext
  fi

  xgotext . -o ./po/en/default.pot --lang en --package-version "$(cat ./meta/version.txt)"
  for dir in ./po/*; do
    if [[ "$dir" == "en" ]]; then
      continue
    fi

    local file=$dir/default.po
    local lang
    lang=$(basename "$(dirname "$file")")
    msgmerge -U --lang "$lang" "$file" ./po/en/default.pot
  done
  find po -name "*.po~" -delete
}

_asciinema() {
  go build -v .
  asciinema rec --title "fsize $(cat ./meta/version.txt)" \
    --command "./fsize /usr/share" ./screenshots/demo.cast \
    --overwrite
}

case "$1" in
"version") version ;;
"locales") locales ;;
"asciinema") _asciinema ;;
*) echo "Unrecognized option ($1)" ;;
esac
