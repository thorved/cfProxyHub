// docker-cloudflare-tunnels.js - Handles Docker-based Cloudflare Tunnel UI operations

// Debug logger
function debugLog(message) {
    console.log("[Docker CF Tunnels] " + message);
    
    // Also update debug section if it exists
    const debugInfo = document.getElementById('debugInfo');
    if (debugInfo) {
        debugInfo.textContent += "\n" + message;
    }
}

// Global initialization flag to prevent multiple executions
let isInitialized = false;

// Initialize when the document is ready or when manually called
function initializeUI() {
    // Check if already initialized
    if (isInitialized) {
        debugLog("UI already initialized, skipping");
        return;
    }
    
    debugLog("Initializing Docker Cloudflare Tunnels UI");
    
    try {
        // Check if jQuery is available
        if (typeof jQuery === 'undefined') {
            debugLog("ERROR: jQuery is not available. Waiting...");
            setTimeout(initializeUI, 500);  // Try again in 500ms
            return;
        }
        
        // Set up event listener for creating a new tunnel
        const createBtn = document.getElementById('createTunnelBtn');
        if (createBtn) {
            createBtn.addEventListener('click', createTunnel);
            debugLog("Create tunnel button listener added");
        } else {
            debugLog("ERROR: Create tunnel button not found!");
        }
        
        // Load existing tunnels
        loadTunnels();
        
        // Mark as initialized
        isInitialized = true;
        
    } catch (err) {
        debugLog("ERROR during initialization: " + err.message);
        showError("Failed to initialize page: " + err.message);
        
        // Try again after a delay
        setTimeout(initializeUI, 1000);
    }
}

// Initialize when the document is ready
document.addEventListener('DOMContentLoaded', function() {
    debugLog("DOM content loaded, initializing UI...");
    setTimeout(initializeUI, 100); // Short delay to ensure other scripts load
});

// Also initialize when the window is fully loaded
window.addEventListener('load', function() {
    debugLog("Window fully loaded, ensuring UI is initialized...");
    setTimeout(initializeUI, 100);
});

// Load existing Cloudflare Tunnel containers
function loadTunnels() {
    debugLog("Loading tunnels...");
    
    // Show loading indicator
    document.getElementById('loadingIndicator').style.display = 'block';
    document.getElementById('errorContainer').style.display = 'none';
    
    fetch('/api/docker/cloudflare/tunnels')
        .then(response => {
            debugLog("Received response with status: " + response.status);
            if (!response.ok) {
                throw new Error(`HTTP error! Status: ${response.status}`);
            }
            return response.json();
        })
        .then(data => {
            debugLog("Parsed data: " + JSON.stringify(data).substring(0, 100) + "...");
            
            // Enhanced logging for debugging response structure
            if (data && data.data && Array.isArray(data.data) && data.data.length > 0) {
                const firstItem = data.data[0];
                debugLog(`First tunnel ID format: ${firstItem.Id ? 'Id' : (firstItem.ID ? 'ID' : 'neither')}`);
                debugLog(`First tunnel ID value: ${getContainerId(firstItem)}`);
            }
            
            if (data && (data.success === true || data.status === "success")) {
                displayTunnels(data.data || []);
            } else {
                const message = data && data.message ? data.message : 'Unknown error';
                showError('Failed to load tunnels: ' + message);
                showAlert('danger', 'Failed to load tunnels: ' + message);
            }
        })
        .catch(error => {
            debugLog("Error loading tunnels: " + error.message);
            
            // Try to get more details from the error response if possible
            if (error.response) {
                try {
                    error.response.json().then(errorData => {
                        const errorMsg = errorData && errorData.message ? errorData.message : error.message;
                        showError('Error loading tunnels: ' + errorMsg);
                    }).catch(() => {
                        showError('Error loading tunnels: ' + error.message);
                    });
                } catch (e) {
                    showError('Error loading tunnels: ' + error.message);
                }
            } else {
                showError('Error loading tunnels: ' + error.message);
            }
            showAlert('danger', 'Error loading tunnels: ' + error.message);
        })
        .finally(() => {
            // Hide loading indicator
            document.getElementById('loadingIndicator').style.display = 'none';
        });
}

// Display tunnels in the table
// Show error message in the error container
function showError(message) {
    const errorContainer = document.getElementById('errorContainer');
    const errorMessage = document.getElementById('errorMessage');
    
    if (errorContainer && errorMessage) {
        errorMessage.textContent = message;
        errorContainer.style.display = 'block';
    }
}

// Display tunnels in the table
function displayTunnels(tunnels) {
    debugLog("Displaying " + (tunnels ? tunnels.length : 0) + " tunnels");
    
    const tableBody = document.getElementById('tunnelsTableBody');
    if (!tableBody) {
        debugLog("ERROR: tunnelsTableBody element not found!");
        return;
    }
    
    // Clear the table body
    tableBody.innerHTML = '';

    if (!tunnels || tunnels.length === 0) {
        tableBody.innerHTML = '<tr><td colspan="5" class="text-center">No Cloudflare tunnel containers found</td></tr>';
        debugLog("No tunnels to display");
        return;
    }

    try {
        tunnels.forEach(tunnel => {
            debugLog("Processing tunnel: " + JSON.stringify(tunnel).substring(0, 100));
            
            const tr = document.createElement('tr');
            
            // Use helper function to get container ID safely
            const fullId = getContainerId(tunnel);
            const id = fullId !== 'Unknown' ? fullId.substring(0, 12) : 'Unknown';
            
            const names = tunnel.Names && Array.isArray(tunnel.Names) ? 
                tunnel.Names.join(', ').replace(/^\//g, '') : 'Unnamed';
            const state = tunnel.State || 'unknown';
            const created = tunnel.Created ? new Date(tunnel.Created * 1000).toLocaleString() : 'Unknown';
            
            tr.innerHTML = `
                <td>${id}</td>
                <td>${names}</td>
                <td><span class="badge ${state === 'running' ? 'badge-success' : 'badge-danger'}">${state}</span></td>
                <td>${created}</td>
                <td>
                    <div class="btn-group" role="group">
                        ${state === 'running' 
                            ? `<button class="btn btn-sm btn-warning" onclick="stopTunnel('${id}')">Stop</button>` 
                            : `<button class="btn btn-sm btn-success" onclick="startTunnel('${id}')">Start</button>`}
                    <button class="btn btn-sm btn-info" onclick="restartTunnel('${id}')">Restart</button>
                    <button class="btn btn-sm btn-danger" onclick="deleteTunnel('${id}')">Delete</button>
                </div>
            </td>
        `;
        tableBody.appendChild(tr);
    });
    } catch (error) {
        debugLog("Error while displaying tunnels: " + error.message);
        showError("Error displaying tunnels: " + error.message);
    }
}

// Create a new tunnel
function createTunnel() {
    const name = document.getElementById('tunnelName').value.trim();
    const token = document.getElementById('tunnelToken').value.trim();
    const restartPolicy = document.getElementById('restartPolicy').value;

    if (!token) {
        showAlert('danger', 'Tunnel token is required');
        return;
    }

    showAlert('info', 'Creating tunnel... Please wait.');
    fetch('/api/docker/cloudflare/tunnels', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({
            name: name || `cloudflared-tunnel-${Date.now().toString().substring(7)}`,
            token: token,
            restart_policy: restartPolicy
        }),
    })
    .then(response => response.json())
    .then(data => {
        if (data.success === true || data.status === 'success') {
            showAlert('success', (data.data && data.data.message) ? data.data.message : 'Tunnel created successfully');
            document.getElementById('tunnelName').value = '';
            document.getElementById('tunnelToken').value = '';
        } else {
            showAlert('warning', 'Tunnel creation status uncertain. Will verify in a moment...');
        }
        setTimeout(loadTunnels, 2000);
    })
    .catch(error => {
        showAlert('warning', 'Error from server, but will verify tunnel creation...');
        setTimeout(loadTunnels, 2000);
    });
}

// Start a tunnel
function startTunnel(id) {
    if (id === 'Unknown' || !id) {
        showAlert('danger', 'Cannot start tunnel: Invalid container ID');
        debugLog("Attempted to start tunnel with invalid ID: " + id);
        return;
    }
    debugLog("Starting tunnel with ID: " + id);
    showAlert('info', 'Starting tunnel... Please wait.');
    fetch(`/api/docker/cloudflare/tunnels/${id}/start`, { method: 'POST' })
    .then(response => response.json())
    .then(data => {
        if (data.success === true || data.status === 'success') {
            showAlert('success', (data.data && data.data.message) ? data.data.message : 'Tunnel started successfully');
        } else {
            showAlert('warning', 'Tunnel start status uncertain. Will verify in a moment...');
        }
        setTimeout(loadTunnels, 2000);
    })
    .catch(error => {
        showAlert('warning', 'Error from server, but will verify tunnel start...');
        setTimeout(loadTunnels, 2000);
    });
}

// Stop a tunnel
function stopTunnel(id) {
    if (id === 'Unknown' || !id) {
        showAlert('danger', 'Cannot stop tunnel: Invalid container ID');
        debugLog("Attempted to stop tunnel with invalid ID: " + id);
        return;
    }
    debugLog("Stopping tunnel with ID: " + id);
    showAlert('info', 'Stopping tunnel... Please wait.');
    fetch(`/api/docker/cloudflare/tunnels/${id}/stop`, { method: 'POST' })
    .then(response => response.json())
    .then(data => {
        if (data.success === true || data.status === 'success') {
            showAlert('success', (data.data && data.data.message) ? data.data.message : 'Tunnel stopped successfully');
        } else {
            showAlert('warning', 'Tunnel stop status uncertain. Will verify in a moment...');
        }
        setTimeout(loadTunnels, 2000);
    })
    .catch(error => {
        showAlert('warning', 'Error from server, but will verify tunnel stop...');
        setTimeout(loadTunnels, 2000);
    });
}

// Restart a tunnel
function restartTunnel(id) {
    if (id === 'Unknown' || !id) {
        showAlert('danger', 'Cannot restart tunnel: Invalid container ID');
        debugLog("Attempted to restart tunnel with invalid ID: " + id);
        return;
    }
    debugLog("Restarting tunnel with ID: " + id);
    showAlert('info', 'Restarting tunnel... Please wait.');
    fetch(`/api/docker/cloudflare/tunnels/${id}/restart`, { method: 'POST' })
    .then(response => response.json())
    .then(data => {
        if (data.success === true || data.status === 'success') {
            showAlert('success', (data.data && data.data.message) ? data.data.message : 'Tunnel restarted successfully');
        } else {
            showAlert('warning', 'Tunnel restart status uncertain. Will verify in a moment...');
        }
        setTimeout(loadTunnels, 2500);
    })
    .catch(error => {
        showAlert('warning', 'Error from server, but will verify tunnel restart...');
        setTimeout(loadTunnels, 2500);
    });
}

// Delete a tunnel
function deleteTunnel(id) {
    if (id === 'Unknown' || !id) {
        showAlert('danger', 'Cannot delete tunnel: Invalid container ID');
        debugLog("Attempted to delete tunnel with invalid ID: " + id);
        return;
    }
    if (!confirm('Are you sure you want to delete this tunnel? This action cannot be undone.')) {
        return;
    }
    debugLog("Deleting tunnel with ID: " + id);
    showAlert('info', 'Deleting tunnel... Please wait.');
    fetch(`/api/docker/cloudflare/tunnels/${id}`, { method: 'DELETE' })
    .then(response => response.json())
    .then(data => {
        if (data.success === true || data.status === 'success') {
            showAlert('success', (data.data && data.data.message) ? data.data.message : 'Tunnel deleted successfully');
        } else {
            showAlert('warning', 'Tunnel deletion status uncertain. Will verify in a moment...');
        }
        setTimeout(loadTunnels, 2000);
    })
    .catch(error => {
        showAlert('warning', 'Error from server, but will verify tunnel deletion...');
        setTimeout(loadTunnels, 2000);
    });
}

// Check Docker connectivity
function checkDockerConnectivity() {
    debugLog("Checking Docker connectivity...");
    
    fetch('/api/docker/debug')
        .then(response => response.json())
        .then(data => {
            if (data.success) {
                debugLog("Docker connectivity: " + data.message);
                showAlert('info', 'Docker connectivity: ' + data.message);
            } else {
                debugLog("Docker connectivity issue: " + data.message);
                showAlert('danger', 'Docker connectivity issue: ' + data.message);
            }
        })
        .catch(error => {
            debugLog("Error checking Docker: " + error.message);
            showAlert('danger', 'Error checking Docker connectivity: ' + error.message);
        });
}

// Get detailed Docker diagnostics
function getDockerDiagnostics() {
    debugLog("Getting Docker diagnostics...");
    
    fetch('/api/docker/diagnostics')
        .then(response => response.json())
        .then(data => {
            if (data.success && data.data) {
                const info = data.data;
                let message = `Docker Status: ${info.docker_status}\n`;
                message += `Container Count: ${info.container_count || 0}\n`;
                message += `Tunnel Count: ${info.tunnel_count || 0}\n`;
                
                if (info.errors && info.errors.length > 0) {
                    message += "\nErrors:\n" + info.errors.join("\n");
                }
                
                debugLog(message);
                
                // Add diagnostics button to the UI for detailed information
                const debugInfo = document.getElementById('debugInfo');
                if (debugInfo) {
                    // Create a pre element with full diagnostics data
                    const pre = document.createElement('pre');
                    pre.textContent = JSON.stringify(info, null, 2);
                    pre.style.maxHeight = '300px';
                    pre.style.overflow = 'auto';
                    pre.style.marginTop = '10px';
                    pre.style.backgroundColor = '#333';
                    pre.style.padding = '10px';
                    pre.style.borderRadius = '4px';
                    
                    // Clear previous diagnostics if any
                    const existingPre = debugInfo.querySelector('pre');
                    if (existingPre) {
                        debugInfo.removeChild(existingPre);
                    }
                    
                    // Add new diagnostics
                    debugInfo.appendChild(pre);
                }
                
                showAlert('info', `Docker diagnostics: ${info.docker_status}, ${info.container_count || 0} containers, ${info.tunnel_count || 0} tunnels`);
            } else {
                debugLog("Error getting diagnostics: " + (data.message || "Unknown error"));
                showAlert('warning', 'Error getting Docker diagnostics');
            }
        })
        .catch(error => {
            debugLog("Error getting diagnostics: " + error.message);
            showAlert('danger', 'Error getting Docker diagnostics: ' + error.message);
        });
}

// Show alert notification to the user
function showAlert(type, message) {
    // Log to console for debugging
    if (type === 'danger' || type === 'error') {
        console.error(message);
    } else {
        console.log(message);
    }
    
    // Create alert element
    const alertDiv = document.createElement('div');
    alertDiv.className = `alert alert-${type} alert-dismissible fade show`;
    alertDiv.role = 'alert';
    
    // Add message
    alertDiv.innerHTML = `
        ${message}
        <button type="button" class="close" data-dismiss="alert" aria-label="Close">
            <span aria-hidden="true">&times;</span>
        </button>
    `;
    
    // Insert at top of page
    const container = document.querySelector('.content-wrapper');
    if (container) {
        // Insert after page header
        const pageHeader = container.querySelector('.page-header');
        if (pageHeader) {
            pageHeader.insertAdjacentElement('afterend', alertDiv);
        } else {
            container.insertAdjacentElement('afterbegin', alertDiv);
        }
        
        // Auto dismiss after 5 seconds
        setTimeout(() => {
            try {
                alertDiv.classList.remove('show');
                setTimeout(() => {
                    alertDiv.remove();
                }, 300);
            } catch (e) {
                // Ignore errors if element already removed
            }
        }, 5000);
    }
    
    // Also update debug info if available
    const debugInfo = document.getElementById('debugInfo');
    if (debugInfo) {
        const timestamp = new Date().toLocaleTimeString();
        debugInfo.textContent += `\n[${timestamp}] ${type.toUpperCase()}: ${message}`;
        
        // Scroll to bottom
        debugInfo.scrollTop = debugInfo.scrollHeight;
    }
}

// Helper function to get the container ID safely (handles both Id and ID formats)
function getContainerId(container) {
    if (!container) return 'Unknown';
    return container.Id ? container.Id : (container.ID ? container.ID : 'Unknown');
}
