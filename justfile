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

# Parameters

skip-compress := env_var_or_default("SKIP_COMPRESS", "0")

default:
    go build -v .

release:
    # Cleaning ./builds/
    just clean
    # Linux
    just build linux 386
    just build linux amd64
    just build linux arm
    just build linux arm64
    # Windows
    just build windows 386
    just build windows amd64
    just build windows arm
    just build windows arm64
    # Darwin
    just build darwin amd64
    just build darwin arm64

build os arch:
    #!/usr/bin/env -S bash -x
    bin=builds/fsize-{{os}}-{{arch}}

    if [[ "{{ os }}" == "windows" ]]; then
        bin="$bin.exe"
    fi

    GOOS={{os}} GOARCH={{arch}} \
    go build -v \
    {{version-flag}} \
    -o $bin

    if [[ {{skip-compress}} == 1 ]]; then
        exit 0
    fi

    if [[ {{os}} == "windows" && {{arch}} == "arm64" ]]; then
        echo ---------------------------------------------
        echo compression not supported for {{os}}-{{arch}}
        echo skipping compression process...
        echo ---------------------------------------------
        exit 0
    fi
 
    just compress $bin

build-local:
    @ go build -v \
    {{version-flag}} .

clean:
    @rm -rf builds completions ./fsize

go-install:
    go install -v {{go-install-version-flag}} github.com/Tom5521/fsize@{{short-latest-tag}}

go-uninstall:
    rm ~/go/bin/fsize

go-reinstall:
    @just go-uninstall
    @just go-install

[private]
compress bin:
    #!/usr/bin/env -S bash -x

    if [[ {{skip-compress}} == 1 ]]; then
        echo skipping compression of {{bin}}...
        exit 0
    fi

    which upx > /dev/null 2>&1
    if [[ $? != 0 ]]; then
        echo ---------------------------------
        echo upx binary not found in PATH
        echo skipping compression process...
        echo ---------------------------------
        exit 0
    fi

    upx --force-macos --8mib-ram -9 {{bin}}
    upx -t {{bin}}

[confirm]
[unix]
install:
    just build-local
    cp fsize {{linux-install-path}}
    fsize --gen-bash-completion {{bash-completion-path}}fsize
    -command -v fish && \
    fsize --gen-fish-completion {{fish-completion-path}}fsize.fish 
    -command -v zsh && \
    fsize --gen-zsh-completion {{zsh-completion-path}}_fsize

[confirm]
[windows]
install:
    just build-local
    cp fsize.exe C:/Windows/System32/

[confirm]
[unix]
uninstall:
    -rm {{linux-install-path}} \
    {{bash-completion-path}}fsize \
    {{fish-completion-path}}fsize.fish
    -rm {{zsh-completion-path}}_fsize

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
    gh release create {{short-latest-tag}} ./builds/* --generate-notes

test:
    go test -v ./*/*_test.go

test-update:
    #!/bin/bash
    go build -v .
    ./fsize --update
    v=$(./fsize --version)
    if [[ "$v" != "fsize version {{short-latest-tag}}" ]]; then
        exit 1
    fi


update-asciinema:
    just build-local
    asciinema rec --title "fsize {{short-latest-tag}}" --command "./fsize /usr/share/"
