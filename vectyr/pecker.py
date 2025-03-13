#!/usr/bin/env python3
import os
import time
import http.server
import socketserver
import webbrowser
from watchdog.observers import Observer
from watchdog.events import FileSystemEventHandler

# Configuration
SVG_DIR = "svg"  # Directory containing SVG icons
HTML_TEMPLATE = "index.html"  # HTML template file
PORT = 8000  # Port for the local server

# Create directory for SVG icons if it doesn't exist
os.makedirs(SVG_DIR, exist_ok=True)

# Create initial HTML template if it doesn't exist
if not os.path.exists(HTML_TEMPLATE):
    with open(HTML_TEMPLATE, "w") as f:
        f.write("""<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>SVG Icon Gallery</title>
    <style>
        body {
            font-family: system-ui, -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            max-width: 1200px;
            margin: 0 auto;
            padding: 20px;
            background-color: #f7f7f7;
        }
        h1 {
            text-align: center;
            margin-bottom: 30px;
            color: #333;
        }
        .gallery {
            display: grid;
            grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
            gap: 20px;
            margin-top: 20px;
        }
        .icon-wrapper {
            background-color: white;
            border-radius: 8px;
            padding: 20px;
            display: flex;
            flex-direction: column;
            align-items: center;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
            transition: transform 0.2s, box-shadow 0.2s;
            cursor: pointer;
        }
        .icon-wrapper:hover {
            transform: translateY(-5px);
            box-shadow: 0 4px 8px rgba(0,0,0,0.15);
        }
        .icon-wrapper svg {
            width: 48px;
            height: 48px;
            margin-bottom: 10px;
        }
        .icon-name {
            font-size: 14px;
            color: #555;
            text-align: center;
            word-break: break-word;
        }
        #codeCopied {
            position: fixed;
            bottom: 20px;
            left: 50%;
            transform: translateX(-50%);
            background-color: #333;
            color: white;
            padding: 10px 20px;
            border-radius: 4px;
            opacity: 0;
            transition: opacity 0.3s;
        }
        #codeCopied.show {
            opacity: 1;
        }
    </style>
</head>
<body>
    <h1>SVG Icon Gallery</h1>
    <p>Click on any icon to copy its SVG code to your clipboard.</p>
    
    <div class="gallery">
        <!-- Icons will be inserted here dynamically -->
    </div>

    <div id="codeCopied" class="code-copied-message">
        SVG Code Copied to Clipboard!
    </div>

    <script>
        function copySvgCode(iconWrapper) {
            const svgCode = iconWrapper.querySelector('svg').outerHTML;
            navigator.clipboard.writeText(svgCode).then(() => {
                showCopiedMessage();
            }).catch(err => {
                console.error('Failed to copy SVG code: ', err);
                alert('Failed to copy SVG code. Please try again.');
            });
        }

        function showCopiedMessage() {
            const messageElement = document.getElementById('codeCopied');
            messageElement.classList.add('show');
            setTimeout(() => {
                messageElement.classList.remove('show');
            }, 1500); // Hide message after 1.5 seconds
        }
    </script>
</body>
</html>""")
    print(f"Created initial {HTML_TEMPLATE}")

class SVGHandler(FileSystemEventHandler):
    def __init__(self):
        self.update_html()
    
    def on_any_event(self, event):
        # Only respond to .svg files
        if event.is_directory or not event.src_path.endswith('.svg'):
            return
        
        print(f"Detected change in {event.src_path}")
        self.update_html()
    
    def update_html(self):
        # Read all SVG files
        svg_files = [f for f in os.listdir(SVG_DIR) if f.endswith('.svg')]
        svg_files.sort()
        
        # Read each SVG file content
        icons = []
        for file_name in svg_files:
            file_path = os.path.join(SVG_DIR, file_name)
            try:
                with open(file_path, 'r') as f:
                    svg_content = f.read()
                icon_name = os.path.splitext(file_name)[0]
                icons.append((icon_name, svg_content))
            except Exception as e:
                print(f"Error reading {file_path}: {e}")
        
        # Generate gallery content
        gallery_content = ""
        for icon_name, svg_content in icons:
            gallery_content += f"""        <div class="icon-wrapper" onclick="copySvgCode(this)">
            {svg_content}
            <div class="icon-name">{icon_name}</div>
        </div>\n"""
        
        # Read template and insert icons
        try:
            with open(HTML_TEMPLATE, 'r') as f:
                html_content = f.read()
            
            # Replace gallery section
            gallery_start = html_content.find('<div class="gallery">')
            gallery_end = html_content.find('</div>', gallery_start)
            
            # Find the position of the closing tag for the gallery div
            closing_tag_pos = gallery_end + 6
            
            new_html = (
                html_content[:gallery_start + len('<div class="gallery">')] +
                "\n" + gallery_content +
                "    " + html_content[gallery_end:]
            )
            
            # Write updated HTML
            with open(HTML_TEMPLATE, 'w') as f:
                f.write(new_html)
            
            print(f"Updated {HTML_TEMPLATE} with {len(icons)} icons")
        except Exception as e:
            print(f"Error updating HTML: {e}")

class SimpleHTTPRequestHandler(http.server.SimpleHTTPRequestHandler):
    def log_message(self, format, *args):
        # Customize logging to be less verbose
        if args[0].startswith('GET') and args[1] == '200':
            return
        super().log_message(format, *args)

def start_server():
    # Set up HTTP server
    handler = SimpleHTTPRequestHandler
    httpd = socketserver.TCPServer(("", PORT), handler)
    
    print(f"Serving at http://localhost:{PORT}")
    print(f"To view the gallery, open a browser and navigate to http://localhost:{PORT}/{HTML_TEMPLATE}")
    print(f"Place SVG files in the '{SVG_DIR}' directory to display them in the gallery")
    print("Press Ctrl+C to stop the server")
    
    # Open browser automatically
    webbrowser.open(f"http://localhost:{PORT}/{HTML_TEMPLATE}")
    
    # Start server
    try:
        httpd.serve_forever()
    except KeyboardInterrupt:
        print("\nShutting down server...")
        httpd.shutdown()
        httpd.server_close()

def main():
    # Set up file system watcher
    event_handler = SVGHandler()
    observer = Observer()
    observer.schedule(event_handler, SVG_DIR, recursive=False)
    observer.start()
    
    # Start HTTP server
    start_server()
    
    # Clean up
    observer.stop()
    observer.join()

if __name__ == "__main__":
    main()