#!/usr/bin/env -S bash -x

go-install() {
  go install -v github.com/Tom5521/fsize@$(git describe --tags --abbrev=0)
}

_install() {
  prefix=$1

  ./do generate completions
  install -D "./fsize" "$prefix/bin/"
  install -D "./completions/fsize.fish" \
    "$prefix/share/fish/vendor_completions.d/fsize.fish"
  install -D "./completions/fsize.sh" \
    "$prefix/share/bash-completion/completions/fsize"
  install -D "./completions/_fsize" \
    "$prefix/share/zsh/site-functions/_fsize"
}

_xgotext() {
  filename=gotext-tools_$(uname -s)_$(uname -m).tar.gz
  cd /tmp
  wget https://github.com/Tom5521/gotext-tools/releases/latest/download/"$filename"
  tar -xzf "$filename"
  install -D ./xgotext "$HOME/.local/bin/xgotext"
}

case "$1" in
"go") go-install ;;
"local") _install "$HOME/.local" ;;
"system") _install "/usr/local" ;;
"xgotext") _xgotext ;;
esac
