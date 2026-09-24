package cmd

import (
	"fmt"
	"os"

	"github.com/HimanshuSardana/kite/internal/build"
)

func runCheck(args []string) {
	broken, err := build.CheckSite("./output")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if len(broken) == 0 {
		fmt.Println("No broken links ✓")
		return
	}

	fmt.Printf("Found %d broken link(s):\n", len(broken))
	for _, b := range broken {
		fmt.Printf("  %s -> %s\n", b.Page, b.Target)
	}
	os.Exit(1)
}
