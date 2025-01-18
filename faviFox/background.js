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
