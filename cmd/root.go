package cmd

import (
	"fmt"
	"os"

	"github.com/HimanshuSardana/kite/internal/build"
)

const (
	DefaultTheme = "modern-light"
	DefaultPort  = "8000"
)

func Execute() {
	args := os.Args
	if len(args) < 2 {
		build.ShowHelpMessage()
		return
	}

	switch args[1] {
	case "build":
		runBuild(args)
	case "serve":
		runServe(args)
	case "list-themes":
		runListThemes(args)
	default:
		build.ShowHelpMessage()
	}
}

func ShowHelp() {
	fmt.Println(`
Kite — A lightweight static site generator

USAGE:
  kite <command> [options]

COMMANDS:
  build         Build the static site into the output directory
  serve         Start a local development server with live reload
  list-themes   List all available themes

OPTIONS:
  -h, --help    Show this help message

EXAMPLES:
  kite build
  kite build gruvbox
  kite serve
  kite list-themes

DESCRIPTION:
  Kite converts your content into a static website using themes and templates.
`)
}
