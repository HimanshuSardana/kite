package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/HimanshuSardana/kite/internal/build"
)

func runListThemes(args []string) {
	themeList, err := build.ListThemes("")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(strings.Join(themeList, "\n"))
}
