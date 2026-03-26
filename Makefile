all: run

run:
	go run .

build:
	go build .

build-release:
	go build -ldflags="-s -w" -o kite-release
	upx --best --lzma kite-release
	stat -c %s ./kite-release

serve:
	go run . serve
