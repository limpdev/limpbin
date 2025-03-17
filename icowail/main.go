package main

import (
	"context"
	"embed"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed frontend/dist
var assets embed.FS

// App struct
type App struct {
	ctx         context.Context
	svgDir      string
	watcher     *fsnotify.Watcher
	svgFiles    []string
	svgContents map[string]string
}

// NewApp creates a new App application struct
func NewApp(svgDir string) *App {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	// Ensure svgDir is absolute
	absDir, err := filepath.Abs(svgDir)
	if err != nil {
		log.Printf("Warning: NEED absolute path for %s: %v", svgDir, err)
		absDir = svgDir
	}

	return &App{
		svgDir:      absDir,
		watcher:     watcher,
		svgContents: make(map[string]string),
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// Log startup info
	runtime.LogInfof(ctx, "Starting application with SVG directory: %s", a.svgDir)

	// Create SVG directory if it doesn't exist
	if _, err := os.Stat(a.svgDir); os.IsNotExist(err) {
		runtime.LogInfo(ctx, "SVG directory doesn't exist, creating it...")
		if err := os.MkdirAll(a.svgDir, 0755); err != nil {
			runtime.LogErrorf(ctx, "Failed to create SVG directory: %v", err)
		}
	}

	// Start watching the SVG directory
	runtime.LogInfo(ctx, "Setting up file watcher...")
	err := a.watcher.Add(a.svgDir)
	if err != nil {
		runtime.LogErrorf(ctx, "Failed to watch directory: %v", err)
	}

	// Initial scan of SVG files
	runtime.LogInfo(ctx, "Performing initial SVG scan...")
	a.updateSvgFiles()

	// Start file watcher in a goroutine
	go a.watchFiles()

	// Send update event after a short delay to ensure frontend is ready
	go func() {
		time.Sleep(1 * time.Second)
		runtime.EventsEmit(ctx, "svg-files-updated")
	}()
}

// watchFiles watches for changes in the SVG directory
func (a *App) watchFiles() {
	for {
		select {
		case event, ok := <-a.watcher.Events:
			if !ok {
				return
			}

			if event.Op&(fsnotify.Write|fsnotify.Create|fsnotify.Remove) != 0 &&
				strings.HasSuffix(event.Name, ".svg") {
				runtime.LogInfof(a.ctx, "SVG file changed: %s", event.Name)
				a.updateSvgFiles()
				runtime.EventsEmit(a.ctx, "svg-files-updated")
			}

		case err, ok := <-a.watcher.Errors:
			if !ok {
				return
			}
			runtime.LogErrorf(a.ctx, "Watcher error: %v", err)
		}
	}
}

// updateSvgFiles scans the SVG directory and updates the file list
func (a *App) updateSvgFiles() {
	a.svgFiles = []string{}
	a.svgContents = make(map[string]string)

	// Check if directory exists before walking
	if _, err := os.Stat(a.svgDir); os.IsNotExist(err) {
		runtime.LogErrorf(a.ctx, "SVG directory does not exist: %s", a.svgDir)
		return
	}

	err := filepath.Walk(a.svgDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			runtime.LogErrorf(a.ctx, "Error accessing path %s: %v", path, err)
			return nil // Continue despite errors
		}

		if !info.IsDir() && strings.HasSuffix(strings.ToLower(info.Name()), ".svg") {
			relPath, err := filepath.Rel(a.svgDir, path)
			if err != nil {
				runtime.LogErrorf(a.ctx, "Error getting relative path for %s: %v", path, err)
				relPath = info.Name() // Fallback to just the filename
			}

			a.svgFiles = append(a.svgFiles, relPath)

			// Read SVG content
			content, err := os.ReadFile(path)
			if err != nil {
				runtime.LogErrorf(a.ctx, "Error reading SVG file %s: %v", path, err)
			} else {
				a.svgContents[relPath] = string(content)
			}
		}

		return nil
	})

	if err != nil {
		runtime.LogErrorf(a.ctx, "Error walking SVG directory: %v", err)
	}

	// Sort files alphabetically
	sort.Strings(a.svgFiles)

	runtime.LogInfof(a.ctx, "Found %d SVG files", len(a.svgFiles))
	for i, file := range a.svgFiles {
		if i < 5 { // Log first 5 files only to avoid spam
			runtime.LogInfof(a.ctx, "SVG file: %s", file)
		}
	}
}

// GetSvgFiles returns the list of SVG files
func (a *App) GetSvgFiles() []string {
	runtime.LogInfof(a.ctx, "GetSvgFiles called, returning %d files", len(a.svgFiles))
	return a.svgFiles
}

// GetSvgContent returns the content of a specific SVG file
func (a *App) GetSvgContent(filename string) string {
	content, exists := a.svgContents[filename]
	if !exists {
		runtime.LogWarningf(a.ctx, "SVG content not found for %s", filename)
		return ""
	}
	return content
}

// GetIconName extracts the icon name from the filename
func (a *App) GetIconName(filename string) string {
	return strings.TrimSuffix(filepath.Base(filename), ".svg")
}

// Debug helper method to check directory contents
func (a *App) DebugListDirectory() string {
	if a.ctx == nil {
		return "Context is nil"
	}

	var result strings.Builder
	result.WriteString(fmt.Sprintf("Listing contents of directory: %s\n", a.svgDir))

	files, err := os.ReadDir(a.svgDir)
	if err != nil {
		return fmt.Sprintf("Error reading directory: %v", err)
	}

	for _, file := range files {
		info, _ := file.Info()
		if info == nil {
			result.WriteString(fmt.Sprintf("- %s (unknown size)\n", file.Name()))
			continue
		} else {
			result.WriteString(fmt.Sprintf("- %s (%d bytes)\n", file.Name(), info.Size()))
		}
	}

	return result.String()
}

// In your main.go or app.go file
func (a *App) Quit(ctx context.Context) {
	if a.watcher != nil {
		a.watcher.Close()
	}
	runtime.Quit(ctx)
}

func main() {
	// Default to a directory in the current working directory
	svgDir := "./svg"

	if len(os.Args) > 1 {
		svgDir = os.Args[1]
	}

	// Log starting info
	log.Printf("Starting application with SVG directory: %s", svgDir)

	// Create application with the SVG directory
	app := NewApp(svgDir)

	// Create application with options
	err := wails.Run(&options.App{
		Title:     "Icons n' Wails",
		Width:     1024,
		Height:    1024,
		MinWidth:  750,
		MinHeight: 750,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		Frameless:                true,
		CSSDragProperty:          "widows",
		CSSDragValue:             "1",
		BackgroundColour:         &options.RGBA{R: 18, G: 18, B: 18, A: 120},
		OnStartup:                app.startup,
		OnShutdown:               app.Quit,
		EnableDefaultContextMenu: true,
		Bind: []interface{}{
			app,
		},
		Debug: options.Debug{
			OpenInspectorOnStartup: false,
		},
		Windows: &windows.Options{
			WebviewIsTransparent: true,
			WindowIsTranslucent:  false,
			DisablePinchZoom:     true,
		},
		Mac: &mac.Options{
			WebviewIsTransparent: true,
			WindowIsTranslucent:  false,
		},
	})

	if err != nil {
		log.Fatal(err)
	}
}
