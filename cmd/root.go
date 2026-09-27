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
	case "-v", "--version", "version":
		runVersion(args)
	case "update":
		runUpdate(args)
	case "build":
		runBuild(args)
	case "serve":
		runServe(args)
	case "list-themes":
		runListThemes(args)
	case "init":
		runInit(args)
	case "new":
		runNew(args)
	case "check":
		runCheck(args)
	default:
		build.ShowHelpMessage()
	}
}

func runInit(args []string) {
	if err := RunInit(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func ShowHelp() {
	fmt.Print(`
Kite — A lightweight static site generator

USAGE:
  kite <command> [options]

COMMANDS:
  build         Build the static site into the output directory
  serve         Start a local development server with live reload
  list-themes   List all available themes
  init          Initialize a new blog project
  new           Create a new post with frontmatter
  check         Check built site for broken internal links
  version       Show the kite version
  update        Update kite to the latest GitHub release

OPTIONS:
  -h, --help    Show this help message
  -v, --version Show the kite version

EXAMPLES:
  kite build
  kite build gruvbox
  kite serve
  kite list-themes
  kite init
  kite new hello-world.md
  kite check
  kite version
  kite update

DESCRIPTION:
  Kite converts your content into a static website using themes and templates.
  Use 'kite init' to start a new blog project.
`)
}
