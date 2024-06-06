

short-latest-tag := `git describe --tags --abbrev=0`
long-latest-tag := `git describe --tags`


version-flag := '-ldflags "-X github.com/Tom5521/fsize/meta.Version=' / long-latest-tag / '"'
go-install-version-flag := '-ldflags "-X github.com/Tom5521/fsize/meta.Version=' / short-latest-tag / '"'

fish-completion-path := "/usr/share/fish/vendor_completions.d/"
bash-completion-path := "/usr/share/bash-completion/completions/"
zsh-completion-path := "/usr/share/zsh/site-functions/"

fish-local-completion-path := "~/.config/fish/completions/"
bash-local-completion-path := "~/local/share/bash-completion/completions/"
zsh-local-completion-path := "$fpath/"


linux-install-path := "/usr/bin/fsize"
linux-local-install-path := "~/.local/bin/fsize"

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
  {{version-flag}} \
  -o builds/fsize-{{os}}-{{arch}}\
  $([[ "{{os}}" == "windows" ]] && echo ".exe")
build-local:
  @ go build -v \
  {{version-flag}} .
build-linux arch:
  @just build linux {{arch}}
build-windows arch:
  @just build windows {{arch}}
build-darwin arch:
  @just build darwin {{arch}}
clean:
  @rm -rf builds completions ./fsize
go-install:
  go install -v {{go-install-version-flag}} github.com/Tom5521/fsize@{{short-latest-tag}}
go-uninstall:
  rm ~/go/bin/fsize
go-reinstall:
  @just go-uninstall
  @just go-install
linux-local-install:
  just build-local
  cp fsize {{linux-local-install-path}}
  -[ -d "{{bash-local-completion-path}}" ] && \
  fsize --gen-bash-completion {{bash-local-completion-path}}fsize
  -which fish && \
  fsize --gen-fish-completion {{fish-local-completion-path}}fsize.fish 
  -which zsh && \
  fsize --gen-zsh-completion {{zsh-local-completion-path}}_fsize
linux-local-uninstall:
  -rm {{linux-local-install-path}} \
  {{bash-local-completion-path}}fsize \
  {{fish-local-completion-path}}fsize.fish
  -rm {{zsh-local-completion-path}}_fsize
linux-local-reinstall:
  just linux-local-uninstall
  just linux-local-install
linux-install:
  just build-local
  cp fsize {{linux-install-path}}
  fsize --gen-bash-completion {{bash-completion-path}}fsize
  -which fish && \
  fsize --gen-fish-completion {{fish-completion-path}}fsize.fish 
  -which zsh && \
  fsize --gen-zsh-completion {{zsh-completion-path}}_fsize
linux-uninstall:
  -rm {{linux-install-path}} \
  {{bash-completion-path}}fsize \
  {{fish-completion-path}}fsize.fish
  -rm {{zsh-completion-path}}_fsize
linux-reinstall:
  just linux-uninstall
  just linux-install
generate-completions:
  mkdir -p completions
  just build-local
  ./fsize --gen-bash-completion ./completions/fsize.sh
  ./fsize --gen-fish-completion ./completions/fsize.fish
  ./fsize --gen-zsh-completion ./completions/_fsize
