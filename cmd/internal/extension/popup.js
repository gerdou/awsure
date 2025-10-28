// Popup script for extension settings

// Load current port
chrome.storage.local.get(['cliPort'], (result) => {
  if (result.cliPort) {
    document.getElementById('port').value = result.cliPort;
  }
});

// Save port button
const btn = document.getElementById('savePort');
btn.addEventListener('click', () => {
  const port = parseInt(document.getElementById('port').value);
  if (port < 1024 || port > 65535) {
    alert('Port must be between 1024 and 65535');
    return;
  }
  chrome.storage.local.set({ cliPort: port }, () => {
    chrome.tabs.query({}, (tabs) => {
      tabs.forEach(tab => {
        chrome.tabs.sendMessage(tab.id, { type: 'UPDATE_PORT', port }).catch(() => {});
      });
    });
    btn.textContent = 'Saved!';
    setTimeout(() => (btn.textContent = 'Save Port'), 1500);
  });
});
