#!/usr/bin/env -S bash -x

completions() {
  mkdir -p completions
  go build -v .
  ./fsize --gen-bash-completion ./completions/fsize.sh
  ./fsize --gen-fish-completion ./completions/fsize.fish
  ./fsize --gen-zsh-completion ./completions/_fsize
}
binaries() {
  ./do clean
  ./do update version

  local oses=(
    linux
    windows
    darwin
  )
  local archs=(
    386
    amd64
    arm
    arm64
  )

  local valid
  valid="$(go tool dist list)"

  for os in "${oses[@]}"; do
    for arch in "${archs[@]}"; do
      if ! echo "$valid" | grep -qw "$os/$arch"; then
        continue
      fi
      ./do build "$os" "$arch"
    done
  done
}

case "$1" in
"completions") completions ;;
"binaries") binaries ;;
*) echo "Unrecognized option ($1)" ;;
esac
