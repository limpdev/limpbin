# Izysays - Chromium Markdown Renderer

## `dist\manifest.json`

```json
{
	"manifest_version": 3,
	"name": "Izysays",
	"version": "1.0",
	"description": "Automatically renders local Markdown files as HTML using remark/rehype.",
	"permissions": [],
	"host_permissions": ["file://*/*"],
	"content_scripts": [
		{
			"matches": ["file://*/*.md", "file://*/*.markdown", "file://*/*.mdown"],
			"js": ["content.bundle.js"],
			"run_at": "document_start"
		}
	],
	"web_accessible_resources": [
		{
			"resources": ["style.css", "fonts/*"],
			"matches": ["file://*/*"]
		}
	],
	"icons": {
		"32": "icon32.png",
		"48": "icon48.png",
		"64": "icon64.png",
		"96": "icon96.png",
		"128": "icon128.png"
	}
}
```

## `package.json`

```json
{
    "name": "izysays",
    "version": "1.0.0",
    "main": "content.js",
    "private": true,
    "scripts": {
        "build": "webpack --mode=production",
        "watch": "webpack --mode=development --watch"
    },
    "keywords": [],
    "author": "Limp Cheney",
    "license": "MIT",
    "description": "Renders local Markdown files in Chrome.",
    "dependencies": {
        "rehype-highlight": "^7.0.2",
        "rehype-stringify": "^10.0.1",
        "remark-breaks": "^4.0.0",
        "remark-code-frontmatter": "^1.0.0",
        "remark-gfm": "^4.0.1",
        "remark-parse": "^11.0.0",
        "remark-rehype": "^11.1.2",
        "unified": "^11.0.5"
    },
    "devDependencies": {
        "@fec/remark-a11y-emoji": "^4.0.2",
        "copy-webpack-plugin": "^13.0.0",
        "webpack": "^5.99.9",
        "webpack-cli": "^6.0.1"
    }
}
```

## `public\manifest.json`

```json
{
	"manifest_version": 3,
	"name": "Izysays",
	"version": "1.0",
	"description": "Automatically renders local Markdown files as HTML using remark/rehype.",
	"permissions": [],
	"host_permissions": ["file://*/*"],
	"content_scripts": [
		{
			"matches": ["file://*/*.md", "file://*/*.markdown", "file://*/*.mdown"],
			"js": ["content.bundle.js"],
			"run_at": "document_start"
		}
	],
	"web_accessible_resources": [
		{
			"resources": ["style.css", "fonts/*"],
			"matches": ["file://*/*"]
		}
	],
	"icons": {
		"32": "icon32.png",
		"48": "icon48.png",
		"64": "icon64.png",
		"96": "icon96.png",
		"128": "icon128.png"
	}
}
```

## `src\content.js`

```javascript
import { unified } from "unified";
import remarkParse from "remark-parse";
import remarkGfm from "remark-gfm";
import remarkRehype from "remark-rehype";
import rehypeHighlight from "rehype-highlight";
import rehypeStringify from "rehype-stringify";
import a11yEmoji from '@fec/remark-a11y-emoji';

async function renderMarkdown() {
    try {
        // More robust check for markdown content
        const preElement = document.querySelector("pre");

        if (!preElement || !preElement.textContent) {
            console.log("No markdown content found");
            return;
        }

        // Additional check to ensure we're dealing with a markdown file
        const isMarkdownFile = window.location.pathname.endsWith(".md") || window.location.pathname.endsWith(".markdown");

        if (!isMarkdownFile) {
            console.log("Not a markdown file");
            return;
        }

        // Hide the body to prevent flash of unstyled content
        document.body.style.display = "none";

        const rawMarkdown = preElement.textContent;

        // Set up the remark/rehype processor with error handling
        const processor = unified()
            .use(remarkParse) // Parse markdown
            .use(remarkGfm) // Support GFM (tables, etc.)
            .use(remarkRehype, { allowDangerousHtml: true }) // Turn markdown into HTML
            .use(rehypeHighlight) // Add syntax highlighting
            .use(rehypeStringify) // Convert HTML to string
            .use(a11yEmoji) // Enables Emojis

        const file = await processor.process(rawMarkdown);
        const renderedHtml = String(file);

        // Stop the browser's default rendering
        if (window.stop) {
            window.stop();
        }

        // Get the original title or create one from the filename
        const title = document.title || window.location.pathname.split("/").pop() || "Markdown Preview";

        // Prepare the new document structure
        document.documentElement.innerHTML = `
            <head>
                <meta charset="UTF-8">
                <meta name="viewport" content="width=device-width, initial-scale=1.0">
                <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.7.2/css/all.min.css" integrity="sha512-Evv84Mr4kqVGRNSgIGL/F/aIDqQb7xQ2vcrdIwxfjThSH8CSR7PBEakCr51Ck+w+/U6swU2Im1vVX0SVk9ABhg==" crossorigin="anonymous" referrerpolicy="no-referrer" />
                <title>${title}</title>
            </head>
            <body>
                <div id="markdown-content-container">
                    ${renderedHtml}
                </div>
                <script src="https://cdnjs.cloudflare.com/ajax/libs/mermaid/11.5.0/mermaid.min.js" integrity="sha512-3EZqKCkk3nMLmbrI7mfry81KH7dkzy/BoDfQrodwLQnS/RbsVlERdYP6J0oiJegRUxSOmx7Y35WNbVKSw7mipw==" crossorigin="anonymous" referrerpolicy="no-referrer"></script>
            </body>
        `;

        // Inject our stylesheet with error handling
        try {
            const link = document.createElement("link");
            link.rel = "stylesheet";
            link.href = chrome.runtime.getURL("style.css");
            link.onerror = () => console.warn("Failed to load stylesheet");
            document.head.appendChild(link);
        } catch (styleError) {
            console.warn("Error loading stylesheet:", styleError);
        }

        // Reveal the body now that it's styled
        document.body.style.display = "block";
    } catch (error) {
        console.error("Error rendering markdown:", error);

        // Restore original content on error
        document.body.style.display = "block";

        // Show error message to user
        const errorDiv = document.createElement("div");
        errorDiv.style.cssText = `
            position: fixed;
            top: 10px;
            right: 10px;
            background: #000000;
            color: white;
            padding: 10px;
            border-radius: 4px;
            z-index: 10000;
            font-family: monospace;
        `;
        errorDiv.textContent = `Markdown rendering failed: ${error.message}`;
        document.body.appendChild(errorDiv);

        // Auto-hide error after 5 seconds
        setTimeout(() => {
            if (errorDiv.parentNode) {
                errorDiv.parentNode.removeChild(errorDiv);
            }
        }, 5000);
    }
}

// Wait for DOM to be ready
if (document.readyState === "loading") {
    document.addEventListener("DOMContentLoaded", renderMarkdown);
} else {
    renderMarkdown();
}
```

## `webpack.config.js`

```javascript
const path = require("path");
const CopyPlugin = require("copy-webpack-plugin");

module.exports = {
    // Entry point for our content script
    entry: {
        content: "./src/content.js",
    },
    // Output configuration
    output: {
        path: path.resolve(__dirname, "dist"),
        filename: "[name].bundle.js",
        clean: true, // Clean the dist folder before each build
    },
    // Plugins
    plugins: [
        new CopyPlugin({
            patterns: [
                { from: "public", to: "." }, // Copies files from public/ to dist/
            ],
        }),
    ],
    // Optional: configuration for resolving modules, etc.
    resolve: {
        extensions: [".js"],
    },
};
```

