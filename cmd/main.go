package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	internal "github.com/HimanshuSardana/kite/internal/build"
)

var themeName = "gruvbox"

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

func main() {
	args := os.Args
	if len(args) > 1 {
		switch args[1] {
		case "serve":
			copyFile("./themes/"+themeName+"/style.css", "./output/style.css")

			fs := http.FileServer(http.Dir("./output/"))
			http.Handle("/", fs)

			log.Println("Serving on http://localhost:8000")

			err := http.ListenAndServe(":8000", nil)
			if err != nil {
				log.Fatalf("Error occured %s\n", err)
			}
		case "build":
			internal.Build()
		case "list-themes":
			themeList := internal.ListThemes()
			fmt.Println(strings.Join(themeList, "\n"))
		default:
			internal.ShowHelpMessage()
		}
	} else {
		internal.ShowHelpMessage()
	}
}
