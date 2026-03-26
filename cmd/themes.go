package cmd

import (
	"fmt"
	"strings"

	"github.com/HimanshuSardana/kite/internal/build"
)

func runListThemes(args []string) {
	themeList := build.ListThemes("")
	fmt.Println(strings.Join(themeList, "\n"))
}
