all: run

# Version metadata injected into internal/version. Release builds (tags)
# get e.g. v0.1.0; local builds get "dev" + commit hash.
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo none)
DATE := $(shell date -u +%Y-%m-%dT%H:%M:%SZ 2>/dev/null || echo unknown)
VERSION_PKG := github.com/HimanshuSardana/kite/internal/version
LDFLAGS := -X $(VERSION_PKG).Version=$(VERSION) -X $(VERSION_PKG).Commit=$(COMMIT) -X $(VERSION_PKG).Date=$(DATE)
RELEASE_LDFLAGS := -s -w $(LDFLAGS)

run:
	go run -ldflags "$(LDFLAGS)" .

build:
	go build -ldflags "$(LDFLAGS)" .

build-release:
	go build -ldflags "$(RELEASE_LDFLAGS)" -o kite-release .
	upx --best --lzma kite-release
	stat -c %s ./kite-release

serve:
	go run -ldflags "$(LDFLAGS)" . serve
