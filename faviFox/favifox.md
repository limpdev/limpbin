# faviFox
> a Firefox extension to handle favicon customization. Let's break this down into parts:
- On FireFox, disable extension signing in `about:config` -> `xpinstall.signatures.required` = `false`
- Then, add a signature `id`, which can you can generate using `web-ext`... why you ask? because god-forbid you make the user experience simple and easy - No, you have to do all this bullshit just to load up some spaghettiscript
> Once you've downloaded a fucking 3rd-party app to "sIgN yOuR eXtEnSiOn", go to mozilla.whatever to obtain a fucking API key so the crackpot team over there can ensure the security of my locally hosted, spaghettiscript extension... because why the fuck not, it must be necessary!
>> I wonder why developers are so reluctant to move over from chromium's ecosystem, I guess we'll never know!

```json
"browser_specific_settings": {
  "gecko": {
    "id": "limpdev"
  }
}
```

## manifest.json
```json
{
  "manifest_version": 2,
  "name": "Favicon Customizer",
  "version": "1.0",
  "description": "Customize favicons for websites and bookmarks",
  "permissions": [
    "<all_urls>",
    "bookmarks",
    "storage",
    "tabs"
  ],
  "background": {
    "scripts": ["background.js"]
  },
  "content_scripts": [{
    "matches": ["<all_urls>"],
    "js": ["content.js"],
    "run_at": "document_start"
  }]
}
```

## background.js
```js
let faviconMap = {};

// Load saved favicon mappings
browser.storage.local.get('faviconMap').then(result => {
  if (result.faviconMap) {
    faviconMap = result.faviconMap;
  }
});

// Handle tab updates
browser.tabs.onUpdated.addListener((tabId, changeInfo, tab) => {
  if (changeInfo.status === 'complete') {
    const customIcon = faviconMap[tab.url];
    if (customIcon) {
      browser.tabs.sendMessage(tabId, {
        type: 'updateFavicon',
        iconUrl: customIcon
      });
    }
  }
});

// Handle bookmark creation and updates
browser.bookmarks.onCreated.addListener((id, bookmark) => {
  updateBookmarkFavicon(id, bookmark);
});

browser.bookmarks.onChanged.addListener((id, changeInfo) => {
  browser.bookmarks.get(id).then(bookmarks => {
    updateBookmarkFavicon(id, bookmarks[0]);
  });
});

function updateBookmarkFavicon(id, bookmark) {
  const customIcon = faviconMap[bookmark.url];
  if (customIcon) {
    // Firefox doesn't provide direct API to modify bookmark icons
    // We can store the custom icon URL in extension storage
    browser.storage.local.set({
      [`bookmark-${id}`]: customIcon
    });
  }
}

// API to set custom favicon
function setCustomFavicon(url, iconUrl) {
  faviconMap[url] = iconUrl;
  browser.storage.local.set({ faviconMap });
  
  // Update any open tabs with this URL
  browser.tabs.query({ url }).then(tabs => {
    tabs.forEach(tab => {
      browser.tabs.sendMessage(tab.id, {
        type: 'updateFavicon',
        iconUrl
      });
    });
  });
}
```

## content.js
```js
browser.runtime.onMessage.addListener((message) => {
  if (message.type === 'updateFavicon') {
    // Remove existing favicon links
    const existingLinks = document.querySelectorAll('link[rel*="icon"]');
    existingLinks.forEach(link => link.remove());
    
    // Create new favicon link
    const link = document.createElement('link');
    link.type = 'image/x-icon';
    link.rel = 'shortcut icon';
    link.href = message.iconUrl;
    document.head.appendChild(link);
  }
});

```

This extension handles both the tab favicons and bookmark favicons, though with some limitations due to Firefox's architecture. Here's how it works:

1. Tab Favicon Management:
   - The content script injects the custom favicon by replacing existing favicon links
   - Changes are applied immediately when a tab loads or updates
   - Uses the `tabs` permission to monitor tab changes

2. Bookmark Favicon Management:
   - Firefox doesn't provide direct API access to modify bookmark icons
   - The extension stores custom favicon mappings in local storage
   - When bookmarks are created or updated, it associates the custom icon URL with the bookmark ID

3. Storage:
   - Uses `browser.storage.local` to persist favicon mappings
   - Maintains a mapping of URLs to custom favicon URLs

To use this extension, you would:
1. Create these files in your extension directory
2. Load it as a temporary extension in Firefox for testing
3. Call `setCustomFavicon(url, iconUrl)` to set custom icons

Some limitations to be aware of:
1. Firefox's bookmark system caches favicons internally, and there's no direct API to modify them
2. You might need to clear Firefox's favicon cache to see changes in bookmarks
3. The bookmark icons might revert to original ones when Firefox refreshes its cache
---

>I'll enhance the extension with image support and add a management UI.

```json
// manifest.json
{
  "manifest_version": 2,
  "name": "Favicon Customizer",
  "version": "1.0",
  "description": "Customize favicons for websites and bookmarks",
  "permissions": [
    "<all_urls>",
    "bookmarks",
    "storage",
    "tabs"
  ],
  "background": {
    "scripts": ["background.js"]
  },
  "content_scripts": [{
    "matches": ["<all_urls>"],
    "js": ["content.js"],
    "run_at": "document_start"
  }],
  "browser_action": {
    "default_popup": "popup.html"
  }
}
```

```js
// background.js
let faviconRules = {};

// Load saved favicon rules
browser.storage.local.get('faviconRules').then(result => {
  if (result.faviconRules) {
    faviconRules = result.faviconRules;
  }
});

// Helper to check if URL matches any rule
function findMatchingRule(url) {
  for (const [pattern, rule] of Object.entries(faviconRules)) {
    try {
      const regex = new RegExp(pattern);
      if (regex.test(url)) {
        return rule;
      }
    } catch (e) {
      // Invalid regex, try exact match
      if (url.includes(pattern)) {
        return rule;
      }
    }
  }
  return null;
}

// Handle tab updates
browser.tabs.onUpdated.addListener((tabId, changeInfo, tab) => {
  if (changeInfo.status === 'complete') {
    const rule = findMatchingRule(tab.url);
    if (rule) {
      browser.tabs.sendMessage(tabId, {
        type: 'updateFavicon',
        iconData: rule.iconData,
        iconType: rule.iconType
      });
    }
  }
});

// Handle bookmark updates
browser.bookmarks.onCreated.addListener((id, bookmark) => {
  updateBookmarkFavicon(id, bookmark);
});

browser.bookmarks.onChanged.addListener((id, changeInfo) => {
  browser.bookmarks.get(id).then(bookmarks => {
    updateBookmarkFavicon(id, bookmarks[0]);
  });
});

function updateBookmarkFavicon(id, bookmark) {
  const rule = findMatchingRule(bookmark.url);
  if (rule) {
    browser.storage.local.set({
      [`bookmark-${id}`]: {
        iconData: rule.iconData,
        iconType: rule.iconType
      }
    });
  }
}

// API to set custom favicon rule
function setFaviconRule(urlPattern, iconData, iconType) {
  faviconRules[urlPattern] = { iconData, iconType };
  browser.storage.local.set({ faviconRules });
  
  // Update any open tabs that match the pattern
  browser.tabs.query({}).then(tabs => {
    tabs.forEach(tab => {
      if (findMatchingRule(tab.url)) {
        browser.tabs.sendMessage(tab.id, {
          type: 'updateFavicon',
          iconData,
          iconType
        });
      }
    });
  });
}

// content.js
browser.runtime.onMessage.addListener((message) => {
  if (message.type === 'updateFavicon') {
    // Remove existing favicon links
    const existingLinks = document.querySelectorAll('link[rel*="icon"]');
    existingLinks.forEach(link => link.remove());
    
    // Create new favicon link
    const link = document.createElement('link');
    link.rel = 'shortcut icon';
    
    if (message.iconType === 'base64') {
      link.href = message.iconData;
    } else if (message.iconType === 'svg') {
      link.type = 'image/svg+xml';
      link.href = message.iconData;
    } else if (message.iconType === 'png') {
      link.type = 'image/png';
      link.href = message.iconData;
    }
    
    document.head.appendChild(link);
  }
});

```

```html
<!-- popup.html -->
<!DOCTYPE html>
<html>
<head>
  <style>
    body {
      width: 400px;
      padding: 16px;
      font-family: system-ui, -apple-system, sans-serif;
    }
    .rule-form {
      display: flex;
      flex-direction: column;
      gap: 12px;
      margin-bottom: 20px;
    }
    .input-group {
      display: flex;
      flex-direction: column;
      gap: 4px;
    }
    .input-group label {
      font-weight: 500;
    }
    input, textarea {
      padding: 8px;
      border: 1px solid #ccc;
      border-radius: 4px;
    }
    button {
      padding: 8px 16px;
      background: #0060df;
      color: white;
      border: none;
      border-radius: 4px;
      cursor: pointer;
    }
    button:hover {
      background: #003eaa;
    }
    .rules-list {
      margin-top: 20px;
    }
    .rule-item {
      display: flex;
      align-items: center;
      gap: 8px;
      padding: 8px;
      border: 1px solid #eee;
      margin-bottom: 8px;
    }
    .rule-item img {
      width: 16px;
      height: 16px;
    }
    .delete-btn {
      background: #dd0000;
      margin-left: auto;
    }
  </style>
</head>
<body>
  <div class="rule-form">
    <div class="input-group">
      <label for="urlPattern">URL Pattern (supports regex)</label>
      <input type="text" id="urlPattern" placeholder="e.g., .*\.example\.com or example.com">
    </div>
    
    <div class="input-group">
      <label for="iconInput">Icon (URL, Base64, or SVG)</label>
      <textarea id="iconInput" rows="3" placeholder="Paste icon URL, Base64 data, or SVG code"></textarea>
    </div>
    
    <button id="addRule">Add Rule</button>
  </div>
  
  <div class="rules-list" id="rulesList"></div>
  
  <script src="popup.js"></script>
</body>
</html>
```

```js
// popup.js
document.addEventListener('DOMContentLoaded', () => {
  // Load and display existing rules
  loadRules();
  
  document.getElementById('addRule').addEventListener('click', addNewRule);
});

function loadRules() {
  browser.storage.local.get('faviconRules').then(result => {
    const rules = result.faviconRules || {};
    const rulesList = document.getElementById('rulesList');
    rulesList.innerHTML = '';
    
    Object.entries(rules).forEach(([pattern, rule]) => {
      const ruleElement = createRuleElement(pattern, rule);
      rulesList.appendChild(ruleElement);
    });
  });
}

function createRuleElement(pattern, rule) {
  const div = document.createElement('div');
  div.className = 'rule-item';
  
  // Preview icon
  const img = document.createElement('img');
  img.src = rule.iconData;
  div.appendChild(img);
  
  // Pattern text
  const span = document.createElement('span');
  span.textContent = pattern;
  div.appendChild(span);
  
  // Delete button
  const deleteBtn = document.createElement('button');
  deleteBtn.className = 'delete-btn';
  deleteBtn.textContent = 'Delete';
  deleteBtn.onclick = () => deleteRule(pattern);
  div.appendChild(deleteBtn);
  
  return div;
}

function addNewRule() {
  const pattern = document.getElementById('urlPattern').value.trim();
  const iconInput = document.getElementById('iconInput').value.trim();
  
  if (!pattern || !iconInput) {
    alert('Please fill in both fields');
    return;
  }
  
  // Detect icon type
  let iconType;
  if (iconInput.startsWith('data:image/')) {
    iconType = 'base64';
  } else if (iconInput.startsWith('<svg')) {
    iconType = 'svg';
  } else if (iconInput.match(/\.(png)$/i)) {
    iconType = 'png';
  } else {
    iconType = 'base64'; // Default to base64 for pasted data
  }
  
  // Save rule
  browser.runtime.getBackgroundPage().then(background => {
    background.setFaviconRule(pattern, iconInput, iconType);
    loadRules(); // Refresh the list
    
    // Clear inputs
    document.getElementById('urlPattern').value = '';
    document.getElementById('iconInput').value = '';
  });
}

function deleteRule(pattern) {
  browser.storage.local.get('faviconRules').then(result => {
    const rules = result.faviconRules || {};
    delete rules[pattern];
    browser.storage.local.set({ faviconRules: rules }).then(() => {
      loadRules(); // Refresh the list
    });
  });
}

```

I've made several enhancements to the extension:

1. Image Support:
   - Added support for PNG, SVG, and Base64 images
   - Automatic detection of image type
   - Proper MIME type handling for different image formats

2. Rule-based Matching:
   - Uses regex or simple string matching for URL patterns
   - Stores rules instead of direct URL mappings
   - More flexible matching system for wildcards and patterns

3. Management UI:
   - Clean, simple popup interface
   - Fields for URL pattern and icon input
   - Live preview of existing rules
   - Ability to delete rules
   - Error handling for invalid inputs

To use the enhanced extension:

1. Create a new rule:
   - Click the extension icon
   - Enter a URL pattern (can be regex or simple string)
   - Paste your Base64 image data, SVG code, or image URL
   - Click "Add Rule"

2. The extension will:
   - Apply the favicon immediately to matching tabs
   - Store the rule for future use
   - Show a preview in the rules list

Would you like me to:
1. Add validation for Base64/SVG input?
2. Add import/export functionality for rules?
3. Add a feature to capture and convert existing favicons to Base64?