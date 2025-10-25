SHELL=/usr/bin/bash

GOOS := $(shell go env GOOS)
GOARCH := $(shell go env GOARCH)
GOPATH := $(shell go env GOPATH)
GOENV := GOOS=$(GOOS) GOARCH=$(GOARCH)
BUILD_TAGS :=
BUILD_ARGS :=

ROOT_PREFIX := /usr/local
LOCAL_PREFIX := $(HOME)/.local
PREFIX = $(if $(filter local,$(1)),$(LOCAL_PREFIX),$(ROOT_PREFIX))

SUDO = $(if $(and $(filter root,$(1)),$(filter-out root,$(USER))),sudo)

LATEST_TAG := $(shell git describe --tags)
LATEST_TAG_SHORT := $(shell git describe --tags --abbrev=0)

# These are the platforms that I want to maintain & support.
# These are the ones that will be used as release binaries.
define SUPPORTED_PLATFORMS
windows/amd64
windows/386
windows/arm
windows/arm64
linux/amd64
linux/386
linux/arm64
linux/arm
darwin/amd64
darwin/arm64
android/arm64
android/arm
endef

# These are just the platforms that fsize is compatible with.
define COMPATIBLE_PLATFORMS
windows/amd64
windows/386
windows/arm
windows/arm64
dragonfly/amd64
freebsd/386
freebsd/amd64
freebsd/arm
freebsd/arm64
freebsd/riscv64
linux/amd64
linux/386
linux/arm64
linux/arm
linux/loong64
linux/mips
linux/mips64
linux/mips64le
linux/mipsle
linux/ppc64
linux/ppc64le
linux/riscv64
linux/s390x
netbsd/386
netbsd/amd64
netbsd/arm
netbsd/arm64
openbsd/386
openbsd/amd64
openbsd/arm
openbsd/arm64
openbsd/ppc64
openbsd/riscv64
solaris/amd64
darwin/amd64
darwin/arm64
android/arm64
android/arm
endef

TARGETS := COMPATIBLE_PLATFORMS
VERBOSE ?= 0

OS_EXT = $(if $(filter windows,$(1)),.exe)

RED = $(shell tput setaf 1)
GREEN = $(shell tput setaf 2)
YELLOW = $(shell tput setaf 3)
BLUE = $(shell tput setaf 4)
BOLD = $(shell tput bold)
NC = $(shell tput sgr0)

ERROR = @echo "$(BOLD)$(RED)ERROR:$(NC) $(1)"
WARN = @echo "$(YELLOW)WARNING:$(NC) $(1)"
INFO = @echo "$(BOLD)$(GREEN)INFO:$(NC) $(1)"
TITLE = @echo "$(BOLD)$(GREEN)$(1)$(NC)"

override CMD := $(GOENV) go
override V_FLAG := $(if $(filter 1,$(VERBOSE)),-v)
override WIN_EXT := $(call OS_EXT,$(GOOS))
override GO_PACKAGE := github.com/Tom5521/fsize

MANAGE_UPDATE := 1

LD_FLAGS := -s -w
LD_FLAGS += -X '$(GO_PACKAGE)/meta.LongVersion=$(LATEST_TAG)'
LD_FLAGS += -X '$(GO_PACKAGE)/meta.Version=$(LATEST_TAG_SHORT)'
ifneq ($(MANAGE_UPDATE),1)
LD_FLAGS += -X 'main.releaseTarget=package-manager'
endif

XGOTEXT_SUFFIX := $(if $(WIN_EXT),.exe,.bin)
XGOTEXT_PATH := $(if $(shell command -v xgotext),\
	xgotext,./xgotext$(XGOTEXT_SUFFIX))

define CLEAN_TARGETS
*.mo
*.po~
*.log
*.diff
*.rej
*.orig
completions
builds
fsize
changelog.md
endef

BIN = ./builds/fsize-$(1)-$(2)$(call WIN_EXT,$(1))
override CURRENT_BIN := $(call BIN,$(GOOS),$(GOARCH))
override NATIVE_GOOS := $(shell GOOS="" go env GOOS)
override NATIVE_GOARCH := $(shell GOARCH="" go env GOARCH)
override NATIVE_BIN := $(call BIN,$(NATIVE_GOOS),$(NATIVE_GOARCH))

NDK_TOOLCHAIN ?= /opt/android-ndk/toolchains/llvm/prebuilt/linux-x86_64/bin

.PHONY: test all clean default run build build-targets install uninstall \
	%-completions-install %-completions-uninstall go-install \
	go-uninstall %-install %-uninstall release update-assets po/en/default.pot \
	xgotext$(XGOTEXT_SUFFIX) build-supported build-compatible
.DEFAULT_GOAL := default

default: test
test:
	go test $(V_FLAG) ./...

run:
	$(CMD) run $(V_FLAG) .

.ONESHELL:
clean:
	targets=($(CLEAN_TARGETS))
	for target in $${targets[@]}; do
		find . -name "$$target" | xargs rm -rf
	done

screenshots/demo.gif:
	vhs ./screenshots/demo.tape

.ONESHELL:
.SILENT:
xgotext$(XGOTEXT_SUFFIX):
	if stat ./xgotext$(XGOTEXT_SUFFIX) &> /dev/null; then
		exit
	fi
	if command -v xgotext &> /dev/null; then
		exit
	fi
	$(call WARN,xgotext isn't installed.)
	$(call WARN,downloading xgotext...)
	wget -O "$(XGOTEXT_PATH)" \
		"https://github.com/Tom5521/gotext-tools/releases/latest/download/xgotext-$(NATIVE_GOOS)-$(NATIVE_GOARCH)$(call OS_EXT,$(NATIVE_GOOS))" 2> /dev/null
	chmod +x "$(XGOTEXT_PATH)"
	$(call TITLE,xgotext downloaded successfully!)

po/en/default.pot: xgotext$(XGOTEXT_SUFFIX)
	$(call INFO,Updating english template...)
	$(XGOTEXT_PATH) . -o ./po/en/default.pot --lang en --package-version \
		$(LATEST_TAG)

.ONESHELL:
.SILENT:
po: po/en/default.pot
	$(call TITLE,Updating translations...)
	for dir in ./po/*; do
		if [[ "$$dir" == "./po/en" ]]; then
			continue
		fi
		file=$$dir/default.po
		lang=$$(basename "$$(dirname "$$file")")
		$(call INFO,Updating $$lang...)
		msgmerge -U --lang "$$lang" "$$file" ./po/en/default.pot 2> /dev/null
	done
	$(call INFO,Deleting intermediate files...)
	$(MAKE) -s clean
	$(call INFO,po generations finished successfully!)

log.diff:
	git diff --staged > log.diff

.ONESHELL:
changelog.md:
	echo '## Changelog' > changelog.md
	echo >> changelog.md

	latest_tag="$(LATEST_TAG_SHORT)"
	penultimate_tag=$$(git describe --tags --abbrev=0 "$$latest_tag^")

	git log --pretty=format:'- [%h](https://github.com/Tom5521/fsize/commit/%H): %s' \
		$$penultimate_tag..$$latest_tag >> changelog.md

.SILENT:
.ONESHELL:
build:
	if [[ "$(GOOS)" == "android"  && "$(NATIVE_GOOS)" != "android" ]];then
		export CGO_ENABLED=1
		if [[ "$(GOARCH)" == "arm64" ]];then
			export CC=$(NDK_TOOLCHAIN)/aarch64-linux-android35-clang
		elif [[ "$(GOARCH)" == "arm" ]]; then
			export CC=$(NDK_TOOLCHAIN)/armv7a-linux-androideabi35-clang
		fi
	fi
	$(CMD) build $(V_FLAG) \
	$(BUILD_ARGS) \
	$(if $(BUILD_TAGS),-tags "$(BUILD_TAGS)") \
	-o $(CURRENT_BIN) \
	-ldflags="$(LD_FLAGS)" .

.ONESHELL:
.SILENT:
build-targets: clean
	platforms=(
		$(value $(TARGETS))
	)
	$(call TITLE,Building...)
	for platform in $${platforms[@]}; do
		os=$$(echo "$$platform"| cut -d'/' -f1)
		arch=$$(echo "$$platform"| cut -d'/' -f2-)

		$(call INFO,$(BLUE)$$os$(NC)/$(BOLD)$$arch$(NC))
		$(MAKE) -s build GOOS=$$os GOARCH=$$arch
	done

build-supported: 
	$(MAKE) build-targets -s TARGETS='SUPPORTED_PLATFORMS'
build-compatible:
	$(MAKE) build-targets -s TARGETS='COMPATIBLE_PLATFORMS'

release: build-targets changelog.md
	gh release create $(LATEST_TAG_SHORT) \
		--notes-file ./changelog.md \
		--title $(LATEST_TAG_SHORT) \
		--fail-on-no-commits builds/*

update-assets: build-targets
	gh release upload --clobber "$(LATEST_TAG_SHORT)" ./builds/*

completions: $(NATIVE_BIN)
	@mkdir -p completions
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

install: local-install
uninstall: local-uninstall

go-install:
	go install $(V_FLAG) .
	$(MAKE) -s local-completions-install

go-uninstall:
	rm -f $(GOPATH)/bin/fsize
	$(MAKE) -s local-completions-uninstall

%-install:
	$(MAKE) -s $*-completions-install
	$(call SUDO,$*) install -D $(NATIVE_BIN) $(call PREFIX,$*)/bin/fsize$(call OS_EXT,$(NATIVE_GOOS))

%-uninstall:
	$(call SUDO,$*) rm -f $(call PREFIX,$*)/bin/fsize$(call OS_EXT,$(NATIVE_GOOS))
	$(MAKE) -s $*-completions-uninstall

.ONESHELL:
builds/fsize-%$(WIN_EXT):
	os=$(word 1,$(subst -, ,$*))
	arch=$(word 1,$(subst ., ,$(word 2,$(subst -, ,$*))))
	$(MAKE) build -s GOOS=$$os GOARCH=$$arch
