#!/usr/bin/env -S bash -x

go-uninstall() {
  rm -f "$(go env GOPATH)/bin/fsize"
}

uninstall() {
  prefix=$1

  rm -f "$prefix/bin/fsize"
  rm -f "$prefix/share/fish/vendor_completions.d/fsize.fish"
  rm -f "$prefix/share/bash-completion/completions/fsize"
  rm -f "$prefix/share/zsh/site-functions/_fsize"
}

case "$1" in
"go") go-uninstall ;;
"local") uninstall "$HOME/.local" ;;
"system") uninstall "/usr/local" ;;
esac
