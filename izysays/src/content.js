import { unified } from "unified";
import remarkParse from "remark-parse";
import remarkGfm from "remark-gfm";
import remarkRehype from "remark-rehype";
import rehypeHighlight from "rehype-highlight";
import rehypeStringify from "rehype-stringify";
import a11yEmoji from '@fec/remark-a11y-emoji';
import remarkFlexibleContainers from "remark-flexible-containers";
import remarkFlexibleParagraphs from "remark-flexible-paragraphs";
import remarkFrontmatter from 'remark-frontmatter'
import remarkDirective from "remark-directive";
import remarkTextr from 'remark-textr'
import remarkToc from 'remark-toc'
import sectionize from 'remark-sectionize'
import remarkGithubAdmonitionsToDirectives from "remark-github-admonitions-to-directives";

// Language to Font Awesome icon mapping
const languageIcons = {
    javascript: 'fab fa-js-square',
    js: 'fab fa-js-square',
    typescript: 'fab fa-js-square', // TypeScript uses JS icon
    ts: 'fab fa-js-square',
    python: 'fab fa-python',
    py: 'fab fa-python',
    java: 'fab fa-java',
    php: 'fab fa-php',
    html: 'fab fa-html5',
    css: 'fab fa-css3-alt',
    scss: 'fab fa-sass',
    sass: 'fab fa-sass',
    less: 'fab fa-less',
    react: 'fab fa-react',
    jsx: 'fab fa-react',
    vue: 'fab fa-vuejs',
    angular: 'fab fa-angular',
    node: 'fab fa-node-js',
    nodejs: 'fab fa-node-js',
    npm: 'fab fa-npm',
    yarn: 'fab fa-yarn',
    docker: 'fab fa-docker',
    git: 'fab fa-git-alt',
    github: 'fab fa-github',
    gitlab: 'fab fa-gitlab',
    bitbucket: 'fab fa-bitbucket',
    go: 'fas fa-code', // Go doesn't have a specific FA icon
    golang: 'fas fa-code',
    rust: 'fas fa-cog',
    cpp: 'fas fa-code',
    'c++': 'fas fa-code',
    c: 'fas fa-code',
    csharp: 'fas fa-code',
    'c#': 'fas fa-code',
    swift: 'fab fa-swift',
    kotlin: 'fas fa-code',
    ruby: 'fas fa-gem',
    rb: 'fas fa-gem',
    shell: 'fas fa-terminal',
    bash: 'fas fa-terminal',
    sh: 'fas fa-terminal',
    powershell: 'fas fa-terminal',
    sql: 'fas fa-database',
    mysql: 'fas fa-database',
    postgresql: 'fas fa-database',
    mongodb: 'fas fa-database',
    json: 'fas fa-file-code',
    xml: 'fas fa-code',
    yaml: 'fas fa-file-code',
    yml: 'fas fa-file-code',
    toml: 'fas fa-file-code',
    ini: 'fas fa-file-code',
    markdown: 'fab fa-markdown',
    md: 'fab fa-markdown',
    text: 'fas fa-file-alt',
    txt: 'fas fa-file-alt',
    // Add more as needed
    default: 'fas fa-code'
};

// Custom rehype plugin to add language icons to code blocks
function rehypeCodeLanguageIcons() {
    return (tree) => {
        const visit = (node, callback) => {
            callback(node);
            if (node.children) {
                node.children.forEach(child => visit(child, callback));
            }
        };

        visit(tree, (node) => {
            if (node.type === 'element' && node.tagName === 'pre') {
                const codeElement = node.children.find(child =>
                    child.type === 'element' && child.tagName === 'code'
                );

                if (codeElement && codeElement.properties && codeElement.properties.className) {
                    // Extract language from class name (format: language-xxx)
                    const languageClass = codeElement.properties.className.find(cls =>
                        cls.startsWith('language-')
                    );

                    if (languageClass) {
                        const language = languageClass.replace('language-', '').toLowerCase();
                        const iconClass = languageIcons[language] || languageIcons.default;

                        // Add language class to pre element for styling
                        if (!node.properties) node.properties = {};
                        if (!node.properties.className) node.properties.className = [];
                        node.properties.className.push('has-language');
                        node.properties['data-language'] = language;

                        // Create language icon element
                        const languageIcon = {
                            type: 'element',
                            tagName: 'div',
                            properties: {
                                className: ['language-icon']
                            },
                            children: [
                                {
                                    type: 'element',
                                    tagName: 'i',
                                    properties: {
                                        className: iconClass.split(' '),
                                        title: language.toUpperCase()
                                    },
                                    children: []
                                }
                            ]
                        };

                        // Add the icon as the first child of pre
                        node.children.unshift(languageIcon);
                    }
                }
            }
        });
    };
}

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
			.use(remarkParse)
			.use(rehypeStringify, { allowDangerousHtml: true })
			.use(remarkGfm)
			.use(a11yEmoji)
			.use(remarkRehype)
			.use(rehypeHighlight)
			.use(rehypeCodeLanguageIcons)
			.use(remarkFlexibleContainers)
			.use(remarkFrontmatter, ['yaml', 'toml'])
			.use(remarkFlexibleParagraphs)
			.use(remarkGithubAdmonitionsToDirectives)
			.use(remarkDirective)
            .use(remarkTextr, {plugins: [ellipses]})
            .use(sectionize)
            .use(remarkToc)

        const file = await processor.process(rawMarkdown);
        const renderedHtml = String(file);

        // Stop the browser's default rendering
        if (window.stop) {
            window.stop();
        }

        /**
        * Replace triple dots with ellipses.
        *
        * @type {TextrPlugin}
        */
        function ellipses(input) {
         return input.replace(/\.{3}/gim, '…')
        }

        // Get the original title or create one from the filename
        const title = document.title || window.location.pathname.split("/").pop() || "Markdown Preview";

        // Prepare the new document structure
        document.documentElement.innerHTML = `
            <head>
                <meta charset="UTF-8">
                <link rel="shortcut-icon" type="image/png" href="data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAACAAAAAgCAYAAABzenr0AAAFZ2lUWHRYTUw6Y29tLmFkb2JlLnhtcAAAAAAAPD94cGFja2V0IGJlZ2luPSLvu78iIGlkPSJXNU0wTXBDZWhpSHpyZVN6TlRjemtjOWQiPz4KPHg6eG1wbWV0YSB4bWxuczp4PSJhZG9iZTpuczptZXRhLyIgeDp4bXB0az0iWE1QIENvcmUgNS41LjAiPgogPHJkZjpSREYgeG1sbnM6cmRmPSJodHRwOi8vd3d3LnczLm9yZy8xOTk5LzAyLzIyLXJkZi1zeW50YXgtbnMjIj4KICA8cmRmOkRlc2NyaXB0aW9uIHJkZjphYm91dD0iIgogICAgeG1sbnM6ZXhpZj0iaHR0cDovL25zLmFkb2JlLmNvbS9leGlmLzEuMC8iCiAgICB4bWxuczpwaG90b3Nob3A9Imh0dHA6Ly9ucy5hZG9iZS5jb20vcGhvdG9zaG9wLzEuMC8iCiAgICB4bWxuczp0aWZmPSJodHRwOi8vbnMuYWRvYmUuY29tL3RpZmYvMS4wLyIKICAgIHhtbG5zOnhtcD0iaHR0cDovL25zLmFkb2JlLmNvbS94YXAvMS4wLyIKICAgIHhtbG5zOnhtcE1NPSJodHRwOi8vbnMuYWRvYmUuY29tL3hhcC8xLjAvbW0vIgogICAgeG1sbnM6c3RFdnQ9Imh0dHA6Ly9ucy5hZG9iZS5jb20veGFwLzEuMC9zVHlwZS9SZXNvdXJjZUV2ZW50IyIKICAgZXhpZjpDb2xvclNwYWNlPSIxIgogICBleGlmOlBpeGVsWERpbWVuc2lvbj0iMzIiCiAgIGV4aWY6UGl4ZWxZRGltZW5zaW9uPSIzMiIKICAgcGhvdG9zaG9wOkNvbG9yTW9kZT0iMyIKICAgcGhvdG9zaG9wOklDQ1Byb2ZpbGU9InNSR0IgSUVDNjE5NjYtMi4xIgogICB0aWZmOkltYWdlTGVuZ3RoPSIzMiIKICAgdGlmZjpJbWFnZVdpZHRoPSIzMiIKICAgdGlmZjpSZXNvbHV0aW9uVW5pdD0iMiIKICAgdGlmZjpYUmVzb2x1dGlvbj0iNzIvMSIKICAgdGlmZjpZUmVzb2x1dGlvbj0iNzIvMSIKICAgeG1wOkNyZWF0b3JUb29sPSJBZG9iZSBF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"/>
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
