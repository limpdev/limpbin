// frontend/src/app.js
// Connect to Go backend
import {GetSvgFiles, GetSvgContent, GetIconName} from '../wailsjs/go/main/App';
import {EventsOn, LogInfo, LogError} from '../wailsjs/runtime/runtime';

// DOM elements
let gallery;
let codeCopiedMessage;

// Initialize DOM elements
function initDOM() {
    gallery = document.getElementById('gallery');
    codeCopiedMessage = document.getElementById('codeCopied');
    
    if (!gallery) {
        LogError("Gallery element not found!");
        // Create it if missing
        gallery = document.createElement('div');
        gallery.id = 'gallery';
        gallery.className = 'gallery';
        document.body.appendChild(gallery);
    }
    
    if (!codeCopiedMessage) {
        LogError("CodeCopied element not found!");
        // Create it if missing
        codeCopiedMessage = document.createElement('div');
        codeCopiedMessage.id = 'codeCopied';
        codeCopiedMessage.className = 'code-copied-message';
        codeCopiedMessage.textContent = 'Copied';
        document.body.appendChild(codeCopiedMessage);
    }
}

// Update gallery with SVG files
async function updateGallery() {
    try {
        LogInfo("Updating gallery...");
        const files = await GetSvgFiles();
        LogInfo(`Found ${files.length} SVG files`);
        
        if (!gallery) {
            LogError("Gallery element is null, reinitializing DOM");
            initDOM();
        }
        
        gallery.innerHTML = '';
        
        if (files.length === 0) {
            LogInfo("No SVG files found");
            gallery.innerHTML = '<div class="no-files">No SVG files found. Add some .svg files to your directory.</div>';
            return;
        }
        
        for (const file of files) {
            try {
                const iconName = await GetIconName(file);
                const svgContent = await GetSvgContent(file);
                
                if (!svgContent) {
                    LogError(`Empty SVG content for ${file}`);
                    continue;
                }
                
                const iconWrapper = document.createElement('div');
                iconWrapper.className = 'icon-wrapper';
                iconWrapper.innerHTML = `
                    ${svgContent}
                    <div class="icon-name">${iconName}</div>
                `;
                
                iconWrapper.addEventListener('click', function(event) {
                    copySvgCode(this, svgContent);
                    
                    // Add ripple effect
                    const ripple = document.createElement('span');
                    ripple.classList.add('ripple');
                    const rect = this.getBoundingClientRect();
                    const size = Math.max(rect.width, rect.height);
                    ripple.style.width = ripple.style.height = `${size}px`;
                    ripple.style.left = `${event.clientX - rect.left - size / 2}px`;
                    ripple.style.top = `${event.clientY - rect.top - size / 2}px`;
                    this.appendChild(ripple);
                    
                    setTimeout(() => ripple.remove(), 1500);
                });
                
                gallery.appendChild(iconWrapper);
            } catch (fileErr) {
                LogError(`Error processing file ${file}: ${fileErr}`);
            }
        }
    } catch (err) {
        LogError(`Error updating gallery: ${err}`);
        gallery.innerHTML = `<div class="error">Error: ${err.message}</div>`;
    }
}

// Copy SVG code to clipboard
function copySvgCode(_iconWrapper, svgCode) {
    navigator.clipboard.writeText(svgCode).then(() => {
        showCopiedMessage();
    }).catch(err => {
        LogError(`Failed to copy SVG code: ${err}`);
        alert('Failed to copy SVG code. Please try again.');
    });
}

// Show copied message
function showCopiedMessage() {
    codeCopiedMessage.classList.add('show');
    setTimeout(() => {
        codeCopiedMessage.classList.remove('show');
    }, 1500);
}

// Listen for updates from Go backend
EventsOn('svg-files-updated', () => {
    LogInfo("Received svg-files-updated event");
    updateGallery();
});

// Initialize
function init() {
    LogInfo("Initializing app.js");
    initDOM();
    updateGallery();
}

// Initial load
if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', init);
} else {
    // DOM already loaded
    init();
}

// Also listen for Wails ready event
window.runtime = {
    ...window.runtime,
    WindowLoad: () => {
        LogInfo("Wails WindowLoad event triggered");
        setTimeout(updateGallery, 500); // Add a small delay
    }
};

const iconsPerPage = 100; // Adjust based on your preference
  const icons = document.querySelectorAll(".icon-wrapper");
  let currentPage = 1;
  const totalPages = Math.ceil(icons.length / iconsPerPage);

  function showPage(page) {
      icons.forEach((icon, index) => {
          icon.classList.toggle("hidden", index < (page - 1) * iconsPerPage || index >= page * iconsPerPage);
      });

      document.getElementById("page-info").textContent = `Page ${page} of ${totalPages}`;
      document.getElementById("prev").disabled = page === 1;
      document.getElementById("next").disabled = page === totalPages;
  }

  document.getElementById("prev").addEventListener("click", () => {
      if (currentPage > 1) {
          currentPage--;
          showPage(currentPage);
      }
  });

  document.getElementById("next").addEventListener("click", () => {
      if (currentPage < totalPages) {
          currentPage++;
          showPage(currentPage);
      }
  });

  // Initial display
  showPage(currentPage);