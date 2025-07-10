/* Fix for missing tooltip functionality */
if (typeof $.fn.tooltip !== 'function') {
    $.fn.tooltip = function(options) {
        console.log('Tooltip polyfill used');
        // Simple no-op polyfill to prevent errors
        return this;
    };
}

/* Docker Cloudflare Tunnels JavaScript */
// This script will run once the main docker-cloudflare-tunnels.js is loaded
// and will ensure that any missing dependencies are addressed

document.addEventListener('DOMContentLoaded', function() {
    // Fallback for missing jQuery plugins
    if (typeof $.fn.tooltip !== 'function') {
        $.fn.tooltip = function() { return this; };
        console.log("Added tooltip polyfill");
    }
    
    // Patch account manager to handle errors
    setTimeout(function() {
        if (window.AccountManager) {
            const originalSetAccount = AccountManager.setAccount;
            AccountManager.setAccount = function(account) {
                try {
                    return originalSetAccount.call(AccountManager, account);
                } catch (err) {
                    console.error("Error in setAccount:", err);
                    // Still update UI elements directly
                    document.getElementById('current-account-name').textContent = 
                        account ? account.name : 'No account selected';
                    document.getElementById('current-account-id').textContent = 
                        account ? account.id : 'Select an account from the list';
                    return account;
                }
            };
            
            const originalUpdateTooltips = AccountManager.updateTooltips;
            AccountManager.updateTooltips = function() {
                try {
                    return originalUpdateTooltips.call(AccountManager);
                } catch (err) {
                    console.error("Error in updateTooltips:", err);
                    // Do nothing, tooltips are non-essential
                }
            };
            
            console.log("Patched AccountManager for error handling");
        }
    }, 500);
});
