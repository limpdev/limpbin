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

