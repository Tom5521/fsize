SHELL=/usr/bin/bash

GOOS := $(shell go env GOOS)
GOARCH := $(shell go env GOARCH)
GOPATH := $(shell go env GOPATH)
GOFLAGS := GOOS=$(GOOS) GOARCH=$(GOARCH)

ROOT_PREFIX := /usr/local
LOCAL_PREFIX := $(HOME)/.local
PREFIX = $(if $(filter local,$(1)),$(LOCAL_PREFIX),$(ROOT_PREFIX))

SUDO = $(if $(and $(filter root,$(1)),$(filter-out root,$(USER))),sudo)

LATEST_TAG := $(shell git describe --tags)
LATEST_TAG_SHORT := $(shell git describe --tags --abbrev=0)

SUPPORTED_OSES := windows linux darwin
SUPPORTED_ARCHITECTURES := 386 amd64 arm arm64

VERBOSE ?= 0

MWIN_EXT = $(if $(filter windows,$(1)),.exe)

override CMD := $(GOFLAGS) go
override V_FLAG := $(if $(filter 1,$(VERBOSE)),-v)
override WIN_EXT := $(call MWIN_EXT,$(GOOS))

BIN = ./builds/fsize-$(1)-$(2)$(call WIN_EXT,$(1))
override NATIVE_GOOS := $(shell go env GOOS)
override NATIVE_GOARCH := $(shell go env GOARCH)
override NATIVE_BIN := $(call BIN,$(NATIVE_GOOS),$(NATIVE_GOARCH))

.PHONY: test all clean default run build build-all \
	%-completions-install %-completions-uninstall go-install \
	go-uninstall %-install %-uninstall release update-assets
.ONESHELL: po changelog.md build-all
.DEFAULT_GOAL := default

default: test
test:
	go test $(V_FLAG) ./...
run:
	$(CMD) run $(V_FLAG) .
clean:
	rm -rf completions builds fsize changelog.md
	find . -name "*.mo" -delete
	find . -name "*.po~" -delete
	find . -name "*.log" -delete
	find . -name "*.diff" -delete
screenshots/demo.cast: build
	LANG=en asciinema rec --title "fsize $(LATEST_TAG)" \
		--command "$(NATIVE_BIN) /usr/share" ./screenshots/demo.cast \
		--overwrite
meta/version.txt:
	$(LATEST_TAG) > meta/version.txt

po:
	if ! command -v xgotext; then
		bin=$(HOME)/.local/bin/xgotext$(call MWIN_EXT,$(NATIVE_GOOS))
		mkdir -p builds "$$(dirname $$bin)"
		wget -O "$$bin" \
			"https://github.com/Tom5521/gotext-tools/releases/latest/download/xgotext-$(NATIVE_GOOS)-$(NATIVE_GOARCH)$(call MWIN_EXT,$(NATIVE_GOOS))"
		chmod +x "$$bin"
	fi
	xgotext . -o ./po/en/default.pot --lang en --package-version \
		$(LATEST_TAG)
	for dir in ./po/*; do
		if [[ "$$dir" == "./po/en" ]]; then
			continue
		fi
		
		file=$$dir/default.po
		lang=$$(basename "$$(dirname "$$file")")
		
		msgmerge -U --lang "$$lang" "$$file" ./po/en/default.pot
	done
	find po -name "*.po~" -delete

log.diff:
	git diff --staged > log.diff

changelog.md:
	echo '## Changelog' > changelog.md
	echo >> changelog.md

	latest_tag=$$(git describe --tags --abbrev=0)
	penultimate_tag=$$(git describe --tags --abbrev=0 "$$latest_tag^")

	git log --pretty=format:'- [%h](https://github.com/Tom5521/fsize/commit/%H): %s' \
		$$penultimate_tag..$$latest_tag >> changelog.md
build:
	$(CMD) build $(V_FLAG) -o ./builds/fsize-$(GOOS)-$(GOARCH) \
	-ldflags '-s -w' .

build-all: clean
	valid=$$($(CMD) tool dist list)
	for os in $(SUPPORTED_OSES); do
		for arch in $(SUPPORTED_ARCHITECTURES); do
			if ! echo $$valid | grep -qw "$$os/$$arch"; then 
				continue
			fi
			$(MAKE) build GOOS=$$os GOARCH=$$arch
		done
	done
release: build-all changelog.md
	gh release create $(LATEST_TAG_SHORT) --notes-file \
		./changelog.md --fail-on-no-commits builds/*
update-assets: build-all changelog.md
	gh release upload "$(LATEST_TAG_SHORT)" --notes-file \
		./changelog.md ./builds/*
completions: build
	mkdir -p completions
	$(NATIVE_BIN) --gen-bash-completion ./completions/fsize.bash
	$(NATIVE_BIN) --gen-fish-completion ./completions/fsize.fish
	$(NATIVE_BIN) --gen-zsh-completion ./completions/fsize.zsh
%-completions-install: completions
	$(call SUDO,$*) install -D ./completions/fsize.fish \
		$(call PREFIX,$*)/share/fish/vendor_completions.d/fsize.fish
	$(call SUDO,$*) install -D ./completions/fsize.bash \
		$(call PREFIX,$*)/share/bash-completion/completions/fsize
	$(call SUDO,$*) install -D ./completions/fsize.zsh \
		$(call PREFIX,$*)/share/zsh/site-functions/_fsize
%-completions-uninstall:
	$(call SUDO,$*) rm -f $(call PREFIX,$*)/share/fish/vendor_completions.d/fsize.fish
	$(call SUDO,$*) rm -f $(call PREFIX,$*)/share/bash-completion/completions/fsize
	$(call SUDO,$*) rm -f $(call PREFIX,$*)/share/zsh/site-functions/_fsize
install:
	$(MAKE) local-install
go-install:
	go install $(V_FLAG) .
	$(MAKE) local-completions-install
go-uninstall:
	rm -f $(GOPATH)/bin/fsize
	$(MAKE) local-completions-uninstall
%-install:
	$(MAKE) $*-completions-install
	$(call SUDO,$*) install -D $(NATIVE_BIN) $(call PREFIX,$*)/bin/fsize$(call MWIN_EXT,$(NATIVE_GOOS))
%-uninstall:
	$(call SUDO,$*) rm -f $(call PREFIX,$*)/bin/fsize$(call MWIN_EXT,$(NATIVE_GOOS))
	$(MAKE) $*-completions-uninstall
