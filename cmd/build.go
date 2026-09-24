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
	includeDrafts := false

	if cfg, err := config.Load("config.yaml"); err == nil && cfg.DefaultTheme != "" {
		themeName = cfg.DefaultTheme
	}

	for _, arg := range args[2:] {
		if arg == "--drafts" {
			includeDrafts = true
		} else {
			themeName = arg
		}
	}

	opts := build.BuildOptions{
		ThemeName:     themeName,
		IncludeDrafts: includeDrafts,
	}

	fmt.Printf("Building with theme: %s\n", themeName)
	if includeDrafts {
		fmt.Println("Including drafts")
	}

	if err := build.Build(opts); err != nil {
		log.Fatalf("Build failed: %v", err)
		os.Exit(1)
	}

	fmt.Println("Build completed successfully!")
}
