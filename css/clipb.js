// clipb.js

// Module for adding copy buttons to code blocks
const ClipbModule = (() => {
    // Function to create and append the copy button
    const createCopyButton = () => {
        const copyButton = document.createElement('button');
        copyButton.innerText = 'Copy';
        copyButton.className = 'copy-button';
        return copyButton;
    };

    // Function to add copy button to a single code block
    const addCopyButton = (codeBlock) => {
        const wrapper = document.createElement('div');
        wrapper.className = 'code-wrapper';

        const button = createCopyButton();
        wrapper.appendChild(button);

        if (codeBlock.tagName === 'CODE' && codeBlock.parentNode.tagName === 'PRE') {
            const preElement = codeBlock.parentNode;
            preElement.parentNode.insertBefore(wrapper, preElement);
            wrapper.appendChild(preElement);
        } else if (codeBlock.tagName === 'PRE') {
            codeBlock.parentNode.insertBefore(wrapper, codeBlock);
            wrapper.appendChild(codeBlock);
        } else {
            codeBlock.parentNode.insertBefore(wrapper, codeBlock);
            wrapper.appendChild(codeBlock);
        }
    };

    // Function to add copy buttons to all existing code blocks
    const addCopyButtonsToAll = () => {
        const codeBlocks = document.querySelectorAll('pre, code');
        codeBlocks.forEach(addCopyButton);
    };

    // MutationObserver callback to handle dynamically added code blocks
    const handleMutations = (mutations) => {
        mutations.forEach((mutation) => {
            mutation.addedNodes.forEach((node) => {
                if (node.nodeType === 1) { // Element node
                    if (node.matches('pre, code')) {
                        addCopyButton(node);
                    } else {
                        const nestedCodeBlocks = node.querySelectorAll('pre, code');
                        nestedCodeBlocks.forEach(addCopyButton);
                    }
                }
            });
        });
    };

    // Initialize the module
    const init = () => {
        addCopyButtonsToAll();

        const observer = new MutationObserver(handleMutations);
        observer.observe(document.body, {
            childList: true,
            subtree: true
        });
    };

    // Public API
    return {
        init
    };
})();

// Ensure that Prism.js has completed its initialization before running ClipbModule
document.addEventListener('DOMContentLoaded', () => {
    // Check if Prism is loaded
    if (window.Prism) {
        ClipbModule.init();
    } else {
        // Wait for Prism to load
        const prismScript = document.querySelector('script[src*="prismHL.js"]');
        if (prismScript) {
            prismScript.addEventListener('load', ClipbModule.init);
        } else {
            console.error('Prism.js script not found. ClipbModule may not function correctly.');
            ClipbModule.init(); // Attempt to initialize anyway
        }
    }
});