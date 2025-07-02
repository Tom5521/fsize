#!/usr/bin/env -S bash -x

os=$1
arch=$2

if [[ "$os" == "" ]]; then
  os=$(go env GOOS)
fi

if [[ "$arch" == "" ]]; then
  arch=$(go env GOARCH)
fi

bin="builds/fsize-$os-$arch"

if [[ "$os" == "windows" ]]; then
  bin="$bin.exe"
fi

go build -v -o "$bin" || exit $?
