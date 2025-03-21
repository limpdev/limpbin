I'll convert this Python script to Go using the Wails framework. Wails is a good choice for this project as it allows you to build desktop applications with Go and web technologies.

```go
package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

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
	ctx           context.Context
	svgDir        string
	watcher       *fsnotify.Watcher
	svgFiles      []string
	svgContents   map[string]string
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
	svgDir := "."
	if len(os.Args) > 1 {
		svgDir = os.Args[1]
	}
	
	// Create application with the SVG directory
	app := NewApp(svgDir)

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "vGallery",
		Width:  1024,
		Height: 768,
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

Now, let's create the frontend part using HTML, CSS, and JavaScript:

```html
<!-- frontend/src/index.html -->
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>vGallery</title>
    <link rel="stylesheet" href="./styles.css">
</head>
<body>
    <h1>vGallery</h1>
    <p>Place SVG files in your configured directory</p>
    <div id="gallery" class="gallery"></div>
    <div id="codeCopied" class="code-copied-message">Copied</div>
    
    <script src="./app.js"></script>
</body>
</html>

```

```css
/* frontend/src/styles.css */
body {
    display: block;
    color: white;
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
    max-width: 1200px;
    margin: 0 auto;
    padding: 20px;
    background-color: #0d1117;
}

h1 {
    text-align: center;
    margin-bottom: 30px;
    color: #91C8FF;
}

svg {
    display: block;
    margin: 0 auto;
    fill: offwhite;
    transition: 0.8s ease;
    opacity: 0.75;
}

svg:hover {
    transform: scale(1.25);
    opacity: 1;
}

.gallery {
    display: inline-flex;
    flex-wrap: wrap;
    justify-content: center;
    gap: 20px;
    margin-top: 20px;
}

.icon-wrapper {
    position: relative;
    overflow: hidden;
    padding: 20px;
    display: inline-flex;
    justify-content: center;
    box-shadow: 0 2px 4px rgba(0,0,0,0.1);
    transition: transform 0.2s, box-shadow 0.2s;
    cursor: pointer;
    background: #09090a;
    border-radius: 14px;
}

.ripple {
    position: absolute;
    border-radius: 50%;
    background: rgba(62, 62, 62, 0.3);
    transform: scale(0);
    animation: ripple 1.5s;
    pointer-events: none;
}

@keyframes ripple {
    from {
        transform: scale(0);
        opacity: 1;
    }
    to {
        transform: scale(4);
        opacity: 0;
    }
}

.icon-wrapper:hover {
    transform: translateY(-5px);
    box-shadow: 0 4px 8px rgba(0,0,0,0.15);
}

.icon-wrapper svg {
    width: 48px;
    height: 48px;
}

.icon-name {
    display: none;
    font-size: 14px;
    color: #555;
    text-align: center;
    word-break: break-word;
}

#codeCopied {
    position: fixed;
    top: 20px;
    left: 50%;
    transform: translateX(-50%);
    background-color: #367E39;
    color: white;
    padding: 10px;
    border-radius: 12px;
    opacity: 0;
    transition: opacity 0.9s ease;
}

#codeCopied.show {
    opacity: 0.75;
}

```

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
function copySvgCode(iconWrapper, svgCode) {
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

```json
{
  "name": "vGallery",
  "outputfilename": "vgallery",
  "frontend:install": "npm install",
  "frontend:build": "npm run build",
  "frontend:dev:watcher": "npm run dev",
  "frontend:dev:serverUrl": "auto",
  "author": {
    "name": "vGallery",
    "email": "user@example.com"
  }
}

```



# vGallery Wails Setup Instructions

This document outlines how to set up and run the vGallery application using Wails and Go.

## Prerequisites

1. Install Go (1.16 or later): [https://golang.org/dl/](https://golang.org/dl/)
2. Install Wails CLI:
   ```
   go install github.com/wailsapp/wails/v2/cmd/wails@latest
   ```
3. Install dependencies for Wails:
   - For Windows: GCC compiler (can be installed through MSYS2)
   - For macOS: Xcode Command Line Tools
   - For Linux: Various packages (refer to Wails documentation)

## Project Structure

Create the following structure for your project:

```
vgallery/
├── frontend/
│   ├── src/
│   │   ├── index.html
│   │   ├── styles.css
│   │   └── app.js
│   └── package.json
├── main.go
└── wails.json
```

## Setup Steps

1. Create a new Wails project:
   ```
   wails init -n vgallery -t vanilla
   ```

2. Replace the generated files with the ones provided in this guide:
   - `main.go` - The Go backend application
   - `frontend/src/index.html` - The HTML template
   - `frontend/src/styles.css` - CSS styles
   - `frontend/src/app.js` - Frontend JavaScript
   - `wails.json` - Wails configuration

3. Install frontend dependencies:
   ```
   cd vgallery/frontend
   npm install
   ```

## Running the Application

1. Build and run the application:
   ```
   cd vgallery
   wails dev
   ```

2. To build for production:
   ```
   wails build
   ```

## Using the Application

1. The application will start and automatically open in a window
2. Place SVG files in the specified directory (passed as command-line argument or defaults to current directory)
3. The gallery will automatically update when SVG files are added, modified, or removed
4. Click on any icon to copy its SVG code to the clipboard

## Command-line Arguments

You can specify the SVG directory when running the built application:
```
./vgallery /path/to/svg/directory
```


## Explanation of the Conversion

I've recreated your SVG gallery application using Go with the Wails framework. Here's how the components work together:

1. **Go Backend (`main.go`)**:
   - Uses the Wails framework to create a desktop application
   - Implements file watching with `fsnotify` (equivalent to Python's `watchdog`)
   - Provides methods for the frontend to access SVG files and their contents
   - Emits events when SVG files change

2. **Frontend**:
   - HTML structure maintaining the same look and feel as the original
   - CSS styling copied from your original template
   - JavaScript to handle UI interactions and communication with the Go backend

3. **Key Differences from the Python Version**:
   - Instead of generating HTML from a template, the Go version uses a frontend/backend architecture
   - The Go backend watches for file changes and provides data
   - The JavaScript frontend renders the SVG gallery and handles user interactions
   - Uses Wails' built-in event system instead of HTTP server

The application maintains all the functionality from the original Python version:
- Automatic update when SVG files change
- Click-to-copy SVG code
- Ripple effect on click
- Visual feedback when copying

To get started, follow the setup instructions in the provided markdown file. The Wails framework will handle bundling everything into a single executable for easy distribution.