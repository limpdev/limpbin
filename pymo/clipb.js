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

    const processLanguageLabels = () => {
        const pres = document.querySelectorAll('pre[class*="language-"]');
        console.log(`Processing language labels for ${pres.length} pre elements`);

        pres.forEach((pre, index) => {
            // Get the language from the class name
            let language = null;
            const classNames = pre.className.split(' ');

            for (const className of classNames) {
                if (className.startsWith("language-")) {
                    language = className.split("-")[1];
                    break;
                }
            }

            if (language) {
                console.log(`Found language ${language} for pre ${index + 1}`);

                // Create the language label
                const label = document.createElement("div");
                label.className = "language-label";
                label.textContent = language;

                // Check if pre is already wrapped
                let wrapper;
                if (pre.parentElement && pre.parentElement.classList.contains("code-wrapper")) {
                    wrapper = pre.parentElement;
                } else {
                    // Create and insert wrapper
                    wrapper = document.createElement("div");
                    wrapper.className = "code-wrapper";
                    pre.parentNode.insertBefore(wrapper, pre);
                    wrapper.appendChild(pre);
                }

                // Add the language label to the wrapper if not already present
                if (!wrapper.querySelector('.language-label')) {
                    wrapper.appendChild(label);
                }
            }
        });
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