package cmd

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/HimanshuSardana/kite/internal/build"
	"github.com/HimanshuSardana/kite/pkg/config"
)

var (
	reloadChan = make(chan struct{})
	reloadMutex sync.RWMutex
)

func runServe(args []string) {
	themeName := DefaultTheme
	port := DefaultPort

	if cfg, err := config.Load("config.yaml"); err == nil && cfg.DefaultTheme != "" {
		themeName = cfg.DefaultTheme
	}

	for i := 2; i < len(args); i++ {
		if args[i] == "--port" && i+1 < len(args) {
			port = args[i+1]
			i++ // skip the port value
		} else if args[i] != "--help" && args[i] != "-h" {
			themeName = args[i]
		}
	}

	// Initial build
	buildSite(themeName)

	// Start file watcher for hot-reload
	go watchAndRebuild(themeName)

	// Serve static files
	fs := http.FileServer(http.Dir("./output/"))
	http.Handle("/", fs)

	// SSE endpoint for live-reload
	http.HandleFunc("/livereload", liveReloadHandler)

	log.Printf("Serving on http://localhost:%s (with hot-reload)", port)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Server error: %s\n", err)
	}
}

func buildSite(themeName string) {
	log.Println("Building site...")

	themeCSS := fmt.Sprintf("./themes/%s/style.css", themeName)
	outputCSS := "./output/style.css"

	if err := copyFile(themeCSS, outputCSS); err != nil {
		log.Printf("Warning: Could not copy theme CSS: %v", err)
	}

	opts := build.BuildOptions{
		ThemeName: themeName,
	}

	if err := build.Build(opts); err != nil {
		log.Printf("Build error: %v", err)
	} else {
		log.Println("Build completed successfully!")
		// Inject live-reload script into HTML files
		injectLiveReloadScript()
		// Notify browsers to reload
		triggerReload()
	}
}

func watchAndRebuild(themeName string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Printf("Failed to create watcher: %v", err)
		return
	}
	defer watcher.Close()

	contentDir := "./content"
	if err := addDirRecursive(watcher, contentDir); err != nil {
		log.Printf("Failed to watch content directory: %v", err)
		return
	}

	// Also watch themes directory for theme changes
	themesDir := "./themes"
	if err := addDirRecursive(watcher, themesDir); err != nil {
		log.Printf("Warning: Could not watch themes directory: %v", err)
	}

	if err := watcher.Add("./config.yaml"); err != nil {
		log.Printf("Warning: Could not watch config.yaml: %v", err)
	}

	log.Println("Watching for changes in content/, themes/, and config.yaml...")

	var debounceTimer *time.Timer

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			// Only trigger on write or create events
			if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create {
				log.Printf("Detected change: %s", event.Name)

				// If a new directory is created, watch it
				if event.Op&fsnotify.Create == fsnotify.Create {
					if info, err := os.Stat(event.Name); err == nil && info.IsDir() {
						addDirRecursive(watcher, event.Name)
					}
				}

				// Debounce: wait 500ms before rebuilding
				if debounceTimer != nil {
					debounceTimer.Stop()
				}
				debounceTimer = time.AfterFunc(500*time.Millisecond, func() {
					buildSite(themeName)
				})
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Printf("Watcher error: %v", err)
		}
	}
}

func addDirRecursive(watcher *fsnotify.Watcher, dir string) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			log.Printf("Watching directory: %s", path)
			return watcher.Add(path)
		}
		return nil
	})
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

func injectLiveReloadScript() {
	outputDir := "./output"
	liveReloadScript := `\n<script>
(function() {
    var evtSource = new EventSource('/livereload');
    evtSource.addEventListener('reload', function() {
        console.log('Reloading page...');
        location.reload();
    });
})();
</script>
`

	err := filepath.Walk(outputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(path, ".html") {
			// Read file
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			html := string(content)

			// Only inject if not already present
			if !strings.Contains(html, "EventSource('/livereload')") {
				// Inject before </body>
				if idx := strings.Index(html, "</body>"); idx != -1 {
					html = html[:idx] + liveReloadScript + html[idx:]
					if err := os.WriteFile(path, []byte(html), info.Mode()); err != nil {
						log.Printf("Warning: Could not inject live-reload script into %s: %v", path, err)
					}
				}
			}
		}
		return nil
	})

	if err != nil {
		log.Printf("Warning: Error injecting live-reload script: %v", err)
	}
}

func triggerReload() {
	go func() {
		reloadChan <- struct{}{}
	}()
}

func liveReloadHandler(w http.ResponseWriter, r *http.Request) {
	// Set SSE headers
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Keep connection open and send reload events
	for {
		select {
		case <-reloadChan:
			fmt.Fprintf(w, "event: reload\ndata: {}\n\n")
			if f, ok := w.(http.Flusher); ok {
				f.Flush()
			}
		case <-r.Context().Done():
			return
		}
	}
}
