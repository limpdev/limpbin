import os
import sys
import tempfile
import markdown2
import webview

def md_to_html(markdown_file):
    """
    Converts a markdown file to HTML content with a custom title bar.
    """
    html_template = """
    <!DOCTYPE html>
    <html lang="en">
    <head>
        <meta charset="UTF-8">
        <title>%s</title>
        <link rel="stylesheet" type="text/css" href="https://cdn.jsdelivr.net/gh/limpdev/limpbin@main/css/GithubAPI.css" />
        <link rel="stylesheet" type="text/css" href="https://cdn.jsdelivr.net/gh/limpdev/limpbin@main/css/prismHL.css" />
        <link rel="stylesheet" type="text/css" href="https://cdn.jsdelivr.net/gh/limpdev/limpbin@main/css/clipb.css" />
        <style>
            /* Custom styles for the frameless window */
            body {
                margin: 0;
                overflow: auto;
                height: 100%%;
                background-color: #0d1117;
                color: #d4d4d4;
                -webkit-app-region: drag;
            }
            /* Title Bar */
            .title-bar {
                display: flex;
                justify-content: space-between;
                align-items: center;
                background-color: #0d1117;
                color: #fff;
                padding: 0 10px;
                -webkit-app-region: drag;
                height: 26px;
            }
            .title-text {
                font-weight: 500;
                font-family: Author, Open Sans;
                white-space: nowrap;
                overflow: hidden;
                text-overflow: ellipsis;
                -webkit-app-region: drag;
            }
            .window-controls {
                display: flex;
                -webkit-app-region: drag;
            }
            .window-control-button {
                width: 18px;
                height: 18px;
                display: flex;
                justify-content: center;
                align-items: center;
                margin-left: 5px;
                border-radius: 4px;
                cursor: pointer;
                transition: background-color 0.8s;
                -webkit-app-region: drag;
            }
            .window-control-button:hover {
                background-color: rgba(255, 255, 255, 0.2);
            }
            .close-button:hover {
                background-color: #e74c3c;
                opacity: 0.6;
            }
            /* Content Area */
            .content-container {
                height: calc(100%% - 30px);
                overflow: auto;
                padding: 20px;
            }
        </style>
    </head>
    <body>
        <!-- Custom Title Bar -->
        <div class="title-bar">
            <div class="title-text">%s</div>
            <div class="window-controls">
                <div class="window-control-button minimize-button" onclick="window.pywebview.api.minimize()">
                    &#9472;
                </div>
                <div class="window-control-button maximize-button" onclick="window.pywebview.api.toggleMaximize()">
                    &#9633;
                </div>
                <div class="window-control-button close-button" onclick="window.pywebview.api.close()">
                    &#10005;
                </div>
            </div>
        </div>
        <!-- Content Area -->
        <div class="content-container">
            <article class="markdown-body">
                %s
            </article>
        </div>
        <script src="https://cdn.jsdelivr.net/gh/limpdev/limpbin@main/css/prismHL.js"></script>
        <script src="https://cdn.jsdelivr.net/gh/limpdev/limpbin@main/css/clipb.js"></script>
    </body>
    </html>
    """

    try:
        with open(markdown_file, "r", encoding="utf-8") as f:
            markdown_content = f.read()
    except Exception as e:
        return f"Error: Could not read markdown file: {e}"

    try:
        markdowner = markdown2.Markdown(extras=[
            "fenced-code-blocks",
            "tables",
            "header-ids",
            "footnotes",
            "code-friendly",
            "link-shortrefs",
            "html-classes",
            "highlightjs-lang"
        ])
        html_content = markdowner.convert(markdown_content)
    except Exception as e:
        return f"Error: Could not convert markdown to html: {e}"

    title = os.path.basename(markdown_file).split('.')[0]
    filename = os.path.basename(markdown_file)
    # Pass the title twice - once for the HTML title and once for the title bar display
    complete_html = html_template % (title, filename, html_content)
    return complete_html

# Create a simple class to handle window operations
class WindowButtons:
    def minimize(self):
        webview.windows[0].minimize()

    def toggleMaximize(self):
        if webview.windows[0].is_maximized:
            webview.windows[0].restore()
        else:
            webview.windows[0].maximize()

    def close(self):
        webview.windows[0].destroy()

def main():
    """
    Main function that handles the application flow.
    """
    if len(sys.argv) < 2:
        print("Error: No markdown file specified.")
        print("Usage: python script.py <markdown_file>")
        sys.exit(1)

    markdown_file = sys.argv[1]

    # Check if file exists
    if not os.path.isfile(markdown_file):
        print(f"Error: File '{markdown_file}' does not exist.")
        sys.exit(1)

    # Get HTML content
    html_content = md_to_html(markdown_file)

    # Create a temporary HTML file
    temp_fd, temp_path = tempfile.mkstemp(suffix='.html', prefix='md_viewer_')

    # Create buttons handler instance before any exceptions might occur
    buttons_api = WindowButtons()

    try:
        with os.fdopen(temp_fd, 'w', encoding='utf-8') as f:
            f.write(html_content)

        # Create the window with the temp HTML file and pass the js_api
        window = webview.create_window(
            title=os.path.basename(markdown_file),
            url=temp_path,
            width=1000,
            height=850,
            frameless=True,
            min_size=(500, 400),
            js_api=buttons_api,
            text_select=True,
            easy_drag=True,
            background_color='#09090a'
        )

        # Start the webview
        webview.start()
    finally:
        # Clean up temp file
        if os.path.exists(temp_path):
            try:
                os.remove(temp_path)
            except:
                pass

if __name__ == "__main__":
    main()