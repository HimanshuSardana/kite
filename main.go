package main

import (
	_ "net/http/pprof"

	cmd "github.com/HimanshuSardana/kite/cmd"
)

func main() {
	cmd.Execute()
}
