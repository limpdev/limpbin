```go
package main

import (
	"embed"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"context"

	"github.com/fsnotify/fsnotify"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
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

	return &App{
		svgDir:      svgDir,
		watcher:     watcher,
		svgContents: make(map[string]string),
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// Create SVG directory if it doesn't exist
	if _, err := os.Stat(a.svgDir); os.IsNotExist(err) {
		os.MkdirAll(a.svgDir, 0755)
	}

	// Start watching the SVG directory
	err := a.watcher.Add(a.svgDir)
	if err != nil {
		log.Fatal(err)
	}

	// Initial scan of SVG files
	a.updateSvgFiles()

	// Start file watcher in a goroutine
	go a.watchFiles()
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
				a.updateSvgFiles()
				runtime.EventsEmit(a.ctx, "svg-files-updated")
			}

		case err, ok := <-a.watcher.Errors:
			if !ok {
				return
			}
			log.Printf("Watcher error: %v", err)
		}
	}
}

// updateSvgFiles scans the SVG directory and updates the file list
func (a *App) updateSvgFiles() {
	a.svgFiles = []string{}
	a.svgContents = make(map[string]string)

	filepath.Walk(a.svgDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(strings.ToLower(info.Name()), ".svg") {
			relPath, _ := filepath.Rel(a.svgDir, path)
			a.svgFiles = append(a.svgFiles, relPath)

			// Read SVG content
			content, err := os.ReadFile(path)
			if err == nil {
				a.svgContents[relPath] = string(content)
			}
		}

		return nil
	})

	// Sort files alphabetically
	sort.Strings(a.svgFiles)
}

// GetSvgFiles returns the list of SVG files
func (a *App) GetSvgFiles() []string {
	return a.svgFiles
}

// GetSvgContent returns the content of a specific SVG file
func (a *App) GetSvgContent(filename string) string {
	return a.svgContents[filename]
}

// GetIconName extracts the icon name from the filename
func (a *App) GetIconName(filename string) string {
	return strings.TrimSuffix(filepath.Base(filename), ".svg")
}

func main() {
	svgDir := "./frontend/dist/svg/devicons"
	if len(os.Args) > 1 {
		svgDir = os.Args[1]
	}

	// Create application with the SVG directory
	app := NewApp(svgDir)

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "theGallery",
		Width:  1024,
		Height: 1024,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 13, G: 17, B: 23, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		log.Fatal(err)
	}
}

```

- and app.js

```javascript
// frontend/src/app.js
// Connect to Go backend
import {GetSvgFiles, GetSvgContent, GetIconName} from '../wailsjs/go/main/App';
import {EventsOn} from '../wailsjs/runtime/runtime';

// DOM elements
const gallery = document.getElementById('gallery');
const codeCopiedMessage = document.getElementById('codeCopied');

// Update gallery with SVG files
async function updateGallery() {
    const files = await GetSvgFiles();
    
    gallery.innerHTML = '';
    
    for (const file of files) {
        const iconName = await GetIconName(file);
        const svgContent = await GetSvgContent(file);
        
        const iconWrapper = document.createElement('div');
        iconWrapper.className = 'icon-wrapper';
        iconWrapper.innerHTML = `
            ${svgContent}
            <div class="icon-name">${iconName}</div>
        `;
        
        iconWrapper.addEventListener('click', function(event) {
            copySvgCode(this, svgContent);
            
            // Add ripple effect
            const ripple = document.createElement('span');
            ripple.classList.add('ripple');
            const rect = this.getBoundingClientRect();
            const size = Math.max(rect.width, rect.height);
            ripple.style.width = ripple.style.height = `${size}px`;
            ripple.style.left = `${event.clientX - rect.left - size / 2}px`;
            ripple.style.top = `${event.clientY - rect.top - size / 2}px`;
            this.appendChild(ripple);
            
            setTimeout(() => ripple.remove(), 1500);
        });
        
        gallery.appendChild(iconWrapper);
    }
}

// Copy SVG code to clipboard
function copySvgCode(_iconWrapper, svgCode) {
    navigator.clipboard.writeText(svgCode).then(() => {
        showCopiedMessage();
    }).catch(err => {
        console.error('Failed to copy SVG code: ', err);
        alert('Failed to copy SVG code. Please try again.');
    });
}

// Show copied message
function showCopiedMessage() {
    codeCopiedMessage.classList.add('show');
    setTimeout(() => {
        codeCopiedMessage.classList.remove('show');
    }, 1500);
}

// Listen for updates from Go backend
EventsOn('svg-files-updated', () => {
    updateGallery();
});

// Initial load
document.addEventListener('DOMContentLoaded', () => {
    updateGallery();
});

```