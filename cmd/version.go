package cmd

import (
	"fmt"

	"github.com/HimanshuSardana/kite/internal/version"
)

func runVersion(args []string) {
	// `kite -v` / `kite --version` stay terse for scripts;
	// the `version` subcommand prints full build metadata.
	if len(args) > 1 && args[1] == "version" {
		fmt.Printf("kite %s\n", version.String())
		return
	}
	fmt.Printf("kite %s\n", version.Version)
}
