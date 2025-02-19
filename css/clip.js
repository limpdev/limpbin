/* CLIPBOARD LOGIC - Intended for codeblock copy buttons */
document.addEventListener('DOMContentLoaded', function() {

    document.querySelectorAll('.pre-wrapper').forEach(container => {
        const copyButton = container.querySelector('.copybtn');
        const codeBlock = container.querySelector('pre code') || container.querySelector('pre');

        copyButton.addEventListener('click', () => {
            const codeToCopy = codeBlock.textContent;

            navigator.clipboard.writeText(codeToCopy)
              .then(() => {
                copyButton.textContent = 'Copied!';
                setTimeout(() => {
                    copyButton.textContent = 'Copy Code';
                }, 2000);
               })
              .catch(err => {
                console.error('Failed to copy', err);
                copyButton.textContent = 'Copy Failed';
                setTimeout(() => {
                    copyButton.textContent = 'Copy Code';
                }, 2000);
               });
        });
    });
});

// ADDING BUTTON RECURSIVELY !!!!!!!!!
// Function to create and insert the copy button (now with SVG)
function addCopyButtonToCodeBlocks() {
  const codeBlocks = document.querySelectorAll('pre > code, pre');

  codeBlocks.forEach((codeBlock) => {
    if (codeBlock.parentNode.querySelector('.copybtn')) {
      return;
    }

    // 1. Create the SVG element
    const svg = document.createElementNS("http://www.w3.org/2000/svg", "svg");
    svg.setAttribute("viewBox", "0 0 32 32"); // Set the viewbox (adjust as needed)
    svg.setAttribute("width", "24");  // Initial size (adjust as needed)
    svg.setAttribute("height", "24"); // Initial size (adjust as needed)
    svg.setAttribute("fill", "currentColor"); // Use the current text color
    svg.style.pointerEvents = "none"; // Prevent SVG from capturing clicks


    // 2. Create the path element(s) within the SVG
    const path1 = document.createElementNS("http://www.w3.org/2000/svg", "path");
    path1.setAttribute("d", "M18.605 2.022v0zM18.605 2.022l-2.256 11.856 8.174 0.027-11.127 16.072 2.257-13.043-8.174-0.029zM18.606 0.023c-0.054 0-0.108 0.002-0.161 0.006-0.353 0.028-0.587 0.147-0.864 0.333-0.154 0.102-0.295 0.228-0.419 0.373-0.037 0.043-0.071 0.088-0.103 0.134l-11.207 14.832c-0.442 0.607-0.508 1.407-0.168 2.076s1.026 1.093 1.779 1.099l5.773 0.042-1.815 10.694c-0.172 0.919 0.318 1.835 1.18 2.204 0.257 0.11 0.527 0.163 0.793 0.163 0.629 0 1.145-0.294 1.533-0.825l11.22-16.072c0.442-0.607 0.507-1.408 0.168-2.076-0.34-0.669-1.026-1.093-1.779-1.098l-5.773-0.010 1.796-9.402c0.038-0.151 0.057-0.308 0.057-0.47 0-1.082-0.861-1.964-1.939-1.999-0.024-0.001-0.047-0.001-0.071-0.001v0z");
    svg.appendChild(path1);

    // Add a second path element to give the classic double paper look.  If you only want the one you described, comment this part out.
    const path2 = document.createElementNS("http://www.w3.org/2000/svg", "path");
    path2.setAttribute("d", "M18.605 2.022v0zM18.605 2.022l-2.256 11.856 8.174 0.027-11.127 16.072 2.257-13.043-8.174-0.029zM18.606 0.023c-0.054 0-0.108 0.002-0.161 0.006-0.353 0.028-0.587 0.147-0.864 0.333-0.154 0.102-0.295 0.228-0.419 0.373-0.037 0.043-0.071 0.088-0.103 0.134l-11.207 14.832c-0.442 0.607-0.508 1.407-0.168 2.076s1.026 1.093 1.779 1.099l5.773 0.042-1.815 10.694c-0.172 0.919 0.318 1.835 1.18 2.204 0.257 0.11 0.527 0.163 0.793 0.163 0.629 0 1.145-0.294 1.533-0.825l11.22-16.072c0.442-0.607 0.507-1.408 0.168-2.076-0.34-0.669-1.026-1.093-1.779-1.098l-5.773-0.010 1.796-9.402c0.038-0.151 0.057-0.308 0.057-0.47 0-1.082-0.861-1.964-1.939-1.999-0.024-0.001-0.047-0.001-0.071-0.001v0z");  // Adjust the path as needed
    svg.appendChild(path2);

    // 3. Create the button element
    const copyButton = document.createElement('button');
    copyButton.className = 'copybtn';
    copyButton.type = 'button';
    copyButton.setAttribute('aria-label', 'Copy code');
    copyButton.appendChild(svg); // Append the SVG to the button


    // 4. Create a wrapper div
    const wrapper = document.createElement('div');
    wrapper.style.position = 'relative';


    //Helper function for updating svg fill
    function updateSvgFill(color) {
      const paths = copyButton.querySelectorAll('path');
        paths.forEach(path => {
            path.setAttribute('fill', color);
      });
    }

    // Mouseover/mouseout events for visual feedback (adjust styles as needed)
    copyButton.addEventListener('mouseover', () => {
      //copyButton.style.backgroundColor = 'lightgray';
      updateSvgFill('blue'); // Example: Change fill color on hover
    });
    copyButton.addEventListener('mouseout', () => {
       if(copyButton.dataset.copied !== 'true'){
          //copyButton.style.backgroundColor = ''; // Reset background color
          updateSvgFill('currentColor');  // Reset fill color
       }
    });


    // 5. Click event handler (same as before, but with visual feedback adapted for SVG)
    copyButton.addEventListener('click', () => {
        const codeToCopy = codeBlock.textContent;
        if (navigator.clipboard) {
            navigator.clipboard.writeText(codeToCopy)
                .then(() => {
                    // Success! Provide visual feedback (change SVG color, for example)
                    updateSvgFill('green');
                    copyButton.dataset.copied = 'true'; // Store copied state

                    setTimeout(() => {
                      updateSvgFill('currentColor');
                      copyButton.dataset.copied = 'false';
                    }, 2000);
                })
                .catch((err) => {
                  console.error('Failed to copy:', err);
                  fallbackCopyTextToClipboard(codeToCopy, copyButton);
                });
        } else
        {
            fallbackCopyTextToClipboard(codeToCopy, copyButton);
        }
    });



    // 6. Insert elements (same as before)
      if (codeBlock.tagName === 'CODE' && codeBlock.parentNode.tagName === 'PRE') {
        const preElement = codeBlock.parentNode;
        preElement.parentNode.insertBefore(wrapper, preElement); // Insert wrapper *before* the <pre>
        wrapper.appendChild(preElement);      // Move the <pre> *into* the wrapper
    }
    //If the code block is just a pre tag
    else if(codeBlock.tagName === 'PRE'){
        codeBlock.parentNode.insertBefore(wrapper, codeBlock); // Insert wrapper *before* the <pre>
        wrapper.appendChild(codeBlock);      // Move the <pre> *into* the wrapper
    }
     // Otherwise, if it's a standalone <code>, insert the wrapper around the <code>
    else {
        codeBlock.parentNode.insertBefore(wrapper, codeBlock); // Insert wrapper *before* the <code>
        wrapper.appendChild(codeBlock);          // Move the <code> *into* the wrapper
    }

    wrapper.appendChild(copyButton); // Add the button to the wrapper
  });
}

// Call the function when the DOM is fully loaded
document.addEventListener('DOMContentLoaded', addCopyButtonToCodeBlocks);

//  Optionally, if you're using a framework that dynamically loads content (e.g., React, Vue, Angular),
//  you might need to call `addCopyButtonToCodeBlocks` after the content is loaded.  A MutationObserver
//  is the best general-purpose solution for this:

const observer = new MutationObserver((mutations) => {
  addCopyButtonToCodeBlocks();
});

// Start observing the document body for changes.  Adjust the target node as needed.
observer.observe(document.body, {
  childList: true,  // Watch for additions/removals of child nodes
  subtree: true     // Watch all descendants (not just direct children)
});