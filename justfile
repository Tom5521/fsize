

default:
  go build -v .
release:
  just clean
  # Linux
  just build-linux amd64
  just build-linux arm64
  just build-linux 386
  # Windows
  just build-windows amd64
  just build-windows arm64
  just build-windows 386
  # Darwin
  just build-darwin amd64
  just build-darwin arm64
build os arch:
  @ GOOS={{os}} GOARCH={{arch}} \
  go build -v \
  -ldflags "-X github.com/Tom5521/fsize/meta.Version=$(git describe --tags)" \
  -o builds/fsize-{{os}}-{{arch}}\
  $([[ "{{os}}" == "windows" ]] && echo ".exe")
build-linux arch:
  @just build linux {{arch}}
build-windows arch:
  @just build windows {{arch}}
build-darwin arch:
  @just build darwin {{arch}}
clean:
  @rm -rf builds
install:
  go install -v github.com/Tom5521/fsize@latest
