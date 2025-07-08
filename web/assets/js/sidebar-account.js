// Shared functionality for loading selected account in sidebar across all pages

// Load selected account in sidebar
function loadSelectedAccountInSidebar() {
  const savedAccount = localStorage.getItem('selectedAccount');
  if (savedAccount) {
    try {
      const account = JSON.parse(savedAccount);
      // Update the account name display
      const accountNameElement = document.getElementById('current-account-name');
      const accountIdElement = document.getElementById('current-account-id');
      
      if (accountNameElement && account.name) {
        // For very long names, show a truncated version
        const displayName = account.name.length > 25 ? 
          account.name.substring(0, 22) + '...' : account.name;
        
        accountNameElement.textContent = displayName;
        accountNameElement.setAttribute('title', account.name);
        accountNameElement.setAttribute('data-original-title', account.name);
        
        // Update Bootstrap tooltip if it exists
        if ($(accountNameElement).data('bs.tooltip')) {
          $(accountNameElement).tooltip('dispose');
        }
        $(accountNameElement).tooltip();
      }
      
      if (accountIdElement && account.id) {
        // For IDs, show only the last part if it's very long
        let displayId = account.id;
        if (account.id.length > 20) {
          displayId = '...' + account.id.substring(account.id.length - 17);
        }
        const idText = `ID: ${displayId}`;
        const fullIdText = `ID: ${account.id}`;
        
        accountIdElement.textContent = idText;
        accountIdElement.setAttribute('title', fullIdText);
        accountIdElement.setAttribute('data-original-title', fullIdText);
        
        // Update Bootstrap tooltip if it exists
        if ($(accountIdElement).data('bs.tooltip')) {
          $(accountIdElement).tooltip('dispose');
        }
        $(accountIdElement).tooltip();
      }
    } catch (error) {
      console.error('Error parsing selected account from localStorage:', error);
    }
  }
}

// Initialize tooltips
function initializeTooltips() {
  if (typeof $ !== 'undefined' && $.fn.tooltip) {
    $('[data-toggle="tooltip"]').tooltip();
  }
}

// Initialize sidebar account loading when DOM is ready
$(document).ready(function() {
  // Load selected account after HTMX loads the sidebar component
  $(document).on('htmx:afterSettle', function(event) {
    if (event.target.getAttribute && event.target.getAttribute('hx-get') === '/components/sidebar') {
      loadSelectedAccountInSidebar();
      initializeTooltips();
    }
  });
  
  // Also listen for storage changes to update in real-time when account is changed
  window.addEventListener('storage', function(e) {
    if (e.key === 'selectedAccount') {
      loadSelectedAccountInSidebar();
    }
  });

  // Listen for custom events when account is selected (for same-page updates)
  window.addEventListener('accountSelected', function() {
    loadSelectedAccountInSidebar();
  });
});
