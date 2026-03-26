package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/HimanshuSardana/kite/internal/build"
)

func runBuild(args []string) {
	themeName := DefaultTheme

	if len(args) > 2 {
		themeName = args[2]
	}

	opts := build.BuildOptions{
		ThemeName: themeName,
	}

	fmt.Printf("Building with theme: %s\n", themeName)

	if err := build.Build(opts); err != nil {
		log.Fatalf("Build failed: %v", err)
		os.Exit(1)
	}

	fmt.Println("Build completed successfully!")
}
