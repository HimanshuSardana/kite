all: run

run:
	go run .

build:
	echo "Building Kite release"
	go build -ldflags="-s -w" -o kite-release
	upx --best --lzma kite-release
	echo "Successfully built release binary"
	stat -c %s ./kite-release

serve:
	go run main.go serve
