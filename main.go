package main

import (
	"fmt"

	_ "net/http/pprof"
)

func showHelpMessage() {
	fmt.Println(`
Usage: 	kite <SUBCOMMAND>

SUBCOMMANDS:
build
serve
`)
}
