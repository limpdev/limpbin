// clipb.js
// Module for adding copy buttons to code blocks
const ClipbModule = (() => {
  // Function to create and return the copy button
  const createCopyButton = () => {
    const copyButton = document.createElement("button");
    copyButton.innerText = "📋";
    copyButton.className = "copy-button";
    return copyButton;
  };

  // Function to add a copy button to a single code block
  const addCopyButton = (codeBlock) => {
    // Avoid adding multiple buttons to the same code block
    if (codeBlock.parentElement.classList.contains("code-wrapper")) return;

    // Create a wrapper div for the code block
    const wrapper = document.createElement("div");
    wrapper.className = "code-wrapper";

    // Create the copy button
    const button = createCopyButton();
    wrapper.appendChild(button);

    // Check if the code block is inside a <pre> element
    if (codeBlock.tagName === "CODE" && codeBlock.closest("pre")) {
      const preElement = codeBlock.closest("pre");
      // Insert the wrapper before the <pre> element
      preElement.parentNode.insertBefore(wrapper, preElement);
      // Append the <pre> element to the wrapper
      wrapper.appendChild(preElement);
    } else {
      // Handle cases where the <code> is not inside a <pre>
      // For now, skip
    }
    //

    function getCodeLanguage(element) {
      // Find the nearest <pre> parent
      const preElement = document.querySelector("pre");
      if (!preElement) {
        return ""; // No <pre> found
      }
      // Get all classes of the <pre> element
      const classes = preElement.classList;
      for (let i = 0; i < classes.length; i++) {
        const className = classes[i];
        if (className.startsWith("language-")) {
          return className.slice("language-".length);
        }
      }
      return ""; // No language class found
    }

    // Add event listener for copy functionality
    button.addEventListener("click", () => {
      const textToCopy = codeBlock.innerText;
      navigator.clipboard.writeText(textToCopy).then(
        () => {
          button.innerText = "📎";
          setTimeout(() => {
            button.innerText = "📋";
          }, 2000);
        },
        (err) => {
          console.error("Could not copy text: ", err);
        },
      );
    });
  };

  // Function to add copy buttons to all code blocks on the page
  const addCopyButtons = () => {
    const codeBlocks = document.querySelectorAll("pre code");
    codeBlocks.forEach(addCopyButton);
  };

  // MutationObserver callback to handle dynamically added code blocks
  const handleMutations = (mutations) => {
    mutations.forEach((mutation) => {
      mutation.addedNodes.forEach((node) => {
        if (node.nodeType === 1) {
          // Element node
          if (node.matches("pre code")) {
            addCopyButton(node);
          } else {
            const nestedCodeBlocks = node.querySelectorAll("pre code");
            nestedCodeBlocks.forEach(addCopyButton);
          }
        }
      });
    });
  };

  // Initialize the module
  const init = () => {
    addCopyButtons();
    const observer = new MutationObserver(handleMutations);
    observer.observe(document.body, {
      childList: true,
      subtree: true,
    });
  };

  // Public API
  return { init };
})();

// Ensure that Prism.js has completed its initialization before running ClipbModule
document.addEventListener("DOMContentLoaded", () => {
  // Check if Prism is loaded
  if (window.Prism) {
    ClipbModule.init();
  } else {
    // Wait for Prism to load
    const prismScript = document.querySelector('script[src*="prismHL.js"]');
    if (prismScript) {
      prismScript.addEventListener("load", ClipbModule.init);
    } else {
      console.error("Prism.js script not found. ClipbModule may not function correctly.");
      ClipbModule.init(); // Attempt to initialize anyway
    }
  }
});

// Capture the code language specified in the pre-tags
document.addEventListener("DOMContentLoaded", () => {
  const pres = document.querySelectorAll('pre[class^="language-"]');
  pres.forEach((pre) => {
    // Get the language from the class name
    const languageClass = Array.from(pre.classList).find((className) =>
      className.startsWith("language-"),
    );
    if (languageClass) {
      // Extract the language name (everything after 'language-')
      const language = languageClass.split("-")[1];
      // Create and insert the language label
      const label = document.createElement("div");
      label.className = "language-label";
      label.textContent = language;
      // Wrap pre in container if not already wrapped
      if (!pre.parentElement.classList.contains("code-container")) {
        const container = document.createElement("div");
        container.className = "code-container";
        pre.parentNode.insertBefore(container, pre);
        container.appendChild(pre);
      }
      pre.parentElement.appendChild(label);
    }
  });
});
