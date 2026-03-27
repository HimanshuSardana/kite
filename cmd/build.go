package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/HimanshuSardana/kite/internal/build"
	"github.com/HimanshuSardana/kite/pkg/config"
)

func runBuild(args []string) {
	themeName := DefaultTheme

	if cfg, err := config.Load("config.yaml"); err == nil && cfg.DefaultTheme != "" {
		themeName = cfg.DefaultTheme
	}

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
