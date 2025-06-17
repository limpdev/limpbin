# QuickLook's Markdown Plugin

> [!NOTE]
> This plugin can be almost completely refactored, as it is a miniature `WebView` window. As a result, each markdown page is really a **tmp** HTML file and behaves as such.

```
Folder PATH listing

%USERPROFILE%\%APPDATA%\Local\pooi.moe\QuickLook\QuickLook.Plugin.MarkdownViewer
|   ##.md
|   .version
|   clipb.css
|   clipb.js
|   md2html.html
|   newtree.md
|   original_md2html.html
|   
+---css
|       github-markdown.css
|       original_github-markdown.css
|       
+---highlight.js
|   |   highlight.min.js
|   |   
|   \---styles
|           github-dark.min.css
|           github.min.css
|           
+---js
|       markdown-it.min.js
|       markdownItAnchor.umd.js
|       
\---mathjax
        tex-mml-svg.js
        
```

> Both iterations of **clipb** are also refactored to accomodate the HTML formatting of **md2html.html**.

### Refactoring Copy Buttons

**SOLUTION**: `clipb.js` changed to reliably handle the code block structure in your HTML example. The issue appears to be that the current implementation doesn't properly account for the pre/code structure with highlight classes that appears in your example.

```javascript
// clipb.js
// Module for adding copy buttons to code blocks
const ClipbModule = (() => {
    // Function to create and return the copy button
    const createCopyButton = () => {
        const copyButton = document.createElement("button");
        const dBolt = "M4 14L14 3v7h6L10 21v-7z";
        copyButton.className = "copy-button";
        copyButton.innerHTML = `<svg xmlns="http://www.w3.org/2000/svg" height="24" viewBox="0 0 24 24" width="24"><path d="${dBolt}"/></svg>`;
        return copyButton;
    };

    // Function to add a copy button to a single code block
    const addCopyButton = (codeBlock) => {
        // Get the parent pre element
        const preElement = codeBlock.closest("pre");
        if (!preElement) return; // Skip if not inside a pre element

        // Check if already wrapped
        if (preElement.parentElement.classList.contains("code-wrapper")) return;

        // Create a wrapper div for the code block
        const wrapper = document.createElement("div");
        wrapper.className = "code-wrapper";

        // Create the copy button
        const button = createCopyButton();

        // Insert the wrapper before the pre element
        preElement.parentNode.insertBefore(wrapper, preElement);

        // Move the pre element inside the wrapper
        wrapper.appendChild(preElement);

        // Add the button to the wrapper
        wrapper.insertBefore(button, preElement);

        // Add event listener for copy functionality
        button.addEventListener("click", () => {
            // Get text content, handling highlighted code
            const textToCopy = codeBlock.textContent || codeBlock.innerText;

            navigator.clipboard.writeText(textToCopy).then(
                () => {
                    const successSVG = `
            <svg viewBox="0 0 24 24" width="1.5em" height="1.5em" fill="green">
              <path d="M10 2a3 3 0 0 0-2.83 2H6a3 3 0 0 0-3 3v12a3 3 0 0 0 3 3h12a3 3 0 0 0 3-3V7a3 3 0 0 0-3-3h-1.17A3 3 0 0 0 14 2zM9 5a1 1 0 0 1 1-1h4a1 1 0 1 1 0 2h-4a1 1 0 0 1-1-1m6.78 6.625a1 1 0 1 0-1.56-1.25l-3.303 4.128l-1.21-1.21a1 1 0 0 0-1.414 1.414l2 2a1 1 0 0 0 1.488-.082l4-5z"></path>
            </svg>
          `;
                    button.innerHTML = successSVG; // Set success SVG

                    setTimeout(() => {
                        const defaultSVG = `
              <svg viewBox="0 0 24 24" width="1.5em" height="1.5em" fill="currentColor">
                <path d="M4 14L14 3v7h6L10 21v-7z"></path>
              </svg>
            `;
                        button.innerHTML = defaultSVG; // Revert to default SVG
                    }, 2000);
                },
                (err) => {
                    console.error("Could not copy text: ", err);
                }
            );
        });
    };

    // Function to add copy buttons to all code blocks on the page
    const addCopyButtons = () => {
        // Target all code blocks inside pre elements, including those with hljs class
        const codeBlocks = document.querySelectorAll("pre code, pre.hljs code");
        codeBlocks.forEach(addCopyButton);
    };

    // MutationObserver callback to handle dynamically added code blocks
    const handleMutations = (mutations) => {
        mutations.forEach((mutation) => {
            mutation.addedNodes.forEach((node) => {
                if (node.nodeType === 1) { // Element node
                    // Check if the node itself is a pre with code
                    if (node.matches("pre") && node.querySelector("code")) {
                        addCopyButton(node.querySelector("code"));
                    }
                    // Check for any code blocks within the added node
                    else {
                        const nestedCodeBlocks = node.querySelectorAll("pre code, pre.hljs code");
                        nestedCodeBlocks.forEach(addCopyButton);
                    }
                }
            });
        });
    };

    // Initialize the module
    const init = () => {
        // Wait for the DOM to be fully loaded
        if (document.readyState === "loading") {
            document.addEventListener("DOMContentLoaded", addCopyButtons);
        } else {
            addCopyButtons();
        }

        // Set up observer for dynamically added elements
        const observer = new MutationObserver(handleMutations);
        observer.observe(document.body, {
            childList: true,
            subtree: true,
        });
    };

    // Public API
    return { init };
})();

// Auto-initialize if script is loaded directly
if (typeof window !== "undefined") {
    window.addEventListener("DOMContentLoaded", ClipbModule.init);
}

```

> `Language-Label` javascript has been removed, as highlightJS works differently and required a total rewrite. Not really feeling like doing that right now tbh...

Here are the key improvements I made to the code:

1. **Better selection of code blocks**: Now selecting both regular `pre code` elements and those with syntax highlighting classes (`pre.hljs code`)

2. **Improved DOM manipulation order**:
   - First create the wrapper
   - Insert the wrapper before the pre element
   - Move the pre element inside the wrapper
   - Add the button to the wrapper before the pre element

3. **Better text content extraction**: Using `textContent || innerText` to properly handle the highlighted code with nested span elements

4. **Added DOM ready check**: Ensuring the code runs after the DOM is fully loaded

5. **More robust mutation observer**: Improved handling of dynamically added code blocks, including checking if a newly added node is itself a pre element with code

6. **Auto-initialization**: Added auto-init when the DOM is ready if the script is loaded directly

---

### CSS Improvements: Transitions... How Do They Work?

The key to creating a smooth, centered button growth on hover is to use a combination of `transform: scale()` instead of changing font or element sizes directly. This keeps the button centered during the transition without affecting the layout around it.

Here's how you can achieve that clean, centered growth effect:

```css
.copy-button {
    position: absolute;
    top: 8px;
    right: 8px;
    display: inline-flex;
    align-items: center;
    justify-content: center;
    padding: 6px;
    background: #151b23;
    border: 1px solid #151b23;
    border-radius: 4px;
    color: #7b98da;
    opacity: 0.6;
    cursor: pointer;
    z-index: 10;
    transition: all 0.3s ease;
    /* Set transform-origin to center for centered scaling */
    transform-origin: center;
}

.copy-button:hover {
    color: #f5d49e;
    border-color: transparent;
    opacity: 1.0;
    /* Scale the button up by 20% */
    transform: scale(1.2);
}
```

The important elements here are:

1. **Use `transform: scale()`** - This grows the button from its center point rather than changing dimensions that would affect layout flow

2. **Set `transform-origin: center`** - This ensures the scaling happens from the center of the button (though center is the default for scale)

3. **Remove font-size changes** - Instead of changing font-size (which can cause reflow), scale the entire button

4. **Adjust transition timing** - A value of 0.3s is typically smoother than 1s for hover effects (though you can adjust to taste)

This approach ensures:
- The button grows smoothly in all directions
- The user doesn't need to reposition their mouse
- The growth is centered regardless of button position
- No layout shifts occur during the animation

If you want to scale just the icon inside rather than the entire button, you could apply the transform to the SVG element instead:

```css
.copy-button svg {
    width: 19px;
    height: 19px;
    stroke: #f5d49e;
    opacity: 0.4;
    transition: all 0.3s ease;
    transform-origin: center;
}

.copy-button:hover svg {
    opacity: 0.97;
    fill: #f5d49e;
    transform: scale(1.2);
}
```

This gives you a clean, centered growth animation without affecting the position of the mouse relative to the clickable area.

> [!NOTE]
> I did change the CSS further to my liking; however, the concept of transitions and transforms are still valid.

#### Finalized and Simplified -> clipb.css

Redundancies have also been removed, transitions added...

```css
/* Copy button container */
.code-wrapper {
    position: relative;
}

.copy-button {
    position: absolute;
    top: 8px;
    right: 8px;
    display: inline-flex;
    align-items: center;
    justify-content: center;
    padding: 6px;
    background: #09090a;
    border: 1px solid #09090a;
    border-radius: 4px;
    color: #7b98da;
    opacity: 0.6;
    font-size: large;
    cursor: pointer;
    z-index: 10;
    transition: all 1s ease;
}

.copy-button:hover {
    color: #f5d49e;
    border-color: transparent;
    opacity: 1.0;
    font-size: x-large;
}

/* SVG icon styles */
.copy-button svg {
    width: 19px;
    height: 19px;
    stroke: #f5d49e;
    opacity: 0.4;
    transition: all 0.5s ease;
    transform-origin: center;
}

.copy-button:hover svg {
    opacity: 0.97;
    fill: #f5d49e;
    transform: scale(1.25);
}

/* Tooltip styles */
.copy-button[data-tooltip]::before {
    content: attr(data-tooltip);
    position: absolute;
    bottom: 100%;
    right: 0;
    margin-bottom: 8px;
    padding: 4px 8px;
    background: #1a202c;
    color: #bbbbbb;
    font-size: 9px;
    white-space: nowrap;
    border-radius: 6px;
    opacity: 0.4;
    visibility: hidden;
    transition: all 0.6s ease;
}

.copy-button[data-tooltip]:hover::before {
    opacity: 1;
    visibility: visible;
}

/* Success state */
.copy-button.copied {
    border-color: #b5f700;
    transition: all 0.6s ease;
}

.copy-button.copied svg {
    stroke: #b5f700;
}

/* Ensure pre tags have relative positioning for button placement */
pre {
    position: relative;
    padding-top: 1.25rem !important;
    padding-left: 1.75em !important;
}

.language-label {
    position: absolute;
    top: 15px;
    right: 60px;
    background: #09090a;
    font-size: 12px;
    font-family: Satoshi, Author, "SF Pro Display";
    border-radius: 0 4px 0 4px;
    opacity: 0.4;
    color: #f5d49e;
    cursor: pointer;
}

.language-label:hover {
    opacity: 0.8;
}
```