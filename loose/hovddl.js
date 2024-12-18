// ==UserScript==
// @name         hovddl-a
// @namespace    http://tampermonkey.net/
// @version      1.1
// @description  Downloads the full-size image under the mouse cursor when pressing '0'
// @author       Limp Dev
// @match        *://*/*
// @grant        GM_download
// ==/UserScript==

(function() {
    'use strict';

    // Function to get the element under the mouse cursor
    function getElementUnderMouse(event) {
        return document.elementFromPoint(event.clientX, event.clientY);
    }

    // Function to get full-size image URL
    function getFullSizeImageUrl(element) {
        if (!element) return null;

        // Check various common attributes for full-size image URLs
        const possibleAttributes = [
            'data-src',              // Common for lazy loading
            'data-full-src',         // Used by some sites
            'data-original',         // Used by some sites
            'data-zoom-src',         // Used for zoom/lightbox
            'data-large-src',        // Used for larger versions
            'data-high-res-src',     // High resolution version
            'data-original-src',     // Original source
            'data-1000px',           // Resolution-specific versions
            'data-raw-src',          // Raw/original version
            'src'                    // Default src attribute
        ];

        // For img elements, check all possible attributes
        if (element.tagName === 'IMG') {
            for (const attr of possibleAttributes) {
                const value = element.getAttribute(attr);
                if (value) {
                    // Common patterns for thumbnail URLs and their full-size counterparts
                    const fullSizeUrl = transformToFullSize(value);
                    if (fullSizeUrl) return fullSizeUrl;
                }
            }
            return element.src; // Fallback to regular src if no alternatives found
        }

        // For background images
        const bgImage = window.getComputedStyle(element).backgroundImage;
        if (bgImage && bgImage !== 'none') {
            const url = bgImage.slice(4, -1).replace(/["']/g, "");
            return transformToFullSize(url);
        }

        return null;
    }

    // Function to transform thumbnail URLs to full-size URLs
    function transformToFullSize(url) {
        if (!url) return url;

        // Common URL patterns for thumbnails and their full-size replacements
        const patterns = [
            {
                // Transform thumbnail dimensions to full size
                regex: /(_thumb|_small|_medium|_tiny|\b\d+x\d+\b)/i,
                replacement: ''
            },
            {
                // Replace common size indicators
                regex: /-\d+x\d+\./i,
                replacement: '.'
            },
            {
                // Replace width-based parameters
                regex: /[?&]w=\d+/i,
                replacement: ''
            },
            {
                // Replace size-related parameters
                regex: /[?&]size=\w+/i,
                replacement: ''
            },
            {
                // Replace quality/compression parameters
                regex: /[?&]q=\d+/i,
                replacement: ''
            }
        ];

        let fullSizeUrl = url;
        patterns.forEach(pattern => {
            fullSizeUrl = fullSizeUrl.replace(pattern.regex, pattern.replacement);
        });

        return fullSizeUrl;
    }

    // Function to find the closest image element
    function findClosestImage(element) {
        if (!element) return null;
        
        // If element is an image, return it
        if (element.tagName === 'IMG') {
            return element;
        }
        
        // Check if element has background image
        const bgImage = window.getComputedStyle(element).backgroundImage;
        if (bgImage && bgImage !== 'none') {
            const url = bgImage.slice(4, -1).replace(/["']/g, "");
            return { src: url, isBgImage: true };
        }
        
        // Look for img element within the current element
        const img = element.querySelector('img');
        if (img) {
            return img;
        }
        
        return null;
    }

    // Function to download the image
    function downloadImage(imageElement) {
        if (!imageElement) return;
        
        // Get the full-size image URL
        const imageUrl = getFullSizeImageUrl(imageElement);
        if (!imageUrl) return;
        
        // Extract filename from URL
        const filename = imageUrl.split('/').pop().split('#')[0].split('?')[0];
        
        // Use GM_download if available, otherwise fallback to creating a link
        if (typeof GM_download !== 'undefined') {
            GM_download({
                url: imageUrl,
                name: filename,
                saveAs: false
            });
        } else {
            const link = document.createElement('a');
            link.href = imageUrl;
            link.download = filename;
            document.body.appendChild(link);
            link.click();
            document.body.removeChild(link);
        }
    }

    // Track mouse position
    let mouseX = 0;
    let mouseY = 0;
    document.addEventListener('mousemove', function(e) {
        mouseX = e.clientX;
        mouseY = e.clientY;
    });

    // Listen for '0' key
    document.addEventListener('keydown', function(e) {
        if (e.key === '0') {
            const elementUnderMouse = document.elementFromPoint(mouseX, mouseY);
            const imageElement = findClosestImage(elementUnderMouse);
            if (imageElement) {
                downloadImage(imageElement);
            }
        }
    });
})();
