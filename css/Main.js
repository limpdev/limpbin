// clipb.js
// Module for adding copy buttons to code blocks
const ClipbModule = (() => {
  // Function to create and return the copy button
  const createCopyButton = () => {
    const copyButton = document.createElement('button')
    const dBolt = 'M4 14L14 3v7h6L10 21v-7z'
    copyButton.className = 'copy-button'
    copyButton.innerHTML = `<svg xmlns="http://www.w3.org/2000/svg" height="24" viewBox="0 0 24 24" width="24"><path d="${dBolt}"/></svg>`
    return copyButton
  }

  // Function to add a copy button to a single code block
  const addCopyButton = codeBlock => {
    // Get the parent pre element
    const preElement = codeBlock.closest('pre')
    if (!preElement) return // Skip if not inside a pre element

    // Check if already wrapped
    if (preElement.parentElement.classList.contains('code-wrapper')) return

    // Create a wrapper div for the code block
    const wrapper = document.createElement('div')
    wrapper.className = 'code-wrapper'

    // Create the copy button
    const button = createCopyButton()

    // Insert the wrapper before the pre element
    preElement.parentNode.insertBefore(wrapper, preElement)

    // Move the pre element inside the wrapper
    wrapper.appendChild(preElement)

    // Add the button to the wrapper
    wrapper.insertBefore(button, preElement)

    // Add event listener for copy functionality
    button.addEventListener('click', () => {
      // Get text content, handling highlighted code
      const textToCopy = codeBlock.textContent || codeBlock.innerText

      navigator.clipboard.writeText(textToCopy).then(
        () => {
          const successSVG = `
            <svg viewBox="0 0 24 24" width="1.5em" height="1.5em" fill="green">
              <path d="M10 2a3 3 0 0 0-2.83 2H6a3 3 0 0 0-3 3v12a3 3 0 0 0 3 3h12a3 3 0 0 0 3-3V7a3 3 0 0 0-3-3h-1.17A3 3 0 0 0 14 2zM9 5a1 1 0 0 1 1-1h4a1 1 0 1 1 0 2h-4a1 1 0 0 1-1-1m6.78 6.625a1 1 0 1 0-1.56-1.25l-3.303 4.128l-1.21-1.21a1 1 0 0 0-1.414 1.414l2 2a1 1 0 0 0 1.488-.082l4-5z"></path>
            </svg>
          `
          button.innerHTML = successSVG // Set success SVG

          setTimeout(() => {
            const defaultSVG = `
              <svg viewBox="0 0 24 24" width="1.5em" height="1.5em" fill="currentColor">
                <path d="M4 14L14 3v7h6L10 21v-7z"></path>
              </svg>
            `
            button.innerHTML = defaultSVG // Revert to default SVG
          }, 2000)
        },
        err => {
          console.error('Could not copy text: ', err)
        }
      )
    })
  }

  // Function to add copy buttons to all code blocks on the page
  const addCopyButtons = () => {
    // Target all code blocks inside pre elements, including those with hljs class
    const codeBlocks = document.querySelectorAll('pre code, pre.hljs code')
    codeBlocks.forEach(addCopyButton)
  }

  // MutationObserver callback to handle dynamically added code blocks
  const handleMutations = mutations => {
    mutations.forEach(mutation => {
      mutation.addedNodes.forEach(node => {
        if (node.nodeType === 1) {
          // Element node
          // Check if the node itself is a pre with code
          if (node.matches('pre') && node.querySelector('code')) {
            addCopyButton(node.querySelector('code'))
          }
          // Check for any code blocks within the added node
          else {
            const nestedCodeBlocks = node.querySelectorAll(
              'pre code, pre.hljs code'
            )
            nestedCodeBlocks.forEach(addCopyButton)
          }
        }
      })
    })
  }

  // Initialize the module
  const init = () => {
    // Wait for the DOM to be fully loaded
    if (document.readyState === 'loading') {
      document.addEventListener('DOMContentLoaded', addCopyButtons)
    } else {
      addCopyButtons()
    }

    // Set up observer for dynamically added elements
    const observer = new MutationObserver(handleMutations)
    observer.observe(document.body, {
      childList: true,
      subtree: true
    })
  }

  // Public API
  return { init }
})()

// Auto-initialize if script is loaded directly
if (typeof window !== 'undefined') {
  window.addEventListener('DOMContentLoaded', ClipbModule.init)
}
/**
 * Converts GitHub-style markdown callouts to HTML callouts within existing HTML blockquotes
 */
function convertMarkdownCalloutsToHtml (htmlText) {
  // Define callout types and their icons
  const calloutTypes = {
    NOTE: '<i class="fas fa-info-circle"></i>',
    TIP: '<i class="fas fa-lightbulb"></i>',
    IMPORTANT: '<i class="fas fa-exclamation-circle"></i>',
    WARNING: '<i class="fas fa-exclamation-triangle"></i>',
    CAUTION: '<i class="fas fa-fire"></i>'
  }

  // Regex to match GitHub-style callouts inside blockquotes in HTML
  // Matches: <blockquote><p>[!TYPE] ... </p></blockquote>
  const calloutRegex =
    /<blockquote>\s*<p>\s*\[!(NOTE|TIP|IMPORTANT|WARNING|CAUTION)\]\s*([\s\S]*?)<\/p>\s*<\/blockquote>/gm

  return htmlText.replace(calloutRegex, function (match, type, content) {
    // Normalize the type to handle case variations
    const normalizedType = type.toUpperCase()

    // Make sure we have a valid type, or default to NOTE
    const calloutType = Object.keys(calloutTypes).includes(normalizedType)
      ? normalizedType
      : 'NOTE'

    // Process the content - trim whitespace
    const processedContent = content.trim()

    // Build the HTML replacement
    return `<div class="callout callout-${calloutType.toLowerCase()}">
  <div class="callout-header">
    <span class="callout-icon">${calloutTypes[calloutType]}</span>
    <span class="callout-title">${calloutType}</span>
  </div>
  <div class="callout-content">
    <p>${processedContent}</p>
  </div>
</div>`
  })
}

document.addEventListener('DOMContentLoaded', function () {
  // Get the current HTML content of the body
  let bodyHTML = document.body.innerHTML

  // Convert the markdown callouts in the HTML to proper callout divs
  bodyHTML = convertMarkdownCalloutsToHtml(bodyHTML)

  // Replace the body's HTML with the converted content
  document.body.innerHTML = bodyHTML
})
// Add ripple effect for every mouse click, anywhere on the page using an SVG
document.addEventListener('click', function (e) {
  // Create a container for the ripple effect
  const rippleContainer = document.createElement('div')
  rippleContainer.style.position = 'fixed'
  rippleContainer.style.left = e.clientX - 48 + 'px' // Center the ripple at click position
  rippleContainer.style.top = e.clientY - 48 + 'px'
  rippleContainer.style.pointerEvents = 'none' // Don't interfere with further clicks
  rippleContainer.style.zIndex = '9999'

  // Create SVG element
  const svgNS = 'http://www.w3.org/2000/svg'
  const svg = document.createElementNS(svgNS, 'svg')
  svg.setAttribute('width', '96')
  svg.setAttribute('height', '96')
  svg.setAttribute('viewBox', '0 0 24 24')

  // Create circle element
  const circle = document.createElementNS(svgNS, 'circle')
  circle.setAttribute('cx', '12')
  circle.setAttribute('cy', '12')
  circle.setAttribute('r', '0')
  circle.setAttribute('fill', 'rgba(168, 168, 168, 0.5)')

  // Create animate elements
  const animateRadius = document.createElementNS(svgNS, 'animate')
  animateRadius.setAttribute('attributeName', 'r')
  animateRadius.setAttribute('calcMode', 'spline')
  animateRadius.setAttribute('dur', '0.5s')
  animateRadius.setAttribute('keySplines', '.52,.6,.25,.99')
  animateRadius.setAttribute('values', '0;11')
  animateRadius.setAttribute('fill', 'freeze')

  const animateOpacity = document.createElementNS(svgNS, 'animate')
  animateOpacity.setAttribute('attributeName', 'opacity')
  animateOpacity.setAttribute('calcMode', 'spline')
  animateOpacity.setAttribute('dur', '0.5s')
  animateOpacity.setAttribute('keySplines', '.52,.6,.25,.99')
  animateOpacity.setAttribute('values', '1;0')
  animateOpacity.setAttribute('fill', 'freeze')

  // Assemble the SVG
  circle.appendChild(animateRadius)
  circle.appendChild(animateOpacity)
  svg.appendChild(circle)
  rippleContainer.appendChild(svg)

  // Add to document
  document.body.appendChild(rippleContainer)

  // Remove after animation completes
  setTimeout(() => {
    document.body.removeChild(rippleContainer)
  }, 500) // Match the duration of the animation
})
