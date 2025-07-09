// Cloudflare Zones Management JavaScript

class CloudflareZonesManager {
    constructor() {
        this.currentAccountId = null;
        this.zones = [];
        this.accounts = [];
        this.searchTimeout = null;
        this.pageLoadingStartTime = Date.now();
        this.init();
    }

    init() {
        this.showPageLoading(true);
        this.bindEvents();
        this.loadCurrentAccount();
    }

    bindEvents() {
        // Account selection
        $('#accountSelect').on('change', (e) => {
            this.currentAccountId = e.target.value;
            if (this.currentAccountId) {
                // Update current account name
                const selectedAccount = this.accounts.find(acc => acc.id === this.currentAccountId);
                if (selectedAccount) {
                    this.currentAccountName = selectedAccount.name;
                }
                
                this.saveCurrentAccount();
                this.showCurrentAccountInfo();
                this.loadZones();
                $('#zoneStats, #zoneActions').show();
                $('#refreshBtn').prop('disabled', false);
            } else {
                this.currentAccountName = null;
                this.showCurrentAccountInfo();
                $('#zoneStats, #zoneActions').hide();
                $('#refreshBtn').prop('disabled', true);
            }
        });

        // Refresh button
        $('#refreshBtn').on('click', () => {
            if (this.currentAccountId) {
                this.showRefreshLoading(true);
                this.loadZones().finally(() => {
                    this.showRefreshLoading(false);
                });
            }
        });

        // Search functionality
        $('#searchZones').on('input', (e) => {
            clearTimeout(this.searchTimeout);
            this.searchTimeout = setTimeout(() => {
                this.filterZones(e.target.value);
            }, 300);
        });

        $('#searchBtn').on('click', () => {
            this.filterZones($('#searchZones').val());
        });

        // Create zone modal events
        $('#createZoneBtn').on('click', () => {
            this.populateCreateZoneModal();
        });

        $('#createZoneSubmit').on('click', () => {
            this.createZone();
        });

        // Zone actions
        $(document).on('click', '.view-zone-btn', (e) => {
            const zoneId = $(e.target).closest('button').data('zone-id');
            this.viewZoneDetails(zoneId);
        });

        $(document).on('click', '.pause-resume-zone-btn', (e) => {
            const zoneId = $(e.target).closest('button').data('zone-id');
            const isPaused = $(e.target).closest('button').data('is-paused');
            this.toggleZonePause(zoneId, !isPaused);
        });

        $(document).on('click', '.delete-zone-btn', (e) => {
            const zoneId = $(e.target).closest('button').data('zone-id');
            const zoneName = $(e.target).closest('button').data('zone-name');
            this.showDeleteZoneModal(zoneId, zoneName);
        });

        // Delete zone confirmation
        $('#confirmDeleteZone').on('click', () => {
            this.deleteZone();
        });

        // Copy functionality
        $(document).on('click', '.copy-btn', (e) => {
            const text = $(e.target).data('copy-text');
            this.copyToClipboard(text);
        });
    }

    async loadCurrentAccount() {
        this.showAccountSelectLoading(true);
        
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
                    
                    // Load accounts list and pre-select the current one
                    await this.loadAccounts();
                    $('#accountSelect').val(this.currentAccountId);
                    this.showCurrentAccountInfo();
                    
                    // Auto-load zones for the current account
                    this.loadZones();
                    $('#zoneStats, #zoneActions').show();
                    $('#refreshBtn').prop('disabled', false);
                    this.hidePageLoadingWithDelay();
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
                
                // Load accounts list and pre-select the current one
                await this.loadAccounts();
                $('#accountSelect').val(this.currentAccountId);
                this.showCurrentAccountInfo();
                
                // Auto-load zones for the current account
                this.loadZones();
                $('#zoneStats, #zoneActions').show();
                $('#refreshBtn').prop('disabled', false);
                this.hidePageLoadingWithDelay();
                return;
            }
        } catch (error) {
            console.error('Error parsing account from localStorage:', error);
        }

        // No current account found, just load the accounts list
        await this.loadAccounts();
        this.hidePageLoadingWithDelay();
    }

    showCurrentAccountInfo() {
        if (this.currentAccountId && this.currentAccountName) {
            $('#currentAccountDisplay').text(this.currentAccountName);
            $('#currentAccountInfo').show();
        } else {
            $('#currentAccountInfo').hide();
        }
    }

    saveCurrentAccount() {
        if (!this.currentAccountId) return;

        const selectedAccount = this.accounts.find(acc => acc.id === this.currentAccountId);
        if (!selectedAccount) return;

        try {
            // Save to localStorage
            localStorage.setItem('selectedAccount', JSON.stringify(selectedAccount));
            
            // Update server-side current account
            fetch('/api/current-account', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${localStorage.getItem('token')}`
                },
                body: JSON.stringify({
                    id: selectedAccount.id,
                    name: selectedAccount.name
                })
            }).then(response => {
                if (response.ok) {
                    console.log('Current account saved to server');
                    // Trigger event for sidebar update
                    window.dispatchEvent(new Event('accountSelected'));
                } else {
                    console.warn('Failed to save account to server');
                }
            }).catch(error => {
                console.warn('Error saving account to server:', error);
            });
        } catch (error) {
            console.error('Error saving current account:', error);
        }
    }

    async loadAccounts() {
        try {
            const response = await fetch('/api/cloudflare/accounts', {
                headers: {
                    'Authorization': `Bearer ${localStorage.getItem('token')}`
                }
            });

            if (!response.ok) {
                throw new Error('Failed to fetch accounts');
            }

            const data = await response.json();
            
            // Handle the response structure consistently
            const apiData = data.data || data;
            this.accounts = apiData.accounts || apiData || [];
            this.populateAccountSelect();
        } catch (error) {
            console.error('Error loading accounts:', error);
            this.showAlert('Error loading Cloudflare accounts', 'danger');
        } finally {
            this.showAccountSelectLoading(false);
        }
    }

    populateAccountSelect() {
        const select = $('#accountSelect');
        select.empty();
        
        if (this.accounts.length === 0) {
            select.append('<option value="">No accounts available</option>');
        } else {
            select.append('<option value="">Select Cloudflare Account...</option>');
            this.accounts.forEach(account => {
                select.append(`<option value="${account.id}">${account.name}</option>`);
            });
        }

        // Also populate create zone modal
        const createSelect = $('#createZoneAccount');
        createSelect.empty().append('<option value="">Select account...</option>');
        this.accounts.forEach(account => {
            createSelect.append(`<option value="${account.id}">${account.name}</option>`);
        });
    }

    async loadZones() {
        if (!this.currentAccountId) return;

        this.showLoading(true);
        
        try {
            const response = await fetch(`/api/cloudflare/accounts/${this.currentAccountId}/zones`, {
                headers: {
                    'Authorization': `Bearer ${localStorage.getItem('token')}`
                }
            });

            if (!response.ok) {
                throw new Error('Failed to fetch zones');
            }

            const data = await response.json();
            console.log('Zones API response:', data);
            
            // Handle the response structure: { data: { success: true, data: [...], ... }, status: "success" }
            const apiData = data.data || data;
            
            if (apiData.success) {
                this.zones = apiData.data || [];
                this.updateZoneStats();
                this.renderZones();
            } else {
                throw new Error(apiData.message || 'Unknown error loading zones');
            }
        } catch (error) {
            console.error('Error loading zones:', error);
            this.showAlert('Error loading zones', 'danger');
            this.zones = [];
            this.updateZoneStats();
            this.renderZones();
        } finally {
            this.showLoading(false);
        }
    }

    updateZoneStats() {
        const total = this.zones.length;
        const active = this.zones.filter(z => z.status === 'active').length;
        const paused = this.zones.filter(z => z.status === 'paused').length;
        const pending = this.zones.filter(z => z.status === 'pending').length;

        // Animate the numbers
        this.animateStatValue('totalZones', total);
        this.animateStatValue('activeZones', active);
        this.animateStatValue('pausedZones', paused);
        this.animateStatValue('pendingZones', pending);
    }

    renderZones() {
        const tbody = $('#zonesTableBody');
        tbody.empty();

        if (this.zones.length === 0) {
            $('#noZonesMessage').show();
            return;
        }

        $('#noZonesMessage').hide();

        this.zones.forEach(zone => {
            const row = this.createZoneRow(zone);
            tbody.append(row);
        });
    }

    createZoneRow(zone) {
        const statusBadge = this.getStatusBadge(zone.status);
        const typeBadge = this.getTypeBadge(zone.type);
        const nameServers = zone.name_servers ? zone.name_servers.slice(0, 2).join(', ') : 'N/A';
        const moreNS = zone.name_servers && zone.name_servers.length > 2 ? `... (+${zone.name_servers.length - 2} more)` : '';
        const createdDate = new Date(zone.created_on).toLocaleDateString();
        const isPaused = zone.status === 'paused';

        return `
            <tr>
                <td>
                    <a href="/CloudflareZoneDetails?id=${zone.id}" class="zone-name">
                        ${zone.name}
                    </a>
                    <br>
                    <small class="text-muted">ID: ${zone.id.substring(0, 8)}...</small>
                </td>
                <td>${statusBadge}</td>
                <td>${typeBadge}</td>
                <td>
                    <span title="${zone.name_servers ? zone.name_servers.join(', ') : 'N/A'}">
                        ${nameServers}${moreNS}
                    </span>
                </td>
                <td>${createdDate}</td>
                <td>
                    <div class="zone-actions">
                        <button class="btn btn-sm btn-outline-info view-zone-btn" 
                                data-zone-id="${zone.id}" 
                                title="View Details">
                            <i class="mdi mdi-eye"></i>
                        </button>
                        <button class="btn btn-sm ${isPaused ? 'btn-outline-success' : 'btn-outline-warning'} pause-resume-zone-btn" 
                                data-zone-id="${zone.id}" 
                                data-is-paused="${isPaused}"
                                title="${isPaused ? 'Resume' : 'Pause'} Zone">
                            <i class="mdi ${isPaused ? 'mdi-play' : 'mdi-pause'}"></i>
                        </button>
                        <button class="btn btn-sm btn-outline-danger delete-zone-btn" 
                                data-zone-id="${zone.id}" 
                                data-zone-name="${zone.name}"
                                title="Delete Zone">
                            <i class="mdi mdi-delete"></i>
                        </button>
                    </div>
                </td>
            </tr>
        `;
    }

    getStatusBadge(status) {
        const statusClasses = {
            'active': 'badge-success',
            'pending': 'badge-warning',
            'paused': 'badge-secondary',
            'initializing': 'badge-info',
            'moved': 'badge-dark',
            'deleted': 'badge-danger'
        };

        const badgeClass = statusClasses[status] || 'badge-secondary';
        return `<span class="badge ${badgeClass}">${status.charAt(0).toUpperCase() + status.slice(1)}</span>`;
    }

    getTypeBadge(type) {
        const typeClasses = {
            'full': 'badge-primary',
            'partial': 'badge-info'
        };

        const badgeClass = typeClasses[type] || 'badge-secondary';
        return `<span class="badge ${badgeClass}">${type.charAt(0).toUpperCase() + type.slice(1)}</span>`;
    }

    filterZones(searchTerm) {
        if (!searchTerm) {
            this.renderZones();
            return;
        }

        const filteredZones = this.zones.filter(zone => 
            zone.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
            zone.id.toLowerCase().includes(searchTerm.toLowerCase()) ||
            zone.status.toLowerCase().includes(searchTerm.toLowerCase())
        );

        const tbody = $('#zonesTableBody');
        tbody.empty();

        if (filteredZones.length === 0) {
            tbody.append(`
                <tr>
                    <td colspan="6" class="text-center py-4">
                        <i class="mdi mdi-magnify text-muted" style="font-size: 48px;"></i>
                        <h5 class="mt-3 text-muted">No zones match your search</h5>
                        <p class="text-muted">Try adjusting your search terms</p>
                    </td>
                </tr>
            `);
            return;
        }

        filteredZones.forEach(zone => {
            const row = this.createZoneRow(zone);
            tbody.append(row);
        });
    }

    populateCreateZoneModal() {
        // Set the selected account in the modal
        if (this.currentAccountId) {
            $('#createZoneAccount').val(this.currentAccountId);
        }
        $('#zoneName').val('');
    }

    async createZone() {
        const zoneName = $('#zoneName').val().trim();
        const accountId = $('#createZoneAccount').val();

        if (!zoneName || !accountId) {
            this.showAlert('Please fill in all required fields', 'warning');
            return;
        }

        // Basic domain validation
        const domainRegex = /^[a-zA-Z0-9][a-zA-Z0-9-]{0,61}[a-zA-Z0-9]\.[a-zA-Z]{2,}$/;
        if (!domainRegex.test(zoneName)) {
            this.showAlert('Please enter a valid domain name', 'warning');
            return;
        }

        this.showCreateZoneLoading(true);

        try {
            const response = await fetch(`/api/cloudflare/accounts/${accountId}/zones`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${localStorage.getItem('token')}`
                },
                body: JSON.stringify({
                    name: zoneName,
                    account_id: accountId
                })
            });

            const data = await response.json();

            if (!response.ok) {
                throw new Error(data.message || 'Failed to create zone');
            }

            this.showAlert(`Zone "${zoneName}" created successfully`, 'success');
            $('#createZoneModal').modal('hide');
            this.loadZones(); // Refresh the zones list
        } catch (error) {
            console.error('Error creating zone:', error);
            this.showAlert(error.message || 'Error creating zone', 'danger');
        } finally {
            this.showCreateZoneLoading(false);
        }
    }

    async toggleZonePause(zoneId, shouldPause) {
        try {
            const response = await fetch(`/api/cloudflare/zones/${zoneId}`, {
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
            this.loadZones(); // Refresh the zones list
        } catch (error) {
            console.error('Error updating zone:', error);
            this.showAlert('Error updating zone', 'danger');
        }
    }

    viewZoneDetails(zoneId) {
        // Redirect to zone details page
        window.location.href = `/CloudflareZoneDetails?id=${zoneId}`;
    }

    showDeleteZoneModal(zoneId, zoneName) {
        $('#deleteZoneName').text(zoneName);
        $('#confirmDeleteZone').data('zone-id', zoneId);
        $('#deleteZoneModal').modal('show');
    }

    async deleteZone() {
        const zoneId = $('#confirmDeleteZone').data('zone-id');
        if (!zoneId) return;

        this.showDeleteZoneLoading(true);

        try {
            const response = await fetch(`/api/cloudflare/zones/${zoneId}`, {
                method: 'DELETE',
                headers: {
                    'Authorization': `Bearer ${localStorage.getItem('token')}`
                }
            });

            if (!response.ok) {
                throw new Error('Failed to delete zone');
            }

            this.showAlert('Zone deleted successfully', 'success');
            $('#deleteZoneModal').modal('hide');
            this.loadZones(); // Refresh the zones list
        } catch (error) {
            console.error('Error deleting zone:', error);
            this.showAlert('Error deleting zone', 'danger');
        } finally {
            this.showDeleteZoneLoading(false);
        }
    }

    showLoading(show) {
        if (show) {
            $('#loadingIndicator').show();
            $('#zonesTableBody').hide();
            $('#noZonesMessage').hide();
            this.showStatsLoading(true);
        } else {
            $('#loadingIndicator').hide();
            $('#zonesTableBody').show();
            this.showStatsLoading(false);
        }
    }

    showCreateZoneLoading(show) {
        const spinner = $('#createZoneSpinner');
        const button = $('#createZoneSubmit');
        
        if (show) {
            spinner.show();
            button.prop('disabled', true).text('Creating...');
        } else {
            spinner.hide();
            button.prop('disabled', false).text('Add Zone');
        }
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
            this.showAlert('Failed to copy to clipboard', 'danger');
        });
    }

    // Enhanced loading and UX methods
    showPageLoading(show) {
        const overlay = $('#pageLoadingOverlay');
        const body = $('body');
        
        if (show) {
            body.addClass('loading-active');
            overlay.show();
            // Ensure overlay is on top of everything
            overlay.css('z-index', '999999');
        } else {
            // Smooth fade out
            overlay.fadeOut(300, () => {
                body.removeClass('loading-active');
            });
        }
    }

    hidePageLoadingWithDelay() {
        const minLoadingTime = 1000; // Minimum 1 second loading time
        const elapsedTime = Date.now() - this.pageLoadingStartTime;
        const remainingTime = Math.max(0, minLoadingTime - elapsedTime);
        
        setTimeout(() => {
            this.showPageLoading(false);
        }, remainingTime);
    }

    showAccountSelectLoading(show) {
        const select = $('#accountSelect');
        const loading = $('#accountSelectLoading');
        
        if (show) {
            select.prop('disabled', true);
            loading.show();
        } else {
            select.prop('disabled', false);
            loading.hide();
        }
    }

    showRefreshLoading(show) {
        const spinner = $('.refresh-spinner');
        const button = $('#refreshBtn');
        
        if (show) {
            spinner.show();
            button.prop('disabled', true);
        } else {
            spinner.hide();
            button.prop('disabled', false);
        }
    }

    showStatsLoading(show) {
        const elements = ['#totalZones', '#activeZones', '#pausedZones', '#pendingZones'];
        
        elements.forEach(selector => {
            const element = $(selector);
            const spinner = element.find('.loading-spinner');
            const value = element.find('.stat-value');
            
            if (show) {
                spinner.show();
                value.hide();
            } else {
                spinner.hide();
                value.show();
            }
        });
    }

    animateStatValue(elementId, targetValue) {
        const element = $(`#${elementId} .stat-value`);
        const currentValue = parseInt(element.text()) || 0;
        
        if (currentValue === targetValue) return;
        
        const duration = 800;
        const steps = 20;
        const stepSize = (targetValue - currentValue) / steps;
        const stepDuration = duration / steps;
        
        let currentStep = 0;
        
        const interval = setInterval(() => {
            currentStep++;
            const newValue = Math.round(currentValue + (stepSize * currentStep));
            
            if (currentStep >= steps) {
                element.text(targetValue);
                clearInterval(interval);
            } else {
                element.text(newValue);
            }
        }, stepDuration);
    }
}

// Initialize the zones manager when the page loads
$(document).ready(() => {
    new CloudflareZonesManager();
});
