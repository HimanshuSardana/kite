package cmd

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func runServe(args []string) {
	themeName := DefaultTheme
	port := DefaultPort

	for i := 2; i < len(args); i++ {
		if args[i] == "--port" && i+1 < len(args) {
			port = args[i+1]
		}
		if args[i] != "--port" && args[i] != "--help" && args[i] != "-h" {
			themeName = args[i]
		}
	}

	themeCSS := fmt.Sprintf("./themes/%s/style.css", themeName)
	outputCSS := "./output/style.css"

	if err := copyFile(themeCSS, outputCSS); err != nil {
		log.Printf("Warning: Could not copy theme CSS: %v", err)
	}

	fs := http.FileServer(http.Dir("./output/"))
	http.Handle("/", fs)

	log.Printf("Serving on http://localhost:%s", port)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Server error: %s\n", err)
	}
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	if err := os.MkdirAll(filepath.Dir(dst), os.ModePerm); err != nil {
		return err
	}

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}
