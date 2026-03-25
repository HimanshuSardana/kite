all: run

run:
	go run ./cmd

build:
	go run ./cmd build

build-release:
	echo "Building Kite release"
	go build -ldflags="-s -w" -o kite-release ./cmd
	upx --best --lzma kite-release
	echo "Successfully built release binary"
	stat -c %s ./kite-release

serve:
	go run ./cmd serve
