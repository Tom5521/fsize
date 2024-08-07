# Tags

short-latest-tag := `git describe --tags --abbrev=0`
long-latest-tag := `git describe --tags`

# Flags

version-flag := '-ldflags "-X github.com/Tom5521/fsize/meta.Version=' + short-latest-tag + '"'
go-install-version-flag := '-ldflags "-X github.com/Tom5521/fsize/meta.Version=' + short-latest-tag + '"'

# Paths

fish-completion-path := "/usr/local/share/fish/vendor_completions.d/"
bash-completion-path := "/usr/local/share/bash-completion/completions/"
zsh-completion-path := "/usr/local/share/zsh/site-functions/"
linux-install-path := "/usr/local/bin/fsize"

default:
    go build -v .

release:
    # Cleaning ./builds/
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

[unix]
build os arch:
    @ GOOS={{ os }} GOARCH={{ arch }} \
    go build -v \
    {{ version-flag }} \
    -o builds/fsize-{{ os }}-{{ arch }}\
    $([[ "{{ os }}" == "windows" ]] && echo ".exe")

[windows]
build os arch:
    $extension = if ("{{ os }}" -eq "windows") { ".exe" } else { "" }
    $output = "builds/fsize-{{ os }}-{{ arch }}$extension"
    @ GOOS={{ os }} GOARCH={{ arch }} \
    go build -v \
    {{ version-flag }} \
    -o $output

build-local:
    @ go build -v \
    {{ version-flag }} .

build-linux arch:
    @just build linux {{ arch }}

build-windows arch:
    @just build windows {{ arch }}

build-darwin arch:
    @just build darwin {{ arch }}

[unix]
clean:
    @rm -rf builds completions ./fsize

[windows]
clean:
    @del builds completions .\\fsize.exe

go-install:
    go install -v {{ go-install-version-flag }} github.com/Tom5521/fsize@{{ short-latest-tag }}

go-uninstall:
    rm ~/go/bin/fsize

go-reinstall:
    @just go-uninstall
    @just go-install

[confirm]
[unix]
install:
    just build-local
    cp fsize {{ linux-install-path }}
    fsize --gen-bash-completion {{ bash-completion-path }}fsize
    -which fish && \
    fsize --gen-fish-completion {{ fish-completion-path }}fsize.fish 
    -which zsh && \
    fsize --gen-zsh-completion {{ zsh-completion-path }}_fsize

[confirm]
[windows]
install:
    just build-local
    cp fsize.exe C:/Windows/System32/

[confirm]
[unix]
uninstall:
    -rm {{ linux-install-path }} \
    {{ bash-completion-path }}fsize \
    {{ fish-completion-path }}fsize.fish
    -rm {{ zsh-completion-path }}_fsize

[confirm]
[windows]
uninstall:
    rm -rf C:/Windows/System32/fsize.exe

[confirm]
reinstall:
    just --yes uninstall
    just --yes install

generate-completions:
    mkdir -p completions
    just build-local
    ./fsize --gen-bash-completion ./completions/fsize.sh
    ./fsize --gen-fish-completion ./completions/fsize.fish
    ./fsize --gen-zsh-completion ./completions/_fsize

commit:
    git add .
    meteor
    git push

gh-release:
    just release
    gh release create {{ short-latest-tag }} ./builds/* --generate-notes

test:
    go test -v ./*/*_test.go

update-asciinema:
    just build-local
    asciinema rec --title "fsize {{ short-latest-tag }}" --command "./fsize /usr/share/"
