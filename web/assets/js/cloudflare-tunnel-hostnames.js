$(document).ready(function() {
  let selectedAccountId = null;
  let selectedAccount = null;
  let selectedTunnelId = null;
  let selectedTunnel = null;
  let currentHostname = null;
  let isEditMode = false;
  
  // Initialize page
  initializePage();
  
  // Refresh hostnames
  $('#refreshHostnames').on('click', function() {
    if ($(this).prop('disabled')) return;
    
    const button = $(this);
    button.prop('disabled', true);
    const originalContent = button.html();
    button.html('<i class="mdi mdi-loading mdi-spin"></i> Refreshing...');
    
    loadHostnames().finally(() => {
      button.prop('disabled', false);
      button.html(originalContent);
    });
  });
  
  // Add hostname buttons
  $('#addHostnameBtn, #addFirstHostnameBtn').on('click', function() {
    showHostnameModal();
  });
  
  // Edit hostname
  $(document).on('click', '.edit-hostname-btn', function() {
    const hostname = $(this).data('hostname');
    const serviceType = $(this).data('service-type');
    const serviceUrl = $(this).data('service-url');
    const path = $(this).data('path');
    
    showHostnameModal(hostname, serviceType, serviceUrl, path);
  });
  
  // Delete hostname
  $(document).on('click', '.delete-hostname-btn', function() {
    const hostname = $(this).data('hostname');
    currentHostname = hostname;
    $('#deleteHostnameText').text(hostname);
    $('#deleteHostnameModal').modal('show');
  });
  
  // Save hostname
  $('#saveHostnameBtn').on('click', function() {
    if ($(this).prop('disabled')) return;
    saveHostname();
  });
  
  // Confirm delete hostname
  $('#confirmDeleteHostname').on('click', function() {
    if ($(this).prop('disabled')) return;
    deleteHostname(currentHostname);
  });
  
  // Form validation
  $('#hostnameForm input[required], #hostnameForm select[required]').on('input change', function() {
    validateForm();
  });
  
  // Real-time hostname validation
  $('#hostname').on('input', function() {
    const hostname = $(this).val().trim();
    if (hostname.length > 0) {
      if (hostname.includes('.') && hostname.length > 3) {
        $(this).removeClass('is-invalid').addClass('is-valid');
      } else {
        $(this).removeClass('is-valid').addClass('is-invalid');
      }
    } else {
      $(this).removeClass('is-valid is-invalid');
    }
    validateForm();
  });

  // Real-time subdomain validation and hostname preview
  $('#subdomain, #domain, #manualDomain').on('input change', function() {
    updateHostnamePreview();
    validateHostname();
    validateForm();
  });

  // Real-time subdomain validation
  $('#subdomain').on('input', function() {
    const subdomain = $(this).val().trim();
    const feedback = $(this).siblings('.invalid-feedback');
    
    if (subdomain.length > 0) {
      // DNS subdomain validation
      if (/^[a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?$/.test(subdomain)) {
        $(this).removeClass('is-invalid').addClass('is-valid');
        feedback.hide();
      } else {
        $(this).removeClass('is-valid').addClass('is-invalid');
        let errorMsg = 'Invalid subdomain format. ';
        if (subdomain.length > 63) {
          errorMsg += 'Must be 63 characters or less.';
        } else if (subdomain.startsWith('-') || subdomain.endsWith('-')) {
          errorMsg += 'Cannot start or end with a hyphen.';
        } else if (!/^[a-zA-Z0-9-]+$/.test(subdomain)) {
          errorMsg += 'Only letters, numbers, and hyphens allowed.';
        } else {
          errorMsg += 'Must start and end with a letter or number.';
        }
        
        if (feedback.length === 0) {
          $(this).after('<div class="invalid-feedback"></div>');
        }
        $(this).siblings('.invalid-feedback').text(errorMsg).show();
      }
    } else {
      $(this).removeClass('is-valid is-invalid');
      feedback.hide();
    }
  });
  
  // Service URL validation
  $('#serviceUrl').on('input', function() {
    const serviceUrl = $(this).val().trim();
    const feedback = $(this).siblings('.invalid-feedback');
    
    if (serviceUrl.length > 0) {
      // Basic URL validation
      if (/^[a-zA-Z0-9.-]+:\d+$/.test(serviceUrl) || 
          /^[a-zA-Z0-9.-]+$/.test(serviceUrl) ||
          /^https?:\/\/[a-zA-Z0-9.-]+(:\d+)?$/.test(serviceUrl)) {
        $(this).removeClass('is-invalid').addClass('is-valid');
        feedback.hide();
      } else {
        $(this).removeClass('is-valid').addClass('is-invalid');
        const errorMsg = 'Invalid service URL format. Examples: localhost:8080, 192.168.1.10:3000, http://localhost:8080';
        
        if (feedback.length === 0) {
          $(this).after('<div class="invalid-feedback"></div>');
        }
        $(this).siblings('.invalid-feedback').text(errorMsg).show();
      }
    } else {
      $(this).removeClass('is-valid is-invalid');
      feedback.hide();
    }
    
    validateForm();
  });
  
  // Manual domain input toggle
  $('#manualDomainToggle').on('click', function(e) {
    e.preventDefault();
    const isManualMode = $('#manualDomain').is(':visible');
    
    if (isManualMode) {
      // Switch back to dropdown
      $('#domain').show();
      $('#manualDomain').hide().val('');
      $('#manualDomainHelp').hide();
      $(this).text("Can't find your domain? Enter manually");
    } else {
      // Switch to manual input
      $('#domain').hide();
      $('#manualDomain').show().focus();
      $('#manualDomainHelp').show();
      $(this).text("Use domain dropdown instead");
    }
    
    updateHostnamePreview();
    validateForm();
  });
  
  // Manual domain input validation
  $('#manualDomain').on('input', function() {
    const domain = $(this).val().trim();
    if (domain.length > 0) {
      // Basic domain validation
      if (/^[a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(\.[a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$/.test(domain)) {
        $(this).removeClass('is-invalid').addClass('is-valid');
      } else {
        $(this).removeClass('is-valid').addClass('is-invalid');
      }
    } else {
      $(this).removeClass('is-valid is-invalid');
    }
    
    updateHostnamePreview();
    validateForm();
  });
  
  // Reset modal when closed
  $('#hostnameModal').on('hidden.bs.modal', function() {
    resetHostnameForm();
  });
  
  // Handle Enter key in form inputs
  $('#hostnameForm input').on('keypress', function(e) {
    if (e.which === 13) { // Enter key
      e.preventDefault();
      if (!$('#saveHostnameBtn').prop('disabled')) {
        $('#saveHostnameBtn').click();
      }
    }
  });
});

// Initialize page with tunnel and account info
function initializePage() {
  console.log('Initializing page...');
  
  // Get tunnel info from URL params or localStorage
  const urlParams = new URLSearchParams(window.location.search);
  const tunnelIdFromUrl = urlParams.get('tunnelId');
  const accountIdFromUrl = urlParams.get('accountId');
  
  console.log('URL params - tunnelId:', tunnelIdFromUrl, 'accountId:', accountIdFromUrl);
  
  // Try to get tunnel info from localStorage first
  const savedTunnel = localStorage.getItem('selectedTunnel');
  if (savedTunnel) {
    try {
      selectedTunnel = JSON.parse(savedTunnel);
      selectedTunnelId = selectedTunnel.id;
      selectedAccountId = selectedTunnel.accountId;
      console.log('Loaded tunnel from localStorage:', selectedTunnel);
    } catch (error) {
      console.error('Error parsing selected tunnel:', error);
    }
  }
  
  // Fallback to URL params
  if (!selectedTunnelId && tunnelIdFromUrl) {
    selectedTunnelId = tunnelIdFromUrl;
    console.log('Using tunnelId from URL');
  }
  if (!selectedAccountId && accountIdFromUrl) {
    selectedAccountId = accountIdFromUrl;
    console.log('Using accountId from URL');
  }
  
  // Get account info from localStorage
  const savedAccount = localStorage.getItem('selectedAccount');
  if (savedAccount) {
    try {
      selectedAccount = JSON.parse(savedAccount);
      console.log('Loaded account from localStorage:', selectedAccount);
    } catch (error) {
      console.error('Error parsing selected account:', error);
    }
  }
  
  console.log('Final values - tunnelId:', selectedTunnelId, 'accountId:', selectedAccountId);
  
  if (!selectedTunnelId || !selectedAccountId) {
    showNotification('error', 'Missing tunnel or account information');
    setTimeout(() => {
      window.location.href = '/CloudflareAllTunnels';
    }, 2000);
    return;
  }
  
  // Update display
  updateTunnelDisplay();
  updateAccountDisplay();
  
  // Load hostnames
  loadHostnames();
  
  // Load domains for dropdown
  loadDomains();
}

// Update tunnel display
function updateTunnelDisplay() {
  if (selectedTunnel) {
    $('#tunnelName').text(selectedTunnel.name);
    $('#tunnelId').text(selectedTunnel.id);
  } else {
    $('#tunnelName').text('Unknown Tunnel');
    $('#tunnelId').text(selectedTunnelId || 'Unknown');
  }
}

// Update account display
function updateAccountDisplay() {
  if (selectedAccount) {
    $('#currentAccountName').text(selectedAccount.name);
    $('#currentAccountId').text(`ID: ${selectedAccount.id.substring(0, 8)}...`);
    $('#currentAccountBadge').text(selectedAccount.name);
  } else {
    $('#currentAccountName').text('Unknown Account');
    $('#currentAccountId').text(selectedAccountId || 'Unknown');
    $('#currentAccountBadge').text('Unknown Account');
  }
}

// Load hostnames for the tunnel
function loadHostnames() {
  showLoadingState();
  
  return fetch(`/api/cloudflare/accounts/${selectedAccountId}/tunnels/${selectedTunnelId}/hostnames`, {
    method: 'GET',
    headers: {
      'Content-Type': 'application/json',
    },
    credentials: 'same-origin'
  })
  .then(response => response.json())
  .then(data => {
    hideLoadingState();
    
    if (data.status === 'success' && data.data && data.data.hostnames) {
      if (data.data.hostnames.length > 0) {
        renderHostnames(data.data.hostnames);
        showHostnamesGrid();
        showNotification('success', `Loaded ${data.data.hostnames.length} hostnames`);
      } else {
        showEmptyState();
      }
    } else {
      showEmptyState();
      showNotification('info', 'No hostnames found for this tunnel');
    }
  })
  .catch(error => {
    hideLoadingState();
    console.error('Error loading hostnames:', error);
    showEmptyState();
    showNotification('error', 'Failed to load hostnames');
  });
}

// Load domains for dropdown
function loadDomains() {
  if (!selectedAccountId) {
    console.error('No account selected');
    return Promise.reject('No account selected');
  }
  
  console.log('Loading domains for account:', selectedAccountId);
  
  return fetch(`/api/cloudflare/accounts/${selectedAccountId}/zones/dropdown?active_only=true&limit=50`, {
    method: 'GET',
    credentials: 'same-origin'
  })
  .then(response => {
    console.log('Zone API response status:', response.status);
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    return response.json();
  })
  .then(data => {
    console.log('Domains loaded:', data);
    if (data.status === 'success' && data.data) {
      const domainSelect = $('#domain');
      domainSelect.empty();
      domainSelect.append('<option value="">Select domain...</option>');
      
      if (data.data.zones && data.data.zones.length > 0) {
        console.log('Found zones:', data.data.zones);
        data.data.zones.forEach(zone => {
          console.log('Processing zone:', zone);
          // Only show active domains
          if (zone.status === 'active') {
            domainSelect.append(`<option value="${zone.name}">${zone.name}</option>`);
          }
        });
        
        if (domainSelect.find('option').length === 1) {
          domainSelect.append('<option value="" disabled>No active domains found</option>');
          showNotification('warning', 'No active domains found for this account');
        }
      } else {
        console.log('No zones found in response');
        domainSelect.append('<option value="" disabled>No active domains found</option>');
        showNotification('warning', 'No active domains found for this account');
      }
      
      return data.data.zones || [];
    } else {
      console.error('Invalid response format:', data);
      throw new Error(data.message || 'Invalid response format');
    }
  })
  .catch(error => {
    console.error('Error loading domains:', error);
    const domainSelect = $('#domain');
    domainSelect.empty();
    domainSelect.append('<option value="">Select domain...</option>');
    domainSelect.append('<option value="" disabled>Error loading domains</option>');
    
    let errorMessage = 'Failed to load domains';
    if (error.message.includes('404')) {
      errorMessage = 'Zones API endpoint not found';
    } else if (error.message.includes('401') || error.message.includes('403')) {
      errorMessage = 'Unauthorized access to zones';
    } else if (error.message.includes('500')) {
      errorMessage = 'Server error while loading zones';
    } else if (error.message) {
      errorMessage = 'Failed to load domains: ' + error.message;
    }
    
    showNotification('error', errorMessage);
    return [];
  });
}

// Update hostname preview
function updateHostnamePreview() {
  const subdomain = $('#subdomain').val().trim();
  const domain = getDomainValue();
  const preview = $('#hostnamePreview');
  
  if (subdomain && domain) {
    const fullHostname = `${subdomain}.${domain}`;
    preview.html(`<strong class="text-success">${fullHostname}</strong>`);
  } else if (subdomain) {
    preview.html(`<span class="text-muted">${subdomain}.</span><span class="text-secondary">[select domain]</span>`);
  } else if (domain) {
    preview.html(`<span class="text-secondary">[enter subdomain]</span><span class="text-muted">.${domain}</span>`);
  } else {
    preview.html('<span class="text-muted">Enter subdomain and select domain...</span>');
  }
}

// Get domain value from either dropdown or manual input
function getDomainValue() {
  if ($('#manualDomain').is(':visible')) {
    return $('#manualDomain').val().trim();
  } else {
    return $('#domain').val();
  }
}

// Parse hostname into subdomain and domain
function parseHostname(hostname) {
  if (!hostname || !hostname.includes('.')) return { subdomain: '', domain: '' };
  
  const parts = hostname.split('.');
  if (parts.length < 2) return { subdomain: '', domain: '' };
  
  // For simplicity, assume first part is subdomain and rest is domain
  const subdomain = parts[0];
  const domain = parts.slice(1).join('.');
  
  return { subdomain, domain };
}

// Render hostnames in table
function renderHostnames(hostnames) {
  const tbody = $('#hostnamesList');
  tbody.empty();
  
  hostnames.forEach(hostname => {
    const createdDate = hostname.created_at ? new Date(hostname.created_at).toLocaleDateString('en-US', { 
      year: 'numeric', 
      month: 'short', 
      day: 'numeric' 
    }) : 'N/A';
    
    const status = hostname.status || 'active';
    const statusBadge = getStatusBadge(status);
    
    const serviceInfo = hostname.service || 'N/A';
    
    // Parse service URL for edit functionality
    const serviceUrl = hostname.service || '';
    let serviceType = 'http';
    let cleanUrl = serviceUrl;
    
    if (serviceUrl.startsWith('https://')) {
      serviceType = 'https';
      cleanUrl = serviceUrl.substring(8);
    } else if (serviceUrl.startsWith('http://')) {
      cleanUrl = serviceUrl.substring(7);
    } else if (serviceUrl.startsWith('tcp://')) {
      serviceType = 'tcp';
      cleanUrl = serviceUrl.substring(6);
    }
    
    // Add row to table
    tbody.append(`
      <tr>
        <td>
          <div class="d-flex align-items-center">
            <div class="preview-icon bg-gradient-info rounded-circle mr-3" style="width: 40px; height: 40px;">
              <i class="mdi mdi-earth text-white"></i>
            </div>
            <div>
              <h6 class="mb-1 font-weight-medium">${hostname.hostname}</h6>
              <p class="text-muted mb-0 small">${hostname.path || '/'}</p>
            </div>
          </div>
        </td>
        <td>
          <span class="badge badge-pill badge-${serviceType === 'http' ? 'primary' : serviceType === 'https' ? 'success' : 'info'}">
            ${serviceType.toUpperCase()}
          </span>
          <span class="text-muted small" style="margin-left: 8px;">${cleanUrl}</span>
        </td>
        <td>
          ${statusBadge}
        </td>
        <td>
          ${createdDate}
        </td>
        <td>
          <button class="btn btn-sm btn-light edit-hostname-btn" 
                  data-hostname="${hostname.hostname}" 
                  data-service-type="${serviceType}" 
                  data-service-url="${cleanUrl}" 
                  data-path="${hostname.path || '/'}">
            <i class="mdi mdi-pencil"></i> Edit
          </button>
          <button class="btn btn-sm btn-danger delete-hostname-btn" data-hostname="${hostname.hostname}">
            <i class="mdi mdi-delete"></i> Delete
          </button>
        </td>
      </tr>
    `);
  });
  
  // Update hostname count
  $('#hostnameCount').text(`${hostnames.length} Hostnames`);
}

// Get status badge HTML
function getStatusBadge(status) {
  const statusLower = status.toLowerCase();
  let badgeClass = 'badge-secondary';
  let badgeText = 'Unknown';
  
  if (statusLower === 'active') {
    badgeClass = 'badge-success';
    badgeText = 'Active';
  } else if (statusLower === 'paused') {
    badgeClass = 'badge-warning';
    badgeText = 'Paused';
  } else if (statusLower === 'error') {
    badgeClass = 'badge-danger';
    badgeText = 'Error';
  }
  
  return `<span class="badge badge-pill ${badgeClass}">${badgeText}</span>`;
}

// Show hostname modal for adding or editing
function showHostnameModal(hostname, serviceType, serviceUrl, path) {
  isEditMode = !!hostname;
  currentHostname = hostname || null;
  
  // Reset form
  resetHostnameForm();
  
  if (isEditMode) {
    // Edit mode - pre-fill form fields
    const parsed = parseHostname(hostname);
    $('#subdomain').val(parsed.subdomain);
    $('#domain').val(parsed.domain);
    $('#serviceType').val(serviceType);
    $('#serviceUrl').val(serviceUrl);
    $('#path').val(path);
    
    $('#hostnameModalLabel').html(`
      <i class="mdi mdi-pencil mr-2"></i>Edit Public Hostname
    `);
  } else {
    $('#hostnameModalLabel').html(`
      <i class="mdi mdi-plus mr-2"></i>Add Public Hostname
    `);
  }
  
  // Update hostname preview
  updateHostnamePreview();
  
  // Show modal
  $('#hostnameModal').modal('show');
}

// Reset hostname form
function resetHostnameForm() {
  $('#hostnameForm')[0].reset();
  $('#hostnamePreview').html('<span class="text-muted">Enter subdomain and select domain...</span>');
  $('#subdomain').removeClass('is-valid is-invalid');
  $('#domain').removeClass('is-valid is-invalid');
  $('#serviceType').removeClass('is-valid is-invalid');
  $('#serviceUrl').removeClass('is-valid is-invalid');
  $('#path').removeClass('is-valid is-invalid');
  $('.invalid-feedback').remove();
  isEditMode = false;
  currentHostname = null;
  
  $('#hostnameModalLabel').html(`
    <i class="mdi mdi-plus mr-2"></i>Add Public Hostname
  `);
}

// Validate hostname
function validateHostname() {
  const subdomain = $('#subdomain').val().trim();
  const domain = getDomainValue();
  
  if (subdomain && domain) {
    const fullHostname = `${subdomain}.${domain}`;
    
    // Check if hostname is already in use (excluding current hostname in edit mode)
    const existingHostnames = $('#hostnamesList tr').map(function() {
      const hostname = $(this).find('.edit-hostname-btn').data('hostname');
      return hostname;
    }).get();
    
    if (existingHostnames.includes(fullHostname) && (!isEditMode || fullHostname !== currentHostname)) {
      $('#hostnamePreview').html('<span class="text-danger"><i class="mdi mdi-alert-circle mr-1"></i>Hostname already exists</span>');
      return false;
    }
    
    return true;
  }
  
  return false;
}

// Validate form
function validateForm() {
  const subdomain = $('#subdomain').val().trim();
  const domain = getDomainValue();
  const serviceType = $('#serviceType').val();
  const serviceUrl = $('#serviceUrl').val().trim();
  
  // Check all required fields
  const requiredFieldsValid = subdomain && domain && serviceType && serviceUrl;
  
  // Check subdomain format
  const subdomainValid = /^[a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?$/.test(subdomain);
  
  // Check service URL format
  const serviceUrlValid = serviceUrl.length > 0 && !serviceUrl.includes(' ');
  
  // Check hostname uniqueness
  const hostnameValid = validateHostname();
  
  const isValid = requiredFieldsValid && subdomainValid && serviceUrlValid && hostnameValid;
  $('#saveHostnameBtn').prop('disabled', !isValid);
  
  return isValid;
}

// Save hostname
function saveHostname() {
  if (!validateForm()) {
    showNotification('error', 'Please fix the validation errors before saving');
    return;
  }
  
  const button = $('#saveHostnameBtn');
  button.prop('disabled', true);
  const originalContent = button.html();
  button.html('<i class="mdi mdi-loading mdi-spin mr-2"></i>Saving...');
  
  const subdomain = $('#subdomain').val().trim();
  const domain = getDomainValue();
  const serviceType = $('#serviceType').val();
  const serviceUrl = $('#serviceUrl').val().trim();
  const path = $('#path').val().trim() || '/';
  
  const fullHostname = `${subdomain}.${domain}`;
  
  // Build service URL based on type
  let fullServiceUrl = serviceUrl;
  if (serviceType === 'http' && !serviceUrl.startsWith('http://')) {
    fullServiceUrl = `http://${serviceUrl}`;
  } else if (serviceType === 'https' && !serviceUrl.startsWith('https://')) {
    fullServiceUrl = `https://${serviceUrl}`;
  } else if (serviceType === 'tcp' && !serviceUrl.startsWith('tcp://')) {
    fullServiceUrl = `tcp://${serviceUrl}`;
  }
  
  const requestData = {
    hostname: fullHostname,
    service: fullServiceUrl,
    path: path
  };
  
  console.log('Saving hostname:', requestData);
  
  const url = isEditMode ? 
    `/api/cloudflare/accounts/${selectedAccountId}/tunnels/${selectedTunnelId}/hostnames/${encodeURIComponent(fullHostname)}` :
    `/api/cloudflare/accounts/${selectedAccountId}/tunnels/${selectedTunnelId}/hostnames`;
  
  const method = isEditMode ? 'PUT' : 'POST';
  
  fetch(url, {
    method: method,
    headers: {
      'Content-Type': 'application/json',
    },
    credentials: 'same-origin',
    body: JSON.stringify(requestData)
  })
  .then(response => {
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    return response.json();
  })
  .then(data => {
    console.log('Hostname save response:', data);
    if (data.status === 'success') {
      $('#hostnameModal').modal('hide');
      showNotification('success', isEditMode ? 'Hostname updated successfully' : 'Hostname created successfully');
      loadHostnames();
    } else {
      throw new Error(data.message || 'Failed to save hostname');
    }
  })
  .catch(error => {
    console.error('Error saving hostname:', error);
    let errorMessage = 'Failed to save hostname';
    
    if (error.message.includes('already exists')) {
      errorMessage = 'A hostname with this name already exists';
    } else if (error.message.includes('invalid')) {
      errorMessage = 'Invalid hostname configuration';
    } else if (error.message.includes('unauthorized')) {
      errorMessage = 'Unauthorized access to this domain';
    } else if (error.message) {
      errorMessage = error.message;
    }
    
    showNotification('error', errorMessage);
  })
  .finally(() => {
    button.prop('disabled', false);
    button.html(originalContent);
  });
}

// Delete hostname
function deleteHostname(hostname) {
  if (!hostname) return;
  
  const button = $('#confirmDeleteHostname');
  button.prop('disabled', true);
  const originalContent = button.html();
  button.html('<i class="mdi mdi-loading mdi-spin mr-2"></i>Deleting...');
  
  fetch(`/api/cloudflare/accounts/${selectedAccountId}/tunnels/${selectedTunnelId}/hostnames/${encodeURIComponent(hostname)}`, {
    method: 'DELETE',
    headers: {
      'Content-Type': 'application/json',
    },
    credentials: 'same-origin'
  })
  .then(response => {
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    return response.json();
  })
  .then(data => {
    console.log('Delete response:', data);
    if (data.status === 'success') {
      $('#deleteHostnameModal').modal('hide');
      showNotification('success', 'Hostname deleted successfully');
      loadHostnames();
    } else {
      throw new Error(data.message || 'Failed to delete hostname');
    }
  })
  .catch(error => {
    console.error('Error deleting hostname:', error);
    showNotification('error', error.message || 'Failed to delete hostname');
  })
  .finally(() => {
    button.prop('disabled', false);
    button.html(originalContent);
  });
}

// Show loading state
function showLoadingState() {
  $('#loadingState').show();
  $('#hostnamesGrid').hide();
  $('#emptyState').hide();
}

// Hide loading state
function hideLoadingState() {
  $('#loadingState').hide();
}

// Show hostnames grid
function showHostnamesGrid() {
  $('#hostnamesGrid').show();
  $('#emptyState').hide();
}

// Show empty state
function showEmptyState() {
  $('#emptyState').show();
  $('#hostnamesGrid').hide();
}

// Show notification
function showNotification(type, message) {
  // Create notification element
  const alertClass = type === 'success' ? 'alert-success' : 
                    type === 'error' ? 'alert-danger' : 
                    type === 'warning' ? 'alert-warning' : 'alert-info';
  
  const notification = $(`
    <div class="alert ${alertClass} alert-dismissible fade show" role="alert" style="position: fixed; top: 20px; right: 20px; z-index: 9999; min-width: 300px;">
      <i class="mdi mdi-${type === 'success' ? 'check-circle' : type === 'error' ? 'alert-circle' : 'information'} mr-2"></i>
      ${message}
      <button type="button" class="close" data-dismiss="alert" aria-label="Close">
        <span aria-hidden="true">&times;</span>
      </button>
    </div>
  `);
  
  // Add to body
  $('body').append(notification);
  
  // Auto-dismiss after 5 seconds
  setTimeout(() => {
    notification.alert('close');
  }, 5000);
}
