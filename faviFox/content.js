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