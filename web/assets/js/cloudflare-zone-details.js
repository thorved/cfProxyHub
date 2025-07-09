// Cloudflare Zone Details JavaScript

class CloudflareZoneDetailsManager {
    constructor() {
        this.zoneId = null;
        this.zone = null;
        this.init();
    }

    init() {
        this.zoneId = this.getZoneIdFromURL();
        if (this.zoneId) {
            this.bindEvents();
            this.loadCurrentAccount().then(() => {
                this.loadZoneDetails();
            });
        } else {
            this.showAlert('Zone ID not found in URL', 'danger');
        }
    }

    async loadCurrentAccount() {
        try {
            // First try to get current account from API
            const response = await fetch('/api/current-account', {
                headers: {
                    'Authorization': `Bearer ${localStorage.getItem('token')}`
                }
            });

            if (response.ok) {
                const data = await response.json();
                if (data.success && data.account) {
                    this.currentAccountId = data.account.id;
                    this.currentAccountName = data.account.name;
                    console.log('Loaded current account from API:', data.account);
                    return;
                }
            }
        } catch (error) {
            console.warn('Failed to get current account from API:', error);
        }

        // Fallback to localStorage
        try {
            const savedAccount = localStorage.getItem('selectedAccount');
            if (savedAccount) {
                const account = JSON.parse(savedAccount);
                this.currentAccountId = account.id;
                this.currentAccountName = account.name;
                console.log('Loaded current account from localStorage:', account);
                return;
            }
        } catch (error) {
            console.error('Error parsing account from localStorage:', error);
        }
    }

    getZoneIdFromURL() {
        const urlParams = new URLSearchParams(window.location.search);
        return urlParams.get('id');
    }

    bindEvents() {
        // Refresh zone details
        $('#refreshZoneBtn').on('click', () => {
            this.loadZoneDetails();
        });

        // Pause/Resume zone
        $('#pauseResumeBtn').on('click', () => {
            this.toggleZonePause();
        });

        // Delete zone
        $('#deleteZoneBtn').on('click', () => {
            this.showDeleteZoneModal();
        });

        $('#confirmDeleteZone').on('click', () => {
            this.deleteZone();
        });

        // Manage DNS records (placeholder)
        $('#manageDNSBtn').on('click', () => {
            this.showAlert('DNS management will be available in a future update', 'info');
        });

        // Copy functionality
        $(document).on('click', '.copy-btn', (e) => {
            const text = $(e.target).data('copy-text');
            this.copyToClipboard(text);
        });
    }

    async loadZoneDetails() {
        this.showLoading(true);
        this.showComponentLoading(true);

        // Add a timeout fallback to ensure loading stops
        const timeoutId = setTimeout(() => {
            console.warn('Loading timeout - forcing loading to stop');
            this.showLoading(false);
            this.showComponentLoading(false);
        }, 10000); // 10 second timeout

        try {
            const response = await fetch(`/api/cloudflare/zones/${this.zoneId}`, {
                headers: {
                    'Authorization': `Bearer ${localStorage.getItem('token')}`
                }
            });

            if (!response.ok) {
                throw new Error(`Failed to fetch zone details: ${response.status}`);
            }

            const data = await response.json();
            
            // Handle the response structure consistently
            const apiData = data.data || data;
            this.zone = apiData.data || apiData;
            
            if (!this.zone) {
                throw new Error('Zone data not found in response');
            }
            
            // Wrap updateZoneDetails in try-catch to prevent silent failures
            try {
                this.updateZoneDetails();
            } catch (updateError) {
                console.error('Error updating zone details UI:', updateError);
                this.showAlert('Error displaying zone details', 'warning');
            }
        } catch (error) {
            console.error('Error loading zone details:', error);
            this.showAlert('Error loading zone details: ' + error.message, 'danger');
        } finally {
            clearTimeout(timeoutId); // Clear the timeout
            // Force loading to stop regardless of what happened
            setTimeout(() => {
                this.showLoading(false);
                this.showComponentLoading(false);
            }, 100);
        }
    }

    updateZoneDetails() {
        if (!this.zone) return;

        // Update header
        $('#zoneName .zone-name-text').text(this.zone.name || 'Unknown');
        $('#zoneId').text(this.zone.id || '');
        
        const status = this.zone.status || 'unknown';
        $('#zoneStatus').removeClass().addClass(`badge ${this.getStatusBadgeClass(status)}`);
        $('#zoneStatus .status-text').text(status.charAt(0).toUpperCase() + status.slice(1));

        // Update pause/resume button
        const isPaused = status === 'paused';
        const pauseResumeBtn = $('#pauseResumeBtn');
        pauseResumeBtn
            .removeClass('btn-warning btn-success')
            .addClass(isPaused ? 'btn-success' : 'btn-warning')
            .find('i')
            .removeClass('mdi-pause mdi-play')
            .addClass(isPaused ? 'mdi-play' : 'mdi-pause');
        $('#pauseResumeText').text(isPaused ? 'Resume' : 'Pause');

        // Update zone information
        $('#detailZoneName').html(`${this.zone.name} ${this.createCopyButton(this.zone.name)}`);
        $('#detailZoneStatus').removeClass().addClass(`badge ${this.getStatusBadgeClass(this.zone.status)}`).text(this.zone.status.charAt(0).toUpperCase() + this.zone.status.slice(1));
        $('#detailZoneType').text(this.zone.type || 'N/A');
        $('#detailDevMode').html(`
            <i class="mdi ${this.zone.development_mode ? 'mdi-toggle-switch text-warning' : 'mdi-toggle-switch-off text-muted'}"></i>
            ${this.zone.development_mode ? 'Enabled' : 'Disabled'}
        `);
        $('#detailCreatedOn').text(this.formatDate(this.zone.created_on));
        $('#detailModifiedOn').text(this.formatDate(this.zone.modified_on));
        $('#detailActivatedOn').text(this.formatDate(this.zone.activated_on));

        // Update name servers
        this.updateNameServers();

        // Update account information
        this.updateAccountInfo();

        // Update plan information
        this.updatePlanInfo();

        // Update permissions
        this.updatePermissions();

        // Update delete modal
        $('#deleteZoneName').text(this.zone.name);
    }

    updateNameServers() {
        // Cloudflare name servers
        const nameServersList = $('#nameServersList');
        nameServersList.empty();
        
        if (this.zone.name_servers && this.zone.name_servers.length > 0) {
            this.zone.name_servers.forEach(ns => {
                nameServersList.append(`
                    <div class="nameserver-item d-flex justify-content-between align-items-center">
                        <span>${ns}</span>
                        ${this.createCopyButton(ns)}
                    </div>
                `);
            });
        } else {
            nameServersList.append('<p class="text-muted">No name servers available</p>');
        }

        // Original name servers
        const originalNameServersList = $('#originalNameServersList');
        originalNameServersList.empty();
        
        if (this.zone.original_name_servers && this.zone.original_name_servers.length > 0) {
            this.zone.original_name_servers.forEach(ns => {
                originalNameServersList.append(`
                    <div class="nameserver-item d-flex justify-content-between align-items-center">
                        <span>${ns}</span>
                        ${this.createCopyButton(ns)}
                    </div>
                `);
            });
        } else {
            originalNameServersList.append('<p class="text-muted">No original name servers</p>');
        }
    }

    updateAccountInfo() {
        const account = this.zone.account || {};
        $('#detailAccountId').html(`${account.id || 'N/A'} ${account.id ? this.createCopyButton(account.id) : ''}`);
        $('#detailAccountName').text(account.name || 'N/A');
        $('#detailOriginalRegistrar').text(this.zone.original_registrar || 'N/A');
        $('#detailOriginalDNSHost').text(this.zone.original_dnshost || 'N/A');
    }

    updatePlanInfo() {
        const plan = this.zone.plan || {};
        $('#detailPlanName').text(plan.name || 'N/A');
        $('#detailPlanPrice').text(plan.price ? `$${plan.price}` : 'N/A');
        $('#detailPlanCurrency').text(plan.currency || 'N/A');
        $('#detailPlanFrequency').text(plan.frequency || 'N/A');
    }

    updatePermissions() {
        const permissionsList = $('#permissionsList');
        permissionsList.empty();
        
        if (this.zone.permissions && this.zone.permissions.length > 0) {
            this.zone.permissions.forEach(permission => {
                permissionsList.append(`
                    <span class="permission-item">${permission}</span>
                `);
            });
        } else {
            permissionsList.append('<p class="text-muted">No permissions information available</p>');
        }
    }

    async toggleZonePause() {
        const isPaused = this.zone.status === 'paused';
        const shouldPause = !isPaused;

        try {
            const response = await fetch(`/api/cloudflare/zones/${this.zoneId}`, {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${localStorage.getItem('token')}`
                },
                body: JSON.stringify({
                    paused: shouldPause
                })
            });

            if (!response.ok) {
                throw new Error('Failed to update zone');
            }

            const action = shouldPause ? 'paused' : 'resumed';
            this.showAlert(`Zone ${action} successfully`, 'success');
            
            // Reload zone details to get updated status
            setTimeout(() => {
                this.loadZoneDetails();
            }, 1000);
        } catch (error) {
            console.error('Error updating zone:', error);
            this.showAlert('Error updating zone', 'danger');
        }
    }

    showDeleteZoneModal() {
        $('#deleteZoneModal').modal('show');
    }

    async deleteZone() {
        this.showDeleteZoneLoading(true);

        try {
            const response = await fetch(`/api/cloudflare/zones/${this.zoneId}`, {
                method: 'DELETE',
                headers: {
                    'Authorization': `Bearer ${localStorage.getItem('token')}`
                }
            });

            if (!response.ok) {
                throw new Error('Failed to delete zone');
            }

            this.showAlert('Zone deleted successfully', 'success');
            
            // Redirect to zones list after successful deletion
            setTimeout(() => {
                window.location.href = '/cloudflare/zones';
            }, 2000);
        } catch (error) {
            console.error('Error deleting zone:', error);
            this.showAlert('Error deleting zone', 'danger');
        } finally {
            this.showDeleteZoneLoading(false);
        }
    }

    getStatusBadgeClass(status) {
        const statusClasses = {
            'active': 'badge-success',
            'pending': 'badge-warning',
            'paused': 'badge-secondary',
            'initializing': 'badge-info',
            'moved': 'badge-dark',
            'deleted': 'badge-danger'
        };

        return statusClasses[status] || 'badge-secondary';
    }

    formatDate(dateString) {
        if (!dateString) return 'N/A';
        
        try {
            const date = new Date(dateString);
            return date.toLocaleString();
        } catch (error) {
            return 'Invalid Date';
        }
    }

    createCopyButton(text) {
        return `<button class="copy-btn" data-copy-text="${text}" title="Copy to clipboard">
                    <i class="mdi mdi-content-copy"></i>
                </button>`;
    }

    showLoading(show) {
        const overlay = $('#loadingOverlay');
        const topLoadingBar = $('#topLoadingBar');
        const loadingProgress = $('#loadingProgress');
        
        if (show) {
            // Show top loading bar
            topLoadingBar.show();
            loadingProgress.addClass('animate');
            
            // Show minimal content overlay
            overlay.removeClass('d-none').css({
                'display': 'flex',
                'visibility': 'visible',
                'opacity': '1'
            });
            
            // Add loading state to main content
            $('.main-panel').addClass('zone-info-loading');
        } else {
            // Hide top loading bar
            setTimeout(() => {
                topLoadingBar.hide();
                loadingProgress.removeClass('animate');
            }, 500);
            
            // Hide content overlay
            overlay.addClass('d-none').css({
                'display': 'none',
                'visibility': 'hidden',
                'opacity': '0'
            });
            
            // Remove loading state
            $('.main-panel').removeClass('zone-info-loading');
        }
        
        console.log('Loading state changed:', show ? 'shown' : 'hidden');
    }

    showDeleteZoneLoading(show) {
        const spinner = $('#deleteZoneSpinner');
        const button = $('#confirmDeleteZone');
        
        if (show) {
            spinner.show();
            button.prop('disabled', true);
        } else {
            spinner.hide();
            button.prop('disabled', false);
        }
    }

    showComponentLoading(show) {
        const elements = [
            '#zoneName .loading-spinner',
            '#zoneStatus .loading-spinner'
        ];
        
        elements.forEach(selector => {
            const spinner = $(selector);
            if (show) {
                spinner.show();
            } else {
                spinner.hide();
            }
        });

        // Show/hide text elements
        if (show) {
            $('#zoneName .zone-name-text').text('Loading...').css('color', '#6c757d');
            $('#zoneStatus .status-text').text('Loading...').css('color', '#6c757d');
        } else {
            $('#zoneName .zone-name-text').css('color', '');
            $('#zoneStatus .status-text').css('color', '');
        }
    }

    showAlert(message, type = 'info') {
        const alertClass = `alert-${type}`;
        const iconMap = {
            'success': 'mdi-check-circle',
            'danger': 'mdi-alert-circle',
            'warning': 'mdi-alert',
            'info': 'mdi-information'
        };
        const icon = iconMap[type] || 'mdi-information';

        const alertHtml = `
            <div class="alert ${alertClass} alert-dismissible alert-floating" role="alert">
                <i class="mdi ${icon} me-2"></i>
                ${message}
                <button type="button" class="btn-close" data-bs-dismiss="alert" aria-label="Close"></button>
            </div>
        `;

        $('body').append(alertHtml);

        // Auto-dismiss after 5 seconds
        setTimeout(() => {
            $('.alert-floating').alert('close');
        }, 5000);
    }

    copyToClipboard(text) {
        navigator.clipboard.writeText(text).then(() => {
            this.showAlert('Copied to clipboard', 'success');
        }).catch(() => {
            // Fallback for older browsers
            const textArea = document.createElement('textarea');
            textArea.value = text;
            document.body.appendChild(textArea);
            textArea.select();
            document.execCommand('copy');
            document.body.removeChild(textArea);
            this.showAlert('Copied to clipboard', 'success');
        });
    }
}

// Initialize the zone details manager when the page loads
$(document).ready(() => {
    new CloudflareZoneDetailsManager();
});
